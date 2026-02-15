package lexer

import (
	"os"
	"sort"
	"strconv"

	"github.com/olekukonko/tablewriter"
)

type lexTab struct {
	token Token
	line  string
}

func PrintTable(tokens [][]Token) {
	var entries []lexTab
	var cache = make(map[string]lexTab)

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

	for i, tkns := range tokens {
		for _, tkn := range tkns {
			if token, found := cache[string(tkn.L)]; found {
				tab := lexTab{
					token: tkn,
					line:  token.line + "," + strconv.Itoa(i+1),
				}
				cache[string(tkn.L)] = tab
				continue
			}
			tab := lexTab{token: tkn, line: strconv.Itoa(i + 1)}
			cache[string(tkn.L)] = tab
		}
	}

	for _, v := range cache {
		entries = append(entries, v)
	}

	sort.Slice(entries, func(i, j int) bool {
		return entries[i].token.V < entries[j].token.V
	})

	for i, lex := range entries {
		if lex.token.T == 4 {
			identifiers.Append([]string{
				string(lex.token.L),
				strconv.Itoa(int(lex.token.V)),
				lex.line,
			})
		}

		if lex.token.T == 6 {
			constants.Append([]string{
				string(lex.token.L),
				strconv.Itoa(int(lex.token.V)),
				lex.line,
			})
		}

		// "No.", "Line", "Token", "Type", "Code",
		globalTable.Append([]string{
			strconv.Itoa(i + 1),
			lex.line,
			string(lex.token.L),
			strconv.Itoa(int(lex.token.T)),
			strconv.Itoa(int(lex.token.V)),
		})
	}

	identifiers.Render()
	constants.Render()
	globalTable.Render()

}
