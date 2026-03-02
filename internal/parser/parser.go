package parser

import (
	"asql/internal/lexer"
	"fmt"
)

// Start
// Push('$')              ->  stack := []int{199, 300}
// Push(300)
// Append '$' to the end of the Token Table
// PTR = Pointer to the first Token in the Token Table  ->  pos := 0
// Repeat
//     X = Pop()                          ->  top := stack[len(stack)-1]
//     K = TokenTable[PTR]                ->  cur := input[pos]
//
//     If (X = TERMINAL) or (X = '$') then
//         If (X = K) then
//             Advance PTR                ->  pos++  and  pop stack
//         Else
//             ERROR()                    ->  unexpected token, X != K
//
//     Else
//         If (ParseTable[X,K] = PRODUCTION) then
//             If (PRODUCTION <> 'λ') then
//                 Push(PRODUCTION) in reverse order
//         Else
//             ERROR()                    ->  empty cell in parse table
//
// Until X = '$'
// End

var parseTable = map[int]map[int][]int{
	300: {
		10: {10, 301, 11, 306, 310},
	},
	301: {
		4:  {302},
		72: {72},
	},
	302: {
		4: {304, 303},
	},
	303: {
		50:  {50, 302},
		199: {},
	},
	304: {
		4: {4, 305},
	},
	305: {
		8:   {99},
		10:  {99},
		13:  {99},
		14:  {99},
		15:  {99},
		50:  {51, 4},
		51:  {51, 4},
		53:  {99},
		199: {99},
	},
	306: {
		4: {308, 307},
	},
	307: {
		50:  {50, 306},
		99:  {},
		199: {},
	},
	308: {
		4: {4, 309},
	},
	309: {
		4:   {4},
		99:  {},
		199: {},
	},
	310: {
		12:  {12, 311},
		199: {},
	},
	311: {
		4: {313, 312},
	},
	312: {
		14:  {317, 311},
		15:  {317, 311},
		99:  {},
		199: {},
	},
	313: {
		4: {304, 314},
	},
	314: {
		8:  {315, 316},
		13: {13, 52, 300, 53},
	},
	315: {
		8: {8},
	},
	316: {
		4:  {304},
		54: {54, 318, 54},
		61: {319},
	},
	317: {
		14: {14},
		15: {15},
	},
	318: {
		62: {62},
	},
	319: {
		61: {61},
	},
}

func Parser(tokenStream []lexer.Token) error {
	fmt.Println(parseTable[300])
	return nil
}
