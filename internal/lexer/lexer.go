package lexer

import (
	"asql/internal/utils"
	"errors"
	"iter"
	"regexp"
	"strings"
)

type Token struct {
	L    lexeme
	V    value
	T    typ
	Line int
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
		{L: "SELECT", V: __select, T: 1},
		{L: "FROM", V: __from, T: 1},
		{L: "WHERE", V: __where, T: 1},
		{L: "IN", V: __in, T: 1},
		{L: "AND", V: __and, T: 1},
		{L: "OR", V: __or, T: 1},
		{L: "CREATE", V: __create, T: 1},
		{L: "TABLE", V: __table, T: 1},
		{L: "CHAR", V: __char, T: 1},
		{L: "NUMERIC", V: __numeric, T: 1},
		{L: "NOT", V: __not, T: 1},
		{L: "NULL", V: __null, T: 1},
		{L: "CONSTRAINT", V: __constraint, T: 1},
		{L: "KEY", V: __key, T: 1},
		{L: "PRIMARY", V: __primary, T: 1},
		{L: "FOREIGN", V: __foreign, T: 1},
		{L: "REFERENCES", V: __references, T: 1},
		{L: "INSERT", V: __insert, T: 1},
		{L: "INTO", V: __into, T: 1},
		{L: "VALUES", V: __values, T: 1},
		// delimitators (5)
		{L: ",", V: comma, T: 5},
		{L: ".", V: dot, T: 5},
		{L: "(", V: lparentheses, T: 5},
		{L: ")", V: rparentheses, T: 5},
		{L: "'", V: apostrophe, T: 5},
		// Operators (7)
		// {L: "+", V: plus, T: 7},
		// {L: "-", V: minus, T: 7},
		// {L: "*", V: times, T: 7},
		// {L: "/", V: divition, T: 7},
		// Relations (8)
		{L: ">", V: gt, T: 8},
		{L: "<", V: lt, T: 8},
		{L: "=", V: eq, T: 8},
		{L: ">=", V: ge, T: 8},
		{L: "<=", V: le, T: 8},
	}

	for _, token := range Tokens {
		lex[string(token.L)] = Token{
			L: token.L,
			V: token.V,
			T: token.T,
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
	delimitators = `[.,()*]`                 // is '*' a delimitator?
	relations    = `>=|<=|!=|=|>|<`          //
	noLexer      = `[^\w\s.,()#'>=<!]`       // Any to catch errors
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

// Apply the criteria for tokens
func lexerTable() func(t string) (Token, error) {
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
		token := Token{}
		token.L = lexeme(t)

		switch {
		case reConst.MatchString(t):
			token.T = 6
			token.V = value(constants)
			constants++
			cache[t] = token
			return token, nil
		case reKws.MatchString(t):
			token.T = 4
			token.V = value(indentifiers)
			indentifiers++
			cache[t] = token
			return token, nil
		default:
			return token, errors.New("Unknown symbol")
		}
	}
}

type TokenResult struct {
	Token Token
	Err   error
}

func (c *TokenStream) GenTokenStream() iter.Seq2[int, TokenResult] {
	return func(yield func(int, TokenResult) bool) {
		for i, t := range c.Tokens {
			token, err := c.Lexer(t)
			if !yield(i, TokenResult{Token: token, Err: err}) {
				return
			}
		}
	}
}

type ErrorState struct {
	L    lexeme
	Err  error
	Line int
}

func Lexer(lines [][]string) (tokenStream []Token, errors []ErrorState) {

	lex := &TokenStream{Lexer: lexerTable()}

	for i, line := range lines { // <- Entry point
		lex.Tokens = line
		for _, result := range lex.GenTokenStream() {

			if result.Err != nil {
				newError := ErrorState{
					L:    result.Token.L,
					Err:  result.Err,
					Line: i + 1,
				}
				errors = append(errors, newError)
				continue
			}

			token := result.Token
			token.Line = i + 1
			tokenStream = append(tokenStream, token)
		}
	}

	return
}
