package lexer

import (
	"errors"
	"fmt"
)

type LexerErr struct {
	Typ  int
	Code int
	Desc error
}

var UnknownSymbol = LexerErr{
	Typ:  1,
	Code: 101,
	Desc: errors.New("Unknown symbol"),
}

func (e LexerErr) Error() string {
	return fmt.Sprintf("%d|%d %s", e.Typ, e.Code, e.Desc)
}
