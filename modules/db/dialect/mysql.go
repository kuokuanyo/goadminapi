package dialect

type mysql struct {
	commonDialect
}

// GetName return "mysql"
func (mysql) GetName() string {
	return "mysql"
}

// ShowColumns return "show columns in " + table
func (mysql) ShowColumns(table string) string {
	return "show columns in " + table
}

// ShowTables "show tables"
func (mysql) ShowTables() string {
	return "show tables"
}
