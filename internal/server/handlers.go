package server

import (
	"fmt"
	"net/http"
	"sort"
	"strconv"
	"strings"

	"asql/internal/lexer"
	"github.com/gin-gonic/gin"
)

type lexTab struct {
	token lexer.Token
	line  string
}

func tableLexer(tokens [][]lexer.Token) gin.H {
	var entries []lexTab
	var cache = make(map[string]lexTab)

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

	var identifiersRows [][]string
	var constantsRows [][]string
	var globalRows [][]string

	for i, lex := range entries {
		if lex.token.T == 4 {
			identifiersRows = append(identifiersRows, []string{
				string(lex.token.L),
				strconv.Itoa(int(lex.token.V)),
				lex.line,
			})
		}

		if lex.token.T == 6 {
			constantsRows = append(constantsRows, []string{
				string(lex.token.L),
				strconv.Itoa(int(lex.token.V)),
				lex.line,
			})
		}

		globalRows = append(globalRows, []string{
			strconv.Itoa(i + 1),
			lex.line,
			string(lex.token.L),
			strconv.Itoa(int(lex.token.T)),
			strconv.Itoa(int(lex.token.V)),
		})
	}

	return gin.H{
		"IdentifierHeaders": []string{"Identifier", "Value", "Line"},
		"Identifiers":       identifiersRows,
		"ConstantHeaders":   []string{"Constant", "Value", "Line"},
		"Constants":         constantsRows,
		"GlobalHeaders":     []string{"No.", "Line", "Token", "Type", "Code"},
		"GlobalTable":       globalRows,
	}
}

func (app *app) health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "OK!",
	})
}

func (app *app) lexer(c *gin.Context) {

	formText := c.PostForm("text")
	c.Header("HX-Reswap", "innerHTML")

	if len(formText) == 0 {
		c.HTML(http.StatusOK, "error.html", gin.H{
			"message": "No source code were provided",
		})
		return
	}

	lines := [][]string{}
	for str := range strings.SplitSeq(formText, "\n") {
		tokens := lexer.Tokenize(str)
		lines = append(lines, tokens)
	}

	tokenStream := [][]lexer.Token{}
	errors := []string{}
	lex := &lexer.TokenStream{Lexer: lexer.NewLexer()}

	for i, line := range lines { // <- Entry point
		lex.Tokens = line
		tokens := []lexer.Token{}
		for result := range lex.GenTokenStream() {
			if result.Err != nil {
				errors = append(errors, fmt.Sprintf("Line %d: %v", i+1, result.Err))
				continue
			}
			tokens = append(tokens, result.Token)
		}
		tokenStream = append(tokenStream, tokens)

	}

	data := tableLexer(tokenStream)
	if len(errors) > 0 {
		data["Errors"] = errors
	}
	c.HTML(http.StatusOK, "tables.html", data)
}

func (app *app) root(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", gin.H{
		"title": "main",
	})
}
