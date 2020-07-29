package guard

import "goadminapi/context"

// MenuDeleteParam 刪除menu只需要id
type MenuDeleteParam struct {
	Id string
}

func (g *Guard) MenuDelete(ctx *context.Context) {
	// 查詢url中參數id的值
	id := ctx.Query("id")
	if id == "" {
		alertWithTitleAndDesc(ctx, "Menu", "menu", "wrong ID", g.conn, g.navBtns)
		ctx.Abort()
		return
	}
	// 將參數設置至Context.UserValue[delete_menu_param]
	ctx.SetUserValue("delete_menu_param", &MenuDeleteParam{
		Id: id,
	})
	ctx.Next()
}

// 取得Context.UserValue[delete_menu_param]的值並轉換成MenuDeleteParam(struct)
func GetMenuDeleteParam(ctx *context.Context) *MenuDeleteParam {
	// deleteMenuParamKey = delete_menu_param
	// 將Context.UserValue(map[string]interface{})[delete_menu_param]的值轉換成MenuDeleteParam(struct)類別
	return ctx.UserValue["delete_menu_param"].(*MenuDeleteParam)
}