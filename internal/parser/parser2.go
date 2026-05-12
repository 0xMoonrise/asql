package parser

import (
	"strings"

	"github.com/0xMoonrise/asql/internal/lexer"
)

// DDL_EXPR := CREATE TABLE IDENTIFIER ( TABLE_BODY )
func (s *StackParser) ddl_expr() *parseErr {
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

	tableName := strings.ToUpper(s.lexemeAt(0))
	s.appendValue("ct_tables", tableName)
	s.currentTable = tableName

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

	s.currentTable = ""
	return nil
}

// TABLE_BODY := (COLUMN_DEF | CONSTRAINT_DEF) { , (COLUMN_DEF | CONSTRAINT_DEF) }
func (s *StackParser) table_body() *parseErr {
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

func (s *StackParser) table_body_item() *parseErr {
	if s.peekAt(0) == lexer.CONSTRAIN {
		return s.constraint_def()
	}
	return s.column_def()
}

// COLUMN_DEF := IDENTIFIER DATA_TYPE [ NULLABILITY ]
func (s *StackParser) column_def() *parseErr {
	if !s.isIdentifier() {
		return &expectIdentifier
	}

	colName := strings.ToUpper(s.lexemeAt(0))
	if s.currentTable != "" {
		s.appendValue("ct_columns:"+s.currentTable, colName)
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
func (s *StackParser) data_type() *parseErr {
	switch s.peekAt(0) {
	case lexer.NUMERIC:
		if s.currentTable != "" {
			s.appendValue("ct_types:"+s.currentTable, "NUMERIC")
		}
		return s.numeric_type()
	case lexer.CHAR:
		if s.currentTable != "" {
			s.appendValue("ct_types:"+s.currentTable, "CHAR")
		}
		return s.char_type()
	case lexer.DATE:
		if s.currentTable != "" {
			s.appendValue("ct_types:"+s.currentTable, "DATE")
		}
		s.next()
		return nil
	default:
		return &expectedDataType
	}
}

// NUMERIC_TYPE := NUMERIC ( INTEGER [ , INTEGER ] )
func (s *StackParser) numeric_type() *parseErr {
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
	if terminal := s.peekAt(1); terminal == EOF {
		return &expectParenthesisClosed
	}
	s.next()
	return nil
}

// CHAR_TYPE := CHAR ( INTEGER )
func (s *StackParser) char_type() *parseErr {
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
	if terminal := s.peekAt(1); terminal == EOF {
		return &expectParenthesisClosed
	}
	s.next()
	return nil
}

// NULLABILITY := NOT NULL | NULL
func (s *StackParser) nullability() *parseErr {
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
func (s *StackParser) constraint_def() *parseErr {
	if err := s.expect(lexer.CONSTRAIN, expectKeyword); err != nil {
		return err
	}
	s.next()
	if !s.isIdentifier() {
		return &expectIdentifier
	}

	constraintName := strings.ToUpper(s.lexemeAt(0))
	if s.currentTable != "" {
		s.appendValue("ct_constraints:"+s.currentTable, constraintName)
	}

	s.next()
	return s.constraint_type()
}

// CONSTRAINT_TYPE := PRIMARY KEY ( COL_LIST )
//
//	| CHECK ( CONDITION )
//	| FOREIGN KEY IDENTIFIER ( COL_LIST ) REFERENCES IDENTIFIER ( COL_LIST )
func (s *StackParser) constraint_type() *parseErr {
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
		if terminal := s.peekAt(1); terminal == EOF {
			return &expectParenthesisClosed
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
		if terminal := s.peekAt(1); terminal == EOF {
			return &expectParenthesisClosed
		}
		s.next()
		return nil

	case lexer.FOREIGN:
		s.next()
		if err := s.expect(lexer.KEY, expectKeyword); err != nil {
			return err
		}
		s.next()
		// FOREIGN KEY ( COL_LIST ) REFERENCES IDENTIFIER ( COL_LIST )
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
		if terminal := s.peekAt(1); terminal == EOF {
			return &expectParenthesisClosed
		}
		s.next()
		return nil

	default:
		return &expectedConstraintType
	}
}

// COL_LIST := IDENTIFIER { , IDENTIFIER }
func (s *StackParser) col_list() *parseErr {
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
