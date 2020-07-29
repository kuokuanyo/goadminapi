package controller

import (
	"goadminapi/context"
	"goadminapi/modules/auth"
	"goadminapi/modules/db"
	"goadminapi/modules/menu"
	"goadminapi/plugins/admin/models"
	"goadminapi/plugins/admin/modules/guard"
	"goadminapi/plugins/admin/modules/response"
)

// 將MenuNewParam(struct)值新增至資料表(MenuModel.Base.TableName(menu))中
// 如果multipart/form-data有設定roles[]值，檢查條件後將資料加入role_menu資料表
// *********************前端頁面函式還沒處理************************
func (h *Handler) NewMenu(ctx *context.Context) {
	param := guard.GetMenuNewParam(ctx)
	if param.HasAlert() {
		//**************函式還沒寫***********************
		//h.getMenuInfoPanel(ctx, param.Alert)
		ctx.AddHeader("Content-Type", "text/html; charset=utf-8")
		ctx.AddHeader("X-PJAX-Url", h.routePath("menu"))
		return
	}
	// 透過參數ctx回傳目前登入的用戶(Context.UserValue["user"])並轉換成UserModel
	user := auth.Auth(ctx)

	// SetConn將參數h.conn設置至MenuModel.Base.Conn
	// New將參數值新增至資料表(MenuModel.Base.TableName(goadmin_menu))中，最後將參數值都設置在MenuModel中
	menuModel, createErr := models.Menu().SetConn(h.conn).
		// GetGlobalMenu回傳參數user(struct)的Menu(設置menuList、menuOption、MaxOrder)
		New(param.Title, param.Icon, param.Uri, param.Header, param.ParentId, (menu.GetGlobalMenu(user, h.conn)).MaxOrder+1)
	if db.CheckError(createErr, db.INSERT) {
		//**************函式還沒寫***********************
		// h.showNewMenu(ctx, createErr)
		return
	}

	// AddRole 檢查role_menu資料表裡是否有符合role_id = 參數roleId、menu_id = MenuModel.Id資料
	// 如果沒有接則將參數roleId(role_id)與MenuModel.Id(menu_id)加入role_menu資料表
	for _, roleId := range param.Roles {
		_, addRoleErr := menuModel.AddRole(roleId)
		if db.CheckError(addRoleErr, db.INSERT) {
			//**************函式還沒寫***********************
			// h.showNewMenu(ctx, addRoleErr)
			return
		}
	}

	// GetGlobalMenu回傳參數user(struct)的Menu(設置menuList、menuOption、MaxOrder)
	menu.GetGlobalMenu(user, h.conn).AddMaxOrder()

	//**************函式還沒寫***********************
	// h.getMenuInfoPanel(ctx, "")
	ctx.AddHeader("Content-Type", "text/html; charset=utf-8")
	ctx.AddHeader("X-PJAX-Url", h.routePath("menu"))
}

// EditMenu edit the menu of given id.
func (h *Handler) EditMenu(ctx *context.Context) {
	param := guard.GetMenuEditParam(ctx)
	if param.HasAlert() {
		//**************函式還沒寫***********************
		// h.getMenuInfoPanel(ctx, param.Alert)
		ctx.AddHeader("Content-Type", "text/html; charset=utf-8")
		ctx.AddHeader("X-PJAX-Url", h.routePath("menu"))
		return
	}

	menuModel := models.MenuWithId(param.Id).SetConn(h.conn)

	// 先刪除所有角色再將新角色新增至資料表
	deleteRolesErr := menuModel.DeleteRoles()
	if db.CheckError(deleteRolesErr, db.DELETE) {
		//**************函式還沒寫***********************
		// formInfo, _ := h.table("menu", ctx).GetDataWithId(parameter.BaseParam().WithPKs(param.Id))
		// h.showEditMenu(ctx, formInfo, deleteRolesErr)
		ctx.AddHeader("X-PJAX-Url", h.routePath("menu"))
		return
	}
	for _, roleId := range param.Roles {
		_, addRoleErr := menuModel.AddRole(roleId)
		if db.CheckError(addRoleErr, db.INSERT) {
			//**************函式還沒寫***********************
			// formInfo, _ := h.table("menu", ctx).GetDataWithId(parameter.BaseParam().WithPKs(param.Id))
			// h.showEditMenu(ctx, formInfo, addRoleErr)
			ctx.AddHeader("X-PJAX-Url", h.routePath("menu"))
			return
		}
	}

	// 更新資料表資料
	_, updateErr := menuModel.Update(param.Title, param.Icon, param.Uri, param.Header, param.ParentId)
	if db.CheckError(updateErr, db.UPDATE) {
		//**************函式還沒寫***********************
		// formInfo, _ := h.table("menu", ctx).GetDataWithId(parameter.BaseParam().WithPKs(param.Id))
		// h.showEditMenu(ctx, formInfo, updateErr)
		ctx.AddHeader("X-PJAX-Url", h.routePath("menu"))
		return
	}

	// h.getMenuInfoPanel(ctx, "")
	ctx.AddHeader("Content-Type", "text/html; charset=utf-8")
	ctx.AddHeader("X-PJAX-Url", h.routePath("menu"))
}



// DeleteMenu delete the menu of given id.
func (h *Handler) DeleteMenu(ctx *context.Context) {
	models.MenuWithId(guard.GetMenuDeleteParam(ctx).Id).SetConn(h.conn).Delete()
	response.OkWithMsg(ctx, "刪除成功")
}
