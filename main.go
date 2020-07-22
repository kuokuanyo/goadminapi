package main

import (
	_ "github.com/go-sql-driver/mysql"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	// eng := engine.Default()
	// app := context.NewApp()
	// route := app.Group("/admin")

	// route.POST("/signin", admin.Admin.handler.Auth)

	_ = r.Run(":8080")
}

// func Auth(ctx *gin.Context) {
// 	password := ctx.FormValue("password")
// 	username := ctx.FormValue("username")
// 	if password == "" || username == "" {
// 		BadRequest(ctx, "wrong password or username")
// 		return
// 	}

// }
