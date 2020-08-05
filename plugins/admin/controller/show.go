package controller

import (
	"bytes"
	"crypto/md5"
	"fmt"
	"goadminapi/context"
	"goadminapi/plugins/admin/modules/parameter"
	"goadminapi/plugins/admin/modules/table"
	"mime"
	"net/http"
	"path"
	"strconv"
	"strings"

	"goadminapi/template"
)

// ShowInfo 前端資訊介面
func (h *Handler) ShowInfo(ctx *context.Context) {
	prefix := ctx.Query("__prefix")

	// table 先透過參數prefix取得Table(interface)，接著判斷條件後將[]context.Node加入至Handler.operations後回傳
	panel := h.table(prefix, ctx)

	// 如果只有編輯表單的權限，導向編輯表單的頁面
	if panel.GetOnlyUpdateForm() {
		ctx.Redirect(h.routePathWithPrefix("show_edit", prefix))
		return
	}
	// 如果只有新增表單的權限，導向新增表單的頁面
	if panel.GetOnlyNewForm() {
		ctx.Redirect(h.routePathWithPrefix("show_new", prefix))
		return
	}
	// 如果只能取得細節的權限，導向細節的頁面
	if panel.GetOnlyDetail() {
		ctx.Redirect(h.routePathWithPrefix("detail", prefix))
		return
	}

	// 取得頁面size、資料排列方式、選擇欄位...等資訊後設置至Parameters(struct)並回傳
	params := parameter.GetParam(ctx.Request.URL, panel.GetInfo().DefaultPageSize, panel.GetInfo().SortField,
		panel.GetInfo().GetSort())

	buf := h.showTable(ctx, prefix, params, panel)
	ctx.HTML(http.StatusOK, buf.String())
}

// 處理前端檔案
func (h *Handler) Assets(ctx *context.Context) {
	// URLRemovePrefix將URL的前綴(ex:/admin)去除
	filepath := h.config.URLRemovePrefix(ctx.Path())

	// aTemplate判斷templateMap(map[string]Template)的key鍵是否參數globalCfg.Theme，有則回傳Template(interface)
	data, err := aTemplate().GetAsset(filepath)
	if err != nil {
		// 如果沒有設置js、css檔案，則會執行
		data, err = template.GetAsset(filepath)
		if err != nil {
			ctx.Write(http.StatusNotFound, map[string]string{}, "")
			panic("asset err")
		}
	}

	var contentType = mime.TypeByExtension(path.Ext(filepath))
	if contentType == "" {
		contentType = "application/octet-stream"
	}

	etag := fmt.Sprintf("%x", md5.Sum(data))

	if match := ctx.Headers("If-None-Match"); match != "" {
		if strings.Contains(match, etag) {
			ctx.SetStatusCode(http.StatusNotModified)
			return
		}
	}

	// 將code, headers and body(參數)設置至在Context.Response中
	ctx.DataWithHeaders(http.StatusOK, map[string]string{
		"Content-Type":   contentType,
		"Cache-Control":  "max-age=2592000",
		"Content-Length": strconv.Itoa(len(data)),
		"ETag":           etag,
	}, data)
}

func (h *Handler) showTableData(ctx *context.Context, prefix string, params parameter.Parameters,
	panel table.Table, urlNamePrefix string) (table.Table, table.PanelInfo, []string, error) {
	// 先設置table(interface)
	if panel == nil {
		panel = h.table(prefix, ctx)
	}

	panelInfo, err := panel.GetData(params.WithIsAll(false))
}

func (h *Handler) showTable(ctx *context.Context, prefix string, params parameter.Parameters, panel table.Table) *bytes.Buffer {
	panel, panelInfo, urls, err := h.showTableData(ctx, prefix, params, panel, "")
}
