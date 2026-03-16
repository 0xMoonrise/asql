package main

import (
	"fmt"
	"github.com/0xMoonrise/asql/internal/lexer"
	"github.com/0xMoonrise/asql/internal/parser"
	"log"
)

func run() error {

	reader, err := readFile()
	if err != nil {
		return err
	}

	source := readFromFile(reader)

	// Lexer: satage 1
	tokenStream, lexerErrors := lexer.Lexer(source)
	// lexer.PrintLexer(tokenStream)

	for _, state := range lexerErrors {
		fmt.Println("[Lexer]", state.Err.Error(), "line", state.Line, "token", state.L)
	}

	// Parser: stage 2
	state := parser.NewParser(tokenStream).Parse()
	if state != nil {
		fmt.Println("[Parser]", state.Error())
	}
	return err
}

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}
