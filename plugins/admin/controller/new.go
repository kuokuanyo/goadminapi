package controller

import (
	"fmt"
	"goadminapi/context"
	"goadminapi/modules/auth"
	"goadminapi/plugins/admin/modules"
	"goadminapi/plugins/admin/modules/guard"
	"goadminapi/plugins/admin/modules/response"
	"goadminapi/template/types"
	template2 "html/template"
	"net/http"
)

func (h *Handler) ShowNewForm(ctx *context.Context) {
	param := guard.GetShowNewFormParam(ctx)

	h.showNewForm(ctx, "", param.Prefix, param.Param.GetRouteParamStr(), false)
}

// showNewForm 處理前端新增表單介面HTML
func (h *Handler) showNewForm(ctx *context.Context, alert template2.HTML, prefix, paramStr string, isNew bool) {
	// 透過參數ctx回傳目前登入的用戶(Context.UserValue["user"])並轉換成UserModel
	user := auth.Auth(ctx)

	// table 先透過參數prefix取得Table(interface)，接著判斷條件後將[]context.Node加入至Handler.operations後回傳
	panel := h.table(prefix, ctx)

	// GetNewForm 處理並取得表單欄位資訊(設置選項...等)
	formInfo := panel.GetNewForm()

	// routePathWithPrefix透過參數name取得該路徑名稱的URL，將url中的:__prefix改成第二個參數(prefix)
	// 用在返回按鍵
	infoUrl := h.routePathWithPrefix("info", prefix) + paramStr
	// 新建資料功能(post)
	newUrl := h.routePathWithPrefix("new", prefix)
	// 用於如果勾選繼續新增，導向新增資料頁面
	showNewUrl := h.routePathWithPrefix("show_new", prefix) + paramStr

	referer := ctx.Headers("Referer")
	if referer != "" && !isInfoUrl(referer) && !isNewUrl(referer, ctx.Query("__prefix")) {
		infoUrl = referer
	}

	// GetForm 將參數值設置至BaseTable.Form(FormPanel(struct)).primaryKey
	f := panel.GetForm()

	// 如果url沒設置__iframe則為空值
	isNotIframe := ctx.Query("__iframe") != "true"

	// 隱藏資訊(__token_、__previous_)
	hiddenFields := map[string]string{
		"__token_":    h.authSrv().AddToken(),
		"__previous_": infoUrl,
	}

	// 如果url沒設置則都為空
	// IframeKey = __iframe
	if ctx.Query("__iframe") != "" {
		hiddenFields["__iframe"] = ctx.Query("__iframe")
	}
	// IframeIDKey = __iframe_id
	if ctx.Query("__iframe_id") != "" {
		hiddenFields["__iframe_id"] = ctx.Query("__iframe_id")
	}

	content := formContent(aForm().
		SetPrefix(h.config.PrefixFixSlash()).      // ex:/admin
		SetFieldsHTML(f.HTMLContent).              // ex:""
		SetContent(formInfo.FieldList).            // 表單欄位資訊
		SetTabContents(formInfo.GroupFieldList).   // ex:[]
		SetTabHeaders(formInfo.GroupFieldHeaders). // ex:[]
		SetUrl(newUrl).
		SetAjax(f.AjaxSuccessJS, f.AjaxErrorJS).
		SetInputWidth(f.InputWidth).               // ex:0
		SetHeadWidth(f.HeadWidth).                 // ex:0
		SetLayout(f.Layout).                       // ex:LayoutDefault
		SetPrimaryKey(panel.GetPrimaryKey().Name). // ex:id
		SetHiddenFields(hiddenFields).
		SetTitle("New").
		// formFooter 處理繼續新增、繼續編輯、保存、重製....等HTML語法
		SetOperationFooter(formFooter("new", f.IsHideContinueEditCheckBox, f.IsHideContinueNewCheckBox,
						f.IsHideResetButton)).
		SetHeader(f.HeaderHtml). // ex:HeaderHtml、FooterHtml為[]
		SetFooter(f.FooterHtml), len(formInfo.GroupFieldHeaders) > 0, !isNotIframe, f.IsHideBackButton, f.Header)

	// 一般不會執行
	if f.Wrapper != nil {
		content = f.Wrapper(content)
	}

	h.HTML(ctx, user, types.Panel{
		Content:     alert + content,
		Description: template2.HTML(f.Description),
		Title:       modules.AorBHTML(isNotIframe, template2.HTML(f.Title), ""),
	}, alert == "")

	// --------一般不會執行，如果勾選繼續新增則會執行---------------
	if isNew {
		ctx.AddHeader("X-PJAX-Url", showNewUrl)
	}
}

// NewForm 新增表單資料(POST功能)
func (h *Handler) NewForm(ctx *context.Context) {
	param := guard.GetNewFormParam(ctx)

	// 如果有上傳頭像檔案才會執行，否則為空map[]
	// if len(param.MultiForm.File) > 0 {
	// 	err := file.GetFileEngine(h.config.FileUploadEngine.Name).Upload(param.MultiForm)
	// 	if err != nil {
	// 		if ctx.WantJSON() {
	// 			response.Error(ctx, err.Error())
	// 		} else {
	// 			h.showNewForm(ctx, aAlert().Warning(err.Error()), param.Prefix, param.Param.GetRouteParamStr(), true)
	// 		}
	// 		return
	// 	}
	// }

	err := param.Panel.InsertData(param.Value())
	if err != nil {
		if ctx.WantJSON() {
			response.Error(ctx, err.Error())
		} else {
			h.showNewForm(ctx, aAlert().Warning(err.Error()), param.Prefix, param.Param.GetRouteParamStr(), true)
		}
		return
	}

	// ---------------一般下面的if不會執行----------------------
	// type Responder func(ctx *context.Context)
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
			// ------繼續新增會執行此條件後return---------------------
			h.showNewForm(ctx, param.Alert, param.Prefix, param.Param.GetRouteParamStr(), true)
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

	buf := h.showTable(ctx, param.Prefix, param.Param, nil)
	ctx.HTML(http.StatusOK, buf.String())

	// GetRouteParamStr處理url後(?...)的部分(頁面設置、排序方式....等)
	ctx.AddHeader("X-PJAX-Url", h.routePathWithPrefix("info", param.Prefix)+param.Param.GetRouteParamStr())
}
