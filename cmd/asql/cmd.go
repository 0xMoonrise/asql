package main

import (
	"flag"
	"io"
	"os"
)

func cmd() (reader io.ReadCloser, err error) {
	filePath := flag.String("f", "", "Path file input")
	flag.Parse()

	if *filePath != "" {
		file, err := os.Open(*filePath)
		if err != nil {
			return nil, err
		}

		reader = file
	} else {
		reader = os.Stdin
	}

	return
}
