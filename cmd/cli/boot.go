package main

import (
	"asql/internal/lexer"
	"bufio"
	"io"
	"strings"
)

func readFromFile(reader io.ReadCloser) (lines [][]string) {
	buffer := bufio.NewScanner(reader)

	for buffer.Scan() {
		rawText := buffer.Text()
		if strings.HasPrefix(rawText, "#") { // Check special cace for shebang #!
			continue
		}

		tokens := lexer.Tokenize(rawText)
		lines = append(lines, tokens)
	}
	return
}
