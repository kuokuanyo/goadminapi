package components

import (
	"goadminapi/template/types"
	"html/template"
)

// 分頁
type PaginatorAttribute struct {
	Name              string
	CurPageStartIndex string // 該頁起始資料的index
	CurPageEndIndex   string // 該頁結束資料的index
	Total             string // 總資料數
	PreviousClass     string // 如沒上一頁則PreviousClass=disabled，否則為空
	PreviousUrl       string // 前一頁的url參數，例如在第二頁時回傳第一頁的url參數
	Pages             []map[string]string // 每一分頁的資訊，包刮page、active、issplit、url...
	NextClass         string // 如沒有下一頁=disables，否則為空
	NextUrl           string // 下頁的url參數
	PageSizeList      []string
	Option            map[string]template.HTML
	Url               string // 取得第一頁的url
	HideEntriesInfo   bool
	ExtraInfo         template.HTML
	types.Attribute
}

func (compo *PaginatorAttribute) SetCurPageStartIndex(value string) types.PaginatorAttribute {
	compo.CurPageStartIndex = value
	return compo
}

func (compo *PaginatorAttribute) SetCurPageEndIndex(value string) types.PaginatorAttribute {
	compo.CurPageEndIndex = value
	return compo
}

func (compo *PaginatorAttribute) SetTotal(value string) types.PaginatorAttribute {
	compo.Total = value
	return compo
}

func (compo *PaginatorAttribute) SetExtraInfo(value template.HTML) types.PaginatorAttribute {
	compo.ExtraInfo = value
	return compo
}

func (compo *PaginatorAttribute) SetHideEntriesInfo() types.PaginatorAttribute {
	compo.HideEntriesInfo = true
	return compo
}

func (compo *PaginatorAttribute) SetPreviousClass(value string) types.PaginatorAttribute {
	compo.PreviousClass = value
	return compo
}

func (compo *PaginatorAttribute) SetPreviousUrl(value string) types.PaginatorAttribute {
	compo.PreviousUrl = value
	return compo
}

func (compo *PaginatorAttribute) SetPages(value []map[string]string) types.PaginatorAttribute {
	compo.Pages = value
	return compo
}

func (compo *PaginatorAttribute) SetPageSizeList(value []string) types.PaginatorAttribute {
	compo.PageSizeList = value
	return compo
}

func (compo *PaginatorAttribute) SetNextClass(value string) types.PaginatorAttribute {
	compo.NextClass = value
	return compo
}

func (compo *PaginatorAttribute) SetNextUrl(value string) types.PaginatorAttribute {
	compo.NextUrl = value
	return compo
}

func (compo *PaginatorAttribute) SetOption(value map[string]template.HTML) types.PaginatorAttribute {
	compo.Option = value
	return compo
}

func (compo *PaginatorAttribute) SetUrl(value string) types.PaginatorAttribute {
	compo.Url = value
	return compo
}

// 首先將符合TreeAttribute.TemplateList["components/paginator"](map[string]string)的值加入text(string)，接著將參數及功能添加給新的模板並解析模板
func (compo *PaginatorAttribute) GetContent() template.HTML {
	return ComposeHtml(compo.TemplateList, *compo, "paginator")
}