package html

import (
	"fmt"
	"html/template"
	"strings"
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

// BaseEl return default Element(struct)
func BaseEl() Element {
	return Element{Style: make(map[string]string), Attribute: make(map[string]string)}
}

// SetTag set tag
func (b Element) SetTag(tag template.HTML) Element {
	b.Tag = tag
	return b
}

// SetContent set content
func (b Element) SetContent(content template.HTML) Element {
	b.Content = content
	return b
}

// SetAttr set attribute
func (b Element) SetAttr(key, value string) Element {
	b.Attribute[key] = value
	return b
}

// Get return HTML
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

// String return HTML
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

// SetStyleAndAttr return Element(struct)
func (b Element) SetStyleAndAttr(ms []M) Element {
	if len(ms) > 0 {
		for k, v := range ms[0] {
			b.Style[k] = v
		}
	}
	if len(ms) > 1 {
		for k, v := range ms[1] {
			b.Attribute[k] = v
		}
	}
	return b
}

// AEl set tag
func AEl() Element {
	return BaseEl().SetTag("a")
}

// LiEl set tag
func LiEl() Element {
	return BaseEl().SetTag("li")
}

// SetClass set class
func (b Element) SetClass(class ...string) Element {
	if b.Attribute["class"] != "" {
		b.Attribute["class"] += " " + strings.Join(class, " ")
	} else {
		b.Attribute["class"] += strings.Join(class, " ")
	}
	return b
}

// DivEl set tag 
func DivEl() Element {
	return BaseEl().SetTag("div")
}

// Div return html
func Div(content template.HTML, ms ...M) template.HTML {
	return DivEl().SetContent(content).SetStyleAndAttr(ms).Get()
}

// UlEl set tag 
func UlEl() Element {
	return BaseEl().SetTag("ul")
}

// Ul return html
func Ul(content template.HTML, ms ...M) template.HTML {
	return UlEl().SetContent(content).SetStyleAndAttr(ms).Get()
}

// A return html
func A(content template.HTML, ms ...M) template.HTML {
	return AEl().SetContent(content).SetStyleAndAttr(ms).Get()
}

// IEl set tag
func IEl() Element {
	return BaseEl().SetTag("i")
}
