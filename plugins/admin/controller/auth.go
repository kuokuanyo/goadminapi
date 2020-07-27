package controller

import (
	"bytes"
	template2 "html/template"

	"goadminapi/context"
	"goadminapi/modules/auth"
	"goadminapi/plugins/admin/models"
	"goadminapi/plugins/admin/modules/response"
	"goadminapi/template"
	"net/http"
	"net/url"
)

func (h *Handler) Auth(ctx *context.Context) {
	var (
		user   models.UserModel
		errMsg = "fail"
		ok     bool
	)

	password := ctx.FormValue("password")
	username := ctx.FormValue("username")
	if password == "" || username == "" {
		response.BadRequest(ctx, "wrong password or username")
		return
	}
	// 檢查user密碼是否正確之後取得user的role、permission及可用menu，最後更新資料表(goadmin_users)的密碼值(加密)
	user, ok = auth.Check(password, username, h.conn)
	if !ok {
		response.BadRequest(ctx, errMsg)
		return
	}

	// 設置cookie(struct)並儲存在response header Set-Cookie中
	err := auth.SetCookie(ctx, user, h.conn)
	if err != nil {
		response.Error(ctx, err.Error())
		return
	}

	// 藉由參數Referer獲得Header
	if ref := ctx.Headers("Referer"); ref != "" {
		if u, err := url.Parse(ref); err == nil {
			v := u.Query()
			if r := v.Get("ref"); r != "" {
				rr, _ := url.QueryUnescape(r)
				response.OkWithData(ctx, map[string]interface{}{
					"url": rr,
				})
				return
			}
		}
	}

	// 成功，回傳code:200 and msg:ok and data
	response.OkWithData(ctx, map[string]interface{}{
		"url": h.config.GetIndexURL(),
	})
	return
}

// ShowLogin判斷map[string]Component(interface)是否有參數login(key)的值，接著執行template將data寫入buf並輸出HTML
func (h *Handler) ShowLogin(ctx *context.Context) {

	// GetComp判斷map[string]Component是否有參數name(login)的值，有的話則回傳Component(interface)
	// GetTemplate添加login_theme1給新的HTML模板，接著將函式加入模板並解析
	// 最後回傳模板及模板名稱
	tmpl, name := template.GetComp("login").GetTemplate()
	buf := new(bytes.Buffer)

	// ExecuteTemplate為html/template套件
	// 將第三個參數data寫入buf(struct)後輸出HTML
	if err := tmpl.ExecuteTemplate(buf, name, struct {
		UrlPrefix string
		Title     string
		Logo      template2.HTML
		CdnUrl    string
		// System    types.SystemInfo
	}{
		// AssertPrefix取得Config.prefix
		UrlPrefix: h.config.AssertPrefix(),
		Title:     h.config.LoginTitle,
		Logo:      h.config.LoginLogo,
		// System: types.SystemInfo{
		// 	Version: system.Version(),
		// },
		CdnUrl: h.config.AssetUrl,
	}); err == nil {
		// 輸出HTML，參數body保存至Context.response.Body及設置ContentType、StatusCode
		ctx.HTML(http.StatusOK, buf.String())
	} else {
		ctx.HTML(http.StatusOK, "parse template error (；′⌒`)")
		panic(err)
	}
}
