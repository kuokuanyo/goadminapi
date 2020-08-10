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

// GetContent 取得ex:<div class="col-md-2 "></div>
func (compo *ColAttribute) GetContent() template.HTML {
	return ComposeHtml(compo.TemplateList, *compo, "col")
}
