package action

import (
	"goadminapi/context"
	"html/template"
)

var operationHandlerSetter context.NodeProcessor

type BaseAction struct {
	BtnId   string
	BtnData interface{}
	JS      template.JS
}

// *****************Action(interface)的所有方法*******************

// SetBtnId set BaseAction.BtnId
func (base *BaseAction) SetBtnId(btnId string) { base.BtnId = btnId }

// Js get BaseAction.JS
func (base *BaseAction) Js() template.JS { return base.JS }

// BtnClass return ""
func (base *BaseAction) BtnClass() template.HTML { return "" }

// BtnAttribute return ""
func (base *BaseAction) BtnAttribute() template.HTML { return "" }

// GetCallbacks get context.Node{}
func (base *BaseAction) GetCallbacks() context.Node { return context.Node{} }

// ExtContent return ""
func (base *BaseAction) ExtContent() template.HTML { return template.HTML(``) }

// FooterContent return ""
func (base *BaseAction) FooterContent() template.HTML { return template.HTML(``) }

// SetBtnData set BaseAction.BtnData
func (base *BaseAction) SetBtnData(data interface{}) { base.BtnData = data }

// *****************Action(interface)的所有方法*******************

// 將參數p(func(...Node))設置給operationHandlerSetter
func InitOperationHandlerSetter(p context.NodeProcessor) {
	operationHandlerSetter = p
}
