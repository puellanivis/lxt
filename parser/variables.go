package parser

import (
	"context"

	"github.com/puellanivis/lxt/tokenizer"
	"github.com/puellanivis/lxt/xslt"
)

func (r *Reader) parseParamList(ctx context.Context) ([]*xslt.Param, error) {
	tok, err := r.peak(ctx)
	if err != nil {
		return nil, err
	}

	var end tokenizer.Token
	switch tok.Type {
	case tokenizer.TokenTypeBeginGroup:
		r.consume()
		end = endTokenFromStart(tok)

	default:
		return nil, r.parseError("expected group")
	}

	var params []*xslt.Param

	for {
		tok, err := r.peak(ctx)
		if tok.Type == tokenizer.TokenTypeComma {
			tok, err = r.read(ctx)
		}

		if tok.Type == tokenizer.TokenTypeEndGroup {
			if tok != end {
				return nil, r.parseErrorf("unexpected end param list token, was expecting: %s", end)
			}

			r.consume()
			return params, nil
		}

		if err != nil {
			return nil, err
		}

		param, err := r.parseParam(ctx)
		if err != nil {
			return nil, err
		}

		params = append(params, param)
	}
}

func (r *Reader) parseArgumentList(ctx context.Context) ([]*xslt.WithParam, error) {
	tok, err := r.peak(ctx)
	if err != nil {
		return nil, err
	}

	var end tokenizer.Token
	switch tok.Type {
	case tokenizer.TokenTypeBeginGroup:
		r.consume()
		end = endTokenFromStart(tok)

	default:
		return nil, r.parseError("expected group")
	}

	var args []*xslt.WithParam

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
			return args, nil
		}

		if err != nil {
			return nil, err
		}

		arg, err := r.parseArgument(ctx)
		if err != nil {
			return nil, err
		}

		args = append(args, arg)
	}
}

func (r *Reader) parseArgument(ctx context.Context) (*xslt.WithParam, error) {
	v, err := r.parseVariable(ctx)
	return (*xslt.WithParam)(v), err
}

func (r *Reader) parseParam(ctx context.Context) (*xslt.Param, error) {
	v, err := r.parseVariable(ctx)
	return (*xslt.Param)(v), err
}

func (r *Reader) parseVariable(ctx context.Context) (*xslt.Variable, error) {
	ident, err := r.peak(ctx)
	if err != nil {
		return nil, err
	}

	name := ident.Value
	if ident.Type != tokenizer.TokenTypeIdentifier {
		if ident.Type != tokenizer.TokenTypeXPath || !tokenizer.IsIdent(ident.Value) {
			return nil, r.parseError("expected identifier")
		}
		name = ident.Value[1:]
	}

	if name == "" {
		return nil, r.parseError("param cannot have empty name")
	}

	val, err := r.read(ctx)
	if val.Type == tokenizer.TokenTypeComma {
		val, err = r.read(ctx)
	}

	if err != nil {
		return nil, err
	}

	switch val.Type {
	case tokenizer.TokenTypeXPath, tokenizer.TokenTypeNumber:
		r.consume()
		return &xslt.Variable{
			Name:   name,
			Select: val.Value,
		}, nil

	case tokenizer.TokenTypeDoubleQuote, tokenizer.TokenTypeSingleQuote:
		if val.Value == "" {
			r.consume()
			return &xslt.Variable{
				Name: name,
			}, nil
		}
	}

	value, err := r.parseExpression(ctx)
	if err != nil {
		return nil, err
	}

	return &xslt.Variable{
		Name:  name,
		Value: value,
	}, nil
}
