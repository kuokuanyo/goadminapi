package controller

import (
	"fmt"
	"goadminapi/context"
	"goadminapi/modules/auth"
	"goadminapi/modules/file"
	"goadminapi/plugins/admin/modules"
	"goadminapi/plugins/admin/modules/guard"
	"goadminapi/plugins/admin/modules/parameter"
	"goadminapi/plugins/admin/modules/response"
	"goadminapi/template/types"
	"goadminapi/template/types/form"
	template2 "html/template"
	"net/http"
	"net/url"
)

// ShowForm 將選項及預設值設置至FormFields後處理前端表單顯示介面並匯出HTML
func (h *Handler) ShowForm(ctx *context.Context) {
	param := guard.GetShowFormParam(ctx)

	// param.Param的資訊有最大顯示資料數、排列順序、依照什麼欄位排序、選擇顯示欄位....等
	h.showForm(ctx, "", param.Prefix, param.Param, false)
}

// showForm 將選項及預設值設置至FormFields後處理前端表單顯示介面並匯出HTML
func (h *Handler) showForm(ctx *context.Context, alert template2.HTML, prefix string, param parameter.Parameters, isEdit bool, animation ...bool) {
	// table 先透過參數prefix取得Table(interface)，接著判斷條件後將[]context.Node加入至Handler.operations後回傳
	panel := h.table(prefix, ctx)

	// 透過參數ctx回傳目前登入的用戶(Context.UserValue["user"])並轉換成UserModel
	user := auth.Auth(ctx)

	// GetRouteParamStr 處理url後(?...)的部分(頁面設置、排序方式....等)
	// ex: ?__edit_pk=26&__page=1&__pageSize=10&__pk=26&__sort=id&__sort_type=desc
	paramStr := param.GetRouteParamStr()

	// 新增資料頁面
	newUrl := modules.AorEmpty(panel.GetCanAdd(), h.routePathWithPrefix("show_new", prefix)+paramStr)

	footerKind := "edit"
	// 如果沒有新增資料權限，則為"edit_only"
	if newUrl == "" || !user.CheckPermissionByUrlMethod(newUrl, h.route("show_new").Method(), url.Values{}) {
		footerKind = "edit_only"
	}

	// formInfo為所有欄位資訊(包含資料數值)
	// GetDataWithId 透過id取得資料，並且將選項、預設值...等資訊設置至FormFields
	formInfo, err := panel.GetDataWithId(param)

	// DeletePK 刪除Parameters.Fields[__pk]
	// ex:/admin/info/manager/edit?__edit_pk=26&__page=1&__pageSize=10&__sort=id&__sort_type=desc(編輯頁面)
	showEditUrl := h.routePathWithPrefix("show_edit", prefix) + param.DeletePK().GetRouteParamStr()

	if err != nil {
		h.HTML(ctx, user, types.Panel{
			Content:     aAlert().Warning(err.Error()),
			Description: template2.HTML(panel.GetForm().Description),
			Title:       template2.HTML(panel.GetForm().Title),
		}, alert == "" || ((len(animation) > 0) && animation[0]))

		if isEdit {
			ctx.AddHeader("X-PJAX-Url", showEditUrl)
		}
		return
	}

	// 顯示所有資料頁面
	// DeleteField刪除Parameters.Fields[參數(__edit_pk)]
	infoUrl := h.routePathWithPrefix("info", prefix) + param.DeleteField("__edit_pk").GetRouteParamStr()
	// 編輯功能(post)
	editUrl := h.routePathWithPrefix("edit", prefix)

	referer := ctx.Headers("Referer")
	if referer != "" && !isInfoUrl(referer) && !isEditUrl(referer, ctx.Query("__prefix")) {
		infoUrl = referer
	}

	// 所有表單欄位資訊
	f := panel.GetForm()

	// 如果url沒設置_iframe則為空值
	isNotIframe := ctx.Query("__iframe") != "true" // ex: true

	// 隱藏資訊(__token_、__previous_)
	hiddenFields := map[string]string{
		"__token_":    h.authSrv().AddToken(),
		"__previous_": infoUrl,
	}

	// 如果url沒設置則都為空
	if ctx.Query("__iframe") != "" {
		hiddenFields["__iframe"] = ctx.Query("__iframe")
	}
	if ctx.Query("__iframe_id") != "" {
		hiddenFields["__iframe_id"] = ctx.Query("__iframe_id")
	}

	content := formContent(aForm().
		SetContent(formInfo.FieldList).            // 表單欄位資訊
		SetFieldsHTML(f.HTMLContent).              // ex:""
		SetTabContents(formInfo.GroupFieldList).   // ex:[]
		SetTabHeaders(formInfo.GroupFieldHeaders). // ex:[]
		SetPrefix(h.config.PrefixFixSlash()).      // ex:/admin
		SetInputWidth(f.InputWidth).               // ex:0
		SetHeadWidth(f.HeadWidth).                 // ex:0
		SetPrimaryKey(panel.GetPrimaryKey().Name). // ex:id
		SetUrl(editUrl).
		SetAjax(f.AjaxSuccessJS, f.AjaxErrorJS).
		SetLayout(f.Layout).                      // ex:LayoutDefault
		SetTitle("Edit").
		SetHiddenFields(hiddenFields).            // 隱藏資訊
		SetOperationFooter(formFooter(footerKind, // formFooter 處理繼續新增、繼續編輯、保存、重製....等HTML語法
						!f.IsHideContinueEditCheckBox,
						!f.IsHideContinueNewCheckBox,
						f.IsHideResetButton)).
		SetHeader(f.HeaderHtml). // ex:HeaderHtml、FooterHtml為[]
		SetFooter(f.FooterHtml), len(formInfo.GroupFieldHeaders) > 0, !isNotIframe, f.IsHideBackButton, f.Header)

	// 一般不會執行
	if f.Wrapper != nil {
		content = f.Wrapper(content)
	}

	h.HTML(ctx, user, types.Panel{
		Content:     alert + content,
		Description: template2.HTML(formInfo.Description),
		Title:       modules.AorBHTML(isNotIframe, template2.HTML(formInfo.Title), ""),
	}, alert == "" || ((len(animation) > 0) && animation[0]))

	// 一般不會執行
	if isEdit {
		ctx.AddHeader("X-PJAX-Url", showEditUrl)
	}
}

// EditForm 更新資料(POST功能)
func (h *Handler) EditForm(ctx *context.Context) {
	param := guard.GetEditFormParam(ctx)

	// 如果有上傳頭像檔案才會執行，否則為空map[]
	if len(param.MultiForm.File) > 0 {
		err := file.GetFileEngine(h.config.FileUploadEngine.Name).Upload(param.MultiForm)
		if err != nil {
			if ctx.WantJSON() {
				response.Error(ctx, err.Error())
			} else {
				h.showForm(ctx, aAlert().Warning(err.Error()), param.Prefix, param.Param, true)
			}
			return
		}
	}

	// GetForm 將參數值設置至BaseTable.Form(FormPanel(struct)).primaryKey
	// field為編輯頁面每一欄位的資訊FormField(struct)
	for _, field := range param.Panel.GetForm().FieldList {
		// FormField.FormType為表單型態，ex: select、text、file
		// 頭像會執行此動作
		if field.FormType == form.File &&
			len(param.MultiForm.File[field.Field]) == 0 &&
			param.MultiForm.Value[field.Field+"__delete_flag"][0] != "1" {
			// 刪除param.MultiForm.Value[field.Field]值
			delete(param.MultiForm.Value, field.Field)
		}
	}

	// 更新資料
	// Value()取得multipart/form-data所設定的參數
	err := param.Panel.UpdateData(param.Value())
	if err != nil {
		// 判斷header裡包含accept:json
		if ctx.WantJSON() {
			response.Error(ctx, err.Error())
		} else {
			h.showForm(ctx, aAlert().Warning(err.Error()), param.Prefix, param.Param, true)
		}
		return
	}

	// -------下面四個條件式用戶、角色、權限都不會執行---------
	if param.Panel.GetForm().Responder != nil {
		param.Panel.GetForm().Responder(ctx)
		return
	}
	if ctx.WantJSON() && !param.IsIframe {
		response.OkWithData(ctx, map[string]interface{}{
			"url": param.PreviousPath,
		})
		return
	}

	// --------------在介面中選擇繼續新增會執行，執行後直接return---------------
	if !param.FromList {
		if isNewUrl(param.PreviousPath, param.Prefix) {
			// ---------------------繼續編輯新增此條件後return----------------
			// 新增表單介面不需要__edit_pk
			h.showNewForm(ctx, param.Alert, param.Prefix, param.Param.DeleteEditPk().GetRouteParamStr(), true)
			return
		}
		if isEditUrl(param.PreviousPath, param.Prefix) {
			// ---------------------繼續編輯執行此條件後return----------------
			h.showForm(ctx, param.Alert, param.Prefix, param.Param, true, false)
			return
		}
		ctx.HTML(http.StatusOK, fmt.Sprintf(`<script>location.href="%s"</script>`, param.PreviousPath))
		ctx.AddHeader("X-PJAX-Url", param.PreviousPath)
		return
	}
	if param.IsIframe {
		ctx.HTML(http.StatusOK, fmt.Sprintf(`<script>
		swal('%s', '', 'success');
		setTimeout(function(){
			$("#%s", window.parent.document).hide();
			$('.modal-backdrop.fade.in', window.parent.document).hide();
		}, 1000)
</script>`, "success", param.IframeID))
		return
	}

	buf := h.showTable(ctx, param.Prefix, param.Param.DeletePK().DeleteEditPk(), nil)
	ctx.HTML(http.StatusOK, buf.String())

	ctx.AddHeader("X-PJAX-Url", param.PreviousPath)
}
