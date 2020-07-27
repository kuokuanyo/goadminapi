package components

import (
	"goadminapi/template/types"
	"html/template"
)

type ColAttribute struct {
	Name    string
	Content template.HTML
	Size    string
	types.Attribute
}

// -----------------為types.ColAttribute(interface)的所有方法--------------------
func (compo *ColAttribute) SetContent(value template.HTML) types.ColAttribute {
	compo.Content = value
	return compo
}

func (compo *ColAttribute) AddContent(value template.HTML) types.ColAttribute {
	compo.Content += value
	return compo
}

func (compo *ColAttribute) SetSize(value types.S) types.ColAttribute {
	compo.Size = ""
	for key, size := range value {
		compo.Size += "col-" + key + "-" + size + " "
	}
	return compo
}

// 將符合ColAttribute.TemplateList["box"](map[string]string)的值加入text(string)，接著將參數及功能添加給新的模板並解析模板
func (compo *ColAttribute) GetContent() template.HTML {
	// 首先將符合ColAttribute.TemplateList["col"](map[string]string)的值加入text(string)，接著將參數及功能添加給新的模板並解析模板
	// 將參數compo寫入buffer(bytes.Buffer)中最後輸出HTML
	return ComposeHtml(compo.TemplateList, *compo, "col")
}
