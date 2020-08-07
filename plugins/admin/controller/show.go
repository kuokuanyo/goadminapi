package controller

import (
	"bytes"
	"crypto/md5"
	"fmt"
	"goadminapi/context"
	"goadminapi/modules/auth"
	"goadminapi/modules/errors"
	"goadminapi/plugins/admin/modules"
	"goadminapi/plugins/admin/modules/parameter"
	"goadminapi/plugins/admin/modules/table"
	"goadminapi/template/types"
	"goadminapi/template/types/action"
	"mime"
	"net/http"
	"net/url"
	"path"
	"strconv"
	"strings"

	"goadminapi/html"
	"goadminapi/template"
	template2 "html/template"

	"goadminapi/template/icon"
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

// showTableData 取得table(interface)、PanelInfo(有主題、描述名稱、可以篩選條件的欄位、選擇顯示的欄位、分頁、[]TheadItem(欄位資訊)等資訊)
// 以及前端介面會使用到的url路徑
func (h *Handler) showTableData(ctx *context.Context, prefix string, params parameter.Parameters,
	panel table.Table, urlNamePrefix string) (table.Table, table.PanelInfo, []string, error) {
	// 先設置table(interface)
	if panel == nil {
		panel = h.table(prefix, ctx)
	}

	// WithIsAll 添加Parameters.Fields["__is_all"]
	// GetData 透過參數處理後取得前端介面顯示資料，將值設置至PanelInfo(struct)
	// PanelInfo裡的資訊有主題、描述名稱、可以篩選條件的欄位、選擇顯示的欄位、分頁、[]TheadItem(欄位資訊)等資訊
	panelInfo, err := panel.GetData(params.WithIsAll(false))
	if err != nil {
		return panel, panelInfo, nil, err
	}

	// DeleteIsAll 刪除Parameters.Fields["__is_all"]
	// 處理url後的參數
	paramStr := params.DeleteIsAll().GetRouteParamStr()

	// AorEmpty判斷第一個(condition)參數，如果true則回傳第二個參數，否則回傳""
	// ex: /admin/info/manager/edit
	editUrl := modules.AorEmpty(!panel.GetInfo().IsHideEditButton, h.routePathWithPrefix(urlNamePrefix+"show_edit", prefix)+paramStr)
	// ex: /admin/info/manager/new
	newUrl := modules.AorEmpty(!panel.GetInfo().IsHideNewButton, h.routePathWithPrefix(urlNamePrefix+"show_new", prefix)+paramStr)
	// ex: /admin/delete/manager
	deleteUrl := modules.AorEmpty(!panel.GetInfo().IsHideDeleteButton, h.routePathWithPrefix(urlNamePrefix+"delete", prefix))
	// ex: /admin/detail/manager
	detailUrl := modules.AorEmpty(!panel.GetInfo().IsHideDetailButton, h.routePathWithPrefix(urlNamePrefix+"detail", prefix)+paramStr)
	// ex: /admin/info/manager
	infoUrl := h.routePathWithPrefix(urlNamePrefix+"info", prefix)

	// 取得權限、角色、可用menu
	user := auth.Auth(ctx)

	// 藉由參數檢查權限，如果有權限回傳第一個參數(path)，反之回傳""
	editUrl = user.GetCheckPermissionByUrlMethod(editUrl, h.route(urlNamePrefix+"show_edit").Method())
	newUrl = user.GetCheckPermissionByUrlMethod(newUrl, h.route(urlNamePrefix+"show_new").Method())
	deleteUrl = user.GetCheckPermissionByUrlMethod(deleteUrl, h.route(urlNamePrefix+"delete").Method())
	detailUrl = user.GetCheckPermissionByUrlMethod(detailUrl, h.route(urlNamePrefix+"detail").Method())

	return panel, panelInfo, []string{editUrl, newUrl, deleteUrl, detailUrl, infoUrl}, nil
}

func (h *Handler) showTable(ctx *context.Context, prefix string, params parameter.Parameters, panel table.Table) *bytes.Buffer {
	// showTableData 取得table(interface)、PanelInfo(有主題、描述名稱、可以篩選條件的欄位、選擇顯示的欄位、分頁、[]TheadItem(欄位資訊)等資訊)
	// 以及前端介面會使用到的url路徑
	panel, panelInfo, urls, err := h.showTableData(ctx, prefix, params, panel, "")
	if err != nil {
		// 將參數設置至ExecuteParam(struct)，接著將給定的數據寫入buf(struct)並回傳
		return h.Execute(ctx, auth.Auth(ctx), types.Panel{
			Content: aAlert().SetTitle(errors.MsgWithIcon).
				SetTheme("warning").
				SetContent(template2.HTML(err.Error())).
				GetContent(),
			Description: template2.HTML(errors.Msg),
			Title:       template2.HTML(errors.Msg),
		}, params.Animation)
	}

	editUrl, newUrl, deleteUrl, detailUrl, infoUrl := urls[0], urls[1], urls[2], urls[3], urls[4]

	// 取得權限、角色、可用menu
	user := auth.Auth(ctx)

	var (
		body       template2.HTML
		dataTable  types.DataTableAttribute
		info       = panel.GetInfo()
		actionBtns = info.Action
		actionJs   template2.JS
		allBtns    = make(types.Buttons, 0)
	)

	// ------------一般info.Buttons為空---------------
	for _, b := range info.Buttons {
		if b.URL() == "" || b.METHOD() == "" || user.CheckPermissionByUrlMethod(b.URL(), b.METHOD(), url.Values{}) {
			allBtns = append(allBtns, b)
		}
	}

	// 取得HTML及JSON
	btns, btnsJs := allBtns.Content()
	allActionBtns := make(types.Buttons, 0)

	// ------------一般info.ActionButtons為空---------------
	for _, b := range info.ActionButtons {
		if b.URL() == "" || b.METHOD() == "" || user.CheckPermissionByUrlMethod(b.URL(), b.METHOD(), url.Values{}) {
			allActionBtns = append(allActionBtns, b)
		}
	}

	// ---------如果上面為空，因此這裡不執行-----------------
	if actionBtns == template.HTML("") && len(allActionBtns) > 0 {
		ext := template.HTML("")
		if deleteUrl != "" {
			ext = html.LiEl().SetClass("divider").Get()
			allActionBtns = append([]types.Button{types.GetActionButton(template.HTML("delete"),
				types.NewDefaultAction(`data-id='{{.Id}}' style="cursor: pointer;"`,
					ext, "", ""), "grid-row-delete")}, allActionBtns...)
		}
		ext = template.HTML("")
		if detailUrl != "" {
			if editUrl == "" && deleteUrl == "" {
				ext = html.LiEl().SetClass("divider").Get()
			}
			allActionBtns = append([]types.Button{types.GetActionButton(template.HTML("detail"),
				action.Jump(detailUrl+"&"+"__detail_pk"+"={{.Id}}", ext))}, allActionBtns...)
		}
		if editUrl != "" {
			if detailUrl == "" && deleteUrl == "" {
				ext = html.LiEl().SetClass("divider").Get()
			}
			allActionBtns = append([]types.Button{types.GetActionButton(template.HTML("edit"),
				action.Jump(editUrl+"&"+"__edit_pk"+"={{.Id}}", ext))}, allActionBtns...)
		}

		var content template2.HTML
		content, actionJs = allActionBtns.Content()

		actionBtns = html.Div(
			html.A(icon.Icon(icon.EllipsisV),
				html.M{"color": "#676565"},
				html.M{"class": "dropdown-toggle", "href": "#", "data-toggle": "dropdown"},
			)+html.Ul(content,
				html.M{"min-width": "20px !important", "left": "-32px", "overflow": "hidden"},
				html.M{"class": "dropdown-menu", "role": "menu", "aria-labelledby": "dLabel"}),
			html.M{"text-align": "center"}, html.M{"class": "dropdown"})
	}

	// --------------------一般都為false-------------------
	if info.TabGroups.Valid() {
		dataTable = aDataTable().
		SetThead(panelInfo.Thead).
		SetDeleteUrl(deleteUrl).
		SetNewUrl(newUrl)

		var (
			tabsHtml    = make([]map[string]template2.HTML, len(info.TabHeaders))
			infoListArr = panelInfo.InfoList.GroupBy(info.TabGroups)
			theadArr    = panelInfo.Thead.GroupBy(info.TabGroups)
		)

		for key, header := range info.TabHeaders {
			tabsHtml[key] = map[string]template2.HTML{
				"title": template2.HTML(header),
				"content": aDataTable().
					SetInfoList(infoListArr[key]).
					SetInfoUrl(infoUrl).
					SetButtons(btns).
					SetActionJs(btnsJs + actionJs).
					SetHasFilter(len(panelInfo.FilterFormData) > 0).
					SetAction(actionBtns).
					SetIsTab(key != 0).
					SetPrimaryKey(panel.GetPrimaryKey().Name).
					SetThead(theadArr[key]).
					SetHideRowSelector(info.IsHideRowSelector).
					SetLayout(info.TableLayout).
					SetNewUrl(newUrl).
					SetSortUrl(params.GetFixedParamStrWithoutSort()).
					SetEditUrl(editUrl).
					SetDetailUrl(detailUrl).
					SetDeleteUrl(deleteUrl).
					GetContent(),
			}
		}
		body = aTab().SetData(tabsHtml).GetContent()
	}
}
