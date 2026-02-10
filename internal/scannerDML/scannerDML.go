package scannerdml

import "regexp"

type keyword struct {
	l lexeme
	s symbol
	v value
	t typ
}

type lexer map[string]keyword

func NewLexer() lexer {
	lex := make(lexer)

	keywords := []keyword{
		// keywords (1)
		{"SELECT", 's', __select, 1},
		{"FROM", 'f', __from, 1},
		{"WHERE", 'w', __where, 1},
		{"IN", 'n', __in, 1},
		{"AND", 'y', __and, 1},
		{"OR", 'o', __or, 1},
		{"CREATE", 'c', __create, 1},
		{"TABLE", 't', __table, 1},
		{"CHAR", 'h', __char, 1},
		{"NUMERIC", 'u', __numeric, 1},
		{"NOT", 'e', __not, 1},
		{"NULL", 'g', __null, 1},
		{"CONSTRAINT", 'b', __constraint, 1},
		{"KEY", 'k', __key, 1},
		{"PRIMARY", 'p', __primary, 1},
		{"FOREIGN", 'j', __foreign, 1},
		{"REFERENCES", 'l', __references, 1},
		{"INSERT", 'm', __insert, 1},
		{"INTO", 'q', __into, 1},
		{"VALUES", 'v', __values, 1},
		// delimitators (5)
		{",", noSym, comma, 5},
		{".", noSym, dot, 5},
		{"(", noSym, lparentheses, 5},
		{")", noSym, rparentheses, 5},
		{"'", noSym, apostrophe, 5},
		// Constants (6)
		{"d", noSym, numeric, 6},
		{"a", noSym, alpha, 6},
		// Operators (7)
		{"+", noSym, plus, 7},
		{"-", noSym, minus, 7},
		{"*", noSym, times, 7},
		{"/", noSym, divition, 7},
		// Relations (8)
		{">", noSym, gt, 8},
		{"<", noSym, lt, 8},
		{"=", noSym, eq, 8},
		{">=", noSym, ge, 8},
		{"<=", noSym, le, 8},
	}

	for _, kw := range keywords {
		lex[string(kw.l)] = keyword{
			l: kw.l,
			s: kw.s,
			v: kw.v,
			t: kw.t,
		}
	}

	return lex
}

var tokenRegex = regexp.MustCompile(
	`'|` +
		`>=|<=|<>|!=|[<>=]|` +
		`[+\-*/]|` +
		`[(),.]|` +
		`[a-zA-Z0-9_]*|` +
		`\d+|` +
		`\s+`,
)

func Tokenize(input string) []string {
	matches := tokenRegex.FindAllString(input, -1)
	var tokens []string
	for _, token := range matches {
		if token != " " && token != "\n" && token != "\t" && token != "\r" {
			tokens = append(tokens, token)
		}
	}

	return tokens
}
