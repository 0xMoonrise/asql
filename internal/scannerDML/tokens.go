package scannerdml

type lexeme string
type symbol rune
type value int
type typ int

const noSym symbol = 0

// keywords
const (
	__select = iota + 10
	__from
	__where
	__in
	__and
	__or
	__create
	__table
	__char
	__numeric
	__not
	__null
	__constraint
	__key
	__primary
	__foreign
	__references
	__insert
	__into
	__values
)

// delimitators
const (
	comma = iota + 50
	dot
	lparentheses
	rparentheses
	apostrophe
)

// operators
const (
	plus = iota + 70
	minus
	times
	divition
)

// constants
const (
	numeric = iota + 61
	alpha
)

// relations
const (
	gt = iota + 81
	lt
	eq
	ge
	le
)
