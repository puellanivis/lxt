package xslt

import (
	"encoding/xml"
	"errors"
)

type ApplyTemplates struct {
	Select string `xml:"select,attr,omitempty"`

	Sort       interface{}
	WithParams []*WithParam
}

func (a *ApplyTemplates) MarshalXML(e *xml.Encoder, _ xml.StartElement) error {
	start := xmlStartElement("xsl:apply-templates",
		xmlAttr("select", a.Select),
	)

	if err := e.EncodeToken(start); err != nil {
		return err
	}

	if a.Sort != nil {
		if err := e.Encode(a.Sort); err != nil {
			return err
		}
	}

	for _, param := range a.WithParams {
		if err := e.Encode(param); err != nil {
			return err
		}
	}

	return e.EncodeToken(start.End())
}

type ForEach struct {
	Select string `xml:"select,attr"`

	Sort interface{}
	Body interface{}
}

func (f *ForEach) MarshalXML(e *xml.Encoder, _ xml.StartElement) error {
	if f.Select == "" {
		return errors.New("xsl:for-each must have a select")
	}

	start := xmlStartElement("xsl:for-each",
		xmlAttr("select", f.Select),
	)

	if err := e.EncodeToken(start); err != nil {
		return err
	}

	if f.Sort != nil {
		if err := e.Encode(f.Sort); err != nil {
			return err
		}
	}

	if err := e.Encode(f.Body); err != nil {
		return err
	}

	return e.EncodeToken(start.End())
}
