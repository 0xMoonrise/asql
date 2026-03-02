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

func tableLexer(tokens []lexer.Token) gin.H {
	var entries []lexTab
	var cache = make(map[string]lexTab)

	for _, token := range tokens {
		if current, found := cache[string(token.L)]; found {
			tab := lexTab{
				token: token,
				line:  current.line + "," + strconv.Itoa(token.Line),
			}
			cache[string(token.L)] = tab
			continue
		}

		tab := lexTab{token: token, line: strconv.Itoa(token.Line)}
		cache[string(token.L)] = tab
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

func (app *app) root(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", gin.H{
		"title": "main",
	})
}

func (app *app) run(c *gin.Context) {
	formText := c.PostForm("text")

	if len(formText) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "No source code were provided",
		})
		return
	}

	lines := [][]string{}
	for str := range strings.SplitSeq(formText, "\n") {
		tokens := lexer.Tokenize(str)
		lines = append(lines, tokens)
	}

	tokenStream, lexerErrors := lexer.Lexer(lines)
	errors := []string{}

	for _, err := range lexerErrors {
		errors = append(errors, fmt.Sprintf("Line %d: %v", err.Line, err.Err.Error()))
	}

	data := tableLexer(tokenStream)
	if len(errors) > 0 {
		data["Errors"] = errors
	}

	c.JSON(http.StatusOK, data)
}
