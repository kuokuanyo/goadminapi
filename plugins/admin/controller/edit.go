package controller

import (
	"fmt"
	"goadminapi/context"
	"goadminapi/modules/file"
	"goadminapi/plugins/admin/modules/guard"
	"goadminapi/plugins/admin/modules/response"
	"goadminapi/template/types/form"
	"net/http"
)

// EditForm 更新資料
func (h *Handler) EditForm(ctx *context.Context) {
	param := guard.GetEditFormParam(ctx)

	// 如果有上傳頭像檔案才會執行，否則為空map[]
	if len(param.MultiForm.File) > 0 {
		err := file.GetFileEngine(h.config.FileUploadEngine.Name).Upload(param.MultiForm)
		if err != nil {
			if ctx.WantJSON() {
				response.Error(ctx, err.Error())
			} else {
				//**************函式還沒寫***********************
				// h.showForm(ctx, aAlert().Warning(err.Error()), param.Prefix, param.Param, true)
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
			//**************函式還沒寫***********************
			// h.showForm(ctx, aAlert().Warning(err.Error()), param.Prefix, param.Param, true)
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
	if !param.FromList {
		if isNewUrl(param.PreviousPath, param.Prefix) {
			//**************函式還沒寫***********************
			//h.showNewForm(ctx, param.Alert, param.Prefix, param.Param.DeleteEditPk().GetRouteParamStr(), true)
			return
		}
		if isEditUrl(param.PreviousPath, param.Prefix) {
			//**************函式還沒寫***********************
			// h.showForm(ctx, param.Alert, param.Prefix, param.Param, true, false)
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
	// buf := h.showTable(ctx, param.Prefix, param.Param.DeletePK().DeleteEditPk(), nil)
	//ctx.HTML(http.StatusOK, buf.String())

	ctx.AddHeader("X-PJAX-Url", param.PreviousPath)
}
