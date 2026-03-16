package parser

import (
	"github.com/0xMoonrise/asql/internal/lexer"
)

const EOF = -1
const maxRecursion = 10

type stackParser struct {
	ptr   int
	depth int
	stack []lexer.Token
}

func (s *stackParser) expect(value lexer.Value, err parseErr) *parseErr {
	if s.stack[s.ptr].V == value {
		return nil
	}
	return &err
}

func (s *stackParser) peekAt(position int) lexer.Value {
	if s.ptr+position < len(s.stack) {
		return s.stack[s.ptr+position].V
	}
	return EOF
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

func (s *stackParser) isRelation() bool {
	return s.tokenAt(0).T == 8
}

func (s *stackParser) isConstant() bool {
	return s.peekAt(0) >= 600
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

	if s.peekAt(0) != lexer.WHERE {
		return nil
	}

	if err := s.where_clause(); err != nil {
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

// NAME_EXPR := IDENTIFIER | IDENTIFIER . IDENTIFIER
// This piece of code is not a pure LL1, because is simplier
// to implement a LL(k) grammar, taking the correct branch
// in case of ambiguity
func (s *stackParser) name_expr() *parseErr {
	if !s.isIdentifier() {
		return &expectIdentifier
	}

	// LL2
	if s.peekAt(1) != lexer.DOT {
		return nil
	}

	s.next()
	// is always true just to clarify the sintax of the production
	if err := s.expect(lexer.DOT, expectDelimiter); err != nil {
		return err
	}

	s.next()
	if !s.isIdentifier() {
		return &expectIdentifier
	}

	// case a.b.c is not expected, this solves the issue
	// with the ambiguity of LL1 just because looking at 2 tokens
	// transforming the grammar to LL2
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

	s.next()

	if err := s.databases_expr(); err != nil {
		return err
	}

	return nil
}

// DATABASES_EXPR := NAME_EXPR | NAME_EXPR , DATABASES_EXPR
func (s *stackParser) databases_expr() *parseErr {

	if err := s.safeRecursion(); err != nil {
		return err
	}

	defer s.unwind()

	if err := s.name_expr(); err != nil {
		return err
	}

	if terminal := s.next(); terminal == EOF || terminal == lexer.WHERE {
		return nil
	}

	if err := s.expect(lexer.COMMA, expectDelimiter); err != nil {
		return err
	}

	s.next()

	if err := s.name_expr(); err != nil {
		return err
	}

	terminal := s.peekAt(0)
	if terminal == EOF || terminal == lexer.WHERE {
		return nil
	}

	return s.databases_expr()
}

// WHERE_CLAUSE := WHERE CONDITION STMT
// STMT         := AND CONDITION STMT | OR CONDITION STMT |  λ
func (s *stackParser) where_clause() *parseErr {

	if err := s.expect(lexer.WHERE, expectKeyword); err != nil {
		return err
	}

	s.next()
	if !s.isIdentifier() {
		return &expectIdentifier
	}

	s.next()
	if err := s.relation_expr(); err != nil {
		return err
	}

	s.next()
	if !s.isConstant() {
		return &expectConstant
	}

	return nil
}

// RELATION := = | < | <= | > | >= | <>
func (s *stackParser) relation_expr() *parseErr {
	if !s.isRelation() {
		return &expectRelational
	}
	return nil
}
