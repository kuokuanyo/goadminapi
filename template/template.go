package template

import (
	"bytes"
	"errors"
	"fmt"
	"goadminapi/plugins/admin/models"
	"goadminapi/template/login"
	"goadminapi/template/types"
	"html/template"
	"path"
	"strconv"
	"strings"
	"sync"

	c "goadminapi/modules/config"
	"goadminapi/modules/menu"
)

type ExecuteParam struct {
	User       models.UserModel
	Tmpl       *template.Template
	TmplName   string
	Panel      types.Panel
	Config     c.Config
	Menu       *menu.Menu
	Animation  bool
	Buttons    types.Buttons
	NoCompress bool
	Iframe     bool
}

type PageType uint8

var (
	templateMu sync.Mutex
	compMu     sync.Mutex
)

var templateMap = make(map[string]Template)

// Login(struct)屬於Component(interface)所有方法
var compMap = map[string]Component{
	// GetLoginComponent設置Login(struct)並回傳
	"login": login.GetLoginComponent(),
}

const (
	NormalPage PageType = iota
	Missing404Page
	Error500Page
	NoPermission403Page
)

const (
	CompCol       = "col"
	CompRow       = "row"
	CompForm      = "form"
	CompTable     = "table"
	CompDataTable = "datatable"
	CompTree      = "tree"
	CompTreeView  = "treeview"
	CompTabs      = "tabs"
	CompAlert     = "alert"
	CompLink      = "link"
	CompPaginator = "paginator"
	CompPopup     = "popup"
	CompBox       = "box"
	CompLabel     = "label"
	CompImage     = "image"
	CompButton    = "button"
)

var DefaultFuncMap = template.FuncMap{
	// "lang":     language.Get,
	// "langHtml": language.GetFromHtml,
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
	"renderRowDataHTML": func(id, content template.HTML, value ...map[string]types.InfoItem) template.HTML {
		return template.HTML(types.ParseTableDataTmplWithID(id, string(content), value...))
	},
	"renderRowDataJS": func(id template.HTML, content template.JS, value ...map[string]types.InfoItem) template.JS {
		return template.JS(types.ParseTableDataTmplWithID(id, string(content), value...))
	},
	"attr": func(s template.HTML) template.HTMLAttr {
		return template.HTMLAttr(s)
	},
	"js": func(s interface{}) template.JS {
		if ss, ok := s.(string); ok {
			return template.JS(ss)
		}
		if ss, ok := s.(template.HTML); ok {
			return template.JS(ss)
		}
		return ""
	},
	"changeValue": func(f types.FormField, index int) types.FormField {
		if len(f.ValueArr) > 0 {
			f.Value = template.HTML(f.ValueArr[index])
		}
		if len(f.OptionsArr) > 0 {
			f.Options = f.OptionsArr[index]
		}
		if f.FormType.IsSelect() {
			f.FieldClass += "_" + strconv.Itoa(index)
		}
		return f
	},
}

type Template interface {
	Name() string

	// Components
	// layout
	Col() types.ColAttribute
	// Row() types.RowAttribute

	// // form and table
	Form() types.FormAttribute
	// Table() types.TableAttribute
	DataTable() types.DataTableAttribute

	// TreeView() types.TreeViewAttribute
	// Tree() types.TreeAttribute
	Tabs() types.TabsAttribute
	Alert() types.AlertAttribute
	// Link() types.LinkAttribute

	Paginator() types.PaginatorAttribute
	// Popup() types.PopupAttribute
	Box() types.BoxAttribute

	Label() types.LabelAttribute
	Image() types.ImgAttribute

	Button() types.ButtonAttribute

	// Builder methods
	// GetTmplList() map[string]string
	GetAssetList() []string
	GetAssetImportHTML(exceptComponents ...string) template.HTML
	GetAsset(string) ([]byte, error)
	GetTemplate(bool) (*template.Template, string)
	// GetVersion() string
	// GetRequirements() []string
	GetHeadHTML() template.HTML
	GetFootJS() template.HTML
	Get404HTML() template.HTML
	Get500HTML() template.HTML
	Get403HTML() template.HTML
}

type Component interface {
	// GetTemplate return a *template.Template and a given key.
	GetTemplate() (*template.Template, string)
	GetAssetList() []string
	GetAsset(string) ([]byte, error)
	GetContent() template.HTML
	IsAPage() bool
	GetName() string
}

func Themes() []string {
	names := make([]string, len(templateMap))
	i := 0
	for k := range templateMap {
		names[i] = k
		i++
	}
	return names
}

// 取得預設的Template(interface)
func Default() Template {
	if temp, ok := templateMap[c.GetTheme()]; ok {
		return temp
	}
	panic("wrong theme name")
}

// 判斷templateMap(map[string]Template)的key鍵是否參數theme，有則回傳Template(interface)
func Get(theme string) Template {
	if temp, ok := templateMap[theme]; ok {
		return temp
	}
	panic("wrong theme name")
}

// 判斷map[string]Component是否有參數name(key)的值，有的話則回傳Component(interface)
func GetComp(name string) Component {
	// Component(interface)
	if comp, ok := compMap[name]; ok {
		return comp
	}
	panic("wrong component name")
}

// 將給定的數據寫入buf(struct)並回傳
func Execute(param ExecuteParam) *bytes.Buffer {
	buf := new(bytes.Buffer)
	err := param.Tmpl.ExecuteTemplate(buf, param.TmplName,
		types.NewPage(types.NewPageParam{
			User:         param.User,
			Menu:         param.Menu,
			Panel:        param.Panel.GetContent(append([]bool{param.Config.IsProductionEnvironment() && (!param.NoCompress)}, param.Animation)...),
			Assets:       GetComponentAssetImportHTML(),
			Buttons:      param.Buttons,
			Iframe:       param.Iframe,
			TmplHeadHTML: Default().GetHeadHTML(),
			TmplFootJS:   Default().GetFootJS(),
		}))
	if err != nil {
		fmt.Println("template execute error")
		panic(err)
	}
	return buf
}

// 新增主題及方法至templateMap(map[string]Template)
func Add(name string, temp Template) {
	templateMu.Lock()
	defer templateMu.Unlock()
	if temp == nil {
		panic("template is nil")
	}
	if _, dup := templateMap[name]; dup {
		panic("add template twice " + name)
	}
	templateMap[name] = temp
}

// 檢查compMap(map[string]Component)的物件後將前端文件路徑加入[]string中
func GetComponentAsset() []string {
	assets := make([]string, 0)
	for _, comp := range compMap {
		// AssetsList為css、js文件路徑
		assets = append(assets, comp.GetAssetList()...)
	}
	return assets
}

// 檢查compMap(map[string]Component)的物件是否符合條件並加入文件路徑到陣列中
func GetComponentAssetWithinPage() []string {
	assets := make([]string, 0)
	for _, comp := range compMap {
		if !comp.IsAPage() {
			assets = append(assets, comp.GetAssetList()...)
		}
	}
	return assets
}

// 透過參數s判斷css或js檔案，取得HTML
func getHTMLFromAssetUrl(s string) template.HTML {
	switch path.Ext(s) {
	case ".css":
		return template.HTML(`<link rel="stylesheet" href="` + c.GetAssetUrl() + c.Url("/assets"+s) + `">`)
	case ".js":
		return template.HTML(`<script src="` + c.GetAssetUrl() + c.Url("/assets"+s) + `"></script>`)
	default:
		return ""
	}
}

// 處理asset後並回傳HTML語法
func GetComponentAssetImportHTML() (res template.HTML) {
	// GetAssetImportHTML(Template(interface)的方法)
	res = Default().GetAssetImportHTML(c.GetExcludeThemeComponents()...)
	// 在頁面中獲取物件asset
	// 檢查map[string]Component物件是否符合條件並加入文件路徑到陣列中
	assets := GetComponentAssetWithinPage()

	for i := 0; i < len(assets); i++ {
		// 透過參數assets[i]判斷css或js檔案，取得HTML
		res += getHTMLFromAssetUrl(assets[i])
	}
	return
}

// 對map[string]Component迴圈，對每一個Component(interface)執行GetAsset方法
func GetAsset(path string) ([]byte, error) {
	for _, comp := range compMap {
		res, err := comp.GetAsset(path)
		if err == nil {
			return res, err
		}
	}
	return nil, errors.New(path + " not found")
}

// 透過參數msg設置Panel(struct)
func WarningPanel(msg string, pts ...PageType) types.Panel {
	pt := Error500Page
	if len(pts) > 0 {
		pt = pts[0]
	}
	pageTitle, description, content := GetPageContentFromPageType(msg, msg, msg, pt)
	return types.Panel{
		Content:     content,
		Description: description,
		Title:       pageTitle,
	}
}

// GetPageContentFromPageType從頁面類型取得頁面內容
func GetPageContentFromPageType(title, desc, msg string, pt PageType) (template.HTML, template.HTML, template.HTML) {
	if c.GetDebug() {
		return template.HTML(title), template.HTML(desc), Default().Alert().Warning(msg)
	}
	if pt == Missing404Page {
		if c.GetCustom404HTML() != template.HTML("") {
			return "", "", c.GetCustom404HTML()
		} else {
			return "", "", Default().Get404HTML()
		}
	} else if pt == NoPermission403Page {
		if c.GetCustom404HTML() != template.HTML("") {
			return "", "", c.GetCustom403HTML()
		} else {
			return "", "", Default().Get403HTML()
		}
	} else {
		if c.GetCustom500HTML() != template.HTML("") {
			return "", "", c.GetCustom500HTML()
		} else {
			return "", "", Default().Get500HTML()
		}
	}
}

func HTML(s string) template.HTML {
	return template.HTML(s)
}
