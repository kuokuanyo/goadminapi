package types

import (
	"bytes"
	"goadminapi/plugins/admin/models"
	"html/template"

	"goadminapi/modules/menu"
)

type Attribute struct {
	TemplateList map[string]string
}

type SystemInfo struct {
	Version string
	Theme   string
}

// 主要內容使用pjax的模板
type Panel struct {
	Title       template.HTML
	Description template.HTML
	Content     template.HTML

	CSS template.CSS
	JS  template.JS
	Url string

	// Whether to toggle the sidebar
	// 是否切換側邊攔
	MiniSidebar bool

	// Auto refresh page switch.
	// 自動刷新頁面轉換
	AutoRefresh bool
	// Refresh page intervals, the unit is second.
	// 刷新頁面間隔，單位為秒
	RefreshInterval []int

	Callbacks Callbacks
}

type Page struct {
	// User is the login user.
	User models.UserModel

	// Menu is the left side menu of the template.
	Menu menu.Menu

	// Panel is the main content of template.
	Panel Panel

	// System contains some system info.
	System SystemInfo

	// UrlPrefix is the prefix of url.
	UrlPrefix string

	// Title is the title of the web page.
	Title string

	// Logo is the logo of the template.
	Logo template.HTML

	// MiniLogo is the downsizing logo of the template.
	MiniLogo template.HTML

	// ColorScheme is the color scheme of the template.
	ColorScheme string

	// IndexUrl is the home page url of the site.
	IndexUrl string

	// AssetUrl is the cdn link of assets
	CdnUrl string

	// Custom html in the tag head.
	CustomHeadHtml template.HTML

	// Custom html after body.
	CustomFootHtml template.HTML

	TmplHeadHTML template.HTML
	TmplFootJS   template.HTML

	// Components assets
	AssetsList template.HTML

	// Footer info
	FooterInfo template.HTML

	// Load as Iframe or not
	Iframe bool

	// Top Nav Buttons
	navButtons     Buttons
	NavButtonsHTML template.HTML
}

type GetPanelFn func(ctx interface{}) (Panel, error)

type TableRowData struct {
	Id    template.HTML
	Ids   template.HTML
	Value map[string]InfoItem
}

// 創建並解析row_data_tmpl模板
func ParseTableDataTmplWithID(id template.HTML, content string, value ...map[string]InfoItem) string {
	t := template.New("row_data_tmpl")
	t, _ = t.Parse(content)
	buf := new(bytes.Buffer)
	v := make(map[string]InfoItem)
	if len(value) > 0 {
		v = value[0]
	}
	_ = t.Execute(buf, TableRowData{
		Id:    id,
		Ids:   `typeof(selectedRows)==="function" ? selectedRows().join() : ""`,
		Value: v,
	})
	return buf.String()
}
