package parser

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/puellanivis/lxt/tokenizer"
	"github.com/puellanivis/lxt/xslt"
)

type Reader struct {
	filename string
	r        *tokenizer.Reader

	tok tokenizer.Token
	err error
}

func (r *Reader) read(ctx context.Context) (tokenizer.Token, error) {
	select {
	case <-ctx.Done():
		return tokenizer.Empty, r.parseError("parsing cancelled", ctx.Err())
	default:
	}

	tok, err := r.r.ReadToken()

	r.tok = tok
	if err != nil {
		r.err = r.parseError("tokenize error", err)
	}

	return r.tok, r.err
}

func (r *Reader) peak(ctx context.Context) (tokenizer.Token, error) {
	if r.tok == tokenizer.Empty {
		return r.read(ctx)
	}

	select {
	case <-ctx.Done():
		return tokenizer.Empty, r.parseError("parsing cancelled", ctx.Err())
	default:
	}

	return r.tok, r.err
}

func (r *Reader) consume() {
	r.tok = tokenizer.Empty
}

func (r *Reader) parseErrorf(f string, args ...interface{}) error {
	return r.parseError(fmt.Sprintf(f, args...))
}

func (r *Reader) parseError(msg string, errs ...error) error {
	if len(errs) > 1 {
		panic("too many errors passed to parseError")
	}

	if len(errs) > 0 && errs[0] != nil {
		return fmt.Errorf("%s: %s:%d: %s: %w", msg, r.filename, r.r.CurrentLine(), r.tok, errs[0])
	}

	return fmt.Errorf("%s: %s:%d: %s", msg, r.filename, r.r.CurrentLine(), r.tok)
}

var endGroupFromStart = map[string]string{
	"(": ")",
	"[": "]",
	"{": "}",
}

func endTokenFromStart(start tokenizer.Token) tokenizer.Token {
	return tokenizer.Token{
		Type:  tokenizer.TokenTypeEndGroup,
		Value: endGroupFromStart[start.Value],
	}
}

func (r *Reader) parseOutput(ctx context.Context, out *xslt.Output) error {
	m, err := r.parseMap(ctx)
	if err != nil {
		return err
	}

	for k, v := range m {
		switch k {
		case "method":
			out.Method = v
		case "version":
			out.Version = v
		case "encoding":
			out.Encoding = v
		case "media-type":
			out.MediaType = v

		case "doctype-public":
			out.DoctypePublic = v
		case "doctype-system":
			out.DoctypeSystem = v

		case "cdata-section-elements":
			var qnames xslt.QNames

			elems := strings.Split(v, " ")
			for _, elem := range elems {
				if elem == "" {
					continue
				}

				qnames = append(qnames, elem)
			}

			out.CDATASectionElements = qnames

		case "omit-xml-declaration":
			b, err := strconv.ParseBool(v)
			if err != nil {
				return r.parseError("bad boolean value", err)
			}
			out.OmitXMLDeclaration = xslt.Bool(b)

		case "standalone":
			b, err := strconv.ParseBool(v)
			if err != nil {
				return r.parseError("bad boolean value", err)
			}
			out.Standalone = xslt.Bool(b)

		case "indent":
			b, err := strconv.ParseBool(v)
			if err != nil {
				return r.parseError("bad boolean value", err)
			}
			out.Indent = xslt.BoolVal(b)

		default:
			return r.parseErrorf("unknown output attribute: %q", k)
		}
	}

	return nil
}

func (r *Reader) parseCall(ctx context.Context) (interface{}, error) {
	name, err := r.read(ctx)
	if err != nil {
		return nil, err
	}

	if name.Type != tokenizer.TokenTypeIdentifier {
		return nil, r.parseError("expected idenitifer")
	}
	r.consume()

	tok, err := r.peak(ctx)
	if err != nil {
		return nil, err
	}

	if tok.Type != tokenizer.TokenTypeBeginGroup || tok.Value != "(" {
		return &xslt.CallTemplate{
			Name: name.Value,
		}, nil
	}

	args, err := r.parseArgumentList(ctx)
	if err != nil {
		return nil, err
	}

	return &xslt.CallTemplate{
		Name:       name.Value,
		WithParams: args,
	}, nil
}

func (r *Reader) parseTemplate(ctx context.Context) (interface{}, error) {
	match, err := r.read(ctx)
	if err != nil {
		return nil, err
	}

	if match.Type != tokenizer.TokenTypeXPath {
		return nil, r.parseError("expected xpath")
	}
	r.consume()

	tok, err := r.peak(ctx)
	if err != nil {
		return nil, err
	}

	var params []*xslt.Param

	if tok.Type == tokenizer.TokenTypeBeginGroup && tok.Value == "(" {
		var err error
		params, err = r.parseParamList(ctx)
		if err != nil {
			return nil, err
		}
	}

	body, err := r.parseExpression(ctx)
	if err != nil {
		return nil, err
	}

	return &xslt.Template{
		Match:  match.Value,
		Params: params,
		Body:   body,
	}, nil
}

func (r *Reader) parseSubfunction(ctx context.Context) (interface{}, error) {
	name, err := r.read(ctx)
	if err != nil {
		return nil, err
	}

	if name.Type != tokenizer.TokenTypeIdentifier {
		return nil, r.parseError("expected idenitifer")
	}
	r.consume()

	tok, err := r.peak(ctx)
	if err != nil {
		return nil, err
	}

	var params []*xslt.Param

	if tok.Type == tokenizer.TokenTypeBeginGroup && tok.Value == "(" {
		var err error
		params, err = r.parseParamList(ctx)
		if err != nil {
			return nil, err
		}
	}

	body, err := r.parseExpression(ctx)
	if err != nil {
		return nil, err
	}

	return &xslt.Template{
		Name:   name.Value,
		Params: params,
		Body:   body,
	}, nil
}

func (r *Reader) parseStatement(ctx context.Context, xsl *xslt.Stylesheet) error {
	tok, err := r.peak(ctx)
	if tok.Type == tokenizer.TokenTypeComma {
		tok, err = r.read(ctx)
	}

	if err != nil {
		return err
	}

	switch tok.Type {
	case tokenizer.TokenTypeOperator:
		switch tok.Value {
		case ";":
			r.consume()
			return nil
		}

	case tokenizer.TokenTypeIdentifier:
		switch tok.Value {
		case "output":
			r.consume()
			return r.parseOutput(ctx, xsl.Output)

		case "sub":
			sub, err := r.parseSubfunction(ctx)
			if err != nil {
				return err
			}

			xsl.Body = append(xsl.Body, sub)
			return nil

		case "template":
			template, err := r.parseTemplate(ctx)
			if err != nil {
				return err
			}

			xsl.Body = append(xsl.Body, template)
			return nil

		case "param":
			r.consume()

			param, err := r.parseParam(ctx)
			if err != nil {
				return err
			}

			xsl.Body = append(xsl.Body, param)
			return nil

		case "var":
			r.consume()

			param, err := r.parseVariable(ctx)
			if err != nil {
				return err
			}

			xsl.Body = append(xsl.Body, param)
			return nil

		}
	}

	return r.parseError("unexpected top-level token")
}

func (r *Reader) parseExpression(ctx context.Context) (interface{}, error) {
	tok, err := r.peak(ctx)
	if tok.Type == tokenizer.TokenTypeComma {
		tok, err = r.read(ctx)
	}

	if err != nil {
		return nil, err
	}

	switch tok.Type {
	case tokenizer.TokenTypeEOF:
		return nil, r.parseError("unexpected EOF")

	case tokenizer.TokenTypeOperator:
		switch tok.Value {
		case ";":
			r.consume()
			return nil, nil
		}

	case tokenizer.TokenTypeError:
		r.consume()

		return nil, r.parseError("unknown token")

	case tokenizer.TokenTypeBeginGroup:
		r.consume()

		return r.parseGroup(ctx, endTokenFromStart(tok))

	case tokenizer.TokenTypeDoubleQuote, tokenizer.TokenTypeSingleQuote:
		r.consume()
		return &xslt.Text{
			Body: tok.Value,
		}, nil

	case tokenizer.TokenTypeXPath, tokenizer.TokenTypeNumber:
		r.consume()
		return &xslt.ValueOf{
			Select: tok.Value,
		}, nil

	case tokenizer.TokenTypeIdentifier:
		switch tok.Value {
		case "text":
			return r.parseText(ctx)
		case "copy-of":
			return r.parseCopyOf(ctx)

		case "var":
			r.consume()
			return r.parseVariable(ctx)

		case "foreach":
			return r.parseForEach(ctx)
		case "apply-templates":
			return r.parseApplyTemplates(ctx)

		case "when":
			return r.parseChoose(ctx)
		case "if":
			return r.parseIf(ctx)
		case "call":
			return r.parseCall(ctx)

		case "tag":
			return r.parseTag(ctx)
		case "attribs":
			return r.parseAttribs(ctx)

		case "span":
			return r.parseSpan(ctx)
		case "div":
			return r.parseDiv(ctx)
		}
	}

	return nil, r.parseError("unexpected token") /*

		r.consume()
		return nil, nil //*/
}

func (r *Reader) parseIf(ctx context.Context) (*xslt.If, error) {
	cond, err := r.read(ctx)
	if err != nil {
		return nil, err
	}

	if cond.Type != tokenizer.TokenTypeXPath {
		return nil, r.parseError("expected xpath")
	}
	r.consume()

	body, err := r.parseExpression(ctx)
	if err != nil {
		return nil, err
	}

	return &xslt.If{
		Test: cond.Value,
		Body: body,
	}, nil
}

func (r *Reader) parseChoose(ctx context.Context) (*xslt.Choose, error) {
	var whens []*xslt.When

	for {
		tok, err := r.peak(ctx)
		if err != nil {
			return nil, err
		}

		if tok.Type != tokenizer.TokenTypeIdentifier {
			return &xslt.Choose{
				Whens: whens,
			}, nil
		}

		var cond tokenizer.Token

		switch tok.Value {
		case "when":
			var err error
			cond, err = r.read(ctx)
			if err != nil {
				return nil, err
			}

			if cond.Type != tokenizer.TokenTypeXPath {
				return nil, r.parseError("expected xpath")
			}
			r.consume()

		case "otherwise":
			r.consume()

		default:
			return &xslt.Choose{
				Whens: whens,
			}, nil
		}

		body, err := r.parseExpression(ctx)
		if err != nil {
			return nil, err
		}

		if tok.Value == "otherwise" {
			return &xslt.Choose{
				Whens: whens,
				Otherwise: &xslt.Otherwise{
					Body: body,
				},
			}, nil
		}

		whens = append(whens, &xslt.When{
			Test: cond.Value,
			Body: body,
		})
	}
}

func (r *Reader) parseMap(ctx context.Context) (map[string]string, error) {
	tok, err := r.peak(ctx)
	if err != nil {
		return nil, err
	}

	var end tokenizer.Token
	switch tok.Type {
	case tokenizer.TokenTypeBeginGroup:
		end = endTokenFromStart(tok)

	case tokenizer.TokenTypeOperator:
		if tok.Value == ";" {
			return nil, nil
		}
		fallthrough

	default:
		return nil, r.parseError("expected group")
	}

	m := make(map[string]string)

	for {
		key, err := r.read(ctx)
		if key.Type == tokenizer.TokenTypeComma {
			key, err = r.read(ctx)
		}

		if err != nil {
			return nil, err
		}

		switch key.Type {
		case tokenizer.TokenTypeEndGroup:
			if key != end {
				return nil, r.parseErrorf("unexpected end map token, was expecting: %s", end)
			}

			r.consume()
			return m, nil

		case tokenizer.TokenTypeIdentifier:
		case tokenizer.TokenTypeDoubleQuote, tokenizer.TokenTypeSingleQuote:

		default:
			return nil, r.parseError("expected ident or string")
		}

		val, err := r.read(ctx)
		if val.Type == tokenizer.TokenTypeComma {
			val, err = r.read(ctx)
		}

		if err != nil {
			return nil, err
		}

		switch val.Type {
		case tokenizer.TokenTypeIdentifier:
		case tokenizer.TokenTypeDoubleQuote, tokenizer.TokenTypeSingleQuote:
		default:
			return nil, r.parseError("expected ident or string")
		}

		m[key.Value] = val.Value
	}
}

func (r *Reader) parseGroup(ctx context.Context, end tokenizer.Token) (xslt.Group, error) {
	var group xslt.Group

	for {
		tok, err := r.peak(ctx)
		if err != nil {
			if err == io.EOF {
				return nil, r.parseError("unexpected EOF")
			}

			return nil, err
		}

		if tok.Type == tokenizer.TokenTypeComma {
			tok, err = r.read(ctx)
		}

		if tok.Type == tokenizer.TokenTypeEndGroup {
			if tok != end {
				return nil, r.parseErrorf("unexpected end group token, was expecting: %s", end)
			}

			r.consume()
			return group, nil
		}

		thing, err := r.parseExpression(ctx)
		if err != nil {
			return nil, err
		}

		if thing != nil {
			group = append(group, thing)
		}
	}
}

func (r *Reader) parseText(ctx context.Context) (*xslt.Text, error) {
	val, err := r.read(ctx)
	if err != nil {
		return nil, err
	}

	switch val.Type {
	case tokenizer.TokenTypeDoubleQuote:
	case tokenizer.TokenTypeSingleQuote:

	default:
		return nil, r.parseError("expected a string")
	}

	r.consume()

	return &xslt.Text{
		Body: val.Value,
	}, nil
}

func (r *Reader) parseCopyOf(ctx context.Context) (*xslt.CopyOf, error) {
	val, err := r.read(ctx)
	if err != nil {
		return nil, err
	}

	switch val.Type {
	case tokenizer.TokenTypeXPath:
	default:
		return nil, r.parseError("expected an xpath")
	}

	r.consume()

	return &xslt.CopyOf{
		Select: val.Value,
	}, nil
}

func ParseFile(ctx context.Context, in io.Reader, filename string, xsl *xslt.Stylesheet) error {
	r := &Reader{
		filename: filename,

		r: &tokenizer.Reader{
			S: bufio.NewScanner(in),
		},
	}

	for {
		tok, err := r.peak(ctx)
		if tok.Type == tokenizer.TokenTypeComma {
			tok, err = r.read(ctx)
		}

		if tok == tokenizer.EOF {
			return nil
		}

		if err != nil {
			if err == io.EOF {
				return r.parseError("unexpected EOF")
			}

			return err
		}

		if err := r.parseStatement(ctx, xsl); err != nil {
			return err
		}
	}
}
