package types

import "html/template"

type Thead []TheadItem
type TheadItem struct {
	Head       string       `json:"head"`
	Sortable   bool         `json:"sortable"`
	Field      string       `json:"field"`
	Hide       bool         `json:"hide"`
	Editable   bool         `json:"editable"`
	EditType   string       `json:"edit_type"`
	EditOption FieldOptions `json:"edit_option"`
	Width      string       `json:"width"`
}

// 分頁
type PaginatorAttribute interface {
	SetCurPageStartIndex(value string) PaginatorAttribute
	SetCurPageEndIndex(value string) PaginatorAttribute
	SetTotal(value string) PaginatorAttribute
	SetHideEntriesInfo() PaginatorAttribute
	SetPreviousClass(value string) PaginatorAttribute
	SetPreviousUrl(value string) PaginatorAttribute
	SetPages(value []map[string]string) PaginatorAttribute
	SetPageSizeList(value []string) PaginatorAttribute
	SetNextClass(value string) PaginatorAttribute
	SetNextUrl(value string) PaginatorAttribute
	SetOption(value map[string]template.HTML) PaginatorAttribute
	SetUrl(value string) PaginatorAttribute
	SetExtraInfo(value template.HTML) PaginatorAttribute
	GetContent() template.HTML
}
