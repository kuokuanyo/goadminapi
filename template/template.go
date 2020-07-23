package template

type PageType uint8

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
