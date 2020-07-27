package components

import (
	"html/template"

	"goadminapi/template/types"
)

type RowAttribute struct {
	Name    string
	Content template.HTML
	types.Attribute
}

// -----------------------types.RowAttribute(interface)所有方法------------------
func (compo *RowAttribute) SetContent(value template.HTML) types.RowAttribute {
	compo.Content = value
	return compo
}

func (compo *RowAttribute) AddContent(value template.HTML) types.RowAttribute {
	compo.Content += value
	return compo
}

// 首先將符合TreeAttribute.TemplateList["components/tree-header"](map[string]string)的值加入text(string)，接著將參數及功能添加給新的模板並解析模板
func (compo *RowAttribute) GetContent() template.HTML {
	// 首先將符合TreeAttribute.TemplateList["components/row"](map[string]string)的值加入text(string)，接著將參數及功能添加給新的模板並解析模板
	return ComposeHtml(compo.TemplateList, *compo, "row")
}
