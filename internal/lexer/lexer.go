package lexer

import (
	"asql/internal/utils"
	"errors"
	"iter"
	"regexp"
	"strings"
)

type Token struct {
	L lexeme
	V value
	T typ
}

type lexer map[string]Token

type TokenStream struct {
	Tokens []string
	Lexer  func(t string) (Token, error)
}

func NewTable() lexer {
	lex := make(lexer)

	Tokens := []Token{
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

		// Dynamics do not need to be declared
		// Constants (6)
		// {"d", numeric, 6},
		// {"a", alpha, 6},

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

	for _, tkn := range Tokens {
		lex[string(tkn.L)] = Token{
			L: tkn.L,
			V: tkn.V,
			T: tkn.T,
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
	keywords     = `[A-Za-z_][A-Za-z0-9_#]*` // kewords and identifiers
	constant     = `\b\d+\b|'[^']*'`         // Strings and numbers
	delimitators = `[\.,()]`
	relations    = `>=|<=|!=|=|>|<`
	spaces       = `\s+`
	noLexer      = `[^\w\s.,()#'>=<!]` // Match any char that is not in the lexer
)

func Tokenize(input string) []string {
	rules := []string{
		keywords,
		constant,
		delimitators,
		spaces,
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

// Apply the criteria for tokens
func NewLexer() func(t string) (Token, error) {
	// Values for dynamic tables
	var indentifiers int = 401
	var constants int = 600
	var cache = make(map[string]Token)
	var lexical lexer = NewTable()
	// Rules for lexer
	reConst := regexp.MustCompile(constant)
	reKws := regexp.MustCompile(keywords)

	return func(t string) (Token, error) {
		// check the static lexcial table
		key := strings.ToUpper(t)
		if token, found := lexical[key]; found {
			return token, nil
		}
		// search the cache
		if token, found := cache[t]; found {
			return token, nil
		}
		// and then it could be a non-lexical character
		// or a valid identifier/constant
		tkn := Token{}
		tkn.L = lexeme(t)

		switch {
		case reConst.MatchString(t):
			tkn.T = 6
			tkn.V = value(constants)
			constants++
			cache[t] = tkn
			return tkn, nil
		case reKws.MatchString(t):
			tkn.T = 4
			tkn.V = value(indentifiers)
			indentifiers++
			cache[t] = tkn
			return tkn, nil
		default:
			return tkn, errors.New("Unknown symbol")
		}
	}
}

type TokenResult struct {
	Token Token
	Err   error
}

func (c *TokenStream) GenTokenStream() iter.Seq[TokenResult] {
	return func(yield func(TokenResult) bool) {
		for _, t := range c.Tokens {
			token, err := c.Lexer(t)
			if !yield(TokenResult{Token: token, Err: err}) {
				return
			}
		}
	}
}
