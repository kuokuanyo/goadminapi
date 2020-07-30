package guard

import (
	"fmt"
	"goadminapi/context"
	"goadminapi/plugins/admin/modules/form"
	"goadminapi/plugins/admin/modules/parameter"
	"goadminapi/plugins/admin/modules/table"
	"html/template"
	"mime/multipart"
	"strings"
)

type NewFormParam struct {
	Panel        table.Table
	Id           string
	Prefix       string
	Param        parameter.Parameters
	Path         string
	MultiForm    *multipart.Form
	PreviousPath string
	FromList     bool
	IsIframe     bool
	IframeID     string
	Alert        template.HTML
}

// return multipart/form-data設定的參數
func (e NewFormParam) Value() form.Values {
	return e.MultiForm.Value
}

func (g *Guard) NewForm(ctx *context.Context) {
	fmt.Println("0000000000000000000000000")
	previous := ctx.FormValue("__previous_")
	fmt.Println("111111111111111111")
	fmt.Println(previous)
	// 取得table(interface)、prefix
	panel, prefix := g.table(ctx)
	fmt.Println("222222222222222222")
	// 取得匹配的service.Service然後轉換成Connection(interface)
	// conn := db.GetConnection(g.services)

	// token := ctx.FormValue("__token")
	// if !auth.GetTokenService(g.services.Get("token_csrf_helper")).CheckToken(token) {
	// 	alert(ctx, panel, "wrong token", conn, g.navBtns)
	// 	ctx.Abort()
	// 	return
	// }

	// GetParamFromURL將頁面size、資料排列方式、選擇欄位...等資訊後設置至Parameters(struct)
	param := parameter.GetParamFromURL(previous, panel.GetInfo().DefaultPageSize,
		// GetPrimaryKey回傳BaseTable.PrimaryKey
		panel.GetInfo().GetSort(), panel.GetPrimaryKey().Name)
		fmt.Println("33333333333333333333")
	// 判斷參數是否是info url(true)
	fromList := isInfoUrl(previous)

	// 取得在multipart/form-data所設定的參數(map[string][]string)
	values := ctx.Request.MultipartForm.Value

	ctx.SetUserValue("new_form_param", &NewFormParam{
		Panel:  panel,
		Id:     "",
		Prefix: prefix,
		Param:  param, // 頁面size、資料排列方式、選擇欄位...等資訊
		// Get透過參數key判斷Values[key]長度是否大於0，如果大於零回傳Values[key][0]，反之回傳""
		IsIframe:     form.Values(values).Get("__iframe") == "true", // ex:false
		IframeID:     form.Values(values).Get("__iframe_id"),        // ex:空
		Path:         strings.Split(previous, "?")[0],               // ex:/admin/info/manager(roles or permissions)
		MultiForm:    ctx.Request.MultipartForm,                     // 在multipart/form-data所設定的參數
		PreviousPath: previous,                                      // ex: /admin/info/manager?__page=1&__pageSize=10&__sort=id&__sort_type=desc
		FromList:     fromList,
	})
	ctx.Next()
	fmt.Println("4444444444444444")
}

func GetNewFormParam(ctx *context.Context) *NewFormParam {
	return ctx.UserValue["new_form_param"].(*NewFormParam)
}
