package models

import (
	"goadminapi/modules/db"
	"goadminapi/modules/db/dialect"
	"strconv"
	"time"
)

// MenuModel is menu model structure.
type MenuModel struct {
	Base

	Id        int64
	Title     string
	ParentId  int64
	Icon      string
	Uri       string
	Header    string
	CreatedAt string
	UpdatedAt string
}

// 設置資料表後後回傳
func Menu() MenuModel {
	return MenuModel{Base: Base{TableName: "menu"}}
}

// 將參數id與tablename(goadmin_menu)設置至MenuModel(struct)後回傳
func MenuWithId(id string) MenuModel {
	idInt, _ := strconv.Atoi(id)
	return MenuModel{Base: Base{TableName: "menu"}, Id: int64(idInt)}
}

// 將參數con設置至MenuModel.Base.Conn
func (t MenuModel) SetConn(con db.Connection) MenuModel {
	t.Conn = con
	return t
}

// 將數值新增至資料表後將數值設置至MenuModel(struct)中
func (t MenuModel) New(title, icon, uri, header string, parentId, order int64) (MenuModel, error) {
	id, err := t.Table(t.TableName).Insert(dialect.H{
		"title":     title,
		"parent_id": parentId,
		"icon":      icon,
		"uri":       uri,
		"order":     order,
		"header":    header,
	})

	t.Id = id
	t.Title = title
	t.ParentId = parentId
	t.Icon = icon
	t.Uri = uri
	t.Header = header
	return t, err
}

// 將menu資料表條件為id = MenuModel.Id的資料透過參數(由multipart/form-data設置)更新
func (t MenuModel) Update(title, icon, uri, header string, parentId int64) (int64, error) {
	return t.Table(t.TableName).
		Where("id", "=", t.Id).
		Update(dialect.H{
			"title":      title,
			"parent_id":  parentId,
			"icon":       icon,
			"uri":        uri,
			"header":     header,
			"updated_at": time.Now().Format("2006-01-02 15:04:05"),
		})
}

// 必須刪除menu以及role_menu資料表的資料，如果是其他菜單的父級也必須刪除
func (t MenuModel) Delete() {
	// 必須刪除menu以及role_menu資料表的資料
	_ = t.Table(t.TableName).Where("id", "=", t.Id).Delete()
	_ = t.Table("role_menu").Where("menu_id", "=", t.Id).Delete()

	// 如果是其他菜單的父級也必須刪除
	items, _ := t.Table(t.TableName).Where("parent_id", "=", t.Id).All()
	if len(items) > 0 {
		ids := make([]interface{}, len(items))
		for i := 0; i < len(ids); i++ {
			ids[i] = items[i]["id"]
		}
		_ = t.Table("role_menu").WhereIn("menu_id", ids).Delete()
	}
	_ = t.Table(t.TableName).Where("parent_id", "=", t.Id).Delete()
}

// 刪除role_menu資料表中menu_id = MenuModel.Id條件的資料
func (t MenuModel) DeleteRoles() error {
	return t.Table("role_menu").
		Where("menu_id", "=", t.Id).
		Delete()
}

// 檢查role_menu資料表裡是否有符合role_id = 參數roleId、menu_id = MenuModel.Id條件
func (t MenuModel) CheckRole(roleId string) bool {
	checkRole, _ := t.Table("role_menu").
		Where("role_id", "=", roleId).
		Where("menu_id", "=", t.Id).
		First()
	return checkRole != nil
}

// 先檢查role_menu條件，接著將參數roleId(role_id)與MenuModel.Id(menu_id)加入role_menu資料表
func (t MenuModel) AddRole(roleId string) (int64, error) {
	if roleId != "" {
		// 檢查goadmin_role_menu資料表裡是否有符合role_id = 參數roleId與menu_id = MenuModel.Id條件
		if !t.CheckRole(roleId) {
			return t.Table("role_menu").
				Insert(dialect.H{
					"role_id": roleId,
					"menu_id": t.Id,
				})
		}
	}
	return 0, nil
}
