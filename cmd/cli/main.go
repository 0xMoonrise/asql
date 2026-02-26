package main

import (
	"asql/internal/lexer"
	"bufio"
	"fmt"
	"log"
	"strings"
)

func run() error {

	lines := [][]string{}
	reader, err := readFile()
	if err != nil {
		return err
	}

	defer reader.Close()
	buffer := bufio.NewScanner(reader)

	for buffer.Scan() {
		rawText := buffer.Text()
		if strings.HasPrefix(rawText, "#") { // Check special cace for shebang #!
			continue
		}

		tokens := lexer.Tokenize(rawText)
		lines = append(lines, tokens)
	}

	//Lexer: satage 1
	tokenStream, lexerErrors := lexer.Lexer(lines)
	lexer.PrintLexer(tokenStream)

	for _, state := range lexerErrors {
		fmt.Println("[Lexer]", state.Err, "line", state.Line, "token", state.L)
	}

	return err
}

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}
