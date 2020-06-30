package tokenizer

import (
	"fmt"
)

type TokenType int

// Token Types
const (
	TokenTypeEmpty = TokenType(iota)
	TokenTypeError
	TokenTypeEOF
	TokenTypeComma
	TokenTypeBeginGroup
	TokenTypeEndGroup
	TokenTypeOperator
	TokenTypeNumber
	TokenTypeIdentifier
	TokenTypeXPath

	TokenTypeBackQuote   = TokenType('`')
	TokenTypeSingleQuote = TokenType('\'')
	TokenTypeDoubleQuote = TokenType('"')
)

func (t TokenType) String() string {
	switch t {
	case TokenTypeEmpty:
		return "EMPTY"
	case TokenTypeError:
		return "ERR"
	case TokenTypeEOF:
		return "EOF"
	case TokenTypeComma:
		return "COMMA"
	case TokenTypeBeginGroup:
		return "BEGIN"
	case TokenTypeEndGroup:
		return "END"
	case TokenTypeOperator:
		return "OP"
	case TokenTypeNumber:
		return "NUM"
	case TokenTypeIdentifier:
		return "IDENT"
	case TokenTypeBackQuote:
		return "BQ"
	case TokenTypeSingleQuote:
		return "SQ"
	case TokenTypeDoubleQuote:
		return "DQ"
	case TokenTypeXPath:
		return "XP"
	}

	return fmt.Sprintf("UNKNOWN%d", int(t))
}
