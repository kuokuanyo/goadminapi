package components

import (
	"goadminapi/template/types"
	"html/template"
)

type LabelAttribute struct {
	Name    string
	Color   template.HTML
	Type    string
	Content template.HTML
	types.Attribute
}

func (compo *LabelAttribute) SetType(value string) types.LabelAttribute {
	compo.Type = value
	return compo
}

func (compo *LabelAttribute) SetColor(value template.HTML) types.LabelAttribute {
	compo.Color = value
	return compo
}

func (compo *LabelAttribute) SetContent(value template.HTML) types.LabelAttribute {
	compo.Content = value
	return compo
}

// 首先將符合TreeAttribute.TemplateList["components/label"](map[string]string)的值加入text(string)，接著將參數及功能添加給新的模板並解析模板
func (compo *LabelAttribute) GetContent() template.HTML {
	return ComposeHtml(compo.TemplateList, *compo, "label")
}