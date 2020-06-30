package xslt

import (
	"encoding/xml"
)

func omitEmptyAttrs(attrs []xml.Attr) []xml.Attr {
	a := make([]xml.Attr, 0, len(attrs))

	for _, attr := range attrs {
		if attr.Value == "" {
			continue
		}

		a = append(a, attr)
	}
	return a
}

func xmlStartElement(name string, attrs ...xml.Attr) xml.StartElement {
	return xml.StartElement{
		Name: xml.Name{
			Local: name,
		},
		Attr: omitEmptyAttrs(attrs),
	}
}
