package models

import (
	"goadminapi/modules/collection"
	"goadminapi/modules/db"
	"goadminapi/modules/db/dialect"
)

const (
	SiteItemOpenState = 1
	SiteItemOffState  = 0
)

// SiteModel is role model structure.
// site資料表欄位
type SiteModel struct {
	Base

	Id    int64
	Key   string
	Value string
	Desc  string
	State int64

	CreatedAt string
	UpdatedAt string
}

// 設置SiteModel(struct)後回傳
func Site() SiteModel {
	return SiteModel{Base: Base{TableName: "site"}}
}

// 將參數(con)設置至SiteModel.Conn後回傳
func (t SiteModel) SetConn(con db.Connection) SiteModel {
	t.Conn = con
	return t
}

func (t SiteModel) Init(cfg map[string]string) {
	items, err := t.Table(t.TableName).All()
	if db.CheckError(err, db.QUERY) {
		panic(err)
	}
	// 將items轉換成type Collection []map[string]interface{}
	itemsCol := collection.Collection(items)
	for key, value := range cfg {
		row := itemsCol.Where("key", "=", key)
		if row.Length() == 0 {
			_, err := t.Table(t.TableName).Insert(dialect.H{
				"key":   key,
				"value": value,
				"state": SiteItemOpenState,
			})
			if db.CheckError(err, db.INSERT) {
				panic(err)
			}
		}
		//else {
		//	if value != "" {
		//		_, err := t.Table(t.TableName).
		//			Where("key", "=", key).Update(dialect.H{
		//			"value": value,
		//		})
		//		if db.CheckError(err, db.UPDATE) {
		//			panic(err)
		//		}
		//	}
		//}
	}
}

func (t SiteModel) AllToMap() map[string]string {

	var m = make(map[string]string, 0)

	items, err := t.Table(t.TableName).Where("state", "=", SiteItemOpenState).All()
	if db.CheckError(err, db.QUERY) {
		return m
	}

	for _, item := range items {
		m[item["key"].(string)] = item["value"].(string)
	}

	return m
}
