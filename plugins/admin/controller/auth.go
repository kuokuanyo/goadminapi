package controller

import (
	"bytes"
	"database/sql"
	template2 "html/template"
	"strings"

	"goadminapi/context"
	"goadminapi/modules/auth"
	"goadminapi/modules/config"
	"goadminapi/modules/db"
	"goadminapi/plugins/admin/models"
	"goadminapi/plugins/admin/modules/response"
	"goadminapi/plugins/admin/modules/table"
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
	phone := ctx.FormValue("phone")
	if password == "" || phone == "" {
		response.BadRequest(ctx, "wrong password or phone")
		return
	}

	// 檢查user密碼是否正確之後取得user的role、permission及可用menu，最後更新資料表(users)的密碼值(加密)
	user, ok = auth.Check(password, phone, h.conn)
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

// Signup 註冊POST功能
func (h *Handler) Signup(ctx *context.Context) {
	var username, picture string
	userid := ctx.FormValue("userid")
	phone := ctx.FormValue("phone")
	email := ctx.FormValue("email")
	password := ctx.FormValue("password")
	passwordCheck := ctx.FormValue("passwordCheck")

	herokuDb, err := sql.Open("mysql", "be94ad46dfd2c5:0986ac8c@tcp(us-cdbr-east-02.cleardb.com:3306)/heroku_340b0d6567ec671")
	defer herokuDb.Close()
	if err != nil {
		panic(err.Error())
	}
	if err := herokuDb.Ping(); err != nil {
		panic(err)
	}

	rows, err := herokuDb.Query("SELECT username, pictureURL FROM line WHERE userid=?", userid)
	if err != nil {
		panic(err)
	}
	for rows.Next() {
		rows.Scan(&username, &picture)
	}

	if userid == "" || username == "" {
		response.BadRequest(ctx, "userid and username can not be empty")
		return
	}

	if phone == "" || email == "" || password == "" || passwordCheck == "" {
		response.BadRequest(ctx, "Value can not be empty")
		return
	}

	if !strings.Contains(phone[:2], "09") {
		response.BadRequest(ctx, "Wrong phone number, ex:09...")
		return
	}

	if password != passwordCheck {
		response.BadRequest(ctx, "password does not match")
		return
	}

	if !strings.Contains(email, "gmail") {
		response.BadRequest(ctx, "Wrong email")
		return
	}

	userbyPhone := models.User("users").SetConn(h.conn).FindByPhone(phone)
	if !userbyPhone.IsEmpty() {
		response.BadRequest(ctx, "Phone number already used")
		return
	}

	userbyUserid := models.User("users").SetConn(h.conn).FindByUserid(userid)
	if !userbyUserid.IsEmpty() {
		response.BadRequest(ctx, "User registered")
		return
	}

	user, err := models.User("users").SetConn(h.conn).New(userid, username, phone, email, table.EncodePassword([]byte(password)), picture)
	if db.CheckError(err, db.INSERT) {
		panic(err)
	}

	_, addRoleErr := user.SetConn(h.conn).AddRole("1")
	if db.CheckError(addRoleErr, db.INSERT) {
		panic(err)
	}
	_, addPermissionErr := user.SetConn(h.conn).AddPermission("1")
	if db.CheckError(addPermissionErr, db.INSERT) {
		panic(err)
	}

	stmt, _ := herokuDb.Prepare("UPDATE line set remarks = true where userid = ?")
	stmt.Exec(userid)

	response.OkWithData(ctx, map[string]interface{}{
		"url": "/admin/login",
	})
	return
}

// ShowSignup 補齊資料前端頁面
func (h *Handler) ShowSignup(ctx *context.Context) {
	tmpl, name := template.GetComp("signup").GetTemplate()
	buf := new(bytes.Buffer)

	// 將第三個參數data寫入buf(struct)後輸出HTML
	if err := tmpl.ExecuteTemplate(buf, name, struct {
		UrlPrefix string
		Logo      template2.HTML
		CdnUrl    string
	}{
		UrlPrefix: h.config.AssertPrefix(),
		Logo:      h.config.LoginLogo,
		CdnUrl:    h.config.AssetUrl,
	}); err == nil {
		ctx.HTML(http.StatusOK, buf.String())
	} else {
		ctx.HTML(http.StatusOK, "parse template error (；′⌒`)")
		panic(err)
	}
}

// ShowLogin 判斷map[string]Component(interface)是否有參數login(key)的值，接著執行template將data寫入buf並輸出HTML
func (h *Handler) ShowLogin(ctx *context.Context) {
	// GetComp判斷map[string]Component是否有參數name(login)的值，有的話則回傳Component(interface)
	// GetTemplate添加login_theme1給新的HTML模板，接著將函式加入模板並解析
	// 最後回傳模板及模板名稱
	tmpl, name := template.GetComp("login").GetTemplate()
	buf := new(bytes.Buffer)

	// 將第三個參數data寫入buf(struct)後輸出HTML
	if err := tmpl.ExecuteTemplate(buf, name, struct {
		UrlPrefix string
		Title     string
		Logo      template2.HTML
		CdnUrl    string
		// System    types.SystemInfo
	}{
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

// Logout delete the cookie.
func (h *Handler) Logout(ctx *context.Context) {
	err := auth.DelCookie(ctx, db.GetConnection(h.services))
	if err != nil {
		panic(err)
	}

	ctx.AddHeader("Location", h.config.Url(config.GetLoginUrl()))
	ctx.SetStatusCode(302)
}
