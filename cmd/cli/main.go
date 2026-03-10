package main

import (
	"asql/internal/lexer"
	"asql/internal/parser"
	"fmt"
	"log"
)

func run() error {

	reader, err := readFile()
	if err != nil {
		return err
	}

	lines := readFromFile(reader)

	// Lexer: satage 1
	tokenStream, lexerErrors := lexer.Lexer(lines)
	// lexer.PrintLexer(tokenStream)

	for _, state := range lexerErrors {
		fmt.Println("[Lexer]", state.Err.Error(), "line", state.Line, "token", state.L)
	}

	// Parser: stage 2
	// parser.Parser(tokenStream)
	p := parser.NewParser(tokenStream)
	state := p.Parse()
	if state != nil {
		fmt.Println("[Parser]", state.Error(),
			"line: ", state.Token.Line,
			"Token:", state.Token.L,
			"Value:", state.Token.V)
	}
	return err
}

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}
