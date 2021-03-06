package admin

import (
	"goadminapi/context"
	"goadminapi/modules/service"
	"goadminapi/plugins"
	"goadminapi/plugins/admin/controller"
	"goadminapi/plugins/admin/modules/guard"
	"goadminapi/plugins/admin/modules/table"
	"goadminapi/template/types/action"

	"goadminapi/modules/config"
)

// Admin也屬於Plugin(interface)所有方法
type Admin struct {
	*plugins.Base
	// GeneratorList類別為map[string]Generator，Generator類別為func(ctx *context.Context) Table
	tableList table.GeneratorList
	guardian  *guard.Guard
	handler   *controller.Handler
}

// NewAdmin 設置一個新的Admin(struct)
func NewAdmin(tableCfg ...table.GeneratorList) *Admin {
	return &Admin{
		tableList: make(table.GeneratorList).CombineAll(tableCfg),
		Base:      &plugins.Base{PlugName: "admin"},
		handler:   controller.New(),
	}
}

// --------------------------plugin(interface)的方法--------------------------

// 初始化router(放置api的地方)
func (admin *Admin) InitPlugin(services service.List) {
	// 將參數services(map[string]Service)設置至Admin.Base(struct)
	admin.InitBase(services)

	// GetService將參數services.Get("config")轉換成Service(struct)後回傳Service.C(Config struct)
	c := config.GetService(services.Get("config"))

	// 將參數設置至SystemTable(struct)
	st := table.NewSystemTable(admin.Conn, c)

	// Combine透過參數判斷GeneratorList已經有該key、value，如果不存在則加入該鍵與值
	// **************用於要判斷:__prefix需要取得map[string]Generator*********************
	admin.tableList.Combine(table.GeneratorList{
		"manager":    st.GetManagerTable,
		"permission": st.GetPermissionTable,
		"roles":      st.GetRolesTable,
		"menu":       st.GetMenuTable,

		// ***************目前先不設置*******************
		// "op":             st.GetOpTable,
		// "normal_manager": st.GetNormalManagerTable,
		// "site":           st.GetSiteTable,
		// "generate":       st.GetGenerateForm,
	})

	// 將參數admin.Services, admin.Conn, admin.tableList設置Admin.guardian(struct)後回傳
	admin.guardian = guard.New(admin.Services, admin.Conn, admin.tableList, admin.UI.NavButtons)

	// 將參數設置至Config(struct)
	handlerCfg := controller.Config{
		Config:     c,
		Services:   services,
		Generators: admin.tableList,
		Connection: admin.Conn,
	}

	// 將參數cfg(struct)裡的值都設置至Handler(struct)
	admin.handler.UpdateCfg(handlerCfg)

	// ***************放置api的地方*****************
	admin.initRouter()

	// 將路徑及method設置至Handler.routes
	admin.handler.SetRoutes(admin.App.Routers)

	admin.handler.AddNavButton(admin.UI.NavButtons)

	// 將參數(services)設置給services(map[string]Service)，services是套件中的全域變數
	table.SetServices(services)

	// InitOperationHandlerSetter 將參數admin.GetAddOperationFn()(func(...Node))設置給operationHandlerSetter
	// GetAddOperationFn回傳Admin.handler.AddOperation(func(...Node))
	action.InitOperationHandlerSetter(admin.GetAddOperationFn())
}

// --------------------------plugin(interface)的方法--------------------------

func (admin *Admin) GetAddOperationFn() context.NodeProcessor {
	return admin.handler.AddOperation
}
