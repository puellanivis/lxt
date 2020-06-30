package parser

import (
	"context"

	"github.com/puellanivis/lxt/tokenizer"
	"github.com/puellanivis/lxt/xslt"
)

func (r *Reader) parseSpan(ctx context.Context) (*xslt.Element, error) {
	className, err := r.read(ctx)
	if err != nil {
		return nil, err
	}

	switch className.Type {
	case tokenizer.TokenTypeIdentifier:
	case tokenizer.TokenTypeDoubleQuote, tokenizer.TokenTypeSingleQuote:
	default:
		return nil, r.parseError("expected a class name")
	}
	r.consume()

	if className.Value == "" {
		return nil, r.parseError("span cannot have an empty class name")
	}

	body, err := r.parseExpression(ctx)
	if err != nil {
		return nil, err
	}

	return &xslt.Element{
		Name: "span",
		Body: xslt.Group{
			&xslt.Attribute{
				Name:  "class",
				Value: className.Value,
			},
			body,
		},
	}, nil
}

func (r *Reader) parseDiv(ctx context.Context) (*xslt.Element, error) {
	className, err := r.read(ctx)
	if err != nil {
		return nil, err
	}

	switch className.Type {
	case tokenizer.TokenTypeIdentifier:
	case tokenizer.TokenTypeDoubleQuote, tokenizer.TokenTypeSingleQuote:
	default:
		return nil, r.parseError("expected a class name")
	}
	r.consume()

	if className.Value == "" {
		return nil, r.parseError("div cannot have an empty class name")
	}

	body, err := r.parseExpression(ctx)
	if err != nil {
		return nil, err
	}

	return &xslt.Element{
		Name: "div",
		Body: xslt.Group{
			&xslt.Attribute{
				Name:  "class",
				Value: className.Value,
			},
			body,
		},
	}, nil
}
