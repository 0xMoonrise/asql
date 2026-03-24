package parser

import (
	"github.com/0xMoonrise/asql/internal/lexer"
)

// DDL_EXPR := CREATE TABLE IDENTIFIER ( TABLE_BODY )
func (s *stackParser) ddl_expr() *parseErr {
	if err := s.expect(lexer.CREATE, expectKeyword); err != nil {
		return err
	}
	s.next()
	if err := s.expect(lexer.TABLE, expectKeyword); err != nil {
		return err
	}
	s.next()
	if !s.isIdentifier() {
		return &expectIdentifier
	}
	s.next()
	if err := s.expect(lexer.LPAR, expectDelimiter); err != nil {
		return err
	}
	s.next()
	if err := s.table_body(); err != nil {
		return err
	}
	if err := s.expect(lexer.RPAR, expectDelimiter); err != nil {
		return err
	}
	s.next()
	return nil
}

// TABLE_BODY := (COLUMN_DEF | CONSTRAINT_DEF) { , (COLUMN_DEF | CONSTRAINT_DEF) }
func (s *stackParser) table_body() *parseErr {
	if err := s.table_body_item(); err != nil {
		return err
	}

	for {
		if s.peekAt(0) != lexer.COMMA {
			return nil
		}
		s.next()
		if err := s.table_body_item(); err != nil {
			return err
		}
	}
}

func (s *stackParser) table_body_item() *parseErr {
	if s.peekAt(0) == lexer.CONSTRAIN {
		return s.constraint_def()
	}
	return s.column_def()
}

// COLUMN_DEF := IDENTIFIER DATA_TYPE [ NULLABILITY ]
func (s *stackParser) column_def() *parseErr {
	if !s.isIdentifier() {
		return &expectIdentifier
	}
	s.next()
	if err := s.data_type(); err != nil {
		return err
	}
	if err := s.nullability(); err != nil {
		return err
	}
	return nil
}

// DATA_TYPE := NUMERIC_TYPE | CHAR_TYPE | DATE
func (s *stackParser) data_type() *parseErr {
	switch s.peekAt(0) {
	case lexer.NUMERIC:
		return s.numeric_type()
	case lexer.CHAR:
		return s.char_type()
	case lexer.DATE:
		s.next()
		return nil
	default:
		return &expectedDataType
	}
}

// NUMERIC_TYPE := NUMERIC ( INTEGER [ , INTEGER ] )
func (s *stackParser) numeric_type() *parseErr {
	if err := s.expect(lexer.NUMERIC, expectKeyword); err != nil {
		return err
	}
	s.next()
	if err := s.expect(lexer.LPAR, expectDelimiter); err != nil {
		return err
	}
	s.next()
	if !s.isConstant() {
		return &expectConstant
	}
	s.next()
	if s.peekAt(0) == lexer.COMMA {
		s.next()
		if !s.isConstant() {
			return &expectConstant
		}
		s.next()
	}
	if err := s.expect(lexer.RPAR, expectDelimiter); err != nil {
		return err
	}
	s.next()
	return nil
}

// CHAR_TYPE := CHAR ( INTEGER )
func (s *stackParser) char_type() *parseErr {
	if err := s.expect(lexer.CHAR, expectKeyword); err != nil {
		return err
	}
	s.next()
	if err := s.expect(lexer.LPAR, expectDelimiter); err != nil {
		return err
	}
	s.next()
	if !s.isConstant() {
		return &expectConstant
	}
	s.next()
	if err := s.expect(lexer.RPAR, expectDelimiter); err != nil {
		return err
	}
	s.next()
	return nil
}

// NULLABILITY := NOT NULL | NULL
func (s *stackParser) nullability() *parseErr {
	switch s.peekAt(0) {
	case lexer.NOT:
		s.next()
		if err := s.expect(lexer.NULL, expectKeyword); err != nil {
			return err
		}
		s.next()
	case lexer.NULL:
		s.next()
	}
	return nil
}

// CONSTRAINT_DEF := CONSTRAINT IDENTIFIER CONSTRAINT_TYPE
func (s *stackParser) constraint_def() *parseErr {
	if err := s.expect(lexer.CONSTRAIN, expectKeyword); err != nil {
		return err
	}
	s.next()
	if !s.isIdentifier() {
		return &expectIdentifier
	}
	s.next()
	return s.constraint_type()
}

// CONSTRAINT_TYPE := PRIMARY KEY ( COL_LIST )
//                  | CHECK ( CONDITION )
//                  | FOREIGN KEY IDENTIFIER ( COL_LIST ) REFERENCES IDENTIFIER ( COL_LIST )

func (s *stackParser) constraint_type() *parseErr {
	switch s.peekAt(0) {
	case lexer.PRIMARY:
		s.next()
		if err := s.expect(lexer.KEY, expectKeyword); err != nil {
			return err
		}
		s.next()
		if err := s.expect(lexer.LPAR, expectDelimiter); err != nil {
			return err
		}
		s.next()
		if err := s.col_list(); err != nil {
			return err
		}
		if err := s.expect(lexer.RPAR, expectDelimiter); err != nil {
			return err
		}
		s.next()
		return nil

	case lexer.CHECK:
		s.next()
		if err := s.expect(lexer.LPAR, expectDelimiter); err != nil {
			return err
		}
		s.next()
		if err := s.condition_expr(); err != nil {
			return err
		}
		if err := s.expect(lexer.RPAR, expectDelimiter); err != nil {
			return err
		}
		s.next()
		return nil

	case lexer.FOREIGN:
		s.next()
		if err := s.expect(lexer.KEY, expectKeyword); err != nil {
			return err
		}
		s.next()
		if !s.isIdentifier() {
			return &expectIdentifier
		}
		s.next()
		if err := s.expect(lexer.LPAR, expectDelimiter); err != nil {
			return err
		}
		s.next()
		if err := s.col_list(); err != nil {
			return err
		}
		if err := s.expect(lexer.RPAR, expectDelimiter); err != nil {
			return err
		}
		s.next()
		if err := s.expect(lexer.REFERENCES, expectKeyword); err != nil {
			return err
		}
		s.next()
		if !s.isIdentifier() {
			return &expectIdentifier
		}
		s.next()
		if err := s.expect(lexer.LPAR, expectDelimiter); err != nil {
			return err
		}
		s.next()
		if err := s.col_list(); err != nil {
			return err
		}
		if err := s.expect(lexer.RPAR, expectDelimiter); err != nil {
			return err
		}
		s.next()
		return nil

	default:
		return &expectedConstraintType
	}
}

// COL_LIST := IDENTIFIER { , IDENTIFIER }
func (s *stackParser) col_list() *parseErr {
	if !s.isIdentifier() {
		return &expectIdentifier
	}
	s.next()

	for {
		if s.peekAt(0) != lexer.COMMA {
			return nil
		}
		s.next()
		if !s.isIdentifier() {
			return &expectIdentifier
		}
		s.next()
	}
}
