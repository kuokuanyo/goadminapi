package controller

import (
	"goadminapi/context"
	"goadminapi/modules/file"
	"goadminapi/plugins/admin/modules/guard"
	"goadminapi/plugins/admin/modules/response"
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
}
