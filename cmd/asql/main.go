package main

import (
	"asql/internal/scanner"
	"bufio"
	"log"
	"log/slog"
	"os"
	"sort"
	"strconv"
	"strings"

	"github.com/olekukonko/tablewriter"
)

type tab struct {
	Key  scanner.Keyword
	Line string
}

func run() error {

	lexTable := make(map[string]tab)
	lines := [][]string{}

	globalTable := tablewriter.NewWriter(os.Stdout)
	identifiers := tablewriter.NewWriter(os.Stdout)
	constants := tablewriter.NewWriter(os.Stdout)

	globalTable.Header([]string{
		"No.", "Line", "Token", "Type", "Code",
	})

	identifiers.Header([]string{
		"Identifer", "Value", "Line",
	})

	constants.Header([]string{
		"Constant", "Value", "Line",
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

	lexer := scanner.NewLexer()
	for i, tokens := range lines {
		for _, token := range tokens {

			tkn, err := lexer(token)

			if err != nil {
				slog.Error(err.Error(), "line", strconv.Itoa(i+1))
				continue
			}

			l, found := lexTable[string(tkn.L)]
			if found {
				l.Line = l.Line + "," + strconv.Itoa(i+1)
				lexTable[string(tkn.L)] = l
				continue
			}

			lexTable[string(tkn.L)] = tab{
				tkn,
				strconv.Itoa(i + 1),
			}

		}
	}

	var entries []tab

	for _, lex := range lexTable {
		entries = append(entries, lex)
	}

	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Key.V < entries[j].Key.V
	})

	for _, lex := range entries {
		if lex.Key.T == 4 {
			identifiers.Append([]string{
				string(lex.Key.L),
				strconv.Itoa(int(lex.Key.V)),
				lex.Line,
			})
		}

		if lex.Key.T == 6 {
			constants.Append([]string{
				string(lex.Key.L),
				strconv.Itoa(int(lex.Key.V)),
				lex.Line,
			})
		}
	}

	identifiers.Render()
	constants.Render()

	return nil
}

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}
