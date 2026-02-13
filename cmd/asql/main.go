package main

import (
	"asql/internal/scanner"
	"bufio"
	"log"
	"log/slog"
	"os"
	"strconv"
	"strings"

	"github.com/olekukonko/tablewriter"
)

func run() error {
	lines := [][]string{}

	table := tablewriter.NewWriter(os.Stdout)
	table.Header([]string{
		"No.", "Line", "Token", "Type", "Code",
	})

	reader, err := cmd()
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
		tokens := scanner.Tokenize(rawText)
		lines = append(lines, tokens)
	}

	line := 1
	lexer := scanner.NewLexer()
	for i, tokens := range lines {
		for _, token := range tokens {
			t, err := lexer(token)
			if err != nil {
				slog.Error(err.Error(), "line", strconv.Itoa(i+1))
				continue
			}
			table.Append([]string{
				strconv.Itoa(line),
				strconv.Itoa(i + 1),
				string(t.L),
				strconv.Itoa(int(t.T)),
				strconv.Itoa(int(t.V)),
			})
			line++
		}
	}

	table.Render()
	return nil
}

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}
