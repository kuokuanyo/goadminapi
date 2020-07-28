package components

import (
	"goadminapi/template/types"
	"html/template"

	"goadminapi/modules/errors"
)

type AlertAttribute struct {
	Name    string
	Theme   string
	Title   template.HTML
	Content template.HTML
	types.Attribute
}

// ---------------------types.AlertAttribute所有方法-------------------
func (compo *AlertAttribute) SetTheme(value string) types.AlertAttribute {
	compo.Theme = value
	return compo
}

func (compo *AlertAttribute) SetTitle(value template.HTML) types.AlertAttribute {
	compo.Title = value
	return compo
}

func (compo *AlertAttribute) SetContent(value template.HTML) types.AlertAttribute {
	compo.Content = value
	return compo
}

// 首先將參數設置至AlertAttribute(struct)後，接著將符合TemplateList["components/alert"]的值加入text(string)，接著將參數及功能添加給新的模板並解析模板
func (compo *AlertAttribute) Warning(msg string) template.HTML {
	return compo.SetTitle(errors.MsgWithIcon).
		SetTheme("warning").
		SetContent(template.HTML(msg)).
		GetContent()
}

// 首先將符合TemplateList["components/alert"]的值加入text(string)，接著將參數及功能添加給新的模板並解析模板
func (compo *AlertAttribute) GetContent() template.HTML {
	return ComposeHtml(compo.TemplateList, *compo, "alert")
}
