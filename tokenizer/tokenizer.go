package tokenizer

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"unicode"
	"unicode/utf8"
)

type Reader struct {
	S *bufio.Scanner

	lineno int
	line   []byte

	off int
}

func (r *Reader) CurrentLine() int {
	return r.lineno
}

func (r *Reader) startNewToken() error {
	for len(r.line) < 1 {
		if !r.S.Scan() {
			if err := r.S.Err(); err != nil {
				return err
			}

			return io.EOF
		}

		r.line = bytes.TrimSpace(r.S.Bytes())
		r.lineno++
	}

	r.line = bytes.TrimSpace(r.line)
	r.off = 0

	return nil
}

const errInvalidCharacter = "invalid character"

func (r *Reader) peak() (rune, int, error) {
	if r.off >= len(r.line) {
		return '\n', 0, nil
	}

	char, sz := utf8.DecodeRune(r.line[r.off:])

	if char == utf8.RuneError && sz == 1 {
		return char, sz, errors.New("invalid UTF-8")
	}

	return char, sz, nil
}

func (r *Reader) advance(n int) {
	r.off += n
}

func any(r rune) bool {
	return true
}

func (r *Reader) next(mustBe func(rune) bool) (rune, int, error) {
	char, sz, err := r.peak()
	r.advance(sz)

	if err != nil {
		return char, sz, err
	}

	if !mustBe(char) {
		return char, sz, errors.New(errInvalidCharacter)
	}

	return char, sz, nil
}

func (r *Reader) bytesSlice(s, e int) []byte {
	length, text := r.off, r.line[s:r.off-e]
	r.line, r.off = r.line[length:], 0
	return text
}

func (r *Reader) bytes() []byte {
	return r.bytesSlice(0, 0)
}

func (r *Reader) text() string {
	return string(r.bytes())
}

func (r *Reader) textSlice(s, e int) string {
	return string(r.bytesSlice(s, e))
}

func unquote(in []byte) string {
	out := make([]byte, 0, len(in))

	for i := 0; i < len(in); i++ {
		b := in[i]

		switch b {
		case '\\':
			i++
			if i >= len(in) {
				return string(out)
			}

			b := in[i]

			switch b {
			case 'n':
				out = append(out, '\n')
			case 'r':
				out = append(out, '\r')
			case 't':
				out = append(out, '\t')
			default:
				out = append(out, b)
			}

		default:
			out = append(out, b)
		}
	}

	return string(out)
}

func (r *Reader) readQuote(quoteChar rune) (int, error) {
	for {
		char, sz, err := r.next(any)
		if err != nil {
			return 0, err
		}

		switch char {
		case quoteChar:
			return sz, nil

		case '\\':
			if _, _, err := r.next(any); err != nil {
				return 0, err
			}
		}
	}
}

func (r *Reader) readBackQuote() (int, error) {
	for {
		char, sz, err := r.next(any)
		if err != nil {
			return 0, fmt.Errorf("reading raw quote: %w", err)
		}

		switch char {
		case '`':
			return sz, nil
		}
	}
}

func identInitial(r rune) bool {
	return unicode.Is(identInitialRange, r)
}

func identFollowing(r rune) bool {
	if identInitial(r) {
		return true
	}

	return unicode.Is(identFollowingRange, r)
}

func IsIdent(s string) bool {
	if len(s) < 1 {
		return false
	}

	if s[0] == '$' || s[0] == '@' {
		s = s[1:]
	}

	if len(s) < 1 {
		return false
	}

	for i, r := range s {
		if i == 0 {
			if !identInitial(r) {
				return false
			}
			continue
		}

		if !identFollowing(r) {
			return false
		}
	}

	return true
}

func (r *Reader) readIdent() error {
	for {
		char, sz, err := r.peak()
		if err != nil {
			r.advance(sz)
			return err
		}

		switch {
		case char == ':':
			r.advance(sz)

			if _, _, err := r.next(identInitial); err != nil {
				return err
			}

		case !identFollowing(char):
			return nil

		default:
			r.advance(sz)
		}
	}
}

func numInitial(r rune) bool {
	return unicode.IsNumber(r)
}

func numFollowing(r rune) bool {
	switch r {
	case '_', '-':
		return true
	}

	return unicode.IsNumber(r)
}

func fixNumber(in []byte) string {
	out := make([]byte, 0, len(in))

	for _, b := range in {
		switch b {
		case '-', '_':
		default:
			out = append(out, b)
		}
	}

	return string(out)
}

func (r *Reader) readNumber() error {
	for {
		char, sz, err := r.peak()
		if err != nil {
			r.advance(sz)
			return err
		}

		if char == '.' {
			r.advance(sz)
			break
		}

		if !numFollowing(char) {
			return nil
		}

		r.advance(sz)
	}

	for {
		char, sz, err := r.peak()
		if err != nil {
			r.advance(sz)
			return err
		}

		if !numFollowing(char) {
			return nil
		}

		r.advance(sz)
	}
}

func simpleXPath(r rune) bool {
	switch r {
	case '/', '*', '.':
		return true
	}

	return identFollowing(r)
}

func is(test rune) func(rune) bool {
	return func(r rune) bool {
		return r == test
	}
}

func (r *Reader) readSimpleXPath() (int, error) {
	for {
		char, sz, err := r.next(any)
		if err != nil {
			return 0, err
		}

		if char == '$' || char == '@' {
			break
		}

		switch {
		case char == '>':
			return sz, nil

		case char == ':':
			if _, _, err := r.next(is(':')); err != nil {
				return 0, err
			}

		case !simpleXPath(char):
			return 0, errors.New(errInvalidCharacter)
		}
	}

	if _, _, err := r.next(identInitial); err != nil {
		return 0, err
	}

	if err := r.readIdent(); err != nil {
		return 0, err
	}

	_, sz, err := r.next(is('>'))
	if err != nil {
		return 0, err
	}

	return sz, nil
}

func (r *Reader) readComplexXPath() (int, error) {
	for {
		char, sz, err := r.next(any)
		if err != nil {
			return 0, err
		}

		switch char {
		case '}':
			char, sz2, err := r.next(any)
			if err != nil {
				return 0, err
			}

			if char == '>' {
				return sz + sz2, nil
			}

		case '\\':
			if _, _, err := r.next(any); err != nil {
				return 0, err
			}
		}
	}
}

func (r *Reader) ReadToken() (Token, error) {
	if err := r.startNewToken(); err != nil {
		if err == io.EOF {
			return EOF, io.EOF
		}

		return Token{
			Type:  TokenTypeError,
			Value: "",
		}, err
	}

	char, sz, err := r.next(any)
	if err != nil {
		return Token{
			Type:  TokenTypeError,
			Value: r.text(),
		}, err
	}

	switch char {
	case '"', '\'':
		e, err := r.readQuote(char)

		return Token{
			Type:  TokenType(char),
			Value: unquote(r.bytesSlice(sz, e)),
		}, err

	case '`':
		e, err := r.readBackQuote()

		return Token{
			Type:  TokenType(char),
			Value: unquote(r.bytesSlice(sz, e)),
		}, err

	case '$', '@':
		if _, _, err := r.next(identInitial); err != nil {
			return Token{
				Type:  TokenTypeXPath,
				Value: r.text(),
			}, err
		}

		err := r.readIdent()

		return Token{
			Type:  TokenTypeXPath,
			Value: r.text(),
		}, err

	case '<':
		if ch, sz2, _ := r.peak(); ch == '{' {
			r.advance(sz2)

			e, err := r.readComplexXPath()
			return Token{
				Type:  TokenTypeXPath,
				Value: string(bytes.TrimSpace(r.bytesSlice(sz+sz2, e))),
			}, err
		}

		e, err := r.readSimpleXPath()
		return Token{
			Type:  TokenTypeXPath,
			Value: r.textSlice(sz, e),
		}, err

	case '=':
		if ch, sz, _ := r.peak(); ch == '>' {
			r.advance(sz)
		}

		return Token{
			Type:  TokenTypeComma,
			Value: r.text(),
		}, nil

	case ':':
		if ch, sz, _ := r.peak(); ch == '=' {
			r.advance(sz)
		}

		return Token{
			Type:  TokenTypeComma,
			Value: r.text(),
		}, nil

	case ',':
		return Token{
			Type:  TokenTypeComma,
			Value: r.text(),
		}, nil

	case '(', '[', '{':
		return Token{
			Type:  TokenTypeBeginGroup,
			Value: r.text(),
		}, nil

	case ')', ']', '}':
		return Token{
			Type:  TokenTypeEndGroup,
			Value: r.text(),
		}, nil

	case ';':
		return Token{
			Type:  TokenTypeOperator,
			Value: r.text(),
		}, nil
	}

	switch {
	case identInitial(char):
		err := r.readIdent()
		return Token{
			Type:  TokenTypeIdentifier,
			Value: r.text(),
		}, err

	case numInitial(char):
		err := r.readNumber()
		return Token{
			Type:  TokenTypeNumber,
			Value: fixNumber(r.bytes()),
		}, err
	}

	return Token{
		Type:  TokenTypeError,
		Value: r.text(),
	}, nil
}
