package main

import (
	"bytes"
	"crypto/sha256"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"time"
)

var src, dst string

func validFolder(s string) (bool, error) {
	if s == "" {
		return false, errors.New("no path provided")
	}
	var fileInfo, err = os.Stat(s)
	if err != nil {
		return false, fmt.Errorf("path '%v' does not exist", s)
	}

	if !fileInfo.IsDir() {
		return false, fmt.Errorf("path '%v' is not a valid folder", s)
	}

	return true, nil
}

func main() {
	flag.StringVar(&src, "src", "", "Source folder containing images.")
	flag.StringVar(&dst, "dst", "", "Destination folder for the copied images images.")

	flag.Parse()

	timestamp := time.Now().UTC().Unix()
	logFileName := filepath.Join(".", fmt.Sprintf("hashr_%v.log", timestamp))
	logFile, err := os.Create(logFileName)
	if err != nil {
		log.Fatalf("error opening log file: %v", err)
	}
	defer logFile.Close()

	log.SetOutput(io.MultiWriter(os.Stdout, logFile))
	log.Printf("log file for this process available at: %s", logFileName)

	if ok, err := validFolder(src); !ok {
		log.Fatalf("invalid argument: please specify a source folder "+
			"containing images to copy: %v", err)
	}

	if ok, err := validFolder(dst); !ok {
		log.Fatalf("invalid argument: please specify a destination folder "+
			"where to copy images: %v", err)
	}

	log.Printf("The source folder of the images is: %v\n", src)
	log.Printf("The destination folder of the images is: %v\n", dst)

	err = filepath.Walk(src, func(srcPath string, info os.FileInfo, err error) error {
		if err != nil {
			log.Fatalf("accessing file %v: %v", srcPath, err)
		}

		if !info.IsDir() {
			log.Printf("file: %s\n", srcPath)
			f, err := os.Open(srcPath)
			if err != nil {
				log.Fatalf("opening file %v: %v", srcPath, err)
			}
			defer f.Close()

			buf := new(bytes.Buffer)
			_, err = io.Copy(buf, f)
			if err != nil {
				log.Fatalf("reading file %v to buffer: %v", srcPath, err)
			}

			h := sha256.New()
			if _, err := io.Copy(h, bytes.NewReader(buf.Bytes())); err != nil {
				log.Fatalf("hashing file %v: %v", srcPath, err)
			}

			destPath := fmt.Sprintf("%v%x%v", dst, h.Sum(nil), filepath.Ext(srcPath))
			_, err = os.Stat(destPath)
			if err != nil {
				df, err := os.OpenFile(destPath, os.O_RDWR|os.O_CREATE, 0644)
				if err != nil {
					log.Fatalf("creating file %v: %v", destPath, err)
				}

				_, err = io.Copy(df, bytes.NewReader(buf.Bytes()))
				if err != nil {
					log.Fatalf("copying to destination file %v: %v", destPath, err)
				}
				log.Printf("created file %v", destPath)
			} else {
				log.Printf("skipping file %v: %v", srcPath, err)
			}
		}
		return nil
	})
	if err != nil {
		log.Fatalf("walking folder structure: %v", err)
	}

}
