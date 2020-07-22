package context

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"
)

// gin.Context(gin-gonic套件)的簡化版本
type Context struct {
	Request   *http.Request
	Response  *http.Response
	UserValue map[string]interface{}
	index     int8
	handlers  Handlers
}

type Handler func(ctx *Context) // 用於middleware
type Handlers []Handler
type HandlerMap map[Path]Handlers

type RouterMap map[string]Router

// Router包含方法模式
type Router struct {
	Methods []string
	Patten  string //url
}

type Path struct {
	URL    string
	Method string
}

type App struct {
	Requests    []Path     //Path包含url、method
	Handlers    HandlerMap // HandlerMap是 map[Path]Handlers，Handlers類別為[]Handler，Handler類別為func(ctx *Context)
	Middlewares Handlers
	Prefix      string

	Routers    RouterMap // RouterMap類別為map[string]Router，Router(struct)裡有methods、patten
	routeIndex int
	routeANY   bool
}

type RouterGroup struct {
	app         *App     //struct
	Middlewares Handlers //Handlers([]Handler)，Handler類別為 func(ctx *Context)
	Prefix      string
}

type NodeProcessor func(...Node)
type Node struct {
	Path     string
	Method   string
	Handlers []Handler
	Value    map[string]interface{}
}

// 取得app(struct)，空的
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

// 將參數設置至App.Routers(RouterMap)中，設定methods及patten(url)
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

// 將參數prefix、middleware新增至RouterGroup(struct)
func (app *App) Group(prefix string, middleware ...Handler) *RouterGroup {
	return &RouterGroup{
		app:         app,
		Middlewares: append(app.Middlewares, middleware...),
		Prefix:      slash(prefix),
	}
}

// 藉由參數key取得multipart/form-data中的值
func (ctx *Context) FormValue(key string) string {
	return ctx.Request.FormValue(key)
}

// 將參數(key、value)添加header中(Context.Response.Header)
func (ctx *Context) AddHeader(key, value string) {
	ctx.Response.Header.Add(key, value)
}

// 藉由參數key獲得Header
func (ctx *Context) Headers(key string) string {
	return ctx.Request.Header.Get(key)
}

// 設置cookie在response header Set-Cookie中
func (ctx *Context) SetCookie(cookie *http.Cookie) {
	if v := cookie.String(); v != "" {
		ctx.AddHeader("Set-Cookie", v)
	}
}

// 將參數添加至Content-Type
func (ctx *Context) SetContentType(contentType string) {
	ctx.AddHeader("Content-Type", contentType)
}

// 轉換成JSON存至Context.Response.body
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

// AppendReqAndResp stores the request info and handle into app.
// support the route parameter. The route parameter will be recognized as
// wildcard store into the RegUrl of Path struct. For example:
//
//         /user/:id      => /user/(.*)
//         /user/:id/info => /user/(.*?)/info
//
// The RegUrl will be used to recognize the incoming path and find
// the handler.
// 在RouterGroup.app(struct)中新增Requests([]Path)路徑及方法、接著在該url中新增參數handler(Handler...)
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

// POST等於在AppendReqAndResp(url, "post", handler)
func (g *RouterGroup) POST(url string, handler ...Handler) *RouterGroup {
	g.app.routeANY = false
	g.AppendReqAndResp(url, "post", handler)
	return g
}

// GET等於在AppendReqAndResp(url, "get", handler)
func (g *RouterGroup) GET(url string, handler ...Handler) *RouterGroup {
	g.app.routeANY = false
	g.AppendReqAndResp(url, "get", handler)
	return g
}

// DELETE等於在AppendReqAndResp(url, "delete", handler)
func (g *RouterGroup) DELETE(url string, handler ...Handler) *RouterGroup {
	g.app.routeANY = false
	g.AppendReqAndResp(url, "delete", handler)
	return g
}

// PUT等於在AppendReqAndResp(url, "put", handler)
func (g *RouterGroup) PUT(url string, handler ...Handler) *RouterGroup {
	g.app.routeANY = false
	g.AppendReqAndResp(url, "put", handler)
	return g
}

// OPTIONS等於在AppendReqAndResp(url, "options", handler)
func (g *RouterGroup) OPTIONS(url string, handler ...Handler) *RouterGroup {
	g.app.routeANY = false
	g.AppendReqAndResp(url, "options", handler)
	return g
}

// HEAD等於在AppendReqAndResp(url, "head", handler)
func (g *RouterGroup) HEAD(url string, handler ...Handler) *RouterGroup {
	g.app.routeANY = false
	g.AppendReqAndResp(url, "head", handler)
	return g
}

// ANY registers a route that matches all the HTTP methods.
// GET, POST, PUT, HEAD, OPTIONS, DELETE.
// 執行所有方法的AppendReqAndResp(url, 方法, handler)
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

// 將參數設置至App.Routers(RouterMap)中，設定methods及patten(url)
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
// 處理斜線(路徑)
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
// join路徑
func join(prefix, suffix string) string {
	if prefix == "/" {
		return suffix
	}
	if suffix == "/" {
		return prefix
	}
	return prefix + suffix
}