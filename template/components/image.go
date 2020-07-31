package components

import (
	"goadminapi/plugins/admin/modules"
	"goadminapi/template/types"
	"html/template"
)

type ImgAttribute struct {
	Name     string
	Width    string
	Height   string
	Uuid     string
	HasModal bool
	Src      template.URL
	types.Attribute
}

func (compo *ImgAttribute) SetWidth(value string) types.ImgAttribute {
	compo.Width = value
	return compo
}

func (compo *ImgAttribute) SetHeight(value string) types.ImgAttribute {
	compo.Height = value
	return compo
}

func (compo *ImgAttribute) WithModal() types.ImgAttribute {
	compo.HasModal = true
	compo.Uuid = modules.Uuid()
	return compo
}

func (compo *ImgAttribute) SetSrc(value template.HTML) types.ImgAttribute {
	compo.Src = template.URL(value)
	return compo
}

// 首先將符合TreeAttribute.TemplateList["components/image"](map[string]string)的值加入text(string)，接著將參數及功能添加給新的模板並解析模板
func (compo *ImgAttribute) GetContent() template.HTML {
	return ComposeHtml(compo.TemplateList, *compo, "image")
}
