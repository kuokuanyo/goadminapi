package guard

import (
	"goadminapi/context"
	"goadminapi/modules/db"
	"goadminapi/modules/errors"
	"goadminapi/modules/service"
	"goadminapi/plugins/admin/modules/response"
	"goadminapi/plugins/admin/modules/table"
	"goadminapi/template"
	"goadminapi/template/types"
)

type Guard struct {
	services  service.List
	conn      db.Connection
	tableList table.GeneratorList
	navBtns   *types.Buttons
}

// 取得table(interface)、prefix
func (g *Guard) table(ctx *context.Context) (table.Table, string) {
	prefix := ctx.Query("__prefix")
	return g.tableList[prefix](ctx), prefix
}

// 將參數s、c、t設置至Guard(struct)後回傳
func New(s service.List, c db.Connection, t table.GeneratorList) *Guard {
	return &Guard{
		services:  s,
		conn:      c,
		tableList: t,
		// navBtns:   b,
	}
}

// 查詢url裡的參數(__prefix)
func (g *Guard) CheckPrefix(ctx *context.Context) {
	prefix := ctx.Query("__prefix")

	if _, ok := g.tableList[prefix]; !ok {
		if ctx.Headers("X-PJAX") == "" && ctx.Method() != "GET" {
			response.BadRequest(ctx, errors.Msg)
		} else {
			response.Alert(ctx, errors.Msg, errors.Msg, "table model not found", g.conn, g.navBtns,
				template.Missing404Page)
		}
		ctx.Abort()
		return
	}

	ctx.Next()
}
