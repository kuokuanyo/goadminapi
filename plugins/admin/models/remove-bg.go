package models

import (
	"database/sql"
	"goadminapi/modules/db"
	"goadminapi/modules/db/dialect"
	"strconv"
)

type RemoveBgModel struct {
	Base
	Id       int64
	UserID   string
	Filename string
}

func RemoveBg() RemoveBgModel {
	return RemoveBgModel{Base: Base{TableName: "remove_background"}}
}

func (t RemoveBgModel) SetConn(con db.Connection) RemoveBgModel {
	t.Conn = con
	return t
}

func (t RemoveBgModel) WithTx(tx *sql.Tx) RemoveBgModel {
	t.Tx = tx
	return t
}

func RemoveBgWithId(id string) RemoveBgModel {
	idInt, _ := strconv.Atoi(id)
	return RemoveBgModel{Base: Base{TableName: "remove_background"}, Id: int64(idInt)}
}

func (t RemoveBgModel) New(userid, filename string) (RemoveBgModel, error) {

	id, err := t.WithTx(t.Tx).Table(t.TableName).Insert(dialect.H{
		"userid":   userid,
		"Filename": filename,
	})

	t.Id = id
	t.UserID = userid
	t.Filename = filename
	return t, err
}

func (t RemoveBgModel) Update(userid, filenmae string) (int64, error) {
	return t.WithTx(t.Tx).Table(t.TableName).
		Where("id", "=", t.Id).
		Update(dialect.H{
			"userid":   userid,
			"Filename": filenmae,
		})
}

func (t RemoveBgModel) MapToModel(m map[string]interface{}) RemoveBgModel {
	t.Id = m["id"].(int64)
	t.UserID, _ = m["userid"].(string)
	t.Filename, _ = m["filename"].(string)
	return t
}
