package context

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
	Requests    []Path //Path包含url、method
	Handlers    HandlerMap // HandlerMap是 map[Path]Handlers，Handlers類別為[]Handler，Handler類別為func(ctx *Context)
	Middlewares Handlers
	Prefix      string

	Routers    RouterMap // RouterMap類別為map[string]Router，Router(struct)裡有methods、patten
	routeIndex int
	routeANY   bool
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

// 將參數prefix、middleware新增至RouterGroup(struct)
func (app *App) Group(prefix string, middleware ...Handler) *RouterGroup {
	return &RouterGroup{
		app:         app,
		Middlewares: append(app.Middlewares, middleware...),
		Prefix:      slash(prefix),
	}
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