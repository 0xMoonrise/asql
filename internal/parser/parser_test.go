package parser

import (
	"testing"

	"github.com/0xMoonrise/asql/internal/lexer"
	"github.com/stretchr/testify/assert"
)

func tokenize(input string) []lexer.Token {
	raw := lexer.Tokenize(input)
	tokens, _ := lexer.Lexer([][]string{raw})
	return tokens
}

func parse(input string) *parseErr {
	tokens := tokenize(input)
	p := NewParser(tokens)
	return p.Parse()
}

type validCase struct {
	query string
	desc  string
}

type invalidCase struct {
	query string
	desc  string
}

type errorCodeCase struct {
	query string
	desc  string
	code  int
}

func runValid(t *testing.T, cases []validCase) {
	t.Helper()
	for _, tc := range cases {
		t.Run(tc.desc, func(t *testing.T) {
			assert.Nil(t, parse(tc.query), tc.query)
		})
	}
}

func runInvalid(t *testing.T, cases []invalidCase) {
	t.Helper()
	for _, tc := range cases {
		t.Run(tc.desc, func(t *testing.T) {
			assert.NotNil(t, parse(tc.query), tc.query)
		})
	}
}

func runErrorCodes(t *testing.T, cases []errorCodeCase) {
	t.Helper()
	for _, tc := range cases {
		t.Run(tc.desc, func(t *testing.T) {
			err := parse(tc.query)
			assert.NotNil(t, err, tc.query)
			if err != nil {
				assert.Equal(t, tc.code, err.Code, tc.query)
			}
		})
	}
}

func TestSelectValid(t *testing.T) {
	runValid(t, []validCase{
		{"SELECT * FROM foo", "star"},
		{"SELECT a FROM foo", "single column"},
		{"SELECT a, b FROM foo", "two columns"},
		{"SELECT a, b, c FROM foo", "three columns"},
		{"SELECT a.b FROM foo", "qualified column"},
		{"SELECT a, a.b FROM foo", "simple and qualified"},
		{"SELECT a, a.b, c FROM foo", "full mix of columns"},
	})
}

func TestSelectInvalid(t *testing.T) {
	runInvalid(t, []invalidCase{
		{"SELECT *", "star without FROM"},
		{"SELECT a, b, FROM foo", "trailing comma"},
		{"SELECT FROM foo", "no columns and no star"},
		{"a, b FROM foo", "missing SELECT keyword"},
		{"", "empty input"},
	})
}

func TestColumnQualificationValid(t *testing.T) {
	runValid(t, []validCase{
		{"SELECT a.b FROM foo", "two-level qualification"},
	})
}

func TestColumnQualificationInvalid(t *testing.T) {
	runInvalid(t, []invalidCase{
		{"SELECT a.b.c FROM foo", "three-level qualification"},
		{"SELECT a.b.c.d FROM foo", "four-level qualification"},
		{"SELECT a.b.c, d FROM foo", "over-qualified first column"},
		{"SELECT d, a.b.c FROM foo", "over-qualified second column"},
	})
}

func TestFromValid(t *testing.T) {
	runValid(t, []validCase{
		{"SELECT * FROM foo", "single table"},
		{"SELECT * FROM foo, bar", "two tables"},
		{"SELECT * FROM foo, bar, baz", "three tables"},
		{"SELECT * FROM db.bar", "qualified table"},
		{"SELECT * FROM db.bar, other", "qualified and simple table"},
	})
}

func TestFromInvalid(t *testing.T) {
	runInvalid(t, []invalidCase{
		{"SELECT *", "missing FROM entirely"},
		{"SELECT * FROM", "FROM without table name"},
		{"SELECT * FROM foo,", "trailing comma after table"},
		{"SELECT * FROM foo, WHERE a = 1", "comma before WHERE"},
	})
}

func TestAliasValid(t *testing.T) {
	runValid(t, []validCase{
		{"SELECT * FROM foo A", "single table with alias"},
		{"SELECT * FROM foo A, bar B", "two tables with alias"},
		{"SELECT * FROM foo A, bar B, baz C", "three tables with alias"},
		{"SELECT * FROM foo A, bar", "mixed alias and no alias"},
		{"SELECT * FROM db.foo A", "qualified table with alias"},
		{"SELECT * FROM foo A WHERE a = 1", "alias then WHERE"},
	})
}

func TestAliasInvalid(t *testing.T) {
	runInvalid(t, []invalidCase{
		{"SELECT * FROM foo A B", "double alias"},
		{"SELECT * FROM foo A, WHERE a = 1", "alias then comma before WHERE"},
	})
}

func TestWhereValid(t *testing.T) {
	runValid(t, []validCase{
		{"SELECT * FROM foo WHERE a = 1", "equality with integer"},
		{"SELECT * FROM foo WHERE a < 1", "less than"},
		{"SELECT * FROM foo WHERE a <= 1", "less than or equal"},
		{"SELECT * FROM foo WHERE a > 1", "greater than"},
		{"SELECT * FROM foo WHERE a >= 1", "greater than or equal"},
		{"SELECT * FROM foo WHERE a <> 1", "not equal"},
		{"SELECT * FROM foo WHERE a = 'hello'", "equality with string"},
		{"SELECT a FROM foo WHERE a = 1", "columns with WHERE"},
		{"SELECT a FROM foo WHERE a = 1 OR b = 2", "boolean OR"},
		{"SELECT a FROM foo WHERE a = 1 AND b = 2", "boolean AND"},
		{"SELECT a FROM foo WHERE a = 1 OR b = 2 AND c = 'x'", "compound conditions"},
	})
}

func TestWhereInvalid(t *testing.T) {
	runInvalid(t, []invalidCase{
		{"SELECT * FROM foo WHERE", "WHERE with nothing after"},
		{"SELECT * FROM foo WHERE = 1", "WHERE missing identifier"},
		{"SELECT * FROM foo WHERE a", "WHERE missing relation and constant"},
		{"SELECT * FROM foo WHERE a =", "WHERE missing constant"},
		{"SELECT * FROM foo WHERE 1 = a", "WHERE constant before identifier"},
	})
}

func TestWhereIdentifierRHSValid(t *testing.T) {
	runValid(t, []validCase{
		{"SELECT * FROM foo, bar WHERE foo.a = bar.a", "join on qualified columns"},
		{"SELECT * FROM foo, bar WHERE foo.a = bar.a AND foo.b = bar.b", "two join conditions"},
		{"SELECT * FROM foo, bar WHERE a = b", "simple identifier rhs"},
		{"SELECT * FROM foo A, bar B WHERE A.a = B.a", "join with aliases"},
		{
			"SELECT ANOMBRE FROM ALUMNOS A, INSCRITOS I, CARRERAS C WHERE A.A#=I.A# AND A.C#=C.C# AND I.SEMESTRE='2010I' AND C.CNOMBRE='ISC'",
			"natural join full query",
		},
	})
}

func TestWhereIdentifierRHSInvalid(t *testing.T) {
	runInvalid(t, []invalidCase{
		{"SELECT * FROM foo WHERE a = b.c.d", "over-qualified rhs"},
		{"SELECT * FROM foo WHERE a = .b", "rhs starts with dot"},
	})
}

func TestSubqueryInWhereValid(t *testing.T) {
	runValid(t, []validCase{
		{
			"SELECT * FROM foo WHERE a IN (SELECT a FROM bar)",
			"simple IN subquery",
		},
		{
			"SELECT * FROM foo WHERE a IN (SELECT a FROM bar WHERE b = 1)",
			"IN subquery with WHERE",
		},
		{
			"SELECT A# FROM ALUMNOS WHERE A# IN (SELECT A# FROM INSCRITOS WHERE P# IN (SELECT P# FROM PROFESORES WHERE GRADO='MAE'))",
			"nested IN subquery two levels",
		},
		{
			"SELECT ANOMBRE FROM ALUMNOS WHERE A# IN (SELECT A# FROM INSCRITOS WHERE P# IN (SELECT P# FROM PROFESORES WHERE GRADO='MAE')) AND C# IN (SELECT C# FROM CARRERAS WHERE CNOMBRE='ISC')",
			"IN subquery with AND after closing paren",
		},
	})
}

func TestSubqueryInWhereInvalid(t *testing.T) {
	runInvalid(t, []invalidCase{
		{
			"SELECT * FROM foo WHERE a IN SELECT a FROM bar",
			"IN subquery missing parenthesis",
		},
		{
			"SELECT * FROM foo WHERE a IN (SELECT a FROM bar",
			"IN subquery unclosed parenthesis",
		},
		{
			"SELECT ANOMBRE FROM ALUMNOS WHERE A# IN (SELECT A# FROM INSCRITOS WHERE P# IN (SELECT P# FROM PROFESORES WHERE GRADO='MAE')) AND C# IN (SELECT C# FROM WHERE CNOMBRE='ISC')",
			"nested subquery FROM without table",
		},
	})
}

func TestSubqueryInFromValid(t *testing.T) {
	runValid(t, []validCase{
		{
			"SELECT * FROM (SELECT a FROM foo) A",
			"simple subquery in FROM",
		},
		{
			"SELECT * FROM foo, (SELECT a FROM bar) B",
			"table and subquery in FROM",
		},
		{
			"SELECT C.M#, MNOMBRE FROM MATERIAS, (SELECT M# FROM INSCRITOS WHERE A# IN (SELECT A# FROM ALUMNOS WHERE ANOMBRE='MESSI LIONEL')) C WHERE MATERIAS.M#=C.M#",
			"integrated natural join full query",
		},
		{
			"SELECT * FROM (SELECT a FROM foo WHERE b = 1) A, bar B WHERE A.a = B.a",
			"subquery in FROM with join",
		},
	})
}

func TestSubqueryInFromInvalid(t *testing.T) {
	runInvalid(t, []invalidCase{
		{
			"SELECT * FROM (SELECT a FROM foo)",
			"subquery in FROM missing alias",
		},
		{
			"SELECT * FROM (SELECT a FROM foo) A B",
			"subquery in FROM double alias",
		},
		{
			"SELECT * FROM (SELECT a FROM foo",
			"subquery in FROM unclosed parenthesis",
		},
		{
			"SELECT * FROM () A",
			"empty subquery in FROM",
		},
	})
}

func TestErrorCodes(t *testing.T) {
	runErrorCodes(t, []errorCodeCase{
		{"FROM foo", "missing SELECT", 201},
		{"SELECT * foo", "missing FROM keyword", 201},

		{"SELECT a, FROM foo", "trailing comma in columns", 204},
		{"SELECT * FROM foo WHERE = 1", "WHERE missing identifier", 204},
		{"SELECT * FROM (SELECT a FROM foo)", "subquery in FROM missing alias", 204},

		{"SELECT a.b.c FROM foo", "three-level qualified column", 205},
		{"SELECT a, a.b.c FROM foo", "over-qualified in column list", 205},
		{"SELECT * FROM foo WHERE a =", "WHERE missing constant", 210},
		{"SELECT * FROM foo WHERE a IN (SELECT a FROM bar", "unclosed IN subquery", 205},

		{"SELECT * FROM foo WHERE a 1", "WHERE missing relational operator", 208},
	})
}
