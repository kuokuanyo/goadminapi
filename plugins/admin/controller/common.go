package controller

import (
	c "goadminapi/modules/config"

	"goadminapi/context"
	"goadminapi/modules/db"
	"goadminapi/modules/service"
	"goadminapi/plugins/admin/modules/table"
	"sync"
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
