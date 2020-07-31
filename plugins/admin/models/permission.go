package models

import (
	"goadminapi/modules/db"
	"strings"
)

type PermissionModel struct {
	Base
	Id         int64
	Name       string
	Slug       string
	HttpMethod []string
	HttpPath   []string
	CreatedAt  string
	UpdatedAt  string
}

// Permission 初始化PermissionModel
func Permission() PermissionModel {
	return PermissionModel{Base: Base{TableName: "permissions"}}
}

func (t PermissionModel) SetConn(con db.Connection) PermissionModel {
	t.Conn = con
	return t
}

// MapToModel 將map設置至permission model
func (t PermissionModel) MapToModel(m map[string]interface{}) PermissionModel {
	t.Id = m["id"].(int64)
	t.Name, _ = m["name"].(string)
	t.Slug, _ = m["slug"].(string)

	methods, _ := m["http_method"].(string)
	if methods != "" {
		t.HttpMethod = strings.Split(methods, ",")
	} else {
		t.HttpMethod = []string{""}
	}

	path, _ := m["http_path"].(string)
	t.HttpPath = strings.Split(path, "\n")
	t.CreatedAt, _ = m["created_at"].(string)
	t.UpdatedAt, _ = m["updated_at"].(string)
	return t
}

// IsSlugExist 檢查標誌是否已經存在
func (t PermissionModel) IsSlugExist(slug string, id string) bool {
	if id == "" {
		check, _ := t.Table(t.TableName).Where("slug", "=", slug).First()
		return check != nil
	}
	check, _ := t.Table(t.TableName).
		Where("slug", "=", slug).
		Where("id", "!=", id).
		First()
	return check != nil
}
