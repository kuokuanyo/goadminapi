package components

import (
	"goadminapi/modules/menu"
	"goadminapi/template/types"
	"html/template"
)

type TreeAttribute struct {
	Name      string
	Tree      []menu.Item
	EditUrl   string
	DeleteUrl string
	UrlPrefix string
	OrderUrl  string
	types.Attribute
}

// ----------------------types.TreeAttribute所有方法-------------------------
func (compo *TreeAttribute) SetTree(value []menu.Item) types.TreeAttribute {
	compo.Tree = value
	return compo
}

func (compo *TreeAttribute) SetEditUrl(value string) types.TreeAttribute {
	compo.EditUrl = value
	return compo
}


func (compo *TreeAttribute) SetUrlPrefix(value string) types.TreeAttribute {
	compo.UrlPrefix = value
	return compo
}


func (compo *TreeAttribute) SetDeleteUrl(value string) types.TreeAttribute {
	compo.DeleteUrl = value
	return compo
}

func (compo *TreeAttribute) SetOrderUrl(value string) types.TreeAttribute {
	compo.OrderUrl = value
	return compo
}

// 首先將符合TemplateList["components/tree"]的值加入text(string)，接著加入方法並解析模板
func (compo *TreeAttribute) GetContent() template.HTML {
	return ComposeHtml(compo.TemplateList, *compo, "tree")
}

// 首先將符合TemplateList["components/tree-header"]的值加入text(string)，接著加入方法並解析模板
func (compo *TreeAttribute) GetTreeHeader() template.HTML {
	return ComposeHtml(compo.TemplateList, *compo, "tree-header")
}
