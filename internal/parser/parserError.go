package parser

import (
	"errors"
	"fmt"
)

type parseErr struct {
	Typ  int
	Code int
	Desc error
}

var ExpectKeyword = parseErr{
	Typ:  2,
	Code: 201,
	Desc: errors.New("A keyword is expected"),
}

var ExpectIdentifier = parseErr{
	Typ:  2,
	Code: 204,
	Desc: errors.New("An identifier is expected"),
}

var ExpectDelimiter = parseErr{
	Typ:  2,
	Code: 205,
	Desc: errors.New("A delimiter is expected"),
}

var ExpectConstant = parseErr{
	Typ:  2,
	Code: 206,
	Desc: errors.New("A constant is expected"),
}

var ExpectOperator = parseErr{
	Typ:  2,
	Code: 207,
	Desc: errors.New("An operator is expected"),
}

var ExpectRelational = parseErr{
	Typ:  2,
	Code: 208,
	Desc: errors.New("A relational operator is expected"),
}

func (e parseErr) Error() string {
	return fmt.Sprintf("%d|%d %s", e.Typ, e.Code, e.Desc)
}
