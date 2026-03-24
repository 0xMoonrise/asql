package parser

import (
	"errors"
	"fmt"
)

type parseErr struct {
	Typ  int
	Code int
	Err  error
}

var expectKeyword = parseErr{
	Typ:  2,
	Code: 201,
	Err:  errors.New("A keyword is expected"),
}

var expectedColOrStar = parseErr{
	Typ:  2,
	Code: 202,
	Err:  errors.New("A start is expected or a column expresion"),
}

var expectIdentifier = parseErr{
	Typ:  2,
	Code: 204,
	Err:  errors.New("An identifier is expected"),
}

var expectDelimiter = parseErr{
	Typ:  2,
	Code: 205,
	Err:  errors.New("A delimiter is expected"),
}

var expectConstant = parseErr{
	Typ:  2,
	Code: 206,
	Err:  errors.New("A constant is expected"),
}

var expectOperator = parseErr{
	Typ:  2,
	Code: 207,
	Err:  errors.New("An operator is expected"),
}

var expectRelational = parseErr{
	Typ:  2,
	Code: 208,
	Err:  errors.New("A relational operator is expected"),
}

var expectParenthesisClosed = parseErr{
	Typ:  2,
	Code: 209,
	Err:  errors.New("Parenthesis not closed"),
}

var unexpectedToken = parseErr{
	Typ:  2,
	Code: 210,
	Err:  errors.New("Unexpcted token"),
}

var expectedConstOrIdent = parseErr{
	Typ:  2,
	Code: 210,
	Err:  errors.New("expected a constant or identifier"),
}

var expectedKeywordAtStart = parseErr{
	Typ:  2,
	Code: 211,
	Err:  errors.New("Expected a select, create or insert"),
}

var expectedDataType = parseErr{
	Typ:  2,
	Code: 212,
	Err:  errors.New("Expected a data type"),
}

var expectedConstraintType = parseErr{
	Typ:  2,
	Code: 213,
	Err:  errors.New("Expected PRIMARY, CHECK or FOREIGN"),
}

// Type 0 kernel error
var maxRecursionReached = parseErr{
	Typ:  0,
	Code: 1,
	Err:  errors.New("Max recursion detected, stoped"),
}

var emptyStack = parseErr{
	Typ:  0,
	Code: 2,
	Err:  errors.New("No token stream was provided"),
}

func (e parseErr) Error() string {
	return fmt.Sprintf("%d|%d %s", e.Typ, e.Code, e.Err)
}
