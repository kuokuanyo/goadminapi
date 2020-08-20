package adapter

import (
	"bytes"
	"net/url"

	"goadminapi/context"
	"goadminapi/modules/auth"
	"goadminapi/modules/config"
	"goadminapi/modules/db"
	"goadminapi/modules/menu"
	"goadminapi/plugins"
	"goadminapi/plugins/admin/models"
	"goadminapi/template"
	"goadminapi/template/types"
)

// WebFrameWork 功能都設定在框架中(使用/adapter/gin/gin.go框架)
type WebFrameWork interface {
	// 回傳使用的web框架名稱
	Name() string

	// 將插件插入框架引擎中
	Use(app interface{}, plugins []plugins.Plugin) error

	// 添加html到框架中
	Content(ctx interface{}, fn types.GetPanelFn, fn2 context.NodeProcessor, navButtons ...types.Button)

	// 從給定的上下文中取得用戶模型
	User(ctx interface{}) (models.UserModel, bool)

	// 將路由(路徑)及處理程式加入框架
	AddHandler(method, path string, handlers context.Handlers)

	DisableLog()

	Static(prefix, path string)

	Run() error

	// Helper functions
	// ================================

	SetApp(app interface{}) error
	SetConnection(db.Connection)
	GetConnection() db.Connection
	SetContext(ctx interface{}) WebFrameWork
	GetCookie() (string, error)
	Path() string
	Method() string
	FormParam() url.Values
	IsPjax() bool
	Redirect()
	SetContentType()
	Write(body []byte)
	CookieKey() string
	HTMLContentType() string
}

// BaseAdapter 是db.Connection(interface)
type BaseAdapter struct {
	db db.Connection
}

// SetConnection 將參數(conn)設置至BaseAdapter.db
func (base *BaseAdapter) SetConnection(conn db.Connection) {
	base.db = conn
}

// GetConnection 回傳BaseAdapter.db
func (base *BaseAdapter) GetConnection() db.Connection {
	return base.db
}

// CookieKey return "session"
func (base *BaseAdapter) CookieKey() string {
	return "session"
}

// HTMLContentType return "text/html; charset=utf-8"
func (base *BaseAdapter) HTMLContentType() string {
	return "text/html; charset=utf-8"
}

// GetUse 首先將參數(app)轉換成gin.Engine(/gin-gonic/gin套件)型態設置至Gin.app
// 接著對參數(plugin []plugins.Plugin)執行迴圈，設置Context(struct)並增加handlers、處理url及寫入header
// -------wf參數是Gin(struct)------------
func (base *BaseAdapter) GetUse(app interface{}, plugin []plugins.Plugin, wf WebFrameWork) error {
	// 將參數轉換成gin.Engine(/gin-gonic/gin套件)型態設置至Gin.app
	if err := wf.SetApp(app); err != nil {
		return err
	}

	// plugin is interface
	for _, plug := range plugin {
		// 回傳Base.App.Handlers(map[Path]Handlers)，path(struct)裡包含URL、method
		for path, handlers := range plug.GetHandler() {
			// 設置Context(struct)並增加handlers、處理url及寫入header
			wf.AddHandler(path.Method, path.URL, handlers)
		}
	}
	return nil
}

// GetUser 透過參數取得cookie後，利用cookie取得用戶角色、權限以及可用menu，最後將UserModel.Conn = nil後回傳UserModel
// -------wf參數是Gin(struct)------------
func (base *BaseAdapter) GetUser(ctx interface{}, wf WebFrameWork) (models.UserModel, bool) {
	// SetContext將參數(ctx)轉換成gin.Context(gin-gonic/gin套件)類別Gin.ctx(struct)
	// 取得cookie
	cookie, err := wf.SetContext(ctx).GetCookie()
	if err != nil {
		return models.UserModel{}, false
	}
	// wf.GetConnection()回傳BaseAdapter.db(interface)
	// 透過cookie、conn可以得到角色、權限以及可使用菜單
	user, exist := auth.GetCurUser(cookie, wf.GetConnection())

	// 設置UserModel.Conn = nil後回傳UserModel
	return user.ReleaseConn(), exist
}

// GetContent 利用cookie驗證使用者，取得role、permission、menu，接著檢查權限，執行模板並導入HTML
// -------wf參數是Gin(struct)------------
func (base *BaseAdapter) GetContent(ctx interface{}, getPanelFn types.GetPanelFn, wf WebFrameWork,
	navButtons types.Buttons, fn context.NodeProcessor) {

	var (
		// 將參數(ctx)轉換成gin.Context(gin-gonic/gin套件)類別Gin.ctx(struct)
		// -------wf參數是Gin(struct)------------
		newBase = wf.SetContext(ctx)
		// 取得session裡設置的cookie
		cookie, hasError = newBase.GetCookie()
	)
	if hasError != nil || cookie == "" {
		newBase.Redirect()
		return
	}

	// wf.GetConnection()回傳BaseAdapter.db(interface)
	// 透過參數sesKey(cookie)取得id並利用id取得該user的role、permission以及可用menu，最後回傳UserModel(struct)
	user, authSuccess := auth.GetCurUser(cookie, wf.GetConnection())
	if !authSuccess {
		newBase.Redirect()
		return
	}

	var (
		panel types.Panel
		err   error
	)

	// CheckPermissions檢查用戶權限(在modules\auth\middleware.go)
	if !auth.CheckPermissions(user, newBase.Path(), newBase.Method(), newBase.FormParam()) {
		panel = template.WarningPanel("no permission", template.NoPermission403Page)
	} else {
		panel, err = getPanelFn(ctx)
		if err != nil {
			panel = template.WarningPanel(err.Error())
		}
	}

	fn(panel.Callbacks...)

	tmpl, tmplName := template.Default().GetTemplate(newBase.IsPjax())

	buf := new(bytes.Buffer)

	hasError = tmpl.ExecuteTemplate(buf, tmplName, types.NewPage(types.NewPageParam{
		User:         user,
		Menu:         menu.GetGlobalMenu(user, wf.GetConnection()).SetActiveClass(config.URLRemovePrefix(newBase.Path())),
		Panel:        panel.GetContent(config.IsProductionEnvironment()),
		Assets:       template.GetComponentAssetImportHTML(),
		Buttons:      navButtons.CheckPermission(user),
		TmplHeadHTML: template.Default().GetHeadHTML(),
		TmplFootJS:   template.Default().GetFootJS(),
	}))

	newBase.SetContentType()

	newBase.Write(buf.Bytes())
}
