package main

import (
	"goadminapi/engine"
	"goadminapi/modules/config"

	"github.com/gin-gonic/gin"

	_ "goadminapi/adapter/gin"     // 框架引擎
	_ "goadminapi/themes/adminlte" // 主題

	_ "github.com/denisenkom/go-mssqldb" // mssql引擎
	_ "github.com/go-sql-driver/mysql"   // mysql引擎
)

func main() {
	r := gin.Default()

	// 設定預設Engine(struct)
	eng := engine.Default()

	cfg := config.Config{
		Databases: config.DatabaseList{
			"default": {
				Host:       "139.162.53.127",
				Port:       "1433",
				User:       "sa",
				Pwd:        "CCO@test53383499",
				Name:       "gotest",
				MaxIdleCon: 50,
				MaxOpenCon: 150,
				Driver:     "mssql",
			},
		},
		UrlPrefix: "admin",
		Store: config.Store{
			Path:   "./uploads",
			Prefix: "uploads",
		},
	}

	// AddConfig 首先將參數cfg(struct)數值處理後設置至globalCfg，接著設置至Engine.config
	// 再來將driver加入Engine.Services，初始化所有資料庫連線並啟動引擎
	_ = eng.AddConfig(cfg).Use(r)

	_ = r.Run(":8080")
}
