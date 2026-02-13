package main

import (
	"asql/internal/scanner"
	"bufio"
	"fmt"
	"log"
)

func run() error {

	lines := make(map[int][]string)
	i := 0

	reader, err := cmd()
	if err != nil {
		return err
	}
	lexScanner := scanner.Lexer()
	defer reader.Close()
	buffer := bufio.NewScanner(reader)

	for buffer.Scan() {
		fileLine := buffer.Text()
		lines[i] = scanner.Tokenize(fileLine)
		i++
	}

	for i = range len(lines) {
		token := lines[i]
		for _, tkn := range token {
			t, err := lexScanner(tkn)
			if err != nil {
				return err
			}
			fmt.Println(t)
		}
	}

	return nil
}

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}
