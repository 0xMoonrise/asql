package parser

import (
	"asql/internal/lexer"
)

const EOF = -1

type stackParser struct {
	ptr   int
	stack []lexer.Token
}

func (s *stackParser) expect(value lexer.Value, err parseErr) *parseErr {
	token := s.stack[s.ptr]
	if token.V == value {
		return nil
	}
	err.Token = token
	return &err
}

func (s *stackParser) peek() lexer.Value {
	return s.stack[s.ptr].V
}

func (s *stackParser) next() lexer.Value {
	if s.ptr+1 < len(s.stack) {
		s.ptr++
		return s.stack[s.ptr].V
	}
	return EOF
}

func NewParser(tokens []lexer.Token) *stackParser {
	return &stackParser{
		ptr:   0,
		stack: tokens,
	}
}

func (s *stackParser) Parse() *parseErr {

	if err := s.select_expr(); err != nil {
		return err
	}

	if err := s.from_expr(); err != nil {
		return err
	}

	return nil
}

// SELECT_EXPR  := SELECT * FROM_EXPR | SELECT COLUMNS_EXPR FROM_EXPR
// COLUMNS_EXPR := NAME_EXPR | NAME_EXPR , COLUMNS_EXPR
// NAME_EXPR    := IDENTIFIER | IDENTIFIER . IDENTIFIER
// FROM_EXPR    := FROM NAME_EXPR | FROM NAME_EXPR WHERE_CLAUSE
// WHERE_CLAUSE := WHERE COLUMN_EXPR OPERATOR CONSTANT
// OPERATOR     := = | < | <= | > | >= | <>

func (s *stackParser) select_expr() *parseErr {
	if err := s.expect(lexer.SELECT, expectKeyword); err != nil {
		return err
	}

	s.next()
	terminal := s.peek()

	switch {
	case terminal == lexer.TIMES:
		s.next()
	case terminal > 400 && terminal < 600:
		return s.name_expr()
	default:
		return &expectedColOrStar
	}

	return nil
}

var cycle = 0

const maxRecursion = 10

func (s *stackParser) name_expr() *parseErr {

	if cycle > maxRecursion {
		return &maxRecursionReached
	}
	cycle++

	s.next()
	if s.peek() == lexer.FROM {
		return nil
	}

	if err := s.expect(lexer.COMMA, expectDelimiter); err != nil {
		return err
	}

	s.next()
	if !(s.peek() > 400 && s.peek() < 600) {
		return &expectIdentifier
	}

	return s.name_expr()
}

func (s *stackParser) from_expr() *parseErr {
	if err := s.expect(lexer.FROM, expectKeyword); err != nil {
		return err
	}

	return nil
}
