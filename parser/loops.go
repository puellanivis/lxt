package parser

import (
	"context"

	"github.com/puellanivis/lxt/tokenizer"
	"github.com/puellanivis/lxt/xslt"
)

func (r *Reader) parseApplyTemplates(ctx context.Context) (*xslt.ApplyTemplates, error) {
	tok, err := r.read(ctx)
	if err != nil {
		return nil, err
	}

	var xpath string
	if tok.Type == tokenizer.TokenTypeXPath {
		xpath = tok.Value

		tok, err = r.read(ctx)
	}

	if tok.Type != tokenizer.TokenTypeBeginGroup || tok.Value != "(" {
		return &xslt.ApplyTemplates{
			Select: xpath,
		}, nil
	}

	args, err := r.parseArgumentList(ctx)
	if err != nil {
		return nil, err
	}

	return &xslt.ApplyTemplates{
		Select:     xpath,
		WithParams: args,
	}, nil
}

func (r *Reader) parseForEach(ctx context.Context) (*xslt.ForEach, error) {
	set, err := r.read(ctx)
	if err != nil {
		return nil, err
	}

	if set.Type != tokenizer.TokenTypeXPath {
		return nil, r.parseError("expected xpath")
	}
	r.consume()

	if set.Value == "" {
		return nil, r.parseError("for-each cannot have empty name")
	}

	body, err := r.parseExpression(ctx)
	if err != nil {
		return nil, err
	}

	return &xslt.ForEach{
		Select: set.Value,
		Body:   body,
	}, nil
}
