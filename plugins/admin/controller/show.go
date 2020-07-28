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

func (h *Handler) ShowInfo(ctx *context.Context) {
	prefix := ctx.Query("__prefix")

	panel := h.table(prefix, ctx)
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
