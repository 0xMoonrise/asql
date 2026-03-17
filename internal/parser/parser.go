package parser

import (
	"fmt"

	"github.com/0xMoonrise/asql/internal/lexer"
)

const EOF = -1
const DEBUG = false
const maxRecursion = 10

type stackParser struct {
	ptr        int
	depth      int
	parCounter int
	stack      []lexer.Token
	State      ErrorState
}

type ErrorState struct {
	Ptr     int
	Message error
	Token   lexer.Token
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

func (s *stackParser) tokenAt(position int) *lexer.Token {
	if s.ptr+position < len(s.stack) {
		return &s.stack[s.ptr+position]
	}
	return nil
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

func (s *stackParser) isOpBool() bool {
	v := s.peekAt(0)
	return v == lexer.AND || v == lexer.OR || v == lexer.NOT
}

func (s *stackParser) pushParStack() *parseErr {
	if s.tokenAt(0).V == lexer.LPAR {
		s.parCounter++
		s.ptr++
		return nil
	}
	return &expectDelimiter
}

func (s *stackParser) popParStack() *parseErr {
	if s.tokenAt(0).V == lexer.RPAR {
		s.parCounter--
		if terminal := s.next(); terminal == EOF && s.parCounter != 0 {
			return &expectParenthesisClosed
		}

		return nil
	}

	return &expectDelimiter
}

func NewParser(tokens []lexer.Token) *stackParser {
	p := &stackParser{
		ptr:        0,
		depth:      0,
		parCounter: 0,
		stack:      tokens,
		State:      ErrorState{},
	}
	return p
}

func (s *stackParser) Parse() *parseErr {

	if len(s.stack) == 0 {
		return &emptyStack
	}

	if DEBUG == true {
		for i, t := range s.stack {
			fmt.Printf("[%d] L=%s V=%d T=%d\n", i, t.L, t.V, t.T)
		}
	}

	if err := s.dml_expr(); err != nil {
		s.State = ErrorState{
			Message: err,
			Token:   *s.tokenAt(0),
			Ptr:     s.ptr,
		}
		fmt.Printf("failed at ptr=%d token=%v\n", s.ptr, s.stack[s.ptr])
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

func (s *stackParser) dml_expr() *parseErr {
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

	if err := s.name_expr(); err != nil {
		return err
	}

	s.next()

	for {
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

		s.next()
	}
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

// DATABASES_EXPR := TABLE_EXPR | TABLE_EXPR , DATABASES_EXPR
func (s *stackParser) databases_expr() *parseErr {
	if err := s.safeRecursion(); err != nil {
		return err
	}
	defer s.unwind()

	if err := s.table_expr(); err != nil {
		return err
	}

	for {
		terminal := s.peekAt(0)
		if terminal == EOF || terminal == lexer.WHERE || terminal == lexer.RPAR {
			return nil
		}
		if s.isIdentifier() {
			return nil
		}

		if err := s.expect(lexer.COMMA, expectDelimiter); err != nil {
			return err
		}

		s.next()

		if err := s.table_expr(); err != nil {
			return err
		}
	}
}

// TABLE_EXPR := NAME_EXPR [ALIAS] | ( DML_EXPR ) ALIAS
func (s *stackParser) table_expr() *parseErr {

	if s.peekAt(0) == lexer.LPAR {
		if err := s.pushParStack(); err != nil {
			return err
		}
		if err := s.dml_expr(); err != nil {
			return err
		}
		if err := s.popParStack(); err != nil {
			return err
		}
		if !s.isIdentifier() {
			return &expectIdentifier
		}
		if next := s.next(); next == EOF {
			return nil
		}
		return s.expectTableTerminator()
	}

	if err := s.name_expr(); err != nil {
		return err
	}

	if next := s.next(); next == EOF {
		return nil
	}

	if s.isIdentifier() {
		if next := s.next(); next == EOF {
			return nil
		}
		return s.expectTableTerminator()
	}

	return nil
}

func (s *stackParser) expectTableTerminator() *parseErr {
	terminal := s.peekAt(0)
	if terminal == EOF ||
		terminal == lexer.WHERE ||
		terminal == lexer.RPAR ||
		terminal == lexer.COMMA {
		return nil
	}
	return &unexpectedToken
}

// WHERE_CLAUSE := WHERE CONDITION STMT
func (s *stackParser) where_clause() *parseErr {
	if err := s.expect(lexer.WHERE, expectKeyword); err != nil {
		return err
	}

	s.next()

	if err := s.condition_expr(); err != nil {
		return err
	}

	return nil
}

// CONDITION := NAME_EXPR RELATION CONSTANT | ( CONDITION STMT )
func (s *stackParser) condition_expr() *parseErr {

	if err := s.safeRecursion(); err != nil {
		return err
	}
	defer s.unwind()

	if s.peekAt(0) == lexer.NOT {
		s.next()
	}

	if err := s.name_expr(); err != nil {
		return err
	}

	s.next()
	if s.peekAt(0) == lexer.IN {
		return s.where_subquery_expr()
	}

	if err := s.relation_expr(); err != nil {
		return &expectRelational
	}

	s.next()
	switch {
	case s.isConstant():
	case s.isIdentifier():
		if err := s.name_expr(); err != nil {
			return err
		}
	default:
		return &expectedConstOrIdent
	}

	if terminal := s.next(); terminal == EOF || terminal == lexer.RPAR {
		return nil
	}

	if !s.isOpBool() {
		return &expectKeyword
	}

	s.next()

	if s.peekAt(0) == lexer.NOT {
		s.next()
	}

	if err := s.condition_expr(); err != nil {
		return err
	}

	return nil
}

// WHERE_SUBQUERY := IN ( DML_EXPR ) | IN ( DML_EXPR ) BOOL_OP CONDITION
func (s *stackParser) where_subquery_expr() *parseErr {
	s.next()
	if err := s.pushParStack(); err != nil {
		return err
	}
	if err := s.dml_expr(); err != nil {
		return err
	}
	if err := s.popParStack(); err != nil {
		return err
	}

	terminal := s.peekAt(0)
	if terminal == EOF || terminal == lexer.RPAR {
		return nil
	}

	if s.isOpBool() {
		s.next()

		if s.peekAt(0) == lexer.NOT {
			s.next()
		}

		return s.condition_expr()
	}

	return &unexpectedToken
}

// RELATION := = | < | <= | > | >= | <>
func (s *stackParser) relation_expr() *parseErr {
	if !s.isRelation() {
		return &expectRelational
	}
	return nil
}
