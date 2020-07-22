package db

import (
	"database/sql"
	"regexp"
	"strings"

	"goadminapi/modules/config"
)

type Mysql struct {
	Base
}

// Connection的方法
func GetMysqlDB() *Mysql {
	return &Mysql{
		Base: Base{
			DbList: make(map[string]*sql.DB),
		},
	}
}

func (db *Mysql) Name() string {
	return "mysql"
}

func (db *Mysql) GetDelimiter() string {
	return "`"
}

// 初始化資料庫連線並啟動引擎
func (db *Mysql) InitDB(cfgs map[string]config.Database) Connection {
	db.Once.Do(func() {
		for conn, cfg := range cfgs {

			if cfg.Dsn == "" {
				cfg.Dsn = cfg.User + ":" + cfg.Pwd + "@tcp(" + cfg.Host + ":" + cfg.Port + ")/" +
					cfg.Name + cfg.ParamStr()
			}

			sqlDB, err := sql.Open("mysql", cfg.Dsn)

			if err != nil {
				if sqlDB != nil {
					_ = sqlDB.Close()
				}
				panic(err)
			} else {
				// Largest set up the database connection reduce time wait
				sqlDB.SetMaxIdleConns(cfg.MaxIdleCon)
				sqlDB.SetMaxOpenConns(cfg.MaxOpenCon)

				db.DbList[conn] = sqlDB
			}
			//啟動資料庫引擎
			if err := sqlDB.Ping(); err != nil {
				panic(err)
			}
		}
	})
	return db
}

// 沒有給定連接(conn)名稱，透過參數查詢db.DbList["default"]資料並回傳
func (db *Mysql) Query(query string, args ...interface{}) ([]map[string]interface{}, error) {
	// CommonQuery查詢資料並回傳
	return CommonQuery(db.DbList["default"], query, args...)
}

// 沒有給定連接(conn)名稱
func (db *Mysql) Exec(query string, args ...interface{}) (sql.Result, error) {
	return CommonExec(db.DbList["default"], query, args...)
}

// 有給定參數連接(conn)名稱，透過參數con查詢db.DbList[con]資料並回傳
func (db *Mysql) QueryWithConnection(con string, query string, args ...interface{}) ([]map[string]interface{}, error) {
	// CommonQuery查詢資料並回傳
	return CommonQuery(db.DbList[con], query, args...)
}

// 有給定連接(conn)名稱，透過參數con執行db.DbList[con]資料並回傳
func (db *Mysql) ExecWithConnection(con string, query string, args ...interface{}) (sql.Result, error) {
	return CommonExec(db.DbList[con], query, args...)
}

// QueryWithTx是sql.Tx的查詢方法(與CommonQuery一樣)
func (db *Mysql) QueryWithTx(tx *sql.Tx, query string, args ...interface{}) ([]map[string]interface{}, error) {
	return CommonQueryWithTx(tx, query, args...)
}

// QueryWithTx是sql.Tx的執行方法(與CommonExec一樣)
func (db *Mysql) ExecWithTx(tx *sql.Tx, query string, args ...interface{}) (sql.Result, error) {
	return CommonExecWithTx(tx, query, args...)
}

// 與CommonQuery一樣(差別在tx執行)
func CommonQueryWithTx(tx *sql.Tx, query string, args ...interface{}) ([]map[string]interface{}, error) {

	rs, err := tx.Query(query, args...)

	if err != nil {
		panic(err)
	}

	defer func() {
		if rs != nil {
			_ = rs.Close()
		}
	}()

	col, colErr := rs.Columns()

	if colErr != nil {
		return nil, colErr
	}

	typeVal, err := rs.ColumnTypes()
	if err != nil {
		return nil, err
	}

	results := make([]map[string]interface{}, 0)

	r, _ := regexp.Compile(`\\((.*)\\)`)
	for rs.Next() {
		var colVar = make([]interface{}, len(col))
		for i := 0; i < len(col); i++ {
			typeName := strings.ToUpper(r.ReplaceAllString(typeVal[i].DatabaseTypeName(), ""))
			SetColVarType(&colVar, i, typeName)
		}
		result := make(map[string]interface{})
		if scanErr := rs.Scan(colVar...); scanErr != nil {
			return nil, scanErr
		}
		for j := 0; j < len(col); j++ {
			typeName := strings.ToUpper(r.ReplaceAllString(typeVal[j].DatabaseTypeName(), ""))
			SetResultValue(&result, col[j], colVar[j], typeName)
		}
		results = append(results, result)
	}
	if err := rs.Err(); err != nil {
		return nil, err
	}
	return results, nil
}

// 與CommonExec一樣(差別在tx執行)
func CommonExecWithTx(tx *sql.Tx, query string, args ...interface{}) (sql.Result, error) {
	rs, err := tx.Exec(query, args...)
	if err != nil {
		return nil, err
	}
	return rs, nil
}

// 查詢資料並回傳
func CommonQuery(db *sql.DB, query string, args ...interface{}) ([]map[string]interface{}, error) {

	//查詢
	rs, err := db.Query(query, args...)

	if err != nil {
		panic(err)
	}

	//最後關閉 *sql.rows
	defer func() {
		if rs != nil {
			_ = rs.Close()
		}
	}()

	//取得欄位名稱
	col, colErr := rs.Columns()

	if colErr != nil {
		return nil, colErr
	}

	// 取得欄位類別
	typeVal, err := rs.ColumnTypes()
	if err != nil {
		return nil, err
	}

	// TODO: regular expressions for sqlite, use the dialect module
	// tell the drive to reduce the performance loss
	results := make([]map[string]interface{}, 0)

	r, _ := regexp.Compile(`\\((.*)\\)`)
	for rs.Next() {
		var colVar = make([]interface{}, len(col))
		//typeName欄位類別名稱
		for i := 0; i < len(col); i++ {
			typeName := strings.ToUpper(r.ReplaceAllString(typeVal[i].DatabaseTypeName(), ""))
			//converter.go中
			//SetColVarType 設定欄位數值類型
			SetColVarType(&colVar, i, typeName)
		}
		result := make(map[string]interface{})
		if scanErr := rs.Scan(colVar...); scanErr != nil {
			return nil, scanErr
		}
		for j := 0; j < len(col); j++ {
			typeName := strings.ToUpper(r.ReplaceAllString(typeVal[j].DatabaseTypeName(), ""))
			// converter.go中
			SetResultValue(&result, col[j], colVar[j], typeName)
		}
		results = append(results, result)
	}
	if err := rs.Err(); err != nil {
		return nil, err
	}
	return results, nil
}

// 執行sql命令
func CommonExec(db *sql.DB, query string, args ...interface{}) (sql.Result, error) {

	rs, err := db.Exec(query, args...)
	if err != nil {
		return nil, err
	}
	return rs, nil
}
