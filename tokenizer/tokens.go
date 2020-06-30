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

var (
	Empty = Token{Type: TokenTypeEmpty}
	EOF   = Token{Type: TokenTypeEOF}
)
