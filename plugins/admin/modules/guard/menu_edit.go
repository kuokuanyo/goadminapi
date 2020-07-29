package guard

import (
	"goadminapi/context"
	"goadminapi/modules/auth"
	"strconv"

	"html/template"
)

type MenuEditParam struct {
	Id       string
	Title    string
	Header   string
	ParentId int64
	Icon     string
	Uri      string
	Roles    []string
	Alert    template.HTML
}

func (g *Guard) MenuEdit(ctx *context.Context) {
	parentId := ctx.FormValue("parent_id")
	if parentId == "" {
		parentId = "0"
	}
	var (
		parentIdInt, _ = strconv.Atoi(parentId)
		token          = ctx.FormValue("__token_")
		alert          template.HTML
	)
	if !auth.GetTokenService(g.services.Get("token_csrf_helper")).CheckToken(token) {
		alert = getAlert("wrong token")
	}
	if alert == "" {
		alert = checkEmpty(ctx, "id", "title", "icon")
	}

	// 將multipart/form-data的key、value值設置至Context.UserValue[edit_menu_param]
	ctx.SetUserValue("edit_menu_param", &MenuEditParam{
		Id:       ctx.FormValue("id"),
		Title:    ctx.FormValue("title"),
		Header:   ctx.FormValue("header"),
		ParentId: int64(parentIdInt),
		Icon:     ctx.FormValue("icon"),
		Uri:      ctx.FormValue("uri"),
		Roles:    ctx.Request.Form["roles[]"],
		Alert:    alert,
	})
	ctx.Next()
}

// 取得Context.UserValue[edit_menu_param]的值並轉換成MenuEditParam(struct)
func GetMenuEditParam(ctx *context.Context) *MenuEditParam {
	return ctx.UserValue["edit_menu_param"].(*MenuEditParam)
}

func (e MenuEditParam) HasAlert() bool {
	return e.Alert != template.HTML("")
}

// 檢查參數(多個key)有在multipart/form-data裡設置(如果值為空則出現錯誤)
func checkEmpty(ctx *context.Context, key ...string) template.HTML {
	for _, k := range key {
		if ctx.FormValue(k) == "" {
			return getAlert("wrong " + k)
		}
	}
	return template.HTML("")
}
