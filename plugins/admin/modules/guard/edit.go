package guard

import (
	"goadminapi/context"
	"goadminapi/modules/config"
	"goadminapi/modules/db"
	"goadminapi/plugins/admin/modules/response"
	"goadminapi/plugins/admin/modules/table"
	"goadminapi/template"
	"goadminapi/template/types"
	"regexp"
	"strings"

	tmpl "html/template"
)

// 判斷參數是否是info url
func isInfoUrl(s string) bool {
	reg, _ := regexp.Compile("(.*?)info/(.*?)$")
	sub := reg.FindStringSubmatch(s)
	return len(sub) > 2 && !strings.Contains(sub[2], "/")
}

func getAlert(msg string) tmpl.HTML {
	return template.Get(config.GetTheme()).Alert().Warning(msg)
}

func alert(ctx *context.Context, panel table.Table, msg string, conn db.Connection, btns *types.Buttons) {
	if ctx.WantJSON() {
		response.BadRequest(ctx, msg)
	} else {
		response.Alert(ctx, panel.GetInfo().Description, panel.GetInfo().Title, msg, conn, btns)
	}
}

// 將給定的數據(types.Page(struct))及參數寫入buf(struct)並回傳，最後輸出HTML
func alertWithTitleAndDesc(ctx *context.Context, title, desc, msg string, conn db.Connection, btns *types.Buttons) {
	response.Alert(ctx, desc, title, msg, conn, btns)
}
