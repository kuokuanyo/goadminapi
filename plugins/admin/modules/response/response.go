package response

import (
	"goadminapi/context"
	"net/http"
)

// 成功，回傳code:200 and msg:ok and data
func OkWithData(ctx *context.Context, data map[string]interface{}) {
	ctx.JSON(http.StatusOK, map[string]interface{}{
		"code": http.StatusOK,
		"msg":  "ok",
		"data": data,
	})
}

// 錯誤請求，回傳code:400 and msg
func BadRequest(ctx *context.Context, msg string) {
	ctx.JSON(http.StatusBadRequest, map[string]interface{}{
		"code": http.StatusBadRequest,
		// Get依照設定的語言給予訊息
		"msg": msg,
	})
}

// 錯誤，回傳code:500 and msg
func Error(ctx *context.Context, msg string) {
	ctx.JSON(http.StatusInternalServerError, map[string]interface{}{
		"code": http.StatusInternalServerError,
		"msg":  msg,
	})
}
