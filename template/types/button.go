package types

import (
	"goadminapi/context"
	"html/template"
)

type Buttons []Button
type Button interface {
	Content() (template.HTML, template.JS)
	GetAction() Action
	URL() string
	METHOD() string
	ID() string
	GetName() string
	SetName(name string)
}
type Action interface {
	Js() template.JS
	BtnAttribute() template.HTML
	BtnClass() template.HTML
	ExtContent() template.HTML
	FooterContent() template.HTML
	SetBtnId(btnId string)
	SetBtnData(data interface{})
	GetCallbacks() context.Node
}
