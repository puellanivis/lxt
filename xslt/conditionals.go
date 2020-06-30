package xslt

import (
	"encoding/xml"
	"errors"
)

type If struct {
	Test string `xml:"test,attr"`

	Body interface{}
}

func (i *If) MarshalXML(e *xml.Encoder, _ xml.StartElement) error {
	if i.Test == "" {
		return errors.New("xsl:if cannot have empty test")
	}

	start := xmlStartElement("xsl:if",
		xmlAttr("test", i.Test),
	)

	if err := e.EncodeToken(start); err != nil {
		return err
	}

	if err := e.Encode(i.Body); err != nil {
		return err
	}

	return e.EncodeToken(start.End())
}

type Choose struct {
	Whens     []*When
	Otherwise *Otherwise `xml:",omitempty"`
}

func (c *Choose) MarshalXML(e *xml.Encoder, _ xml.StartElement) error {
	if len(c.Whens) < 1 {
		return errors.New("xsl:choose must have at least one xsl:when")
	}

	start := xmlStartElement("xsl:choose")

	if err := e.EncodeToken(start); err != nil {
		return err
	}

	for _, when := range c.Whens {
		if err := e.Encode(when); err != nil {
			return err
		}
	}

	if c.Otherwise != nil {
		if err := e.Encode(c.Otherwise); err != nil {
			return err
		}
	}

	return e.EncodeToken(start.End())

}

type When struct {
	Test string `xml:"test,attr"`

	Body interface{}
}

func (w *When) MarshalXML(e *xml.Encoder, _ xml.StartElement) error {
	if w.Test == "" {
		return errors.New("xsl:when cannot have empty test")
	}

	start := xmlStartElement("xsl:when",
		xmlAttr("test", w.Test),
	)

	if err := e.EncodeToken(start); err != nil {
		return err
	}

	if err := e.Encode(w.Body); err != nil {
		return err
	}

	return e.EncodeToken(start.End())
}

type Otherwise struct {
	Body interface{}
}

func (o *Otherwise) MarshalXML(e *xml.Encoder, _ xml.StartElement) error {
	start := xmlStartElement("xsl:otherwise")

	if err := e.EncodeToken(start); err != nil {
		return err
	}

	if err := e.Encode(o.Body); err != nil {
		return err
	}

	return e.EncodeToken(start.End())
}
