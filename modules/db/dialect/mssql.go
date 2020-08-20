package dialect

import "fmt"

type mssql struct {
	commonDialect
}

// GetName return "mssql"
func (mssql) GetName() string {
	return "mssql"
}

// ShowColumns return select column_name, data_type from information_schema.columns where table_name = '%s'
func (mssql) ShowColumns(table string) string {
	return fmt.Sprintf("select column_name, data_type from information_schema.columns where table_name = '%s'", table)
}

// ShowTables return select * from information_schema.TABLES
func (mssql) ShowTables() string {
	return "select * from information_schema.TABLES"
}
