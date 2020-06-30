package xslt

import (
	"encoding/xml"
)

type BoolVal bool

func Bool(b bool) *BoolVal {
	return (*BoolVal)(&b)
}

func (b *BoolVal) String() string {
	if b == nil {
		return ""
	}

	if *b {
		return "yes"
	}

	return "no"
}

func (b *BoolVal) XMLAttr(name xml.Name) xml.Attr {
	return xml.Attr{
		Name:  name,
		Value: b.String(),
	}
}

func (b *BoolVal) MarshalXMLAttr(name xml.Name) (xml.Attr, error) {
	return b.XMLAttr(name), nil
}
