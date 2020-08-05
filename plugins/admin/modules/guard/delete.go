package guard

import (
	"goadminapi/context"
	"goadminapi/modules/errors"
	"goadminapi/plugins/admin/modules/table"
)

type DeleteParam struct {
	Panel  table.Table
	Id     string
	Prefix string
}

// Delete 取得url的id值後將值設置至Context.UserValue[delete_param]
func (g *Guard) Delete(ctx *context.Context) {
	// 取得table(interface)、prefix
	panel, prefix := g.table(ctx)
	if !panel.GetDeletable() {
		alert(ctx, panel, errors.OperationNotAllow, g.conn, g.navBtns)
		ctx.Abort()
		return
	}

	id := ctx.FormValue("id")
	if id == "" {
		alert(ctx, panel, errors.WrongID, g.conn, g.navBtns)
		ctx.Abort()
		return
	}

	ctx.SetUserValue("delete_param", &DeleteParam{
		Panel:  panel,
		Id:     id,
		Prefix: prefix,
	})
	ctx.Next()
}

// GetDeleteParam 取得Context.UserValue[delete_param]的值並轉換成DeleteParam(struct)
func GetDeleteParam(ctx *context.Context) *DeleteParam {
	return ctx.UserValue["delete_param"].(*DeleteParam)
}
