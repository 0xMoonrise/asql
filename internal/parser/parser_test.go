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

func TestSelectStar(t *testing.T) {
	assert.Nil(t, parse("SELECT * FROM foo"),
		"SELECT * FROM foo should be valid")
}

func TestSelectStarMissingFrom(t *testing.T) {
	assert.NotNil(t, parse("SELECT *"),
		"SELECT * without FROM should fail")
}

func TestSelectSingleColumn(t *testing.T) {
	assert.Nil(t, parse("SELECT a FROM foo"),
		"single column should be valid")
}

func TestSelectMultipleColumns(t *testing.T) {
	assert.Nil(t, parse("SELECT a, b, c FROM foo"),
		"multiple columns should be valid")
}

func TestSelectQualifiedColumn(t *testing.T) {
	assert.Nil(t, parse("SELECT a.b FROM foo"),
		"qualified column should be valid")
}

func TestSelectMixedColumns(t *testing.T) {
	assert.Nil(t, parse("SELECT a, a.b FROM foo"),
		"mix of simple and qualified columns should be valid")
}

func TestSelectAllMixedColumns(t *testing.T) {
	assert.Nil(t, parse("SELECT a, a.b, c FROM foo"),
		"full mix of columns should be valid")
}

func TestSelectTrailingComma(t *testing.T) {
	assert.NotNil(t, parse("SELECT a, b, FROM foo"),
		"trailing comma should fail")
}

func TestSelectThreeLevelColumn(t *testing.T) {
	assert.NotNil(t, parse("SELECT a.b.c FROM foo"),
		"three level qualification should fail")
}

func TestSelectFourLevelColumn(t *testing.T) {
	assert.NotNil(t, parse("SELECT a.b.c.d FROM foo"),
		"four level qualification should fail")
}

func TestMissingSelect(t *testing.T) {
	assert.NotNil(t, parse("a, b FROM foo"),
		"query without SELECT should fail")
}

func TestEmptyInput(t *testing.T) {
	assert.NotNil(t, parse(""),
		"empty input should fail")
}

func TestMissingSelectKeywordErrorCode(t *testing.T) {
	err := parse("FROM foo")
	assert.NotNil(t, err)
	assert.Equal(t, 201, err.Code,
		"should return error 201 (keyword expected)")
}

func TestTrailingCommaErrorCode(t *testing.T) {
	err := parse("SELECT a, FROM foo")
	assert.NotNil(t, err)
	assert.Equal(t, 204, err.Code,
		"should return error 204 (identifier expected)")
}
