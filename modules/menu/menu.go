package menu

import (
	"goadminapi/modules/db"
	"goadminapi/plugins/admin/models"
	"regexp"
	"strconv"
)

// Item為資料表goadmin_menu的欄位
type Item struct {
	Name         string
	ID           string
	Url          string
	Icon         string
	Header       string
	Active       string
	ChildrenList []Item
}

type Menu struct {
	List     []Item
	Options  []map[string]string
	MaxOrder int64
}

// 回傳參數user(struct)的Menu(設置menuList、menuOption、MaxOrder)
func GetGlobalMenu(user models.UserModel, conn db.Connection) *Menu {
	var (
		menus      []map[string]interface{}
		menuOption = make([]map[string]string, 0)
	)
	// 查詢角色、權限
	user.WithRoles().WithMenus()
	if user.IsSuperAdmin() {
		// 取得多筆資料(利用where、order等篩選)
		menus, _ = db.WithDriver(conn).Table("menu").
			Where("id", ">", 0).
			OrderBy("order", "asc").
			All()
	} else {
		var ids []interface{}
		for i := 0; i < len(user.MenuIds); i++ {
			ids = append(ids, user.MenuIds[i])
		}
		menus, _ = db.WithDriver(conn).Table("goadmin_menu").
			WhereIn("id", ids).
			OrderBy("order", "asc").
			All()
	}

	var title string
	for i := 0; i < len(menus); i++ {
		title = menus[i]["title"].(string)
		menuOption = append(menuOption, map[string]string{
			"id":    strconv.FormatInt(menus[i]["id"].(int64), 10),
			"title": title,
		})
	}

	// 將參數menus轉換[]Item(Item是資料表menu的欄位)
	menuList := constructMenuTree(menus, 0)

	return &Menu{
		List:     menuList,   // 所有菜單資訊([]Item)，Item為資料表menu的欄位
		Options:  menuOption, // 設置每個菜單的id、title
		MaxOrder: menus[len(menus)-1]["parent_id"].(int64),
	}
}

// 將參數menus轉換成[]Item類別(Item(struct)是資料表menu的欄位)
func constructMenuTree(menus []map[string]interface{}, parentID int64) []Item {

	branch := make([]Item, 0)

	var title string
	for j := 0; j < len(menus); j++ {
		if parentID == menus[j]["parent_id"].(int64) {
			title = menus[j]["title"].(string)

			header, _ := menus[j]["header"].(string)

			child := Item{
				Name:         title,
				ID:           strconv.FormatInt(menus[j]["id"].(int64), 10),
				Url:          menus[j]["uri"].(string),
				Icon:         menus[j]["icon"].(string),
				Header:       header,
				Active:       "",
				ChildrenList: constructMenuTree(menus, menus[j]["id"].(int64)),
			}

			branch = append(branch, child)
		}
	}

	return branch
}

// 設定menu的active
func (menu *Menu) SetActiveClass(path string) *Menu {
	reg, _ := regexp.Compile(`\?(.*)`)
	path = reg.ReplaceAllString(path, "")
	for i := 0; i < len(menu.List); i++ {
		menu.List[i].Active = ""
	}

	for i := 0; i < len(menu.List); i++ {
		if menu.List[i].Url == path && len(menu.List[i].ChildrenList) == 0 {
			menu.List[i].Active = "active"
			return menu
		}
		for j := 0; j < len(menu.List[i].ChildrenList); j++ {
			if menu.List[i].ChildrenList[j].Url == path {
				menu.List[i].Active = "active"
				menu.List[i].ChildrenList[j].Active = "active"
				return menu
			}
			menu.List[i].Active = ""
			menu.List[i].ChildrenList[j].Active = ""
		}
	}
	return menu
}
