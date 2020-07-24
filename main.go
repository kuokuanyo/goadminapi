package main

import (
	"goadminapi/engine"
	"goadminapi/modules/config"
	"net/http"

	_ "goadminapi/adapter/gin"

	_ "github.com/go-sql-driver/mysql"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	r.Static("/admin/assets", "./assets")


	r.LoadHTMLGlob("template/**/*.html")
	//r.LoadHTMLFiles("template/login/index.html")
	r.GET("/admin/login", func(c *gin.Context) {
		c.HTML(http.StatusOK, "login.html", gin.H{
			"UrlPrefix": "/admin",
			"CdnUrl": "",
		})
	})
	// 回傳預設的Engine(struct)
	eng := engine.Default()

	cfg := config.Config{
		// 数据库配置，为一个map，key为连接名，value为对应连接信息
		Databases: config.DatabaseList{
			// 默认数据库连接，名字必须为default
			"default": {
				Host:       "127.0.0.1",
				Port:       "3306",
				User:       "root",
				Pwd:        "asdf4440",
				Name:       "kuo",
				MaxIdleCon: 50,
				MaxOpenCon: 150,
				Driver:     "mysql",
			},
		},
		UrlPrefix: "admin",
		Store: config.Store{
			Path:   "./uploads",
			Prefix: "uploads",
		},
	}
	// 首先將參數cfg(struct)數值處理後設置至globalCfg，接著設置至Engine.config
	// 再來將driver加入Engine.Services，初始化所有資料庫連線並啟動引擎
	_ = eng.AddConfig(cfg).Use(r)

	_ = r.Run(":8080")
}
