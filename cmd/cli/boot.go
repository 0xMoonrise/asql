package main

import (
	"bufio"
	"io"
	"strings"

	"github.com/0xMoonrise/asql/internal/lexer"
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
