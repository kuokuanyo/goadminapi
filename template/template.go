package template

import "html/template"

type PageType uint8


// Login(struct)屬於Component(interface)所有方法
var compMap = map[string]Component{
	// GetLoginComponent設置Login(struct)並回傳
	// "login": login.GetLoginComponent(),
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

// Component is the interface which stand for a ui component.
type Component interface {
	// GetTemplate return a *template.Template and a given key.
	GetTemplate() (*template.Template, string)

	// GetAssetList return the assets url suffix used in the component.
	// example:
	//
	// {{.UrlPrefix}}/assets/login/css/bootstrap.min.css => login/css/bootstrap.min.css
	//
	// See:
	// https://github.com/GoAdminGroup/go-admin/blob/master/template/login/theme1.tmpl#L32
	// https://github.com/GoAdminGroup/go-admin/blob/master/template/login/list.go
	GetAssetList() []string

	// GetAsset return the asset content according to the corresponding url suffix.
	// Asset content is recommended to use the tool go-bindata to generate.
	//
	// See: http://github.com/jteeuwen/go-bindata
	GetAsset(string) ([]byte, error)

	GetContent() template.HTML

	IsAPage() bool

	GetName() string
}

// 判斷map[string]Component是否有參數name(key)的值，有的話則回傳Component(interface)
func GetComp(name string) Component {
	// Component(interface)
	if comp, ok := compMap[name]; ok {
		return comp
	}
	panic("wrong component name")
}

// // 透過參數msg設置Panel(struct)
// func WarningPanel(msg string, pts ...PageType) types.Panel {
// 	pt := Error500Page
// 	if len(pts) > 0 {
// 		pt = pts[0]
// 	}
// 	pageTitle, description, content := GetPageContentFromPageType(msg, msg, msg, pt)
// 	return types.Panel{
// 		// Default()取得預設的template(主題名稱已經通過全局配置)
// 		// Alert為Template(interface)的方法
// 		Content:     content,
// 		Description: description,
// 		Title:       pageTitle,
// 	}
// }

// func GetPageContentFromPageType(title, desc, msg string, pt PageType) (template.HTML, template.HTML, template.HTML) {
// 	// globalCfg.Debug
// 	if c.GetDebug() {
// 		return template.HTML(title), template.HTML(desc), Default().Alert().Warning(msg)
// 	}

// 	if pt == Missing404Page {
// 		if c.GetCustom404HTML() != template.HTML("") {
// 			return "", "", c.GetCustom404HTML()
// 		} else {
// 			return "", "", Default().Get404HTML()
// 		}
// 	} else if pt == NoPermission403Page {
// 		if c.GetCustom404HTML() != template.HTML("") {
// 			return "", "", c.GetCustom403HTML()
// 		} else {
// 			return "", "", Default().Get403HTML()
// 		}
// 	} else {
// 		if c.GetCustom500HTML() != template.HTML("") {
// 			return "", "", c.GetCustom500HTML()
// 		} else {
// 			return "", "", Default().Get500HTML()
// 		}
// 	}
// }
