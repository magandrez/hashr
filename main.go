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

	err := filepath.Walk(src, func(srcPath string, info os.FileInfo, err error) error {
		if err != nil {
			return fmt.Errorf("accessing file: %v : %v", srcPath, err)
		}

		if !info.IsDir() {
			log.Printf("file: %s\n", srcPath)
			f, err := os.Open(srcPath)
			if err != nil {
				log.Fatal(err)
			}
			defer f.Close()

			buf := new(bytes.Buffer)
			_, err = io.Copy(buf, f)
			if err != nil {
				return fmt.Errorf("reading contents of file to buffer: %w", err)
			}

			h := sha256.New()
			if _, err := io.Copy(h, bytes.NewReader(buf.Bytes())); err != nil {
				log.Fatal(err)
			}

			destPath := fmt.Sprintf("%v%x%v", dst, h.Sum(nil), filepath.Ext(srcPath))
			_, err = os.Stat(destPath)
			if err != nil {
				df, err := os.OpenFile(destPath, os.O_RDWR|os.O_CREATE, 0644)
				if err != nil {
					log.Fatal(err)
				}

				_, err = io.Copy(df, bytes.NewReader(buf.Bytes()))
				if err != nil {
					log.Fatal(err)
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
