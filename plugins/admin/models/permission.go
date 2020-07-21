package models

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
