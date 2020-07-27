package controller

import (
	"crypto/md5"
	"fmt"
	"goadminapi/context"
	"mime"
	"net/http"
	"path"
	"strconv"
	"strings"

	"goadminapi/template"
)

// 處理前端檔案
func (h *Handler) Assets(ctx *context.Context) {
	// URLRemovePrefix將URL的前綴(ex:/admin)去除
	filepath := h.config.URLRemovePrefix(ctx.Path())

	// aTemplate判斷templateMap(map[string]Template)的key鍵是否參數globalCfg.Theme，有則回傳Template(interface)
	data, err := aTemplate().GetAsset(filepath)

	if err != nil {
		// 對map[string]Component迴圈，對每一個Component(interface)執行GetAsset方法
		data, err = template.GetAsset(filepath)
		if err != nil {
			// 將狀態碼，標頭(header)及body寫入Context.Response
			ctx.Write(http.StatusNotFound, map[string]string{}, "")
			panic("asset err")
		}
	}

	var contentType = mime.TypeByExtension(path.Ext(filepath))

	if contentType == "" {
		contentType = "application/octet-stream"
	}

	etag := fmt.Sprintf("%x", md5.Sum(data))

	// 藉由參數key獲得Header
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
