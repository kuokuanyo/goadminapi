package action

import (
	"html/template"
	"strings"
)

type JumpAction struct {
	BaseAction
	Url         string
	Ext         template.HTML
	NewTabTitle string
}

func Jump(url string, ext ...template.HTML) *JumpAction {
	url = strings.Replace(url, "{%id}", "{{.Id}}", -1)
	url = strings.Replace(url, "{%ids}", "{{.Ids}}", -1)
	if len(ext) > 0 {
		return &JumpAction{Url: url, Ext: ext[0]}
	}
	return &JumpAction{Url: url, NewTabTitle: ""}
}
