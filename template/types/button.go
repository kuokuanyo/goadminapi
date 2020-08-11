package types

import (
	"goadminapi/context"
	"goadminapi/modules/utils"
	"goadminapi/plugins/admin/models"
	"html/template"
	"net/url"
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

type ActionButton struct {
	*BaseButton
}

type BaseButton struct {
	Id, Url, Method, Name string
	Title                 template.HTML
	Action                Action
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

//************************Button(interface)的所有方法****************************

// Content return "",""
func (b *BaseButton) Content() (template.HTML, template.JS) { return "", "" }

// GetAction get BaseButton.Action
func (b *BaseButton) GetAction() Action { return b.Action }

// ID get BaseButton.ID
func (b *BaseButton) ID() string { return b.Id }

// URL get BaseButton.Url
func (b *BaseButton) URL() string { return b.Url }

// METHOD get BaseButton.Method
func (b *BaseButton) METHOD() string { return b.Method }

// GetName get BaseButton.Name
func (b *BaseButton) GetName() string { return b.Name }

// SetName set theme
func (b *BaseButton) SetName(name string) { b.Name = name }

//************************Button(interface)的所有方法****************************

func GetActionButton(title template.HTML, action Action, ids ...string) *ActionButton {

	id := ""
	if len(ids) > 0 {
		id = ids[0]
	} else {
		id = "action-info-btn-" + utils.Uuid(10)
	}

	action.SetBtnId(id)
	node := action.GetCallbacks()

	return &ActionButton{
		BaseButton: &BaseButton{
			Id:     id,
			Title:  title,
			Action: action,
			Url:    node.Path,
			Method: node.Method,
		},
	}
}

// 檢查權限，回傳Buttons([]Button(interface))
func (b Buttons) CheckPermission(user models.UserModel) Buttons {
	btns := make(Buttons, 0)
	for _, btn := range b {
		// 檢查權限(藉由url、method)
		if user.CheckPermissionByUrlMethod(btn.URL(), btn.METHOD(), url.Values{}) {
			btns = append(btns, btn)
		}
	}
	return btns
}

// 取得HTML及JSON
func (b Buttons) Content() (template.HTML, template.JS) {
	h := template.HTML("")
	j := template.JS("")

	for _, btn := range b {
		hh, jj := btn.Content()
		h += hh
		j += jj
	}
	return h, j
}
