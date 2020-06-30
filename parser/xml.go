package parser

import (
	"context"

	"github.com/puellanivis/lxt/tokenizer"
	"github.com/puellanivis/lxt/xslt"
)

func (r *Reader) parseTag(ctx context.Context) (*xslt.Element, error) {
	name, err := r.read(ctx)
	if err != nil {
		return nil, err
	}

	if name.Type != tokenizer.TokenTypeIdentifier {
		return nil, r.parseError("expected identifier")
	}
	r.consume()

	if name.Value == "" {
		return nil, r.parseError("tag cannot have empty name")
	}

	body, err := r.parseExpression(ctx)
	if err != nil {
		return nil, err
	}

	return &xslt.Element{
		Name: name.Value,
		Body: body,
	}, nil
}

func (r *Reader) parseAttribs(ctx context.Context) ([]*xslt.Attribute, error) {
	tok, err := r.read(ctx)
	if err != nil {
		return nil, err
	}

	if tok.Type != tokenizer.TokenTypeBeginGroup || tok.Value != "(" {
		return nil, r.parseError("unexpected token: expected '('")
	}
	r.consume()
	end := endTokenFromStart(tok)

	var attribs []*xslt.Attribute

	for {
		tok, err := r.peak(ctx)
		if tok.Type == tokenizer.TokenTypeComma {
			tok, err = r.read(ctx)
		}

		if tok.Type == tokenizer.TokenTypeEndGroup {
			if tok != end {
				return nil, r.parseErrorf("unexpected end argument list token, was expecting: %s", end)
			}

			r.consume()
			return attribs, nil
		}

		if err != nil {
			return nil, err
		}

		switch tok.Type {
		case tokenizer.TokenTypeIdentifier:
		case tokenizer.TokenTypeDoubleQuote, tokenizer.TokenTypeSingleQuote:
		default:
			return nil, r.parseError("expected identifier")
		}

		r.consume()

		val, err := r.parseExpression(ctx)
		if err != nil {
			return nil, err
		}

		attribs = append(attribs, &xslt.Attribute{
			Name:  tok.Value,
			Value: val,
		})
	}
}
