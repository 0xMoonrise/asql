package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/0xMoonrise/asql/internal/lexer"
	"github.com/0xMoonrise/asql/internal/parser"
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
	parser := parser.NewParser(tokenStream)
	if err := parser.Parse(); err != nil {
		fmt.Println("[Parser]", parser.State.Message.Error())
		line := parser.State.Token.Line
		fmt.Printf("[Line: %d] %s ", line, strings.Join(source[line-1], " "))
	}
	return err
}

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}
