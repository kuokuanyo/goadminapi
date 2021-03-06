package admin

import (
	"goadminapi/context"
	"goadminapi/modules/auth"
	"goadminapi/modules/config"
	"goadminapi/modules/utils"
	"goadminapi/template"
)

func (admin *Admin) initRouter() *Admin {
	app := context.NewApp()
	route := app.Group("/admin")

	// GetComponentAsset檢查compMap(map[string]Component)的物件後將前端文件路徑加入[]string中
	for _, path := range template.GetComponentAsset() {
		route.GET("/assets"+path, admin.handler.Assets)
	}

	checkRepeatedPath := make([]string, 0)
	for _, themeName := range template.Themes() {
		for _, path := range template.Get(themeName).GetAssetList() {
			if !utils.InArray(checkRepeatedPath, path) {
				checkRepeatedPath = append(checkRepeatedPath, path)
				route.GET("/assets"+path, admin.handler.Assets)
			}
		}
	}

	route.GET("/", admin.handler.Transaction)

	route.GET("/signup", admin.handler.ShowSignup)
	route.GET(config.GetLoginUrl(), admin.handler.ShowLogin)
	route.POST("/signup", admin.handler.Signup)
	route.POST("/signin", admin.handler.Auth)

	// 退費
	route.GET("/refund", admin.handler.ShowRefund)
	route.POST("/refund", admin.handler.Refund)

	authRoute := route.Group("/", auth.Middleware(admin.Conn))
	authRoute.GET("/logout", admin.handler.Logout)

	authRoute.GET("/menu", admin.handler.ShowMenu).Name("menu")
	authRoute.GET("/menu/new", admin.handler.ShowNewMenu).Name("menu_new_show")
	authRoute.GET("/menu/edit/show", admin.handler.ShowEditMenu).Name("menu_edit_show")
	authRoute.POST("/menu/new", admin.guardian.MenuNew, admin.handler.NewMenu).Name("menu_new")
	authRoute.POST("/menu/delete", admin.guardian.MenuDelete, admin.handler.DeleteMenu).Name("menu_delete")
	authRoute.POST("/menu/edit", admin.guardian.MenuEdit, admin.handler.EditMenu).Name("menu_edit")

	authPrefixRoute := route.Group("/", auth.Middleware(admin.Conn), admin.guardian.CheckPrefix)
	authPrefixRoute.GET("/info/:__prefix", admin.handler.ShowInfo).Name("info")
	authPrefixRoute.GET("/info/:__prefix/edit", admin.guardian.ShowForm, admin.handler.ShowForm).Name("show_edit")
	authPrefixRoute.GET("/info/:__prefix/new", admin.guardian.ShowNewForm, admin.handler.ShowNewForm).Name("show_new")
	authPrefixRoute.GET("/info/:__prefix/detail", admin.handler.ShowDetail).Name("detail")
	authPrefixRoute.POST("/new/:__prefix", admin.guardian.NewForm, admin.handler.NewForm).Name("new")
	authPrefixRoute.POST("/delete/:__prefix", admin.guardian.Delete, admin.handler.Delete).Name("delete")
	authPrefixRoute.POST("/edit/:__prefix", admin.guardian.EditForm, admin.handler.EditForm).Name("edit")

	// 功能
	authRoute.POST("/remove-bg", admin.handler.RemoveBg)
	// authRoute.POST("/convert-file", admin.handler.ConvertFile)
	// authRoute.POST("/evaluate", admin.handler.Evaluate)

	admin.App = app
	return admin
}
