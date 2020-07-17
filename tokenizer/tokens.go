package tokenizer

import (
	"fmt"
)

type Token struct {
	Type  TokenType
	Value string
}

func (t Token) String() string {
	return fmt.Sprintf("%s(%q)", t.Type, t.Value)
}

// Sentinel and error values.
var (
	Empty = Token{Type: TokenTypeEmpty}
	EOF   = Token{Type: TokenTypeEOF}
)

// Operator tokens.
var (
	OperatorEquals = Token{
		Type:  TokenTypeOperator,
		Value: "=",
	}
	OperatorArrow = Token{
		Type:  TokenTypeOperator,
		Value: "=>",
	}
	OperatorDefine = Token{
		Type:  TokenTypeOperator,
		Value: ":=",
	}
)

// Comma tokens.
var (
	Comma = Token{
		Type:  TokenTypeComma,
		Value: ",",
	}
)
