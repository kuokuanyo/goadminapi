package controller

import (
	c "goadminapi/modules/config"
	"regexp"

	"goadminapi/context"
	"goadminapi/modules/db"
	"goadminapi/modules/service"
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
	//navButtons    *types.Buttons
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

func isNewUrl(s string, p string) bool {
	reg, _ := regexp.Compile("(.*?)info/" + p + "/new")

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