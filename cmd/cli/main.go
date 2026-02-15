package main

import (
	"asql/internal/lexer"
	"bufio"
	"fmt"
	"log"
	"log/slog"
	"strings"
)

type state struct {
	err error
}

func run() error {

	// lexTable := make(map[string]lexTab)
	lines := [][]string{}
	tokenStream := [][]lexer.Token{}
	state := state{} // global state of all stages

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
	/*
		we should return a lex table an example:
		1:[TOKENS]
		2:[TOKENS]
		3:[TOKENS]
		The order and the lines must be consistent
		to the original source
		[][]string -> [][]tokens
	*/

	lex := &lexer.TokenStream{Lexer: lexer.NewLexer()}
	for i, line := range lines {
		lex.Tokens = line
		tokens := []lexer.Token{}
		for result := range lex.GenTokenStream() {
			if result.Err != nil {
				state.err = result.Err
				slog.Error(fmt.Sprintf("line %d: %v", i+1, result.Err))
				continue
			}
			fmt.Println(result.Token)
			tokens = append(tokens, result.Token)
		}
		tokenStream = append(tokenStream, tokens)
	}

	lexer.PrintTable(tokenStream)
	return nil
}

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}
