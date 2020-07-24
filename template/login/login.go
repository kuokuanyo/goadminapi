package login

// Login也是Component(interface)
type Login struct {
	Name string
}

// 設置Login(struct)並回傳
func GetLoginComponent() *Login {
	return &Login{
		Name: "login",
	}
}

// var DefaultFuncMap = template.FuncMap{
// 	"lang":     language.Get,
// 	"langHtml": language.GetFromHtml,
// 	"link": func(cdnUrl, prefixUrl, assetsUrl string) string {
// 		if cdnUrl == "" {
// 			return prefixUrl + assetsUrl
// 		}
// 		return cdnUrl + assetsUrl
// 	},
// 	"isLinkUrl": func(s string) bool {
// 		return (len(s) > 7 && s[:7] == "http://") || (len(s) > 8 && s[:8] == "https://")
// 	},
// 	"render": func(s, old, repl template.HTML) template.HTML {
// 		return template.HTML(strings.Replace(string(s), string(old), string(repl), -1))
// 	},
// 	"renderJS": func(s template.JS, old, repl template.HTML) template.JS {
// 		return template.JS(strings.Replace(string(s), string(old), string(repl), -1))
// 	},
// 	"divide": func(a, b int) int {
// 		return a / b
// 	},
// }