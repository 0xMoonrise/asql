package parser

import (
	"strings"

	"github.com/0xMoonrise/asql/internal/lexer"
)

// EXPR_INSERT := INSERT INTO IDENTIFIER VALUES ( LIST_ATTR )
func (s *StackParser) expr_insert() *parseErr {
	if err := s.expect(lexer.INSERT, expectKeyword); err != nil {
		return err
	}
	s.next()
	if err := s.expect(lexer.INTO, expectKeyword); err != nil {
		return err
	}
	s.next()
	if !s.isIdentifier() {
		return &expectIdentifier
	}

	tableName := strings.ToUpper(s.lexemeAt(0))
	s.appendValue("insert_tables", tableName)

	s.next()
	if err := s.expect(lexer.VALUES, expectKeyword); err != nil {
		return err
	}
	s.next()
	if s.peekAt(0) != lexer.LPAR {
		return &expectDelimiter
	}
	s.next()
	if err := s.list_attr(); err != nil {
		return err
	}
	return nil
}

// LIST_ATTR := CONSTANT { , CONSTANT } )
func (s *StackParser) list_attr() *parseErr {
	if !s.isConstant() && s.peekAt(0) != lexer.NULL {
		return &expectConstant
	}

	s.appendValue("insert_values", string(s.tokenAt(0).L))
	s.next()

	for {
		if s.peekAt(0) == lexer.RPAR {
			s.next()
			return nil
		}
		if err := s.expect(lexer.COMMA, expectDelimiter); err != nil {
			return err
		}
		s.next()
		if !s.isConstant() && s.peekAt(0) != lexer.NULL {
			return &expectConstant
		}
		s.appendValue("insert_values", string(s.tokenAt(0).L))
		s.next()
	}
}
