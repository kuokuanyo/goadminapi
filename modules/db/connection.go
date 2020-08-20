package db

import (
	"database/sql"
	"fmt"
	"strings"

	"goadminapi/modules/config"
	"goadminapi/modules/service"
)

// Connection 資料庫連接的處理程序
// Connection也屬於Service(interface)
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
	InitDB(cfg map[string]config.Database) Connection

	// GetName get the connection name.
	Name() string

	Close() []error

	// GetDelimiter get the default testDelimiter.
	GetDelimiter() string

	GetDB(key string) *sql.DB

	// BeginTxWithReadUncommitted() *sql.Tx
	// BeginTxWithReadCommitted() *sql.Tx
	// BeginTxWithRepeatableRead() *sql.Tx
	// BeginTx() *sql.Tx
	// BeginTxWithLevel(level sql.IsolationLevel) *sql.Tx

	// BeginTxWithReadUncommittedAndConnection(conn string) *sql.Tx
	// BeginTxWithReadCommittedAndConnection(conn string) *sql.Tx
	// BeginTxWithRepeatableReadAndConnection(conn string) *sql.Tx
	BeginTxAndConnection(conn string) *sql.Tx
	// BeginTxWithLevelAndConnection(conn string, level sql.IsolationLevel) *sql.Tx
}

// GetConnectionByDriver 藉由參數(driver = mysql、mssql...)取得Connection(interface)
func GetConnectionByDriver(driver string) Connection {
	switch driver {
	case "mysql":
		// 取得*Mysql(struct)，也屬於Connection(interface)
		return GetMysqlDB()
	case "mssql":
		// 取得*Mssql(struct)，也屬於Connection(interface)
		return GetMssqlDB()
	// case "sqlite":
	// 	return GetSqliteDB()
	// case "postgresql":
	// 	return GetPostgresqlDB()
	default:
		panic("driver not found!")
	}
}

// GetConnectionFromService 將參數srv轉換為Connection(interface)回傳並回傳
func GetConnectionFromService(srv interface{}) Connection {
	if v, ok := srv.(Connection); ok {
		return v
	}
	panic("wrong service")
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

// CheckError 檢查是否有錯誤
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

// GetConnection 透過資料庫引擎取得匹配的Service然後轉換成Connection(interface)類別
func GetConnection(srvs service.List) Connection {
	// service.List類別為map[string]Service，Service是interface(Name方法)
	// 將所有globalCfg.Databases[key]的driver值設置至DatabaseList(map[string]Database).Database.Driver後回傳
	// GetDefault取得預設資料庫DatabaseList["default"]的值
	// Get透過資料庫driver取得Service(interface)，Get裡的參數ex:mysql...
	if v, ok := srvs.Get(config.GetDatabases().GetDefault().Driver).(Connection); ok {
		return v
	}
	panic("wrong service")
}

// GetAggregationExpression 取得資料庫引擎的Aggregation表達式，將參數值加入表達式
func GetAggregationExpression(driver, field, headField, delimiter string) string {
	switch driver {
	// case "postgresql":
	// 	return fmt.Sprintf("string_agg(%s::character varying, '%s') as %s", field, delimiter, headField)
	case "mysql":
		return fmt.Sprintf("group_concat(%s separator '%s') as %s", field, delimiter, headField)
	// case "sqlite":
	// 	return fmt.Sprintf("group_concat(%s, '%s') as %s", field, delimiter, headField)
	case "mssql":
		return fmt.Sprintf("string_agg(%s, '%s') as [%s]", field, delimiter, headField)
	default:
		panic("wrong driver")
	}
}
