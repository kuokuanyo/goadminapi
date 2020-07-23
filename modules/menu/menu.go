package menu

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