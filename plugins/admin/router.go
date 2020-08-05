package admin

import (
	"goadminapi/context"
	"goadminapi/modules/auth"
	"goadminapi/modules/config"
	"goadminapi/template"
)

func (admin *Admin) initRouter() *Admin {
	app := context.NewApp()
	route := app.Group("/admin")

	// GetComponentAsset檢查compMap(map[string]Component)的物件後將前端文件路徑加入[]string中
	for _, path := range template.GetComponentAsset() {
		route.GET("/assets"+path, admin.handler.Assets)
	}

	route.GET(config.GetLoginUrl(), admin.handler.ShowLogin)
	route.POST("/signin", admin.handler.Auth)

	authRoute := route.Group("/", auth.Middleware(admin.Conn))
	// *************還有前端函式還沒處理getMenuInfoPanel、showNewMenu********************
	authRoute.POST("/menu/new", admin.guardian.MenuNew, admin.handler.NewMenu).Name("menu_new")
	authRoute.POST("/menu/delete", admin.guardian.MenuDelete, admin.handler.DeleteMenu).Name("menu_delete")
	// *************還有前端函式還沒處理showEditMenu、table、getMenuInfoPanel********************
	authRoute.POST("/menu/edit", admin.guardian.MenuEdit, admin.handler.EditMenu).Name("menu_edit")

	authPrefixRoute := route.Group("/", auth.Middleware(admin.Conn), admin.guardian.CheckPrefix)
	authPrefixRoute.GET("/info/:__prefix", admin.handler.ShowInfo).Name("info")
	// *************還有前端函式還沒處理showNewForm、showTable********************
	authPrefixRoute.POST("/new/:__prefix", admin.guardian.NewForm, admin.handler.NewForm).Name("new")
	authPrefixRoute.POST("/delete/:__prefix", admin.guardian.Delete, admin.handler.Delete).Name("delete")
	// *************還有前端函式還沒處理showForm、showNewForm、showTable********************
	authPrefixRoute.POST("/edit/:__prefix", admin.guardian.EditForm, admin.handler.EditForm).Name("edit")


	admin.App = app
	return admin
}
