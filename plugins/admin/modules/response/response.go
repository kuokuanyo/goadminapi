package response

import (
	"goadminapi/context"
	"goadminapi/modules/auth"
	"goadminapi/modules/db"
	"goadminapi/modules/menu"
	"goadminapi/template/types"
	"net/http"

	"goadminapi/modules/config"
	"goadminapi/template"
)

// 成功，回傳code:200 and msg:ok and data
func OkWithData(ctx *context.Context, data map[string]interface{}) {
	ctx.JSON(http.StatusOK, map[string]interface{}{
		"code": http.StatusOK,
		"msg":  "ok",
		"data": data,
	})
}

// 成功，回傳code:200 and msg
func OkWithMsg(ctx *context.Context, msg string) {
	ctx.JSON(http.StatusOK, map[string]interface{}{
		"code": http.StatusOK,
		"msg":  msg,
	})
}

// 錯誤請求，回傳code:400 and msg
func BadRequest(ctx *context.Context, msg string) {
	ctx.JSON(http.StatusBadRequest, map[string]interface{}{
		"code": http.StatusBadRequest,
		// Get依照設定的語言給予訊息
		"msg": msg,
	})
}

// 錯誤，回傳code:500 and msg
func Error(ctx *context.Context, msg string) {
	ctx.JSON(http.StatusInternalServerError, map[string]interface{}{
		"code": http.StatusInternalServerError,
		"msg":  msg,
	})
}

// 將給定的數據(types.Page(struct))寫入buf(struct)並回傳，最後輸出HTML
func Alert(ctx *context.Context, desc, title, msg string, conn db.Connection, btns *types.Buttons,
	pageType ...template.PageType) {
	// 透過參數ctx回傳目前登入的用戶(Context.UserValue["user"])並轉換成UserModel
	user := auth.Auth(ctx)

	pt := template.Error500Page
	if len(pageType) > 0 {
		pt = pageType[0]
	}
	// GetPageContentFromPageType從頁面類型取得頁面內容
	pageTitle, description, content := template.GetPageContentFromPageType(title, desc, msg, pt)
	// GetTemplate為Template(interface)的方法，取得模板
	tmpl, tmplName := template.Default().GetTemplate(ctx.IsPjax())
	buf := template.Execute(template.ExecuteParam{
		User:     user,
		TmplName: tmplName,
		Tmpl:     tmpl,
		Panel: types.Panel{
			Content:     content,
			Description: description,
			Title:       pageTitle,
		},
		Config:    *config.Get(),
		Menu:      menu.GetGlobalMenu(user, conn).SetActiveClass(config.URLRemovePrefix(ctx.Path())),
		Animation: true,
		Buttons:   *btns,
	})
	ctx.HTML(http.StatusOK, buf.String())
}

