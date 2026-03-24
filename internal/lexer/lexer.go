package lexer

import (
	"errors"
	"iter"
	"regexp"
	"strings"

	"github.com/0xMoonrise/asql/internal/utils"
)

type Token struct {
	L    Lexeme
	V    Value
	T    Typ
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
		{L: "SELECT", V: SELECT, T: 1},
		{L: "FROM", V: FROM, T: 1},
		{L: "WHERE", V: WHERE, T: 1},
		{L: "IN", V: IN, T: 1},
		{L: "AND", V: AND, T: 1},
		{L: "OR", V: OR, T: 1},
		{L: "CREATE", V: CREATE, T: 1},
		{L: "TABLE", V: TABLE, T: 1},
		{L: "CHAR", V: CHAR, T: 1},
		{L: "NUMERIC", V: NUMERIC, T: 1},
		{L: "NOT", V: NOT, T: 1},
		{L: "NULL", V: NULL, T: 1},
		{L: "CONSTRAINT", V: CONSTRAIN, T: 1},
		{L: "KEY", V: KEY, T: 1},
		{L: "PRIMARY", V: PRIMARY, T: 1},
		{L: "FOREIGN", V: FOREIGN, T: 1},
		{L: "REFERENCES", V: REFERENCES, T: 1},
		{L: "INSERT", V: INSERT, T: 1},
		{L: "INTO", V: INTO, T: 1},
		{L: "VALUES", V: VALUES, T: 1},
		{L: "DATE", V: DATE, T: 1},
		{L: "CHECK", V: CHECK, T: 1},
		// delimitators (5)
		{L: ",", V: COMMA, T: 5},
		{L: ".", V: DOT, T: 5},
		{L: "(", V: LPAR, T: 5},
		{L: ")", V: RPAR, T: 5},
		{L: "'", V: APOS, T: 5},
		// Operators (7)
		{L: "+", V: PLUS, T: 7},
		{L: "-", V: MINUS, T: 7},
		{L: "*", V: TIMES, T: 7},
		{L: "/", V: DIVS, T: 7},
		// Relations (8)
		{L: ">", V: GT, T: 8},
		{L: "<", V: LT, T: 8},
		{L: "=", V: EQ, T: 8},
		{L: ">=", V: GE, T: 8},
		{L: "<=", V: LE, T: 8},
		{L: "<>", V: NE, T: 8},
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
	delimitators = `[.,()*']`                // is '*' a delimitator?
	relations    = `>=|<=|<>|=|>|<`          //
	noLexer      = `[^\w\s.,()'>=<]`         // Any to catch errors
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
		token.L = Lexeme(t)

		switch {
		case reConst.MatchString(t):
			token.T = 6
			token.V = Value(constants)
			constants++
			cache[t] = token
			return token, nil
		case reKws.MatchString(t):
			token.T = 4
			token.V = Value(indentifiers)
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
	L    Lexeme
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
					Err:  UnknownSymbol,
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
