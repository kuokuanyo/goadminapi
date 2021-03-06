package plugins

import (
	"goadminapi/context"
	"goadminapi/modules/db"
	"goadminapi/modules/service"
	"goadminapi/modules/ui"
)

// Base(struct)包含Plugin(interface)所有方法
type Base struct {
	App       *context.App
	Services  service.List
	Conn      db.Connection
	UI        *ui.Service
	PlugName  string
	URLPrefix string
}

// GetRequest回傳插件中的所有路徑
// InitPlugin初始化插件，類似於初始化資料庫並設置及配置路徑
type Plugin interface {
	GetHandler() context.HandlerMap
	InitPlugin(services service.List)
	Name() string
	Prefix() string
}

// --------------------------plugin(interface)的方法--------------------------

// 回傳Base.App.Handlers(map[Path]Handlers)，path(struct)裡包含URL、method
func (b *Base) GetHandler() context.HandlerMap {
	return b.App.Handlers
}
// 回傳Base.PlugName
func (b *Base) Name() string {
	return b.PlugName
}

// 回傳Base.URLPrefix
func (b *Base) Prefix() string {
	return b.URLPrefix
}

// --------------------------plugin(interface)的方法--------------------------

// 透過參數srv(map[string]Service)設置至Base(struct).Services並且設置Base.Conn、Base.UI
func (b *Base) InitBase(srv service.List) {
	b.Services = srv
	b.Conn = db.GetConnection(b.Services)
	b.UI = ui.GetService(b.Services)
}

