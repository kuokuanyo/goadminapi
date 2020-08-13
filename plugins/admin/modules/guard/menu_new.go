package guard

import (
	"goadminapi/context"
	"goadminapi/modules/auth"
	"strconv"

	"html/template"
)

// MenuNewParam為menu資料表欄位
type MenuNewParam struct {
	Title    string
	Header   string
	ParentId int64
	Icon     string
	Uri      string
	Roles    []string
	Alert    template.HTML
}

func (g *Guard) MenuNew(ctx *context.Context) {
	var (
		alert template.HTML
		token = ctx.FormValue("__token_")
	)
	if !auth.GetTokenService(g.services.Get("token_csrf_helper")).CheckToken(token) {
		alert = getAlert("wrong token")
	}
	
	// title與icon值一定要設置(multipart/form-data)
	if alert == "" {
		alert = checkEmpty(ctx, "title", "icon")
	}
	// 藉由參數取得multipart/form-data中的parent_id值
	parentId := ctx.FormValue("parent_id")
	if parentId == "" {
		parentId = "0"
	}
	parentIdInt, _ := strconv.Atoi(parentId)

	// 將值設置至Context.UserValue[new_menu_param]
	ctx.SetUserValue("new_menu_param", &MenuNewParam{
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

func GetMenuNewParam(ctx *context.Context) *MenuNewParam {
	// 將Context.UserValue(map[string]interface{})[new_menu_param]的值轉換成MenuNewParam(struct)類別
	return ctx.UserValue["new_menu_param"].(*MenuNewParam)
}

func (e MenuNewParam) HasAlert() bool {
	return e.Alert != template.HTML("")
}
