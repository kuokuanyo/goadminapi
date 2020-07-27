package components

import (
	"goadminapi/template/types"
	"html/template"
)

type TabsAttribute struct {
	Name string
	Data []map[string]template.HTML
	types.Attribute
}

// --------------------types.TabsAttribute所有方法----------------------
func (compo *TabsAttribute) SetData(value []map[string]template.HTML) types.TabsAttribute {
	compo.Data = value
	return compo
}

func (compo *TabsAttribute) GetContent() template.HTML {
	return ComposeHtml(compo.TemplateList, *compo, "tabs")
}
