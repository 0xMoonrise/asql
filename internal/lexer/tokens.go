package lexer

type Lexeme string
type Value int
type Typ int

// keywords
const (
	SELECT = iota + 10
	FROM
	WHERE
	IN
	AND
	OR
	CREATE
	TABLE
	CHAR
	NUMERIC
	NOT
	NULL
	CONSTRAIN
	KEY
	PRIMARY
	//foreign
	FOREIGN
	//references
	REFERENCES
	INSERT
	INTO
	VALUES
	DATE
	CHECK
)

// delimitators
const (
	COMMA = iota + 50
	DOT
	LPAR
	RPAR
	APOS
)

// operators
const (
	PLUS = iota + 70
	MINUS
	TIMES
	DIVS
)

// constants
const (
	NUMER = iota + 61
	ALPHA
)

// relations
const (
	GT = iota + 81
	LT
	EQ
	GE
	LE
	NE
)
