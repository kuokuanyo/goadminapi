package controller

import (
	"fmt"
	"goadminapi/context"
	"goadminapi/modules/file"
	"goadminapi/plugins/admin/modules/guard"
	"goadminapi/plugins/admin/modules/response"
	"net/http"
)

func (h *Handler) NewForm(ctx *context.Context) {
	param := guard.GetNewFormParam(ctx)

	// 如果有上傳頭像檔案才會執行，否則為空map[]
	if len(param.MultiForm.File) > 0 {
		err := file.GetFileEngine(h.config.FileUploadEngine.Name).Upload(param.MultiForm)
		if err != nil {
			if ctx.WantJSON() {
				response.Error(ctx, err.Error())
			} else {
				//**************函式還沒寫***********************
				// h.showNewForm(ctx, aAlert().Warning(err.Error()), param.Prefix, param.Param.GetRouteParamStr(), true)
			}
			return
		}
	}

	err := param.Panel.InsertData(param.Value())
	if err != nil {
		if ctx.WantJSON() {
			response.Error(ctx, err.Error())
		} else {
			//**************函式還沒寫***********************
			// h.showNewForm(ctx, aAlert().Warning(err.Error()), param.Prefix, param.Param.GetRouteParamStr(), true)
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
	if !param.FromList {
		if isNewUrl(param.PreviousPath, param.Prefix) {
			//**************函式還沒寫***********************
			// h.showNewForm(ctx, param.Alert, param.Prefix, param.Param.GetRouteParamStr(), true)
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

	//**************函式還沒寫***********************
	// buf := h.showTable(ctx, param.Prefix, param.Param, nil)
	// ctx.HTML(http.StatusOK, buf.String())
	//**************函式還沒寫***********************

	// GetRouteParamStr處理url後(?...)的部分(頁面設置、排序方式....等)
	ctx.AddHeader("X-PJAX-Url", h.routePathWithPrefix("info", param.Prefix)+param.Param.GetRouteParamStr())
}
