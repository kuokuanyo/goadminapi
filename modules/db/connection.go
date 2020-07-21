package db

import (
	"database/sql"
	"strings"
)

// 資料庫連接的處理程序
type Connection interface {
	// 查詢
	Query(query string, args ...interface{}) ([]map[string]interface{}, error)

	// 執行
	Exec(query string, args ...interface{}) (sql.Result, error)
	// 查詢(有給定conn名稱)
	QueryWithConnection(conn, query string, args ...interface{}) ([]map[string]interface{}, error)
	// 執行(有給定conn名稱)
	ExecWithConnection(conn, query string, args ...interface{}) (sql.Result, error)

	QueryWithTx(tx *sql.Tx, query string, args ...interface{}) ([]map[string]interface{}, error)

	ExecWithTx(tx *sql.Tx, query string, args ...interface{}) (sql.Result, error)

	// InitDB initialize the database connections.
	// 初始化資料庫連接
	InitDB(cfg map[string]Database) Connection

	// GetName get the connection name.
	Name() string

	Close() []error

	// GetDelimiter get the default testDelimiter.
	GetDelimiter() string

	GetDB(key string) *sql.DB
}

// 藉由參數(driver = mysql、mssql...)取得Connection(interface)
func GetConnectionByDriver(driver string) Connection {
	switch driver {
	case "mysql":
		return GetMysqlDB()
	// case "mssql":
	// 	return GetMssqlDB()
	// case "sqlite":
	// 	return GetSqliteDB()
	// case "postgresql":
	// 	return GetPostgresqlDB()
	default:
		panic("driver not found!")
	}
}

const (
	INSERT = 0
	DELETE = 1
	UPDATE = 2
	QUERY  = 3
)

var ignoreErrors = [...][]string{
	// insert
	{
		"LastInsertId is not supported",
		"There is no generated identity value",
	},
	// delete
	{
		"no affect",
	},
	// update
	{
		"LastInsertId is not supported",
		"There is no generated identity value",
		"no affect",
	},
	// query
	{
		"LastInsertId is not supported",
		"There is no generated identity value",
		"no affect",
		"out of index",
	},
}

// 檢查是否有錯誤
func CheckError(err error, t int) bool {
	if err == nil {
		return false
	}
	for _, msg := range ignoreErrors[t] {
		if strings.Contains(err.Error(), msg) {
			return false
		}
	}
	return true
}
