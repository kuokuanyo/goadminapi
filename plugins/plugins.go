package plugins

import (
	"goadminapi/context"
	"goadminapi/modules/db"
	"goadminapi/modules/service"
)

// Base(struct)也是Plugin(interface)
type Base struct {
	// context.App在context\context.go中
	App      *context.App
	Services service.List
	Conn     db.Connection
	//UI        *ui.Service
	PlugName  string
	URLPrefix string
}
