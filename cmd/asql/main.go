package main

import (
	"asql/internal/scannerDML"
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

	defer reader.Close()
	scanner := bufio.NewScanner(reader)

	for scanner.Scan() {
		fileLine := scanner.Text()
		lines[i] = scannerdml.Tokenize(fileLine)
		i++
	}

	lex := scannerdml.NewLexer()
	for _, token := range lines[0] {
		fmt.Println(token)
		l := lex[token]
		fmt.Println(l)
	}
	return nil
}

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}
