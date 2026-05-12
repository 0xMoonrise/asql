package parser

import (
	"fmt"
	"strings"

	"github.com/0xMoonrise/asql/internal/lexer"
)

const EOF = -1
const DEBUG = false
const maxRecursion = 10

type StackParser struct {
	ptr          int
	depth        int
	parCounter   int
	stack        []lexer.Token
	State        ErrorState
	Metadata     map[string][]string
	currentTable string
}

type ErrorState struct {
	Ptr     int
	Message error
	Token   lexer.Token
}

func (s *StackParser) lexemeAt(offset int) string {
	if t := s.tokenAt(offset); t != nil {
		return strings.ToUpper(string(t.L))
	}
	return ""
}

func (s *StackParser) appendValue(header string, value string) {
	s.Metadata[header] = append(s.Metadata[header], value)
}

func (s *StackParser) expect(value lexer.Value, err parseErr) *parseErr {
	if s.stack[s.ptr].V == value {
		return nil
	}
	return &err
}

func (s *StackParser) peekAt(position int) lexer.Value {
	if s.ptr+position < len(s.stack) {
		return s.stack[s.ptr+position].V
	}
	return EOF
}

func (s *StackParser) next() lexer.Value {
	if s.ptr+1 < len(s.stack) {
		s.ptr++
		return s.stack[s.ptr].V
	}
	return EOF
}

func (s *StackParser) tokenAt(position int) *lexer.Token {
	idx := s.ptr + position
	if idx >= 0 && idx < len(s.stack) {
		return &s.stack[idx]
	}
	return nil
}

func (s *StackParser) isIdentifier() bool {
	terminal := s.peekAt(0)
	return (terminal > 400 && terminal < 600)
}

func (s *StackParser) isRelation() bool {
	return s.tokenAt(0).T == 8
}

func (s *StackParser) isConstant() bool {
	return s.peekAt(0) >= 600
}

func (s *StackParser) isOpBool() bool {
	v := s.peekAt(0)
	return v == lexer.AND || v == lexer.OR || v == lexer.NOT
}

func (s *StackParser) pushParStack() *parseErr {
	if s.tokenAt(0).V == lexer.LPAR {
		s.parCounter++
		s.ptr++
		return nil
	}
	return &expectDelimiter
}

func (s *StackParser) popParStack() *parseErr {
	if s.tokenAt(0).V == lexer.RPAR {
		s.parCounter--
		if terminal := s.next(); terminal == EOF && s.parCounter != 0 {
			return &expectParenthesisClosed
		}
		return nil
	}
	return &expectDelimiter
}

func NewParser(tokens []lexer.Token) *StackParser {
	return &StackParser{
		ptr:        0,
		depth:      0,
		parCounter: 0,
		stack:      tokens,
		State:      ErrorState{},
		Metadata: map[string][]string{
			"ct_tables":      {},
			"insert_tables":  {},
			"insert_values":  {},
			"select_tables":  {},
			"select_columns": {},
		},
	}
}

func (s *StackParser) NewErrState(err error) ErrorState {
	return ErrorState{
		Message: err,
		Token:   *s.tokenAt(0),
		Ptr:     s.ptr,
	}
}

func (s *StackParser) Parse() *parseErr {
	if len(s.stack) == 0 {
		return &emptyStack
	}

	if DEBUG {
		for i, t := range s.stack {
			fmt.Printf("[%d] L=%s V=%d T=%d\n", i, t.L, t.V, t.T)
		}
	}

	return s.script()
}

func (s *StackParser) script() *parseErr {
	for {
		if err := s.statement(); err != nil {
			s.State = s.NewErrState(err)
			fmt.Printf("failed at ptr=%d token=%v\n", s.ptr, s.stack[s.ptr])
			return err
		}

		if err := s.expect(lexer.SCOLON, expectDelimiter); err != nil {
			s.State = s.NewErrState(err)
			return err
		}

		if next := s.next(); next == EOF {
			return nil
		}
	}
}

func (s *StackParser) statement() *parseErr {
	switch s.stack[s.ptr].V {
	case lexer.SELECT:
		return s.dml_expr()
	case lexer.CREATE:
		return s.ddl_expr()
	case lexer.INSERT:
		return s.expr_insert()
	default:
		return &expectedKeywordAtStart
	}
}

func (s *StackParser) safeRecursion() *parseErr {
	s.depth++
	if s.depth > maxRecursion {
		return &maxRecursionReached
	}
	return nil
}

func (s *StackParser) unwind() {
	s.depth--
}

// DML_EXPR := SELECT_EXPR FROM_EXPR [ WHERE_CLAUSE ]
func (s *StackParser) dml_expr() *parseErr {
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

// SELECT_EXPR := SELECT * | SELECT COLUMNS_EXPR
func (s *StackParser) select_expr() *parseErr {
	if err := s.expect(lexer.SELECT, expectKeyword); err != nil {
		return err
	}

	s.next()
	terminal := s.peekAt(0)

	switch {
	case terminal == lexer.TIMES:
		s.appendValue("select_columns", "*")
		s.next()
	case s.isIdentifier():
		return s.columns_expr()
	default:
		return &expectedColOrStar
	}
	return nil
}

func (s *StackParser) columns_expr() *parseErr {
	if err := s.name_expr(); err != nil {
		return err
	}
	// registrar columna: puede ser "A" o "TABLA.A"
	s.registerSelectColumn()
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
		s.registerSelectColumn()
		s.next()
	}
}

func (s *StackParser) registerSelectColumn() {
	t := s.tokenAt(0)
	if t == nil {
		return
	}
	col := strings.ToUpper(string(t.L))

	prev := s.tokenAt(-1)
	prevPrev := s.tokenAt(-2)
	if prev != nil && prev.V == lexer.DOT &&
		prevPrev != nil && s.isIdentifierValue(prevPrev.V) {
		table := strings.ToUpper(string(prevPrev.L))
		s.appendValue("select_columns", table+"."+col)
		return
	}

	s.appendValue("select_columns", col)
}

func (s *StackParser) isIdentifierValue(v lexer.Value) bool {
	return v > 400 && v < 600
}

// NAME_EXPR := IDENTIFIER | IDENTIFIER . IDENTIFIER
func (s *StackParser) name_expr() *parseErr {
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

	if s.peekAt(1) == lexer.DOT {
		return &expectDelimiter
	}

	return nil
}

// FROM_EXPR := FROM DATABASES_EXPR
func (s *StackParser) from_expr() *parseErr {
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
func (s *StackParser) databases_expr() *parseErr {
	if err := s.safeRecursion(); err != nil {
		return err
	}
	defer s.unwind()

	if err := s.table_expr(); err != nil {
		return err
	}

	for {
		terminal := s.peekAt(0)
		if terminal == EOF ||
			terminal == lexer.WHERE ||
			terminal == lexer.RPAR ||
			terminal == lexer.SCOLON {
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

// TABLE_EXPR := NAME_EXPR [ ALIAS ] | ( DML_EXPR ) ALIAS
func (s *StackParser) table_expr() *parseErr {
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
	s.registerSelectTable()

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

func (s *StackParser) registerSelectTable() {
	t := s.tokenAt(0)
	if t == nil {
		return
	}
	name := strings.ToUpper(string(t.L))

	prev := s.tokenAt(-1)
	prevPrev := s.tokenAt(-2)
	if prev != nil && prev.V == lexer.DOT &&
		prevPrev != nil && s.isIdentifierValue(prevPrev.V) {
		s.appendValue("select_tables", name)
		return
	}

	s.appendValue("select_tables", name)
}

func (s *StackParser) expectTableTerminator() *parseErr {
	terminal := s.peekAt(0)
	if terminal == EOF ||
		terminal == lexer.WHERE ||
		terminal == lexer.RPAR ||
		terminal == lexer.COMMA ||
		terminal == lexer.SCOLON {
		return nil
	}
	return &unexpectedToken
}

// WHERE_CLAUSE := WHERE CONDITION_EXPR
func (s *StackParser) where_clause() *parseErr {
	if err := s.expect(lexer.WHERE, expectKeyword); err != nil {
		return err
	}

	s.next()

	if err := s.condition_expr(); err != nil {
		return err
	}

	return nil
}

// CONDITION_EXPR := [NOT] NAME_EXPR RELATION_EXPR CONSTANT  { (AND | OR) [NOT] CONDITION_EXPR }
//
//	| [NOT] NAME_EXPR RELATION_EXPR NAME_EXPR { (AND | OR) [NOT] CONDITION_EXPR }
//	| NAME_EXPR [NOT] WHERE_SUBQUERY
func (s *StackParser) condition_expr() *parseErr {
	if err := s.safeRecursion(); err != nil {
		return err
	}
	defer s.unwind()

	if err := s.name_expr(); err != nil {
		return err
	}

	whereCol := s.lexemeAt(0)

	s.next()
	if s.peekAt(0) == lexer.NOT {
		s.next()
	}

	if s.peekAt(0) == lexer.IN {
		return s.where_subquery_expr()
	}

	if err := s.relation_expr(); err != nil {
		return &expectRelational
	}

	whereOp := string(s.tokenAt(0).L)

	s.next()
	switch {
	case s.isConstant():
		whereVal := string(s.tokenAt(0).L)
		s.appendValue("where_col", whereCol)
		s.appendValue("where_op", whereOp)
		s.appendValue("where_val", whereVal)
	case s.isIdentifier():
		if err := s.name_expr(); err != nil {
			return err
		}
	default:
		return &expectedConstOrIdent
	}

	if terminal := s.next(); terminal == EOF ||
		terminal == lexer.RPAR ||
		terminal == lexer.SCOLON {
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

// WHERE_SUBQUERY := IN ( DML_EXPR ) [ AND | OR [ NOT ] CONDITION_EXPR ]
func (s *StackParser) where_subquery_expr() *parseErr {
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
	if terminal == EOF || terminal == lexer.RPAR || terminal == lexer.SCOLON {
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

// RELATION_EXPR := = | < | <= | > | >= | <>
func (s *StackParser) relation_expr() *parseErr {
	if !s.isRelation() {
		return &expectRelational
	}
	return nil
}
