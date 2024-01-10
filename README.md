# hashr

hashr is a small script to consolidate files scattered over a folder structure
into a single folder, using their SHA256 signature as filename.

## The problem

Over the years, I have accumulated several dozens of gigabytes of images, both
raw and processed finals. I have stored all of them in different hard drives,
which, after a period of time, for one reason or another have failed either due
to a human error (deleting partition tables) or plain disk failures, leaving me
with the task of recovering them.

Due to this, I have accumulated fragmented and mostly overlapping versions of
the same images, sometimes with the same file names, sometimes with names given
by the recovery tool used, spread over different folder structures.

I decided I would like to consolidate all the images and getting rid of the
duplicates in an automated manner.

TL;DR;: I have accumulated a sparsely populated folder tree of images. Some of
them might be copies of other files in a different folder, some of them might
have different names. I have a folder tree containing image files of various
formats with possibly n copies that might or might not have the same name that I
would like to consolidate.

## A solution

This solution is just another de-duplication tool. I rolled my own as an
exercise but also due to laziness: I did not want to learn the interfaces of
other tool(s) in order to solve the specifics of the problem at hand. Below, an
assorted list of tools that might help:

- find(1) might get you a long way.
- [fdupes](https://github.com/adrianlopezroche/fdupes)
- [exiv2](https://exiv2.org/) is a metadata library tool that might be handy on
  these type of cases, it is well supported and functional.

This is a simple CLI tool that traverses a folder structure and
calculates the SHA256 signature of each file and copies the file in a
destination folder with its calculated signature as its filename. Thus:

1. Relying on the contents of the file to distinguish exact copies from unique
  versions.
2. Getting rid of duplicates as collisions (files with the exact same SHA256
   signature) won't be copied into the same folder.
3. Reconciling a sparsely populated folder structure under a single folder.

## Development

Requirements:

- [Go](https://go.dev/) ^1.11

### Usage

Download, compile and use binary for your architecture like:

```sh
$ git clone https://github.com/magandrez/hashr && cd hashr && go run . -src [folder] -dst [folder]
```

See help with `go run . --help`

