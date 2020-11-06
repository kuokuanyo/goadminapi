package refund

import (
	"bytes"
	"fmt"
	"html/template"
	"strings"
)

// Refund也是Component(interface)
type Refund struct {
	Name string
}

// 設置Refund(struct)並回傳
func GetRefundComponent() *Refund {
	return &Refund{
		Name: "refund",
	}
}

// 加入函式加入模板中
var DefaultFuncMap = template.FuncMap{
	"link": func(cdnUrl, prefixUrl, assetsUrl string) string {
		if cdnUrl == "" {
			return prefixUrl + assetsUrl
		}
		return cdnUrl + assetsUrl
	},
	"isLinkUrl": func(s string) bool {
		return (len(s) > 7 && s[:7] == "http://") || (len(s) > 8 && s[:8] == "https://")
	},
	"render": func(s, old, repl template.HTML) template.HTML {
		return template.HTML(strings.Replace(string(s), string(old), string(repl), -1))
	},
	"renderJS": func(s template.JS, old, repl template.HTML) template.JS {
		return template.JS(strings.Replace(string(s), string(old), string(repl), -1))
	},
	"divide": func(a, b int) int {
		return a / b
	},
}

// ----------------------Refund建立Component的所有方法-----------------------------------
// 添加Refund_theme1給新的HTML模板，接著將函式加入模板並解析
// 最後回傳模板及模板名稱
func (l *Refund) GetTemplate() (*template.Template, string) {
	tmpl, err := template.New("refund_theme1").
		// Funcs將要添加的函式元素加入
		Funcs(DefaultFuncMap).
		Parse(refundTmpl)
	if err != nil {
		panic("Refund GetTemplate Error: ")
	}

	return tmpl, "refund_theme1"
}

// AssetsList為css、js文件路徑
func (l *Refund) GetAssetList() []string { return AssetsList }

// 首先取得模板及模板名稱，取得登入介面的html
func (l *Refund) GetContent() template.HTML {
	buffer := new(bytes.Buffer)
	tmpl, defineName := l.GetTemplate()
	err := tmpl.ExecuteTemplate(buffer, defineName, l)
	if err != nil {
		fmt.Println("ComposeHtml Error:", err)
	}
	return template.HTML(buffer.String())
}

// 取得css、js檔案的[]byte
func (l *Refund) GetAsset(name string) ([]byte, error) { return Asset(name[1:]) }
func (l *Refund) IsAPage() bool                        { return true }
func (l *Refund) GetName() string                      { return "refund" }

// ----------------------login建立Component的所有方法-----------------------------------
