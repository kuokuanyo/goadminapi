package guard

import (
	"goadminapi/modules/db"
	"goadminapi/modules/service"

	"goadminapi/plugins/admin/modules/table"
)

type Guard struct {
	services  service.List
	conn      db.Connection
	tableList table.GeneratorList
	//navBtns   *types.Buttons
}

// 將參數s、c、t設置至Guard(struct)後回傳
func New(s service.List, c db.Connection, t table.GeneratorList) *Guard {
	return &Guard{
		services:  s,
		conn:      c,
		tableList: t,
		// navBtns:   b,
	}
}
