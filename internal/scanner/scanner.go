package scanner

import (
	"asql/internal/utils"
	"errors"
	"regexp"
	"strings"
	"unicode"
)

type Keyword struct {
	L lexeme
	V value
	T typ
}

type lexer map[string]Keyword

func NewTable() lexer {
	lex := make(lexer)

	keywords := []Keyword{
		// keywords (1)
		{"SELECT", __select, 1},
		{"FROM", __from, 1},
		{"WHERE", __where, 1},
		{"IN", __in, 1},
		{"AND", __and, 1},
		{"OR", __or, 1},
		{"CREATE", __create, 1},
		{"TABLE", __table, 1},
		{"CHAR", __char, 1},
		{"NUMERIC", __numeric, 1},
		{"NOT", __not, 1},
		{"NULL", __null, 1},
		{"CONSTRAINT", __constraint, 1},
		{"KEY", __key, 1},
		{"PRIMARY", __primary, 1},
		{"FOREIGN", __foreign, 1},
		{"REFERENCES", __references, 1},
		{"INSERT", __insert, 1},
		{"INTO", __into, 1},
		{"VALUES", __values, 1},

		// delimitators (5)
		{",", comma, 5},
		{".", dot, 5},
		{"(", lparentheses, 5},
		{")", rparentheses, 5},
		{"'", apostrophe, 5},

		// Constants (6)
		{"d", numeric, 6},
		{"a", alpha, 6},

		// Operators (7)
		{"+", plus, 7},
		{"-", minus, 7},
		{"*", times, 7},
		{"/", divition, 7},

		// Relations (8)
		{">", gt, 8},
		{"<", lt, 8},
		{"=", eq, 8},
		{">=", ge, 8},
		{"<=", le, 8},
	}

	for _, kw := range keywords {
		lex[string(kw.L)] = Keyword{
			L: kw.L,
			V: kw.V,
			T: kw.T,
		}
	}

	return lex
}

/*
   Tokenizer procedure:
   1. Keywords
   2. Delimitators
   3. Constants
   4. Operators
   5. Relations

   This could be inefficient because it needs 5 rounds to fill
   the table and/or determine all the tokens in the source code.

   But its scalable for more kinds of rules.
*/
// rules of Tokenizer

var (
	keywords     = `\b[A-Za-z_][A-Za-z0-9_]*\b` // kewords and identifiers
	constant     = `\b\d+\b|'[^']*'`            // Strings and numbers
	delimitators = `[,()]`
	relations    = `>=|<=|!=|=|>|<`
	noLexer      = `\W` // Match any char that is not in the lexer
)

func Tokenize(input string) []string {
	rules := []string{
		keywords,
		constant,
		delimitators,
		relations,
		noLexer,
	}

	lexerRule := strings.Join(rules, "|")
	re := regexp.MustCompile(lexerRule)

	tokens := utils.Filter(
		utils.Map(re.FindAllString(input, -1), strings.TrimSpace),
		func(s string) bool { return s != "" },
	)

	return tokens
}

func isNumeric(s string) bool {
	for _, r := range s {
		if !unicode.IsDigit(r) {
			return false
		}
	}
	return len(s) > 0
}

// Apply the criteria for tokens
func NewLexer() func(t string) (token Keyword, err error) {
	// Values for dynamic tables
	var indentifiers int = 401
	var constants int = 600
	var cache = make(map[string]Keyword)

	return func(t string) (token Keyword, err error) {
		lexical := NewTable()
		token, found := lexical[t]

		if found {
			return
		}

		token, found = cache[t]
		if found {
			return
		}

		token = Keyword{}
		token.L = lexeme(t)

		switch {
		case regexp.MustCompile(constant).MatchString(t):
			token.T = 6
			token.V = value(constants)
			constants++
			cache[t] = token
			return
		case regexp.MustCompile(keywords).MatchString(t):
			token.T = 4
			token.V = value(indentifiers)
			indentifiers++
			cache[t] = token
			return
		default:
			err = errors.New("Unknown symbol")
		}

		return
	}
}
