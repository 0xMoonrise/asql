package parser

import (
	"asql/internal/lexer"
	"errors"
	"fmt"
)

type parseErr struct {
	Typ   int
	Code  int
	Desc  error
	Token lexer.Token
}

var expectKeyword = parseErr{
	Typ:  2,
	Code: 201,
	Desc: errors.New("A keyword is expected"),
}

var expectedColOrStar = parseErr{
	Typ:  2,
	Code: 202,
	Desc: errors.New("A start is expected or a column expresion"),
}

var expectIdentifier = parseErr{
	Typ:  2,
	Code: 204,
	Desc: errors.New("An identifier is expected"),
}

var expectDelimiter = parseErr{
	Typ:  2,
	Code: 205,
	Desc: errors.New("A delimiter is expected"),
}

var expectConstant = parseErr{
	Typ:  2,
	Code: 206,
	Desc: errors.New("A constant is expected"),
}

var expectOperator = parseErr{
	Typ:  2,
	Code: 207,
	Desc: errors.New("An operator is expected"),
}

var expectRelational = parseErr{
	Typ:  2,
	Code: 208,
	Desc: errors.New("A relational operator is expected"),
}

var maxRecursionReached = parseErr{
	Typ:  -1,
	Code: 1,
	Desc: errors.New("Max recursion detected, stoped"),
}

func (e parseErr) Error() string {
	return fmt.Sprintf("%d|%d %s", e.Typ, e.Code, e.Desc)
}
