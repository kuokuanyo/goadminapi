package dialect

import "fmt"

type mssql struct {
	commonDialect
}

func (mssql) GetName() string {
	return "mssql"
}

func (mssql) ShowColumns(table string) string {
	return fmt.Sprintf("select column_name, data_type from information_schema.columns where table_name = '%s'", table)
}

func (mssql) ShowTables() string {
	return "select * from information_schema.TABLES"
}
