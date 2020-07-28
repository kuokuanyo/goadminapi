package page

import (
	"bytes"
	"goadminapi/context"
	"goadminapi/modules/config"
	"goadminapi/modules/db"
	"goadminapi/modules/logger"
	"goadminapi/plugins/admin/models"
	"goadminapi/template"
	"goadminapi/template/types"

	"goadminapi/modules/menu"
)

// 設置並回傳面板內容
func SetPageContent(ctx *context.Context, user models.UserModel, c func(ctx interface{}) (types.Panel, error), conn db.Connection) {
	panel, err := c(ctx)
	if err != nil {
		logger.Error("SetPageContent", err)
		panel = template.WarningPanel(err.Error())
	}

	tmpl, tmplName := template.Get(config.GetTheme()).GetTemplate(ctx.IsPjax())
	ctx.AddHeader("Content-Type", "text/html; charset=utf-8")

	buf := new(bytes.Buffer)
	err = tmpl.ExecuteTemplate(buf, tmplName, types.NewPage(types.NewPageParam{
		User:         user,
		Menu:         menu.GetGlobalMenu(user, conn).SetActiveClass(config.URLRemovePrefix(ctx.Path())),
		Panel:        panel.GetContent(config.IsProductionEnvironment()),
		Assets:       template.GetComponentAssetImportHTML(),
		TmplHeadHTML: template.Default().GetHeadHTML(),
		TmplFootJS:   template.Default().GetFootJS(),
	}))
}
