package controller

// import (
// 	"encoding/base64"
// 	"goadminapi/context"
// 	"goadminapi/plugins/admin/modules"
// 	"goadminapi/plugins/admin/modules/response"
// 	"os"

// 	unipdf "github.com/unidoc/unidoc/pdf/model"

// 	"github.com/ConvertAPI/convertapi-go"
// 	"github.com/ConvertAPI/convertapi-go/config"
// 	"github.com/ConvertAPI/convertapi-go/param"
// )

// // ConvertFile 轉檔
// func (h *Handler) ConvertFile(ctx *context.Context) {
// 	filename := ctx.FormValue("file")
// 	RequestConvertAPI(ctx, filename)
// 	response.OkWithData(ctx, map[string]interface{}{
// 		"file": "http://www.cco.com.tw:8080/convert-file/" + filename + ".xlsx",
// 	})
// 	return
// }

// // Evaluate 取得頁數、估價函式
// func (h *Handler) Evaluate(ctx *context.Context) {
// 	filename := modules.Uuid()
// 	content := ctx.FormValue("file")

// 	dec, err := base64.StdEncoding.DecodeString(content)
// 	if err != nil {
// 		response.BadRequest(ctx, err.Error())
// 		return
// 	}

// 	WriteOriginalPDF(ctx, filename, dec)

// 	response.OkWithData(ctx, map[string]interface{}{
// 		"filename": filename,
// 		"pages":    GetPages(ctx, filename),
// 	})
// 	return
// }

// // RequestConvertAPI 請求轉檔API
// func RequestConvertAPI(ctx *context.Context, file string) {
// 	// 轉excel
// 	config.Default.Secret = "QaPpNniwBlJ7UpeY"
// 	pdfRes := convertapi.ConvDef("pdf", "xlsx", param.NewPath("file", "./themes/adminlte/resource/assets/dist/original/"+file+".pdf", nil))

// 	// save to file
// 	_, err := pdfRes.ToPath("./themes/adminlte/resource/assets/dist/convert-file/" + file + ".xlsx")
// 	if err != nil {
// 		response.BadRequest(ctx, "請檢查檔案是否有問題，如:檔案受密碼保護...等")
// 		return
// 	}
// }

// // GetPages 取得檔案頁數
// func GetPages(ctx *context.Context, filename string) (numPages int) {
// 	file, err := os.Open("./themes/adminlte/resource/assets/dist/original/" + filename + ".pdf")
// 	if err != nil {
// 		response.BadRequest(ctx, err.Error())
// 		return
// 	}

// 	pdfReader, err := unipdf.NewPdfReader(file)
// 	if err != nil {
// 		response.BadRequest(ctx, err.Error())
// 		return
// 	}
// 	numPages, err = pdfReader.GetNumPages()
// 	if err != nil {
// 		response.BadRequest(ctx, "檔案發生錯誤，可能原因:檔案被加密...等")
// 		return
// 	}
// 	return
// }

// // WriteOriginalPDF 儲存原始PDF檔案
// func WriteOriginalPDF(ctx *context.Context, file string, content []byte) {
// 	// 原始pdf
// 	f, err := os.Create("./themes/adminlte/resource/assets/dist/original/" + file + ".pdf")
// 	if err != nil {
// 		response.BadRequest(ctx, err.Error())
// 		return
// 	}
// 	defer f.Close()

// 	if _, err := f.Write(content); err != nil {
// 		response.BadRequest(ctx, err.Error())
// 		return
// 	}
// 	if err := f.Sync(); err != nil {
// 		response.BadRequest(ctx, err.Error())
// 		return
// 	}
// }
