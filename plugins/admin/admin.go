package admin

import (
	"goadminapi/plugins"
	"goadminapi/plugins/admin/controller"
	"goadminapi/plugins/admin/modules/guard"
	"goadminapi/plugins/admin/modules/table"
)

// Admin is a GoAdmin plugin.
type Admin struct {
	*plugins.Base
	// plugins\admin\modules\table\table.go
	// GeneratorList類別為map[string]Generator，Generator類別為func(ctx *context.Context) Table
	tableList table.GeneratorList
	// plugins\admin\modules\guard
	guardian *guard.Guard
	// plugins\admin\controller
	handler *controller.Handler
}
