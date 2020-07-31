package html

import (
	"fmt"
	"html/template"
)

type M map[string]string

type Style M
type Attribute M

type Element struct {
	Tag       template.HTML
	Content   template.HTML
	Style     Style
	Attribute Attribute
}

func BaseEl() Element {
	return Element{Style: make(map[string]string), Attribute: make(map[string]string)}
}

func (b Element) SetTag(tag template.HTML) Element {
	b.Tag = tag
	return b
}

func (b Element) SetContent(content template.HTML) Element {
	b.Content = content
	return b
}

func (b Element) SetAttr(key, value string) Element {
	b.Attribute[key] = value
	return b
}

func (b Element) Get() template.HTML {
	return template.HTML(fmt.Sprintf(`<%s%s%s>%s</%s>`, b.Tag, b.Style.String(), b.Attribute.String(), b.Content, b.Tag))
}

func (s Style) String() template.HTML {
	res := ""
	for k, v := range s {
		res += k + ":" + v + ";"
	}
	if res != "" {
		res = ` style="` + res + `"`
	}
	return template.HTML(res)
}

func (s Attribute) String() template.HTML {
	res := ""
	for k, v := range s {
		res += k + `="` + v + `" `
	}
	if res != "" {
		res = ` ` + res[:len(res)-1]
	}
	return template.HTML(res)
}

func AEl() Element {
	return BaseEl().SetTag("a")
}
