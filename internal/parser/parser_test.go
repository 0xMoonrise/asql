package parser

import (
	"asql/internal/lexer"
	"testing"

	"github.com/stretchr/testify/assert"
)

func tokenize(input string) []lexer.Token {
	raw := lexer.Tokenize(input)
	tokens, _ := lexer.Lexer([][]string{raw})
	return tokens
}

func parse(input string) *parseErr {
	tokens := tokenize(input)
	p := NewParser(tokens)
	return p.Parse()
}

func TestSelect(t *testing.T) {
	t.Run("star", func(t *testing.T) {
		assert.Nil(t, parse("SELECT * FROM foo"),
			"SELECT * FROM foo should be valid")
	})

	t.Run("single column", func(t *testing.T) {
		assert.Nil(t, parse("SELECT a FROM foo"),
			"single column should be valid")
	})

	t.Run("multiple columns", func(t *testing.T) {
		assert.Nil(t, parse("SELECT a, b, c FROM foo"),
			"multiple columns should be valid")
	})

	t.Run("qualified column", func(t *testing.T) {
		assert.Nil(t, parse("SELECT a.b FROM foo"),
			"qualified column should be valid")
	})

	t.Run("mixed columns", func(t *testing.T) {
		assert.Nil(t, parse("SELECT a, a.b FROM foo"),
			"mix of simple and qualified columns should be valid")
	})

	t.Run("all mixed columns", func(t *testing.T) {
		assert.Nil(t, parse("SELECT a, a.b, c FROM foo"),
			"full mix of columns should be valid")
	})
}

func TestSelectErrors(t *testing.T) {
	t.Run("missing FROM after star", func(t *testing.T) {
		assert.NotNil(t, parse("SELECT *"),
			"SELECT * without FROM should fail")
	})

	t.Run("trailing comma", func(t *testing.T) {
		assert.NotNil(t, parse("SELECT a, b, FROM foo"),
			"trailing comma should fail")
	})

	t.Run("missing SELECT keyword", func(t *testing.T) {
		assert.NotNil(t, parse("a, b FROM foo"),
			"query without SELECT should fail")
	})

	t.Run("empty input", func(t *testing.T) {
		assert.NotNil(t, parse(""),
			"empty input should fail")
	})
}

func TestColumnQualification(t *testing.T) {
	t.Run("three level qualification", func(t *testing.T) {
		assert.NotNil(t, parse("SELECT a.b.c FROM foo"),
			"three level qualification should fail")
	})

	t.Run("four level qualification", func(t *testing.T) {
		assert.NotNil(t, parse("SELECT a.b.c.d FROM foo"),
			"four level qualification should fail")
	})
}

func TestErrorCodes(t *testing.T) {
	t.Run("missing SELECT returns 201", func(t *testing.T) {
		err := parse("FROM foo")
		assert.NotNil(t, err)
		assert.Equal(t, 201, err.Code,
			"should return error 201 (keyword expected)")
	})

	t.Run("trailing comma returns 204", func(t *testing.T) {
		err := parse("SELECT a, FROM foo")
		assert.NotNil(t, err)
		assert.Equal(t, 204, err.Code,
			"should return error 204 (identifier expected)")
	})

	t.Run("over-qualified column returns 205", func(t *testing.T) {
		err := parse("SELECT a.b.c FROM foo")
		assert.NotNil(t, err)
		assert.Equal(t, 205, err.Code,
			"should return error 205 (delimiter expected)")
	})
}
