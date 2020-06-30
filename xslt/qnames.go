package xslt

import (
	"encoding/xml"
	"strings"
)

type QNames []string

func (q QNames) MarshalXMLAttr(name xml.Name) (xml.Attr, error) {
	val := strings.Join([]string(q), " ")

	return xml.Attr{
		Name:  name,
		Value: val,
	}, nil
}
