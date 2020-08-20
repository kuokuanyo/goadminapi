package engine

import (
	"goadminapi/adapter"
	"goadminapi/modules/config"
	"goadminapi/modules/db"
	"goadminapi/modules/service"
	"goadminapi/modules/ui"
	"goadminapi/plugins"
	"goadminapi/plugins/admin"
	"goadminapi/plugins/admin/models"
	"goadminapi/template/types"
)

var defaultAdapter adapter.WebFrameWork
var engine *Engine
var navButtons = new(types.Buttons)

// Engine 核心組件，有PluginList及Adapter兩個屬性
type Engine struct {
	PluginList []plugins.Plugin
	Adapter    adapter.WebFrameWork
	Services   service.List //Services類別為map[string]Service，Service為interface(Name方法)
	NavButtons *types.Buttons
	config     *config.Config
}

// Register 建立引擎預設的配適器
func Register(ada adapter.WebFrameWork) {
	if ada == nil {
		panic("adapter is nil")
	}
	defaultAdapter = ada
}

// Default 設定預設Engine(struct)
func Default() *Engine {
	engine = &Engine{
		Adapter:    defaultAdapter,
		Services:   service.GetServices(),
		NavButtons: new(types.Buttons),
	}
	return engine
}

// DB 透過參數(driver)找到匹配的Service(interface)後轉換成Connection(interface)型態
func (eng *Engine) DB(driver string) db.Connection {
	// GetConnectionFromService將參數型態轉換成Connection(interface)後回傳
	// Get藉由參數(driver)取得匹配的Service(interface)
	return db.GetConnectionFromService(eng.Services.Get(driver))
}

// DefaultConnection 透過資料庫引擎(driver)回傳預設的Connection(interface)
func (eng *Engine) DefaultConnection() db.Connection {
	// GetDefault() = DatabaseList["default"]
	// 參數為DatabaseList["default"].driver
	return eng.DB(eng.config.Databases.GetDefault().Driver)
}

// ============================
// Config APIs
// ============================

// setConfig 將參數值處理後設置至全局變數globalCfg，最後設置至Engine.config
// --------此函式處理全局變數globalCfg-------------
func (eng *Engine) setConfig(cfg config.Config) *Engine {
	// 設置Config(struct)title、theme、登入url、前綴url...資訊，如果參數cfg(struct)有些數值為空值，設置預設值
	eng.config = config.Set(cfg)

	//------------------------------------------------------------------------------------------------
	// 檢查版本及主題是否正確
	// CheckRequirements在template\template.go
	// sysCheck, themeCheck := template.CheckRequirements()
	// if !sysCheck {
	// 	panic(fmt.Sprintf("wrong GoAdmin version, theme %s required GoAdmin version are %s",
	// 		eng.config.Theme, strings.Join(template.Default().GetRequirements(), ",")))
	// }
	// if !themeCheck {
	// 	panic(fmt.Sprintf("wrong Theme version, GoAdmin %s required Theme version are %s",
	// 		system.Version(), strings.Join(system.RequireThemeVersion()[eng.config.Theme], ",")))
	// }
	//------------------------------------------------------------------------------------------------
	return eng
}

// InitDatabase 將driver加入Engine.Services，初始化所有資料庫並啟動引擎
func (eng *Engine) InitDatabase() *Engine {
	// GroupByDriver將資料庫依照資料庫引擎分組(ex:mysql一組mssql一組)
	// driver = mysql、mssql等引擎名稱
	for driver, databaseCfg := range eng.config.Databases.GroupByDriver() {
		// Add藉由參數新增List(map[string]Service)，在List加入引擎(driver)
		// GetConnectionByDriver藉由參數(driver = mysql、mssql...)取得Connection(interface)
		// InitDB初始化資料庫並啟動引擎
		eng.Services.Add(driver, db.GetConnectionByDriver(driver).InitDB(databaseCfg))
	}
	if defaultAdapter == nil {
		panic("adapter is nil")
	}
	return eng
}

// AddConfig 首先將參數cfg(struct)數值處理後設置至globalCfg，接著設置至Engine.config
// 最後將資料庫引擎加入Engine.Services，初始化所有資料庫連線並啟動引擎
func (eng *Engine) AddConfig(cfg config.Config) *Engine {
	// setConfig將參數cfg(struct)數值處理後設置至globalCfg，最後設置至Engine.config
	// InitDatabase將driver加入Engine.Services，初始化所有資料庫連線並啟動引擎
	return eng.setConfig(cfg).InitDatabase()
}

// FindPluginByName 在PluginList([]plugins.Plugin)的迴圈中尋找與參數(name)符合的plugin，如果有回傳Plugin,true，反之nil, false
func (eng *Engine) FindPluginByName(name string) (plugins.Plugin, bool) {
	for _, plug := range eng.PluginList {
		if plug.Name() == name {
			return plug, true
		}
	}
	return nil, false
}

// Use 尋找符合的plugin，接著設置context.Context(struct)與設置url與寫入header，取得新的request與middleware
func (eng *Engine) Use(router interface{}) error {
	if eng.Adapter == nil {
		panic("adapter is nil, import the default adapter or use AddAdapter method add the adapter")
	}

	// 在PluginList([]plugins.Plugin)的迴圈中尋找與參數(name)符合的plugin，如果有回傳Plugin,true，反之nil, false
	_, exist := eng.FindPluginByName("admin")
	if !exist {
		eng.PluginList = append(eng.PluginList, admin.NewAdmin())
	}

	// DefaultConnection透過資料庫引擎(driver)回傳預設的Connection(interface)
	site := models.Site().SetConn(eng.DefaultConnection())

	// ToMap將Config的值設置至map[string]string
	site.Init(eng.config.ToMap())

	// 更新Config(struct)值(從site資料表資訊更新)
	_ = eng.config.Update(site.AllToMap())

	// 藉由參數新增List(map[string]Service)，新增config
	eng.Services.Add("config", config.SrvWithConfig(eng.config))

	//------------------------------------------------
	// errors.Init()

	// 隱藏設置入口
	// if !eng.config.HideConfigCenterEntrance {
	// 	*eng.NavButtons = (*eng.NavButtons).AddNavButton(icon.Gear, types.NavBtnSiteName,
	// 		action.JumpInNewTab(config.Url("/info/site/edit"),
	// 			language.GetWithScope("site setting", "config")))
	// }

	// // 隱藏App Info入口
	// if !eng.config.HideAppInfoEntrance {
	// 	*eng.NavButtons = (*eng.NavButtons).AddNavButton(icon.Info, types.NavBtnInfoName,
	// 		action.JumpInNewTab(config.Url("/application/info"),
	// 			language.GetWithScope("system info", "system")))
	// }

	// if !eng.config.HideToolEntrance {
	// 	*eng.NavButtons = (*eng.NavButtons).AddNavButton(icon.Wrench, types.NavBtnToolName,
	// 		action.JumpInNewTab(config.Url("/info/generate/new"),
	// 			language.GetWithScope("tool", "tool")))
	// }
	//------------------------------------------------

	navButtons = eng.NavButtons

	// 藉由參數新增List(map[string]Service)，新增ui
	eng.Services.Add("ui", ui.NewService(eng.NavButtons))

	// 取得匹配的eng.Services然後轉換成Connection(interface)類別
	defaultConnection := db.GetConnection(eng.Services)

	// SetConnection為WebFrameWork(interface)的方法
	//設定連線
	defaultAdapter.SetConnection(defaultConnection)
	eng.Adapter.SetConnection(defaultConnection)

	// Initialize plugins
	for i := range eng.PluginList {
		eng.PluginList[i].InitPlugin(eng.Services)
	}

	// 首先將參數(app)轉換成gin.Engine(/gin-gonic/gin套件)型態設置至Gin.app
	// 接著對參數(plugin []plugins.Plugin)執行迴圈，設置Context(struct)並增加handlers、處理url及寫入header
	return eng.Adapter.Use(router, eng.PluginList)
}
