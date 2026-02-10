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

	for i, token := range lines[0] {
		fmt.Println(i, "token:", token)
		fmt.Println(scannerdml.Lexer(token))
	}
	return nil
}

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}
