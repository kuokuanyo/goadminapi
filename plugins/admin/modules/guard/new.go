package guard

import (
	"goadminapi/context"
	"goadminapi/modules/auth"
	"goadminapi/modules/config"
	"goadminapi/modules/db"
	"goadminapi/plugins/admin/modules/form"
	"goadminapi/plugins/admin/modules/parameter"
	"goadminapi/plugins/admin/modules/table"
	"html/template"
	"mime/multipart"
	"strings"
)

type ShowNewFormParam struct {
	Panel  table.Table
	Prefix string
	Param  parameter.Parameters
}

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

func (g *Guard) ShowNewForm(ctx *context.Context) {
	// panel, prefix := g.table(ctx)
	panel, prefix := g.table(ctx)

	if !panel.GetCanAdd() {
		alert(ctx, panel, "operation not allow", g.conn, g.navBtns)
		ctx.Abort()
		return
	}

	if panel.GetOnlyInfo() {
		ctx.Redirect(config.Url("/info/" + prefix))
		ctx.Abort()
		return
	}

	if panel.GetOnlyDetail() {
		ctx.Redirect(config.Url("/info/" + prefix + "/detail"))
		ctx.Abort()
		return
	}

	if panel.GetOnlyUpdateForm() {
		ctx.Redirect(config.Url("/info/" + prefix + "/edit"))
		ctx.Abort()
		return
	}

	ctx.SetUserValue("show_new_form_param", &ShowNewFormParam{
		Panel:  panel,
		Prefix: prefix,
		Param: parameter.GetParam(ctx.Request.URL, panel.GetInfo().DefaultPageSize, panel.GetInfo().SortField,
			panel.GetInfo().GetSort()),
	})
	ctx.Next()
}

func (g *Guard) NewForm(ctx *context.Context) {
	previous := ctx.FormValue("__previous_")

	// 取得table(interface)、prefix
	panel, prefix := g.table(ctx)

	// 取得匹配的service.Service然後轉換成Connection(interface)
	conn := db.GetConnection(g.services)

	token := ctx.FormValue("__token_")
	if !auth.GetTokenService(g.services.Get("token_csrf_helper")).CheckToken(token) {
		alert(ctx, panel, "wrong token", conn, g.navBtns)
		ctx.Abort()
		return
	}

	// GetParamFromURL將頁面size、資料排列方式、選擇欄位...等資訊後設置至Parameters(struct)
	param := parameter.GetParamFromURL(previous, panel.GetInfo().DefaultPageSize,
		// GetPrimaryKey回傳BaseTable.PrimaryKey
		panel.GetInfo().GetSort(), panel.GetPrimaryKey().Name)

	// 判斷參數是否是info url，如果選擇繼續增加則會是flase
	fromList := isInfoUrl(previous)
	if fromList {
		// GetRouteParamStr取得url.Values後加入__page(鍵)與值，最後編碼並回傳
		previous = config.Url("/info/" + prefix + param.GetRouteParamStr())
	}

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
		FromList:     fromList,                                      // 如果沒有繼續增加則為true，繼續增加則為false
	})
	ctx.Next()
}

func GetNewFormParam(ctx *context.Context) *NewFormParam {
	return ctx.UserValue["new_form_param"].(*NewFormParam)
}

func GetShowNewFormParam(ctx *context.Context) *ShowNewFormParam {
	return ctx.UserValue["show_new_form_param"].(*ShowNewFormParam)
}
