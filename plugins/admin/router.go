package admin

import (
	"goadminapi/context"
)

func (admin *Admin) initRouter() *Admin {
	app := context.NewApp()

	route := app.Group("/admin")

	// GetLoginUrl = globalCfg.LoginUrl
	// route.GET(config.GetLoginUrl(), admin.handler.ShowLogin)

	route.POST("/signin", admin.handler.Auth)

	admin.App = app
	return admin
}
