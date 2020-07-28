package main

import (
	"goadminapi/engine"
	"goadminapi/modules/config"

	_ "goadminapi/adapter/gin"

	_ "goadminapi/themes/adminlte"

	_ "github.com/go-sql-driver/mysql"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	// 回傳預設的Engine(struct)
	eng := engine.Default()

	cfg := config.Config{
		Databases: config.DatabaseList{
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
