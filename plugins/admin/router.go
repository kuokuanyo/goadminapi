package admin

import "goadminapi/context"

func (admin *Admin) initRouter() *Admin {
	app := context.NewApp()

	route := app.Group("/admin")

	route.POST("/signin", admin.handler.Auth)

	admin.App = app
	return admin
}
