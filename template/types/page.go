package types

import (
	"bytes"
	"fmt"
	"goadminapi/context"
	"goadminapi/plugins/admin/models"
	"html/template"
	"strconv"

	"goadminapi/modules/menu"
	"goadminapi/modules/utils"

	"goadminapi/modules/config"
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
	// System SystemInfo

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

type NewPageParam struct {
	User         models.UserModel
	Menu         *menu.Menu
	Panel        Panel
	Assets       template.HTML
	Buttons      Buttons
	Iframe       bool
	TmplHeadHTML template.HTML
	TmplFootJS   template.HTML
}

type GetPanelInfoFn func(ctx *context.Context) (Panel, error)

// 取得row_data_tmpl並解析模版
func ParseTableDataTmpl(content interface{}) string {
	var (
		c  string
		ok bool
	)
	if c, ok = content.(string); !ok {
		if cc, ok := content.(template.HTML); ok {
			c = string(cc)
		} else {
			c = string(content.(template.JS))
		}
	}
	t := template.New("row_data_tmpl")
	t, _ = t.Parse(c)
	buf := new(bytes.Buffer)
	_ = t.Execute(buf, TableRowData{Ids: `typeof(selectedRows)==="function" ? selectedRows().join() : ""`})
	return buf.String()
}

func (param NewPageParam) NavButtonsAndJS() (template.HTML, template.HTML) {
	navBtnFooter := template.HTML("")
	navBtn := template.HTML("")
	btnJS := template.JS("")

	for _, btn := range param.Buttons {
		navBtnFooter += btn.GetAction().FooterContent()
		content, js := btn.Content()
		navBtn += content
		btnJS += js
	}

	return template.HTML(ParseTableDataTmpl(navBtn)),
		navBtnFooter + template.HTML(ParseTableDataTmpl(`<script>`+btnJS+`</script>`))
}

// 設置Page(struct)回傳
func NewPage(param NewPageParam) *Page {

	navBtn, btnJS := param.NavButtonsAndJS()

	return &Page{
		User:  param.User,
		Menu:  *param.Menu,
		Panel: param.Panel,
		// System: SystemInfo{
		// 	Version: system.Version(),
		// 	Theme:   config.GetTheme(),
		// },
		UrlPrefix:      config.AssertPrefix(),
		Title:          config.GetTitle(),
		Logo:           config.GetLogo(),
		MiniLogo:       config.GetMiniLogo(),
		ColorScheme:    config.GetColorScheme(),
		IndexUrl:       config.GetIndexURL(),
		CdnUrl:         config.GetAssetUrl(),
		CustomHeadHtml: config.GetCustomHeadHtml(),
		CustomFootHtml: config.GetCustomFootHtml() + btnJS,
		FooterInfo:     config.GetFooterInfo(),
		AssetsList:     param.Assets,
		navButtons:     param.Buttons,
		Iframe:         param.Iframe,
		NavButtonsHTML: navBtn,
		TmplHeadHTML:   param.TmplHeadHTML,
		TmplFootJS:     param.TmplFootJS,
	}
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

// 取得content設置Panel並回傳
func (p Panel) GetContent(params ...bool) Panel {

	prod := false

	if len(params) > 0 {
		prod = params[0]
	}

	animation := template.HTML("")
	style := template.HTML("") // 處理css
	remove := template.HTML("")
	ani := config.GetAnimation()
	if ani.Type != "" && (len(params) < 2 || params[1]) {
		animation = template.HTML(` class='pjax-container-content animated ` + ani.Type + `'`)
		if ani.Delay != 0 {
			// 設置延遲
			style = template.HTML(fmt.Sprintf(`animation-delay: %fs;-webkit-animation-delay: %fs;`, ani.Delay, ani.Delay))
		}
		if ani.Duration != 0 {
			// 設定動畫持續時間
			style = template.HTML(fmt.Sprintf(`animation-duration: %fs;-webkit-animation-duration: %fs;`, ani.Duration, ani.Duration))
		}
		if style != "" {
			style = ` style="` + style + `"`
		}
		remove = template.HTML(`<script>
		$('.pjax-container-content .modal.fade').on('show.bs.modal', function (event) {
            // Fix Animate.css
			$('.pjax-container-content').removeClass('` + ani.Type + `');
        });
		</script>`)
	}

	p.Content = `<div` + animation + style + ">" + p.Content + "</div>" + remove
	// 切換側邊攔
	if p.MiniSidebar {
		p.Content += `<script>$("body").addClass("sidebar-collapse")</script>`
	}
	// 自動刷新頁面轉換
	if p.AutoRefresh {
		refreshTime := 60
		if len(p.RefreshInterval) > 0 {
			refreshTime = p.RefreshInterval[0]
		}
		// 設定1000秒刷新一次頁面
		p.Content += `<script>
window.setTimeout(function(){
	$.pjax.reload('#pjax-container');	
}, ` + template.HTML(strconv.Itoa(refreshTime*1000)) + `);
</script>`
	}
	if prod {
		// 壓縮內容
		utils.CompressedContent(&p.Content)
	}

	return p
}
