package models

type RoleModel struct {
	Base
	Id        int64
	Name      string
	Slug      string
	CreatedAt string
	UpdatedAt string
}

// 初始化role model
func Role() RoleModel {
	return RoleModel{Base: Base{TableName: "roles"}}
}

// 將取得的值(參數m)設置至rolemodel
func (t RoleModel) MapToModel(m map[string]interface{}) RoleModel {
	t.Id = m["id"].(int64)
	t.Name, _ = m["name"].(string)
	t.Slug, _ = m["slug"].(string)
	t.CreatedAt, _ = m["created_at"].(string)
	t.UpdatedAt, _ = m["updated_at"].(string)
	return t
}