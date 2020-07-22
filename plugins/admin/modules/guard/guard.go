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
