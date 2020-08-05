package controller

import (
	"goadminapi/context"
	"goadminapi/plugins/admin/modules/guard"
	"goadminapi/plugins/admin/modules/response"
)

// Delete 刪除資料後回傳token
func (h *Handler) Delete(ctx *context.Context) {
	param := guard.GetDeleteParam(ctx)

	if err := h.table(param.Prefix, ctx).DeleteData(param.Id); err != nil {
		response.Error(ctx, "delete fail")
		panic(err)
	}

	response.OkWithData(ctx, map[string]interface{}{
		"token": h.authSrv().AddToken(),
	})
}
