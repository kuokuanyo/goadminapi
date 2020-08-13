package controller

import (
	"goadminapi/context"
	"goadminapi/modules/auth"
	"goadminapi/modules/db"
	"goadminapi/modules/menu"
	"goadminapi/plugins/admin/models"
	"goadminapi/plugins/admin/modules/guard"
	"goadminapi/plugins/admin/modules/parameter"
	"goadminapi/plugins/admin/modules/response"
	"goadminapi/plugins/admin/modules/table"
	"goadminapi/template"
	"goadminapi/template/types"
	template2 "html/template"
)

// ShowMenu 前端菜單介面顯示HTML
func (h *Handler) ShowMenu(ctx *context.Context) {
	h.getMenuInfoPanel(ctx, "")
}

// ShowNewMenu 前端顯示新建菜單HTML語法
func (h *Handler) ShowNewMenu(ctx *context.Context) {
	h.showNewMenu(ctx, nil)
}

// ShowNewMenu 前端顯示新建菜單HTML語法
func (h *Handler) showNewMenu(ctx *context.Context, err error) {
	// table 透過參數prefix取得Table(interface)，在generators.go設定的資訊
	panel := h.table("menu", ctx)

	// GetForm 取得table中設置的FormPanel(struct)，在generators.go設置的資訊
	formInfo := panel.GetNewForm()

	user := auth.Auth(ctx)

	var alert template2.HTML
	if err != nil {
		alert = aAlert().Warning(err.Error())
	}

	h.HTML(ctx, user, types.Panel{
		Content: alert + formContent(aForm().
			SetTitle("New").
			SetContent(formInfo.FieldList).
			SetTabContents(formInfo.GroupFieldList).
			SetTabHeaders(formInfo.GroupFieldHeaders).
			SetPrefix(h.config.PrefixFixSlash()).
			SetPrimaryKey(panel.GetPrimaryKey().Name).
			SetUrl(h.routePath("menu_new")).
			SetHiddenFields(map[string]string{
				"__token_":    h.authSrv().AddToken(),
				"__previous_": h.routePath("menu"),
			}).
			// formFooter處理後回傳繼續新增、保存、重製....等HTML語法
			SetOperationFooter(formFooter("new", false, true, false)),
			false, ctx.Query("__iframe") == "true", false, ""),
		Description: template2.HTML(panel.GetForm().Description),
		Title:       template2.HTML(panel.GetForm().Title),
	})
}

// ShowEditMenu 前端菜單編輯介面HTML語法
func (h *Handler) ShowEditMenu(ctx *context.Context) {
	if ctx.Query("id") == "" {
		h.getMenuInfoPanel(ctx, template.Get(h.config.Theme).Alert().Warning("wrong id"))

		ctx.AddHeader("Content-Type", "text/html; charset=utf-8")
		ctx.AddHeader("X-PJAX-Url", h.routePath("menu"))
		return
	}

	// table 透過參數prefix取得Table(interface)，在generators.go設定的資訊
	model := h.table("menu", ctx)

	// BaseParam 設置(頁數及頁數Size)至Parameters(struct)
	formInfo, err := model.GetDataWithId(parameter.BaseParam().WithPKs(ctx.Query("id")))

	user := auth.Auth(ctx)

	if err != nil {
		h.HTML(ctx, user, types.Panel{
			Content:     aAlert().Warning(err.Error()),
			Description: template2.HTML(model.GetForm().Description),
			Title:       template2.HTML(model.GetForm().Title),
		})
		return
	}

	// 將編輯介面的HTML語法匯出
	h.showEditMenu(ctx, formInfo, nil)
}

// showEditMenu 前端菜單編輯介面HTML語法
func (h *Handler) showEditMenu(ctx *context.Context, formInfo table.FormInfo, err error) {
	var alert template2.HTML
	if err != nil {
		alert = aAlert().Warning(err.Error())
	}

	h.HTML(ctx, auth.Auth(ctx), types.Panel{
		Content: alert + formContent(aForm().
			SetContent(formInfo.FieldList).                           // 表單欄位資訊
			SetTabContents(formInfo.GroupFieldList).                  // 空
			SetTabHeaders(formInfo.GroupFieldHeaders).                // 空
			SetPrefix(h.config.PrefixFixSlash()).                     // /admin
			SetPrimaryKey(h.table("menu", ctx).GetPrimaryKey().Name). // id
			SetUrl(h.routePath("menu_edit")).
			// formFooter處理後回傳繼續新增、繼續編輯、保存、重製....等HTML語法
			SetOperationFooter(formFooter("edit", true, true, false)).
			SetHiddenFields(map[string]string{
				"__token_":    h.authSrv().AddToken(),
				"__previous_": h.routePath("menu"),
			}), false, ctx.Query("__iframe") == "true", false, ""),
		Description: template2.HTML(formInfo.Description),
		Title:       template2.HTML(formInfo.Title),
	})

	return
}

// 將MenuNewParam(struct)值新增至資料表(MenuModel.Base.TableName(menu))中
// 如果multipart/form-data有設定roles[]值，檢查條件後將資料加入role_menu資料表
func (h *Handler) NewMenu(ctx *context.Context) {
	param := guard.GetMenuNewParam(ctx)
	if param.HasAlert() {
		h.getMenuInfoPanel(ctx, param.Alert)
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
		h.showNewMenu(ctx, createErr)
		return
	}

	// AddRole 檢查role_menu資料表裡是否有符合role_id = 參數roleId、menu_id = MenuModel.Id資料
	// 如果沒有接則將參數roleId(role_id)與MenuModel.Id(menu_id)加入role_menu資料表
	for _, roleId := range param.Roles {
		_, addRoleErr := menuModel.AddRole(roleId)
		if db.CheckError(addRoleErr, db.INSERT) {
			h.showNewMenu(ctx, addRoleErr)
			return
		}
	}

	// GetGlobalMenu回傳參數user(struct)的Menu(設置menuList、menuOption、MaxOrder)
	menu.GetGlobalMenu(user, h.conn).AddMaxOrder()

	h.getMenuInfoPanel(ctx, "")
	ctx.AddHeader("Content-Type", "text/html; charset=utf-8")

	ctx.AddHeader("X-PJAX-Url", h.routePath("menu"))
}

// EditMenu edit the menu of given id.
func (h *Handler) EditMenu(ctx *context.Context) {
	param := guard.GetMenuEditParam(ctx)
	if param.HasAlert() {
		h.getMenuInfoPanel(ctx, param.Alert)
		ctx.AddHeader("Content-Type", "text/html; charset=utf-8")
		ctx.AddHeader("X-PJAX-Url", h.routePath("menu"))
		return
	}

	menuModel := models.MenuWithId(param.Id).SetConn(h.conn)

	// 先刪除所有角色再將新角色新增至資料表
	deleteRolesErr := menuModel.DeleteRoles()
	if db.CheckError(deleteRolesErr, db.DELETE) {
		formInfo, _ := h.table("menu", ctx).GetDataWithId(parameter.BaseParam().WithPKs(param.Id))
		h.showEditMenu(ctx, formInfo, deleteRolesErr)
		ctx.AddHeader("X-PJAX-Url", h.routePath("menu"))
		return
	}
	for _, roleId := range param.Roles {
		_, addRoleErr := menuModel.AddRole(roleId)
		if db.CheckError(addRoleErr, db.INSERT) {
			formInfo, _ := h.table("menu", ctx).GetDataWithId(parameter.BaseParam().WithPKs(param.Id))
			h.showEditMenu(ctx, formInfo, addRoleErr)
			ctx.AddHeader("X-PJAX-Url", h.routePath("menu"))
			return
		}
	}

	// 更新資料表資料
	_, updateErr := menuModel.Update(param.Title, param.Icon, param.Uri, param.Header, param.ParentId)
	if db.CheckError(updateErr, db.UPDATE) {
		formInfo, _ := h.table("menu", ctx).GetDataWithId(parameter.BaseParam().WithPKs(param.Id))
		h.showEditMenu(ctx, formInfo, updateErr)
		ctx.AddHeader("X-PJAX-Url", h.routePath("menu"))
		return
	}

	h.getMenuInfoPanel(ctx, "")
	ctx.AddHeader("Content-Type", "text/html; charset=utf-8")
	ctx.AddHeader("X-PJAX-Url", h.routePath("menu"))
}

// DeleteMenu delete the menu of given id.
func (h *Handler) DeleteMenu(ctx *context.Context) {
	models.MenuWithId(guard.GetMenuDeleteParam(ctx).Id).SetConn(h.conn).Delete()
	response.OkWithMsg(ctx, "刪除成功")
}

// getMenuInfoPanel 前端菜單顯示介面HTML處理
func (h *Handler) getMenuInfoPanel(ctx *context.Context, alert template2.HTML) {
	user := auth.Auth(ctx)

	tree := aTree().
		SetTree((menu.GetGlobalMenu(user, h.conn)).List). // 回傳菜單([]menu.Item)
		SetEditUrl(h.routePath("menu_edit_show")).        // /admin/menu/edit/show
		SetUrlPrefix(h.config.Prefix()).                  // /admin
		SetDeleteUrl(h.routePath("menu_delete")).         // /admin/menu/delete
		SetOrderUrl(h.routePath("menu_order")).           // /admin/menu/order
		GetContent()

	header := aTree().GetTreeHeader()

	box := aBox().SetHeader(header).SetBody(tree).GetContent()

	col1 := aCol().SetSize(types.SizeMD(6)).SetContent(box).GetContent()

	// table 透過參數prefix取得Table(interface)，在generators.go設定的資訊
	formInfo := h.table("menu", ctx).GetNewForm()

	newForm := menuFormContent(aForm().
		SetPrefix(h.config.PrefixFixSlash()).                     // /admin
		SetUrl(h.routePath("menu_new")).                          // /admin/menu/new
		SetPrimaryKey(h.table("menu", ctx).GetPrimaryKey().Name). // id
		SetHiddenFields(map[string]string{
			"__token_":    h.authSrv().AddToken(),
			"__previous_": h.routePath("menu"),
		}).
		// formFooter處理後回傳繼續保存、重製....等HTML語法
		SetOperationFooter(formFooter("menu", false, false, false)).
		SetTitle("New Menu").
		SetContent(formInfo.FieldList).
		SetTabContents(formInfo.GroupFieldList).   // 空
		SetTabHeaders(formInfo.GroupFieldHeaders)) // 空

	col2 := aCol().SetSize(types.SizeMD(6)).SetContent(newForm).GetContent()

	row := aRow().SetContent(col1 + col2).GetContent()

	h.HTML(ctx, user, types.Panel{
		Content:     alert + row,
		Description: "Menu",
		Title:       "Menu Manage",
	})
}
