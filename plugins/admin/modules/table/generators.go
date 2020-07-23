package table

import (
	"goadminapi/modules/config"
	"goadminapi/modules/db"
)

type SystemTable struct {
	conn db.Connection
	c    *config.Config
}

// 將參數設置至SystemTable(struct)後回傳
func NewSystemTable(conn db.Connection, c *config.Config) *SystemTable {
	return &SystemTable{conn: conn, c: c}
}
