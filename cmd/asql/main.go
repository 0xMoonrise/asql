package main

import (
	"asql/internal/scanner"
	"bufio"
	"log"
	"log/slog"
	"os"
	"strconv"

	"github.com/olekukonko/tablewriter"
)

func run() error {

	lines := make(map[int][]string)
	i := 0

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
		lines[i] = scanner.Tokenize(rawText)
		i++
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
