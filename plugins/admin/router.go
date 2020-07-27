package admin

import (
	"goadminapi/context"
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

	// GetLoginUrl = globalCfg.LoginUrl
	route.GET(config.GetLoginUrl(), admin.handler.ShowLogin)

	route.POST("/signin", admin.handler.Auth)

	admin.App = app
	return admin
}
