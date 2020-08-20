package errors

import (
	"html/template"
)

var (
	Msg         string
	MsgHTML     template.HTML
	MsgWithIcon template.HTML
)

const (
	PermissionDenied     = "permission denied"
	WrongID              = "wrong id"
	OperationNotAllow    = "operation not allow"
	EditFailWrongToken   = "edit fail, wrong token"
	CreateFailWrongToken = "create fail, wrong token"
	NoPermission         = "no permission"
	SiteOff              = "site is off"
)

// WrongPK return "wrong " + 參數
func WrongPK(pk string) string {
	return "wrong " + pk
}

// Init 初始化
func Init() {
	Msg = "error"
	MsgHTML = template.HTML("error")
	//MsgWithIcon = icon.Icon(icon.Warning, 2) + MsgHTML + `!`
}
