package controller

import (
	"bytes"
	"goadminapi/modules/auth"
	c "goadminapi/modules/config"
	"goadminapi/modules/menu"
	"goadminapi/template/icon"
	"goadminapi/template/types"
	template2 "html/template"
	"regexp"

	"goadminapi/context"
	"goadminapi/modules/db"
	"goadminapi/modules/service"
	"goadminapi/plugins/admin/models"
	"goadminapi/plugins/admin/modules/table"
	"sync"

	"goadminapi/template"
)

type Handler struct {
	config        *c.Config
	captchaConfig map[string]string
	services      service.List
	conn          db.Connection
	routes        context.RouterMap
	generators    table.GeneratorList // map[string]Generator
	operations    []context.Node
	navButtons    *types.Buttons
	operationLock sync.Mutex
}

type Config struct {
	Config     *c.Config
	Services   service.List
	Connection db.Connection
	Generators table.GeneratorList
}

// 判斷參數cfg(長度是否大於0)後設置Handler(struct)並回傳
func New(cfg ...Config) *Handler {
	if len(cfg) == 0 {
		return &Handler{
			operations: make([]context.Node, 0),
			// navButtons: new(types.Buttons),
		}
	}
	return &Handler{
		config:     cfg[0].Config,
		services:   cfg[0].Services,
		conn:       cfg[0].Connection,
		generators: cfg[0].Generators,
		operations: make([]context.Node, 0),
		// navButtons: new(types.Buttons),
	}
}

// 將參數設置至ExecuteParam(struct)，接著將給定的數據寫入buf(struct)並回傳
func (h *Handler) Execute(ctx *context.Context, user models.UserModel, panel types.Panel, animation ...bool) *bytes.Buffer {
	tmpl, tmplName := aTemplate().GetTemplate(isPjax(ctx))

	return template.Execute(template.ExecuteParam{
		User:      user,
		TmplName:  tmplName,
		Tmpl:      tmpl,
		Panel:     panel,
		Config:    *h.config,
		Menu:      menu.GetGlobalMenu(user, h.conn).SetActiveClass(h.config.URLRemovePrefix(ctx.Path())),
		Animation: len(animation) > 0 && animation[0] || len(animation) == 0,
		Buttons:   (*h.navButtons).CheckPermission(user),
		Iframe:    ctx.Query("__iframe") == "true",
	})
}

// 藉由參數取得Router(struct)
func (h *Handler) route(name string) context.Router {
	return h.routes.Get(name)
}

func isNewUrl(s string, p string) bool {
	reg, _ := regexp.Compile("(.*?)info/" + p + "/new")

	return reg.MatchString(s)
}

func isEditUrl(s string, p string) bool {
	reg, _ := regexp.Compile("(.*?)info/" + p + "/edit")
	return reg.MatchString(s)
}

// 判斷templateMap(map[string]Template)的key鍵是否參數globalCfg.Theme，有則回傳Template(interface)
func aTemplate() template.Template {
	// 判斷templateMap(map[string]Template)的key鍵是否參數globalCfg.Theme，有則回傳Template(interface)
	// GetTheme return globalCfg.Theme
	return template.Get(c.GetTheme())
}

// 將參數(r)設置至Handler.routes
func (h *Handler) SetRoutes(r context.RouterMap) {
	h.routes = r
}

// 將參數cfg(struct)裡的值都設置至Handler(struct)
func (h *Handler) UpdateCfg(cfg Config) {
	h.config = cfg.Config
	h.services = cfg.Services
	h.conn = cfg.Connection
	h.generators = cfg.Generators
}

// 透過參數name取得該路徑名稱的URL、如果參數value大於零，則處理url中有:__的字串
func (h *Handler) routePath(name string, value ...string) string {
	// Get藉由參數name取得Router(struct)，Router裡有Methods([]string)及Pattern(string)
	// GetURL處理URL後回傳(處理url中有:__的字串)
	return h.routes.Get(name).GetURL(value...)
}

// 透過參數name取得該路徑名稱的URL，將url中的:__prefix改成第二個參數(prefix)
func (h *Handler) routePathWithPrefix(name string, prefix string) string {
	return h.routePath(name, "prefix", prefix)
}

// searchOperation 在Handler.operations([]context.Node)執行迴圈，如果條件符合參數path、method則回傳true(代表已經存在)
func (h *Handler) searchOperation(path, method string) bool {
	for _, node := range h.operations {
		if node.Path == path && node.Method == method {
			return true
		}
	}
	return false
}

// AddOperation 判斷條件後將參數(context.Node)添加至Handler.operations
func (h *Handler) AddOperation(nodes ...context.Node) {
	h.operationLock.Lock()

	addNodes := make([]context.Node, 0)
	for _, node := range nodes {
		// 在Handler.operations([]context.Node)執行迴圈，如果條件符合參數path、method則回傳true
		// 代表Handler.operations裡已經存在，則不添加
		if h.searchOperation(node.Path, node.Method) {
			continue
		}
		addNodes = append(addNodes, node)
	}
	h.operations = append(h.operations, addNodes...)
}

// table 先透過參數prefix取得Table(interface)，接著判斷條件後將[]context.Node加入至Handler.operations後回傳
func (h *Handler) table(prefix string, ctx *context.Context) table.Table {
	// 透過參數prefix執行函式取得table(interface)
	t := h.generators[prefix](ctx)

	// 建立Invoker(Struct)並透過參數ctx取得UserModel，並且取得該user的role、權限與可用menu，最後檢查用戶權限
	// GetConnection取得匹配的service.Service然後轉換成Connection(interface)類別
	authHandler := auth.Middleware(db.GetConnection(h.services))

	// GetInfo 將參數值設置至base.Info(InfoPanel(struct)).primaryKey
	for _, cb := range t.GetInfo().Callbacks {
		if cb.Value["need_auth"] == 1 {
			// AddOperation 判斷條件後將參數(context.Node)添加至Handler.operations
			h.AddOperation(context.Node{
				Path:     cb.Path,
				Method:   cb.Method,
				Handlers: append([]context.Handler{authHandler}, cb.Handlers...),
			})
		} else {
			h.AddOperation(context.Node{Path: cb.Path, Method: cb.Method, Handlers: cb.Handlers})
		}
	}

	// GetForm 將參數值設置至base.Info(FormPanel(struct)).primaryKey
	for _, cb := range t.GetForm().Callbacks {
		if cb.Value["need_auth"] == 1 {
			// 判斷條件後將參數(類別context.Node)添加至Handler.operations
			h.AddOperation(context.Node{
				Path:     cb.Path,
				Method:   cb.Method,
				Handlers: append([]context.Handler{authHandler}, cb.Handlers...),
			})
		} else {
			h.AddOperation(context.Node{Path: cb.Path, Method: cb.Method, Handlers: cb.Handlers})
		}
	}
	return t
}

// filterFormFooter 取得過濾表單中的按鈕(搜尋、重置)...等HTML語法
func filterFormFooter(infoUrl string) template2.HTML {
	col1 := aCol().SetSize(types.SizeMD(2)).GetContent()
	// 搜尋按鈕HTML
	btn1 := aButton().SetType("submit").
		SetContent(icon.Icon("fa-search", 2) + template.HTML("search")).
		SetThemePrimary().
		SetSmallSize().
		SetOrientationLeft().
		SetLoadingText(icon.Icon("fa-search", 1) + template.HTML("search")).
		GetContent()
	// 重置按鈕HTML
	btn2 := aButton().SetType("reset").
		SetContent(icon.Icon("fa-undo", 2) + template.HTML("reset")).
		SetThemeDefault().
		SetOrientationLeft().
		SetSmallSize().
		SetHref(infoUrl).
		SetMarginLeft(12).
		GetContent()
	col2 := aCol().SetSize(types.SizeMD(8)).
		SetContent(btn1 + btn2).GetContent()
	return col1 + col2
}

// 將參數h.services.Get(auth.TokenServiceKey)轉換成TokenService(struct)類別後回傳
func (h *Handler) authSrv() *auth.TokenService {
	return auth.GetTokenService(h.services.Get("token_csrf_helper"))
}

// 判斷是否header X-PJAX:true
func isPjax(ctx *context.Context) bool {
	return ctx.IsPjax()
}

func aAlert() types.AlertAttribute {
	return aTemplate().Alert()
}

func aDataTable() types.DataTableAttribute {
	return aTemplate().DataTable()
}

func aTab() types.TabsAttribute {
	return aTemplate().Tabs()
}

func aBox() types.BoxAttribute {
	return aTemplate().Box()
}

func aForm() types.FormAttribute {
	return aTemplate().Form()
}

func aCol() types.ColAttribute {
	return aTemplate().Col()
}

func aButton() types.ButtonAttribute {
	return aTemplate().Button()
}
