package controller

import (
	c "goadminapi/modules/config"

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

// 判斷templateMap(map[string]Template)的key鍵是否參數globalCfg.Theme，有則回傳Template(interface)
func aTemplate() template.Template {
	// 判斷templateMap(map[string]Template)的key鍵是否參數globalCfg.Theme，有則回傳Template(interface)
	// GetTheme return globalCfg.Theme
	return template.Get(c.GetTheme())
}
