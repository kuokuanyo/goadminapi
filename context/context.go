package context

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

// Context gin.Context(gin-gonic套件)的簡化版本
type Context struct {
	Request   *http.Request
	Response  *http.Response
	UserValue map[string]interface{}
	index     int8
	handlers  Handlers
}

// Handler 用於middleware
type Handler func(ctx *Context) 

// Handlers is []Handler
type Handlers []Handler

// HandlerMap is map[Path]Handlers
type HandlerMap map[Path]Handlers

// RouterMap is map[string]Router
type RouterMap map[string]Router

// Router 包含方法([]string)、模式(url)
type Router struct {
	Methods []string
	Patten  string //url
}

// Path 包含url、method
type Path struct {
	URL    string
	Method string
}

// App struct
type App struct {
	Requests    []Path     //Path包含url、method
	Handlers    HandlerMap // HandlerMap是 map[Path]Handlers，Handlers類別為[]Handler，Handler類別為func(ctx *Context)
	Middlewares Handlers
	Prefix      string

	Routers    RouterMap // RouterMap類別為map[string]Router，Router(struct)裡有methods、patten
	routeIndex int
	routeANY   bool
}

// RouterGroup struct
type RouterGroup struct {
	app         *App     //struct
	Middlewares Handlers //Handlers([]Handler)，Handler類別為 func(ctx *Context)
	Prefix      string
}

// NodeProcessor is func(...Node)
type NodeProcessor func(...Node)

// Node struct
type Node struct {
	Path     string
	Method   string
	Handlers []Handler
	Value    map[string]interface{}
}

// NewApp 取得app(struct)，空的
func NewApp() *App {
	return &App{
		Requests:    make([]Path, 0),
		Handlers:    make(HandlerMap),
		Prefix:      "/",
		Middlewares: make([]Handler, 0),
		routeIndex:  -1,
		Routers:     make(RouterMap),
	}
}

// NewContext 設置新Context(struct)，將參數(req)設置至Context.Request
func NewContext(req *http.Request) *Context {
	return &Context{
		Request:   req,
		UserValue: make(map[string]interface{}),
		Response: &http.Response{
			StatusCode: http.StatusOK,
			Header:     make(http.Header),
		},
		index: -1,
	}
}

// SetHandlers 將參數(handlers)設置至至Context.Handlers
func (ctx *Context) SetHandlers(handlers Handlers) *Context {
	ctx.handlers = handlers
	return ctx
}

// Name 將參數設置至App.Routers(RouterMap)中，設定methods及patten(url)
func (app *App) Name(name string) {
	if app.routeANY {
		// Routers(RouterMap)類別為map[string]Router，Router(struct)裡有methods、patten
		//新增所有方法
		app.Routers[name] = Router{
			Methods: []string{"POST", "GET", "DELETE", "PUT", "OPTIONS", "HEAD"},
			Patten:  app.Requests[app.routeIndex].URL,
		}
	} else {
		app.Routers[name] = Router{
			Methods: []string{app.Requests[app.routeIndex].Method},
			Patten:  app.Requests[app.routeIndex].URL,
		}
	}
}

// User 回傳目前登入的用戶(Context.UserValue["user"])
func (ctx *Context) User() interface{} {
	return ctx.UserValue["user"]
}

// SetUserValue 藉由參數key、value設定Context.UserValue
func (ctx *Context) SetUserValue(key string, value interface{}) {
	ctx.UserValue[key] = value
}

// Path 回傳Request.url.path
func (ctx *Context) Path() string {
	return ctx.Request.URL.Path
}

// Method return method
func (ctx *Context) Method() string {
	return ctx.Request.Method
}

// Group 將參數prefix、middleware新增至RouterGroup(struct)
func (app *App) Group(prefix string, middleware ...Handler) *RouterGroup {
	return &RouterGroup{
		app:         app,
		Middlewares: append(app.Middlewares, middleware...),
		Prefix:      slash(prefix),
	}
}

// Abort abort the context.
func (ctx *Context) Abort() {
	ctx.index = 63
}

// Next 執行迴圈Context.handlers[ctx.index](ctx)
func (ctx *Context) Next() {
	ctx.index++
	// Context.Handlers類別為[]Handler，Handler類別為func(ctx *Context)
	for s := int8(len(ctx.handlers)); ctx.index < s; ctx.index++ {
		ctx.handlers[ctx.index](ctx)
	}
}

// Query 取得Request url(在url中)裡的參數(key)
func (ctx *Context) Query(key string) string {
	return ctx.Request.URL.Query().Get(key)
}

// FormValue 藉由參數key取得multipart/form-data中的值
func (ctx *Context) FormValue(key string) string {
	return ctx.Request.FormValue(key)
}

// AddHeader 將參數(key、value)添加header中(Context.Response.Header)
func (ctx *Context) AddHeader(key, value string) {
	ctx.Response.Header.Add(key, value)
}

// Headers 藉由參數key獲得Header
func (ctx *Context) Headers(key string) string {
	return ctx.Request.Header.Get(key)
}

// SetCookie 設置cookie在response header Set-Cookie中
func (ctx *Context) SetCookie(cookie *http.Cookie) {
	if v := cookie.String(); v != "" {
		ctx.AddHeader("Set-Cookie", v)
	}
}

// SetContentType 將參數添加至Content-Type
func (ctx *Context) SetContentType(contentType string) {
	ctx.AddHeader("Content-Type", contentType)
}

// SetStatusCode 將參數設置至Context.Response.StatusCode
func (ctx *Context) SetStatusCode(code int) {
	ctx.Response.StatusCode = code
}

// Redirect 添加重新導向的url(參數path)至header
func (ctx *Context) Redirect(path string) {
	ctx.Response.StatusCode = http.StatusFound
	ctx.SetContentType("text/html; charset=utf-8")
	ctx.AddHeader("Location", path)
}

// WantHTML 判斷method是否為get以及header裡包含accept:html
func (ctx *Context) WantHTML() bool {
	return ctx.Method() == "GET" && strings.Contains(ctx.Headers("Accept"), "html")
}

// WantJSON 判斷header裡包含accept:json
func (ctx *Context) WantJSON() bool {
	return strings.Contains(ctx.Headers("Accept"), "json")
}

// PostForm 取得表單的值(所有)，參數放於multipart/form-data.
func (ctx *Context) PostForm() url.Values {
	_ = ctx.Request.ParseMultipartForm(32 << 20)
	return ctx.Request.PostForm
}

// Write 將狀態碼、標頭(header)及body寫入Context.Response
func (ctx *Context) Write(code int, header map[string]string, Body string) {
	ctx.Response.StatusCode = code
	for key, head := range header {
		// 加入header
		ctx.AddHeader(key, head)
	}
	ctx.Response.Body = ioutil.NopCloser(strings.NewReader(Body))
}

// WriteString 將參數body保存至Context.response.Body中
func (ctx *Context) WriteString(body string) {
	ctx.Response.Body = ioutil.NopCloser(strings.NewReader(body))
}

// DataWithHeaders 將code, headers and body(參數)設置至在Context.Response中
func (ctx *Context) DataWithHeaders(code int, header map[string]string, data []byte) {
	ctx.Response.StatusCode = code
	for key, head := range header {
		//添加標頭
		ctx.AddHeader(key, head)
	}
	ctx.Response.Body = ioutil.NopCloser(bytes.NewBuffer(data))
}

// HTML 輸出HTML，參數body保存至Context.response.Body及設置ContentType、StatusCode
func (ctx *Context) HTML(code int, body string) {
	ctx.SetContentType("text/html; charset=utf-8")
	ctx.SetStatusCode(code)
	// 將參數body保存至Context.response.Body中
	ctx.WriteString(body)
}

// JSON 轉換成JSON存至Context.Response.body
func (ctx *Context) JSON(code int, Body map[string]interface{}) {
	ctx.Response.StatusCode = code
	//設定SetContentType
	ctx.SetContentType("application/json")
	// Marshal將struct轉成json
	BodyStr, err := json.Marshal(Body)
	if err != nil {
		panic(err)
	}
	ctx.Response.Body = ioutil.NopCloser(bytes.NewReader(BodyStr))
}

// IsPjax 判斷是否header X-PJAX:true
func (ctx *Context) IsPjax() bool {
	return ctx.Headers("X-PJAX") == "true"
}

// Method 取得Method
func (r Router) Method() string {
	return r.Methods[0]
}

// GetURL 處理URL後回傳(處理url中有:__的字串)
func (r Router) GetURL(value ...string) string {
	u := r.Patten

	// 未處理前ex:/admin/info/:__prefix/edit
	for i := 0; i < len(value); i += 2 {
		// 處理url
		u = strings.Replace(u, ":__"+value[i], value[i+1], -1)
	}
	// 處理後 ex:/admin/info/roles/edit
	return u
}

// Get 藉由參數name取得Router(struct)，Router裡有Methods([]string)及Pattern(string)
func (r RouterMap) Get(name string) Router {
	return r[name]
}

// AppendReqAndResp stores the request info and handle into app.
// support the route parameter. The route parameter will be recognized as
// wildcard store into the RegUrl of Path struct. For example:
//
//         /user/:id      => /user/(.*)
//         /user/:id/info => /user/(.*?)/info
//
// The RegUrl will be used to recognize the incoming path and find
// the handler.
// AppendReqAndResp 在RouterGroup.app(struct)中新增Requests([]Path)路徑及方法、接著在該url中新增參數handler(Handler...)
func (g *RouterGroup) AppendReqAndResp(url, method string, handler []Handler) {

	g.app.Requests = append(g.app.Requests, Path{
		URL:    join(g.Prefix, url),
		Method: method,
	})
	g.app.routeIndex++

	var h = make([]Handler, len(g.Middlewares))
	copy(h, g.Middlewares)

	g.app.Handlers[Path{
		URL:    join(g.Prefix, url),
		Method: method,
	}] = append(h, handler...)
}

// Group 將參數prefix、middleware新增至RouterGroup(struct)
func (g *RouterGroup) Group(prefix string, middleware ...Handler) *RouterGroup {
	return &RouterGroup{
		app:         g.app,
		Middlewares: append(g.Middlewares, middleware...),
		Prefix:      join(slash(g.Prefix), slash(prefix)),
	}
}

// POST 等於在AppendReqAndResp(url, "post", handler)
func (g *RouterGroup) POST(url string, handler ...Handler) *RouterGroup {
	g.app.routeANY = false
	g.AppendReqAndResp(url, "post", handler)
	return g
}

// GET 等於在AppendReqAndResp(url, "get", handler)
func (g *RouterGroup) GET(url string, handler ...Handler) *RouterGroup {
	g.app.routeANY = false
	g.AppendReqAndResp(url, "get", handler)
	return g
}

// DELETE 等於在AppendReqAndResp(url, "delete", handler)
func (g *RouterGroup) DELETE(url string, handler ...Handler) *RouterGroup {
	g.app.routeANY = false
	g.AppendReqAndResp(url, "delete", handler)
	return g
}

// PUT 等於在AppendReqAndResp(url, "put", handler)
func (g *RouterGroup) PUT(url string, handler ...Handler) *RouterGroup {
	g.app.routeANY = false
	g.AppendReqAndResp(url, "put", handler)
	return g
}

// OPTIONS 等於在AppendReqAndResp(url, "options", handler)
func (g *RouterGroup) OPTIONS(url string, handler ...Handler) *RouterGroup {
	g.app.routeANY = false
	g.AppendReqAndResp(url, "options", handler)
	return g
}

// HEAD 等於在AppendReqAndResp(url, "head", handler)
func (g *RouterGroup) HEAD(url string, handler ...Handler) *RouterGroup {
	g.app.routeANY = false
	g.AppendReqAndResp(url, "head", handler)
	return g
}

// ANY registers a route that matches all the HTTP methods.
// GET, POST, PUT, HEAD, OPTIONS, DELETE.
// ANY 執行所有方法的AppendReqAndResp(url, 方法, handler)
func (g *RouterGroup) ANY(url string, handler ...Handler) *RouterGroup {
	g.app.routeANY = true
	g.AppendReqAndResp(url, "post", handler)
	g.AppendReqAndResp(url, "get", handler)
	g.AppendReqAndResp(url, "delete", handler)
	g.AppendReqAndResp(url, "put", handler)
	g.AppendReqAndResp(url, "options", handler)
	g.AppendReqAndResp(url, "head", handler)
	return g
}

// Name 將參數設置至App.Routers(RouterMap)中，設定methods及patten(url)
func (g *RouterGroup) Name(name string) {
	// RouterGroup.App(struct)的Name方法
	g.app.Name(name)
}

// slash fix the path which has wrong format problem.
//
// 	 ""      => "/"
// 	 "abc/"  => "/abc"
// 	 "/abc/" => "/abc"
// 	 "/abc"  => "/abc"
// 	 "/"     => "/"
//
// slash 處理斜線(路徑)
func slash(prefix string) string {
	prefix = strings.TrimSpace(prefix)
	if prefix == "" || prefix == "/" {
		return "/"
	}
	if prefix[0] != '/' {
		if prefix[len(prefix)-1] == '/' {
			return "/" + prefix[:len(prefix)-1]
		}
		return "/" + prefix
	}
	if prefix[len(prefix)-1] == '/' {
		return prefix[:len(prefix)-1]
	}
	return prefix
}

// join join the path.
// join 路徑
func join(prefix, suffix string) string {
	if prefix == "/" {
		return suffix
	}
	if suffix == "/" {
		return prefix
	}
	return prefix + suffix
}
