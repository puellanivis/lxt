package xslt

import (
	"encoding/xml"
	"errors"
)

func xmlName(name string) xml.Name {
	return xml.Name{
		Local: name,
	}
}

type Stylesheet struct {
	XMLName xml.Name   `name:"xsl:stylesheet"`
	Attr    []xml.Attr `xml:",attr"`

	Start    Group
	Imports  Group
	Includes Group

	Output *Output

	Body Group
}

func NewStylesheet() *Stylesheet {
	return &Stylesheet{
		XMLName: xmlName("xsl:stylesheet"),
		Attr: []xml.Attr{
			{Name: xmlName("version"), Value: "1.0"},
			{Name: xmlName("xmlns:install"), Value: "http://www.microsoft.com/support"},
			{Name: xmlName("xmlns:msxsl"), Value: "urn:schemas-microsoft-com:xslt"},
			{Name: xmlName("xmlns:xs"), Value: "http://www.w3.org/2001/XMLSchema"},
			{Name: xmlName("xmlns:xsl"), Value: "http://www.w3.org/1999/XSL/Transform"},
		},

		Output: NewOutput(),
	}
}

type Output struct {
	XMLName xml.Name `name:"xsl:output"`

	Method    string `xml:"method,attr"`
	Version   string `xml:"version,attr,omitempty"`
	Encoding  string `xml:"encoding,attr"`
	MediaType string `xml:"media-type,attr"`

	OmitXMLDeclaration *BoolVal `xml:"omit-xml-declaration,attr,omitempty"`
	Standalone         *BoolVal `xml:"standalone,attr,omitempty"`
	Indent             BoolVal  `xml:"indent,attr,omitempty"`

	DoctypePublic string `xml:"doctype-public,attr,omitempty"`
	DoctypeSystem string `xml:"doctype-system,attr,omitempty"`

	CDATASectionElements QNames `xml:"cdata-section-elements,attr,omitempty"`
}

func NewOutput() *Output {
	return &Output{
		XMLName: xmlName("xsl:output"),

		Method:    "xml",
		Version:   "1.0",
		Encoding:  "UTF-8",
		Indent:    true,
		MediaType: "application/xml",
	}
}

type Group []interface{}

type Template struct {
	Name  string `xml:"name,attr,omitempty"`
	Match string `xml:"match,attr,omitempty"`

	Params []*Param

	Body interface{}
}

func (t *Template) MarshalXML(e *xml.Encoder, _ xml.StartElement) error {
	if t.Name == "" && t.Match == "" {
		return errors.New("xsl:template must have at least a name or a match")
	}

	start := xmlStartElement("xsl:template",
		xmlAttr("name", t.Name),
		xmlAttr("match", t.Match),
	)

	if err := e.EncodeToken(start); err != nil {
		return err
	}

	for _, param := range t.Params {
		if err := e.Encode(param); err != nil {
			return err
		}
	}

	if err := e.Encode(t.Body); err != nil {
		return err
	}

	return e.EncodeToken(start.End())
}

type CallTemplate struct {
	Name string `xml:"name,attr"`

	WithParams []*WithParam
}

func (c *CallTemplate) MarshalXML(e *xml.Encoder, _ xml.StartElement) error {
	if c.Name == "" {
		return errors.New("xsl:call-template must have a name")
	}

	start := xmlStartElement("xsl:call-template",
		xmlAttr("name", c.Name),
	)

	if err := e.EncodeToken(start); err != nil {
		return err
	}

	for _, param := range c.WithParams {
		if err := e.Encode(param); err != nil {
			return err
		}
	}

	return e.EncodeToken(start.End())
}

type Attribute struct {
	Name string `xml:"name,attr"`

	Value interface{}
}

func (a *Attribute) MarshalXML(e *xml.Encoder, _ xml.StartElement) error {
	if a.Name == "" {
		return errors.New("xsl:attribute must have a name")
	}

	start := xmlStartElement("xsl:attribute",
		xmlAttr("name", a.Name),
	)

	if err := e.EncodeToken(start); err != nil {
		return err
	}

	if err := e.Encode(a.Value); err != nil {
		return err
	}

	return e.EncodeToken(start.End())
}

type Element struct {
	Name string `xml:"name,attr"`

	Body interface{}
}

func (el *Element) MarshalXML(e *xml.Encoder, _ xml.StartElement) error {
	if el.Name == "" {
		return errors.New("xsl:element must have a name")
	}

	start := xmlStartElement("xsl:element",
		xmlAttr("name", el.Name),
	)

	if err := e.EncodeToken(start); err != nil {
		return err
	}

	if err := e.Encode(el.Body); err != nil {
		return err
	}

	return e.EncodeToken(start.End())
}
