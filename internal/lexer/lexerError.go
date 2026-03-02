package lexer

import "errors"

type LexerErr struct {
	Typ  int
	Code int
	Desc error
}

var UnknownSymbol = LexerErr{
	Typ:  1,
	Code: 101,
	Desc: errors.New("unknown symbol"),
}

func (e LexerErr) Error() string {
	return e.Desc.Error()
}
