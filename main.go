package main

import (
	"goadminapi/modules/config"
	"log"
	"net/http"

	"goadminapi/engine"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/acme/autocert"

	_ "goadminapi/adapter/gin"     // 框架引擎
	_ "goadminapi/themes/adminlte" // 主題

	_ "github.com/denisenkom/go-mssqldb" // mssql引擎
	_ "github.com/go-sql-driver/mysql"   // mysql引擎
)

func main() {
	r := gin.Default()
	r.Static("/thumbnail", "./themes/adminlte/resource/assets/dist/thumbnail")
	r.Static("/original", "./themes/adminlte/resource/assets/dist/original")

	// 設定預設Engine(struct)
	eng := engine.Default()

	cfg := config.Config{
		Databases: config.DatabaseList{
			"default": {
				Host:       "35.194.236.160",
				Port:       "3306",
				User:       "yo",
				Pwd:        "yo123456",
				Name:       "gotest",
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
	//	更新
	// AddConfig 首先將參數cfg(struct)數值處理後設置至globalCfg，接著設置至Engine.config
	// 再來將driver加入Engine.Services，初始化所有資料庫連線並啟動引擎
	_ = eng.AddConfig(cfg).Use(r)

	log.Fatal(http.Serve(autocert.NewListener("hilive.com.tw"), r))
}
