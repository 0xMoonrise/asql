package parser

import (
	"errors"
)

type ParseErr struct {
	Typ  int
	Code int
	Desc error
}

var ExpectKeyword = ParseErr{
	Typ:  2,
	Code: 201,
	Desc: errors.New("A keyword is expected"),
}

var ExpectIdentifier = ParseErr{
	Typ:  2,
	Code: 204,
	Desc: errors.New("An identifier is expected"),
}

var ExpectDelimiter = ParseErr{
	Typ:  2,
	Code: 205,
	Desc: errors.New("A delimiter is expected"),
}

var ExpectConstant = ParseErr{
	Typ:  2,
	Code: 206,
	Desc: errors.New("A constant is expected"),
}

var ExpectOperator = ParseErr{
	Typ:  2,
	Code: 207,
	Desc: errors.New("An operator is expected"),
}

var ExpectRelational = ParseErr{
	Typ:  2,
	Code: 208,
	Desc: errors.New("A relational operator is expected"),
}
