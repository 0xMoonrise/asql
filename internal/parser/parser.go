package parser

import (
	"asql/internal/lexer"
)

const EOF = -1
const maxRecursion = 10

type stackParser struct {
	ptr   int
	depth int
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

func (s *stackParser) peekAt(position int) lexer.Value {
	return s.stack[s.ptr+position].V
}

func (s *stackParser) next() lexer.Value {
	if s.ptr+1 < len(s.stack) {
		s.ptr++
		return s.stack[s.ptr].V
	}
	return EOF
}

func (s *stackParser) tokenAt(position int) lexer.Token {
	return s.stack[s.ptr+position]
}

func (s *stackParser) isIdentifier() bool {
	terminal := s.peekAt(0)
	return (terminal > 400 && terminal < 600)
}

func NewParser(tokens []lexer.Token) *stackParser {
	p := &stackParser{
		ptr:   0,
		depth: 0,
		stack: tokens,
	}
	return p
}

func (s *stackParser) Parse() *parseErr {

	if len(s.stack) == 0 {
		return &emptyStack
	}

	if err := s.select_expr(); err != nil {
		return err
	}

	if err := s.from_expr(); err != nil {
		return err
	}

	return nil
}

func (s *stackParser) safeRecursion() *parseErr {
	s.depth++
	if s.depth > maxRecursion {
		return &maxRecursionReached
	}
	return nil
}

func (s *stackParser) unwind() {
	s.depth--
}

// SELECT_EXPR  := SELECT * FROM_EXPR | SELECT COLUMNS_EXPR FROM_EXPR
func (s *stackParser) select_expr() *parseErr {

	if err := s.expect(lexer.SELECT, expectKeyword); err != nil {
		return err
	}

	s.next()
	terminal := s.peekAt(0)

	switch {
	case terminal == lexer.TIMES:
		s.next()
	case s.isIdentifier():
		return s.columns_expr()
	default:
		return &expectedColOrStar
	}

	return nil
}

// columns_expr: ptr must point to the first identifier on entry
// COLUMNS_EXPR := NAME_EXPR | NAME_EXPR , COLUMNS_EXPR
func (s *stackParser) columns_expr() *parseErr {

	if err := s.safeRecursion(); err != nil {
		return err
	}

	defer s.unwind()

	if err := s.name_expr(); err != nil {
		return err
	}

	s.next()

	if s.peekAt(0) == lexer.FROM {
		return nil
	}

	if err := s.expect(lexer.COMMA, expectDelimiter); err != nil {
		return err
	}

	s.next()

	if err := s.name_expr(); err != nil {
		return err
	}

	return s.columns_expr()
}

// NAME_EXPR    := IDENTIFIER | IDENTIFIER . IDENTIFIER
func (s *stackParser) name_expr() *parseErr {
	if !s.isIdentifier() {
		return &expectIdentifier
	}

	if s.peekAt(1) != lexer.DOT {
		return nil
	}

	s.next()
	if err := s.expect(lexer.DOT, expectDelimiter); err != nil {
		return err
	}

	s.next()
	if !s.isIdentifier() {
		return &expectIdentifier
	}

	//case a.b.c is not expected
	if s.peekAt(1) == lexer.DOT {
		return &expectDelimiter
	}

	return nil
}

// FROM_EXPR    := FROM NAME_EXPR | FROM NAME_EXPR WHERE_CLAUSE
func (s *stackParser) from_expr() *parseErr {
	if err := s.expect(lexer.FROM, expectKeyword); err != nil {
		return err
	}

	return nil
}

// WHERE_CLAUSE := WHERE COLUMN_EXPR OPERATOR CONSTANT
// OPERATOR     := = | < | <= | > | >= | <>
