package xslt

import (
	"encoding/xml"
	"sort"
)

type Attribs map[string]string

func (a Attribs) ToXMLAttrs() []xml.Attr {
	var keys []string
	for key, val := range a {
		if val == "" {
			continue
		}

		keys = append(keys, key)
	}
	sort.Strings(keys)

	var attribs []xml.Attr
	for _, key := range keys {
		attribs = append(attribs, xmlAttr(key, a[key]))
	}

	return attribs
}

func xmlAttr(name, value string) xml.Attr {
	return xml.Attr{
		Name: xml.Name{
			Local: name,
		},
		Value: value,
	}
}
