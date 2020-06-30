package xslt

import (
	"encoding/xml"
)

type Text struct {
	DisableOutputEscaping *BoolVal `xml:"disable-output-escaping,attr,omitempty"`

	Body string `xml:",innerxml"`
}

func (t *Text) MarshalXML(e *xml.Encoder, _ xml.StartElement) error {
	start := xmlStartElement("xsl:text",
		xmlAttr("disable-output-escaping", t.DisableOutputEscaping.String()),
	)

	if err := e.EncodeToken(start); err != nil {
		return err
	}

	if err := e.EncodeToken(xml.CharData(t.Body)); err != nil {
		return err
	}

	return e.EncodeToken(start.End())
}

type ValueOf struct {
	DisableOutputEscaping *BoolVal `xml:"disable-output-escaping,attr,omitempty"`

	Select string `xml:"select,attr"`
}

func (t *ValueOf) MarshalXML(e *xml.Encoder, _ xml.StartElement) error {
	start := xmlStartElement("xsl:value-of",
		xmlAttr("disable-output-escaping", t.DisableOutputEscaping.String()),
		xmlAttr("select", t.Select),
	)

	if err := e.EncodeToken(start); err != nil {
		return err
	}

	return e.EncodeToken(start.End())
}

type CopyOf struct {
	Select string `xml:"select,attr"`
}

func (t *CopyOf) MarshalXML(e *xml.Encoder, _ xml.StartElement) error {
	start := xmlStartElement("xsl:copy-of",
		xmlAttr("select", t.Select),
	)

	if err := e.EncodeToken(start); err != nil {
		return err
	}

	return e.EncodeToken(start.End())
}
