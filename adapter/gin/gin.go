package gin

import (
	"bytes"
	"errors"
	"goadminapi/adapter"
	"goadminapi/engine"
	"goadminapi/plugins"
	"goadminapi/plugins/admin/models"
	"goadminapi/template/types"
	"net/http"
	"net/url"
	"strings"

	"goadminapi/context"

	"goadminapi/modules/config"

	"github.com/gin-gonic/gin"
)

// Gin同時也符合adapter.WebFrameWork(interface)
type Gin struct {
	// adapter.BaseAdapter(struct)裡面為db.Connection(interface)
	adapter.BaseAdapter
	// gin.Context(struct)為gin最重要的部分，允許在middleware傳遞變數(例如驗證請求、管理流程)
	ctx *gin.Context
	// app為框架中的實例，包含muxer,middleware ,configuration，藉由New() or Default()建立Engine
	app *gin.Engine
}

// 初始化
func init() {
	// 建立引擎預設的配適器
	engine.Register(new(Gin))
}

//-------------------------------------
// 下列為adapter.WebFrameWork(interface)的方法
// Gin(struct)也是adapter.WebFrameWork(interface)
//------------------------------------

// 回傳框架名稱
func (gins *Gin) Name() string {
	return "gin"
}

// 首先將參數(app)轉換成gin.Engine(/gin-gonic/gin套件)型態設置至Gin.app
// 接著對參數(plugin []plugins.Plugin)執行迴圈，設置Context(struct)並增加handlers、處理url及寫入header
func (gins *Gin) Use(app interface{}, plugs []plugins.Plugin) error {
	// 首先將參數(app)轉換成gin.Engine(/gin-gonic/gin套件)型態設置至Gin.app
	// 接著對參數(plugin []plugins.Plugin)執行迴圈，設置Context(struct)並增加handlers、處理url及寫入header
	return gins.GetUse(app, plugs, gins)
}

// 利用cookie驗證使用者，取得role、permission、menu，接著檢查權限，執行模板並導入HTML
func (gins *Gin) Content(ctx interface{}, getPanelFn types.GetPanelFn, fn context.NodeProcessor, btns ...types.Button) {
	// 利用cookie驗證使用者，取得role、permission、menu，接著檢查權限，執行模板並導入HTML
	gins.GetContent(ctx, getPanelFn, gins, btns, fn)
}

// 透過參數取得cookie後，利用cookie取得用戶角色、權限以及可用menu，最後將UserModel.Conn = nil後回傳UserModel
func (gins *Gin) User(ctx interface{}) (models.UserModel, bool) {
	// 取得用戶角色、權限以及可用menu
	return gins.GetUser(ctx, gins)
}

// 設置Context(struct)並增加handlers、處理url及寫入header
func (gins *Gin) AddHandler(method, path string, handlers context.Handlers) {
	// Handle第三個參數(主要處理程序)為funcion(*gin.Context)，gin.Context為struct(gin-gonic套件)
	gins.app.Handle(strings.ToUpper(method), path, func(c *gin.Context) {

		// 設置新Context(struct)，將參數(c.Request)設置至Context.Request
		ctx := context.NewContext(c.Request)

		// Context.Params類型為[]Context.Param，Param裡有key以及value(他是url參數的鍵與值)
		// 將參數設置在url中
		for _, param := range c.Params {
			if c.Request.URL.RawQuery == "" {
				c.Request.URL.RawQuery += strings.Replace(param.Key, ":", "", -1) + "=" + param.Value
			} else {
				c.Request.URL.RawQuery += "&" + strings.Replace(param.Key, ":", "", -1) + "=" + param.Value
			}
		}

		// SetHandlers將參數(handlers)設置至至Context.Handlers
		// 執行迴圈Context.handlers[ctx.index](ctx)
		ctx.SetHandlers(handlers).Next()

		for key, head := range ctx.Response.Header {
			c.Header(key, head[0])
		}

		if ctx.Response.Body != nil {
			buf := new(bytes.Buffer)
			_, _ = buf.ReadFrom(ctx.Response.Body)
			c.String(ctx.Response.StatusCode, buf.String())
		} else {
			c.Status(ctx.Response.StatusCode)
		}
	})
}

func (gins *Gin) DisableLog()                { panic("not implement") }
func (gins *Gin) Static(prefix, path string) { panic("not implement") }
func (gins *Gin) Run() error                 { panic("not implement") }

// 將參數(app)轉換成gin.Engine(gin-gonic/gin套件)型態設置至Gin.app
func (gins *Gin) SetApp(app interface{}) error {
	var (
		eng *gin.Engine
		ok  bool
	)
	// app.(*gin.Engine)將interface{}轉換為gin.Engine型態
	if eng, ok = app.(*gin.Engine); !ok {
		return errors.New("gin adapter SetApp: wrong parameter")
	}
	gins.app = eng
	return nil
}

// 將參數(contextInterface)轉換成gin.Context(gin-gonic/gin套件)類別Gin.ctx(struct)
func (gins *Gin) SetContext(contextInterface interface{}) adapter.WebFrameWork {
	var (
		ctx *gin.Context
		ok  bool
	)
	// 將contextInterface類別變成gin.Context(struct)
	if ctx, ok = contextInterface.(*gin.Context); !ok {
		panic("gin adapter SetContext: wrong parameter")
	}
	return &Gin{ctx: ctx}
}

// 取得session裡設置的cookie
func (gins *Gin) GetCookie() (string, error) {
	// Cookie()回傳cookie(藉由參數裡的命名回傳的)
	return gins.ctx.Cookie(gins.CookieKey())
}

// return  Gin.ctx.Request.URL.Path
func (gins *Gin) Path() string {
	return gins.ctx.Request.URL.Path
}

// return gins..ctx.Request.Method
func (gins *Gin) Method() string {
	return gins.ctx.Request.Method
}

// 解析參數(multipart/form-data裡的)
func (gins *Gin) FormParam() url.Values {
	_ = gins.ctx.Request.ParseMultipartForm(32 << 20)
	return gins.ctx.Request.PostForm
}

// 設置標頭 X-PJAX = true
func (gins *Gin) IsPjax() bool {
	return gins.ctx.Request.Header.Get("X-PJAX") == "true"
}

func (gins *Gin) SetContentType() {
	return
}

// 重新導向至登入頁面(出現錯誤)
func (gins *Gin) Redirect() {
	gins.ctx.Redirect(302, config.Url(config.GetLoginUrl()))
	gins.ctx.Abort()
}

// 將參數(body)寫入並更新http代碼
func (gins *Gin) Write(body []byte) {
	// Data將資料寫入body並更新http代碼
	// gins.HTMLContentType() return "text/html; charset=utf-8"
	gins.ctx.Data(http.StatusOK, gins.HTMLContentType(), body)
}
