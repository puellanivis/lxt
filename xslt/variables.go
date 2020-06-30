package xslt

import (
	"encoding/xml"
	"errors"
)

type variable = struct {
	Name   string
	Select string
	Value  interface{}
}

func marshalVariable(e *xml.Encoder, tagName string, v variable) error {
	if v.Name == "" {
		return errors.New("Variable cannot have empty name")
	}

	start := xmlStartElement(tagName,
		xmlAttr("name", v.Name),
		xmlAttr("select", v.Select),
	)

	if err := e.EncodeToken(start); err != nil {
		return err
	}

	if v.Value != nil {
		if err := e.Encode(v.Value); err != nil {
			return err
		}
	}

	return e.EncodeToken(start.End())
}

type Param struct {
	Name   string `xml:"name,attr"`
	Select string `xml:"select,attr,omitempty"`

	Value interface{} `xml:",omitempty`
}

func (p *Param) MarshalXML(e *xml.Encoder, _ xml.StartElement) error {
	return marshalVariable(e, "xsl:param", variable(*p))
}

type Variable struct {
	Name   string `xml:"name,attr"`
	Select string `xml:"select,attr,omitempty"`

	Value interface{} `xml:",omitempty`
}

func (v *Variable) MarshalXML(e *xml.Encoder, _ xml.StartElement) error {
	return marshalVariable(e, "xsl:variable", variable(*v))
}

type WithParam struct {
	Name   string `xml:"name,attr"`
	Select string `xml:"select,attr,omitempty"`

	Value interface{} `xml:",omitempty`
}

func (p *WithParam) MarshalXML(e *xml.Encoder, _ xml.StartElement) error {
	return marshalVariable(e, "xsl:with-param", variable(*p))
}
