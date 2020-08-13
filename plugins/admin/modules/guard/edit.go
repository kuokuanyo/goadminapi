package guard

import (
	"goadminapi/context"
	"goadminapi/modules/auth"
	"goadminapi/modules/config"
	"goadminapi/modules/db"
	"goadminapi/modules/errors"
	"goadminapi/plugins/admin/modules/form"
	"goadminapi/plugins/admin/modules/parameter"
	"goadminapi/plugins/admin/modules/response"
	"goadminapi/plugins/admin/modules/table"
	"goadminapi/template"
	"goadminapi/template/types"
	"mime/multipart"
	"regexp"
	"strings"

	tmpl "html/template"
)

type ShowFormParam struct {
	Panel  table.Table
	Id     string
	Prefix string
	Param  parameter.Parameters
}

type EditFormParam struct {
	Panel        table.Table
	Id           string
	Prefix       string
	Param        parameter.Parameters
	Path         string
	MultiForm    *multipart.Form
	PreviousPath string
	Alert        tmpl.HTML
	FromList     bool
	IsIframe     bool
	IframeID     string
}

// 取得EditFormParam.MultiForm.Value(map[string][]string)
func (e EditFormParam) Value() form.Values {
	return e.MultiForm.Value
}

func (g *Guard) ShowForm(ctx *context.Context) {
	// 取得table(interface)、prefix
	panel, prefix := g.table(ctx)

	if !panel.GetEditable() {
		alert(ctx, panel, errors.OperationNotAllow, g.conn, g.navBtns)
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

	if panel.GetOnlyNewForm() {
		ctx.Redirect(config.Url("/info/" + prefix + "/new"))
		ctx.Abort()
		return
	}

	id := ctx.Query("__edit_pk")
	if id == "" && prefix != "site" {
		alert(ctx, panel, errors.WrongPK(panel.GetPrimaryKey().Name), g.conn, g.navBtns)
		ctx.Abort()
		return
	}
	ctx.SetUserValue("show_form_param", &ShowFormParam{
		Panel:  panel,
		Id:     id,
		Prefix: prefix,
		Param: parameter.GetParam(ctx.Request.URL, panel.GetInfo().DefaultPageSize, panel.GetInfo().SortField,
			panel.GetInfo().GetSort()).WithPKs(id),
	})
	ctx.Next()
}

// EditForm 編輯表單(POST功能)
func (g *Guard) EditForm(ctx *context.Context) {
	previous := ctx.FormValue("__previous_")

	// 取得table(interface)、prefix
	panel, prefix := g.table(ctx)
	// 判斷是否有編輯功能(日誌沒有編輯功能)
	if !panel.GetEditable() {
		alert(ctx, panel, errors.OperationNotAllow, g.conn, g.navBtns)
		ctx.Abort()
		return
	}

	// 藉由參數取得multipart/form-data中的__token_值並判斷
	token := ctx.FormValue("__token_")
	if !auth.GetTokenService(g.services.Get("token_csrf_helper")).CheckToken(token) {
		alert(ctx, panel, errors.EditFailWrongToken, g.conn, g.navBtns)
		ctx.Abort()
		return
	}

	// GetParamFromURL將頁面size、資料排列方式、選擇欄位...等資訊後設置至Parameters(struct)
	param := parameter.GetParamFromURL(previous, panel.GetInfo().DefaultPageSize,
		panel.GetInfo().GetSort(), panel.GetPrimaryKey().Name)

	// 判斷參數是否是info url，如果選擇繼續增加則會是flase
	fromList := isInfoUrl(previous)
	if fromList {
		previous = config.Url("/info/" + prefix + param.GetRouteParamStr())
	}

	// 取得在multipart/form-data所設定的參數(struct)
	multiForm := ctx.Request.MultipartForm
	// 取得id
	id := multiForm.Value[panel.GetPrimaryKey().Name][0]
	// 取得在multipart/form-data所設定的參數(map[string][]string)
	values := ctx.Request.MultipartForm.Value

	ctx.SetUserValue("edit_form_param", &EditFormParam{
		Panel:        panel,
		Id:           id,
		Prefix:       prefix,                                        // manage or roles or permissions
		Param:        param.WithPKs(id),                             // 將參數(id)結合並設置至Parameters.Fields["__pk"]後回傳
		Path:         strings.Split(previous, "?")[0],               // ex:/admin/info/manager(roles or permissions)
		MultiForm:    multiForm,                                     // 在multipart/form-data所設定的參數
		IsIframe:     form.Values(values).Get("__iframe") == "true", // ex:false
		IframeID:     form.Values(values).Get("__iframe_id"),
		PreviousPath: previous, // ex: /admin/info/manager?__page=1&__pageSize=10&__sort=id&__sort_type=desc
		FromList:     fromList, // 如果沒有繼續增加則為true，繼續增加則為false
	})
	ctx.Next()
}

// 回傳Context.UserValue[edit_form_param]的值(struct)
func GetEditFormParam(ctx *context.Context) *EditFormParam {
	return ctx.UserValue["edit_form_param"].(*EditFormParam)
}

func GetShowFormParam(ctx *context.Context) *ShowFormParam {
	return ctx.UserValue["show_form_param"].(*ShowFormParam)
}

// 判斷參數是否是info url
func isInfoUrl(s string) bool {
	reg, _ := regexp.Compile("(.*?)info/(.*?)$")
	sub := reg.FindStringSubmatch(s)
	return len(sub) > 2 && !strings.Contains(sub[2], "/")
}

func getAlert(msg string) tmpl.HTML {
	return template.Get(config.GetTheme()).Alert().Warning(msg)
}

func alert(ctx *context.Context, panel table.Table, msg string, conn db.Connection, btns *types.Buttons) {
	if ctx.WantJSON() {
		response.BadRequest(ctx, msg)
	} else {
		response.Alert(ctx, panel.GetInfo().Description, panel.GetInfo().Title, msg, conn, btns)
	}
}

// 將給定的數據(types.Page(struct))及參數寫入buf(struct)並回傳，最後輸出HTML
func alertWithTitleAndDesc(ctx *context.Context, title, desc, msg string, conn db.Connection, btns *types.Buttons) {
	response.Alert(ctx, desc, title, msg, conn, btns)
}
