package db

import (
	dbsql "database/sql"
	"errors"
	"fmt"
	"goadminapi/modules/db/dialect"
	"regexp"
	"strings"
	"sync"
)

// 包裝Connection、driver與dialect方法
type SQL struct {
	dialect.SQLComponent //sql過濾條件
	diver                Connection
	dialect              dialect.Dialect //sql CRUD等方法(不同資料庫引擎的方法)
	conn                 string
	tx                   *dbsql.Tx
}

// 回傳sql元件
var SQLPool = sync.Pool{
	New: func() interface{} {
		return &SQL{
			SQLComponent: dialect.SQLComponent{
				Fields:     make([]string, 0),
				TableName:  "",
				Args:       make([]interface{}, 0),
				Wheres:     make([]dialect.Where, 0),
				Leftjoins:  make([]dialect.Join, 0),
				UpdateRaws: make([]dialect.RawUpdate, 0),
				WhereRaws:  "",
				Order:      "",
				Group:      "",
				Limit:      "",
			},
			diver:   nil,
			dialect: nil,
		}
	},
}

// 取得新的SQL(struct)
func newSQL() *SQL {
	return SQLPool.Get().(*SQL)
}

// 將SQL(struct)資訊清除後將參數設置至SQL.TableName回傳
func (sql *SQL) Table(table string) *SQL {
	sql.clean()
	sql.TableName = table
	return sql
}

// 將參數(table)設置並回傳sql(struct)
func Table(table string) *SQL {
	sql := newSQL()
	sql.TableName = table //sql.dialect.SQLComponent.TableName
	sql.conn = "default"
	return sql
}

// 將參數設置(conn)並回傳sql(struct)
func WithDriver(conn Connection) *SQL {
	sql := newSQL()
	sql.diver = conn
	sql.dialect = dialect.GetDialectByDriver(conn.Name())
	sql.conn = "default"
	return sql
}

// 返回sql struct藉由給定的conn
func (sql *SQL) WithDriver(conn Connection) *SQL {
	sql.diver = conn
	//GetDialectByDriver 不同資料庫引擎有不同的使用符號
	sql.dialect = dialect.GetDialectByDriver(conn.Name())
	return sql
}

// 將參數設置(connName、conn)並回傳sql(struct)
func WithDriverAndConnection(connName string, conn Connection) *SQL {
	sql := newSQL()
	sql.diver = conn
	sql.dialect = dialect.GetDialectByDriver(conn.Name())
	sql.conn = connName
	return sql
}

// 取得所有欄位資訊
func (sql *SQL) ShowColumns() ([]map[string]interface{}, error) {
	defer RecycleSQL(sql)
	return sql.diver.QueryWithConnection(sql.conn, sql.dialect.ShowColumns(sql.TableName))
}

// 將參數設置至SQL(struct).Fields並且設置SQL(struct).Functions
func (sql *SQL) Select(fields ...string) *SQL {
	sql.Fields = fields
	sql.Functions = make([]string, len(fields))
	reg, _ := regexp.Compile("(.*?)\\((.*?)\\)")
	for k, field := range fields {
		res := reg.FindAllStringSubmatch(field, -1)
		if len(res) > 0 && len(res[0]) > 2 {
			sql.Functions[k] = res[0][1]
			sql.Fields[k] = res[0][2]
		}
	}
	return sql
}

// 藉由id取的符合資料
func (sql *SQL) Find(arg interface{}) (map[string]interface{}, error) {
	return sql.Where("id", "=", arg).First()
}

// 插入給定的參數資料(values(map[string]interface{}))後，最後回傳加入值的id
func (sql *SQL) Insert(values dialect.H) (int64, error) {
	// 清空的sql 資訊放入SQLPool中
	defer RecycleSQL(sql)

	// 新增頁面中設定的數值(ex:map[http_method:GET http_path:s name:ssssssssss slug:ssssssssss])
	sql.Values = values

	sql.dialect.Insert(&sql.SQLComponent)

	var (
		res    dbsql.Result
		err    error
		resMap []map[string]interface{}
	)

	// postgresql引擎才會執行
	if sql.diver.Name() == "postgresql" {
		if sql.TableName == "menu" ||
			sql.TableName == "permissions" ||
			sql.TableName == "roles" ||
			sql.TableName == "users" {

			if sql.tx != nil {
				resMap, err = sql.diver.QueryWithTx(sql.tx, sql.Statement+" RETURNING id", sql.Args...)
			} else {
				resMap, err = sql.diver.QueryWithConnection(sql.conn, sql.Statement+" RETURNING id", sql.Args...)
			}

			if err != nil {
				return 0, err
			}

			if len(resMap) == 0 {
				return 0, errors.New("no affect row")
			}
			return resMap[0]["id"].(int64), nil
		}
	}

	if sql.tx != nil {
		// QueryWithTx是transaction的執行方法
		res, err = sql.diver.ExecWithTx(sql.tx, sql.Statement, sql.Args...)
	} else {
		// ExecWithConnection有給定連接(conn)名稱
		res, err = sql.diver.ExecWithConnection(sql.conn, sql.Statement, sql.Args...)
	}
	if err != nil {
		return 0, err
	}

	if affectRow, _ := res.RowsAffected(); affectRow < 1 {
		return 0, errors.New("no affect row")
	}

	return res.LastInsertId()
}

func (sql *SQL) Update(values dialect.H) (int64, error) {
	defer RecycleSQL(sql)

	sql.Values = values

	sql.dialect.Update(&sql.SQLComponent)

	var (
		res dbsql.Result
		err error
	)

	if sql.tx != nil {
		res, err = sql.diver.ExecWithTx(sql.tx, sql.Statement, sql.Args...)
	} else {
		res, err = sql.diver.ExecWithConnection(sql.conn, sql.Statement, sql.Args...)
	}

	if err != nil {
		return 0, err
	}

	if affectRow, _ := res.RowsAffected(); affectRow < 1 {
		return 0, errors.New("no affect row")
	}

	return res.LastInsertId()
}

func (sql *SQL) Delete() error {
	defer RecycleSQL(sql)

	sql.dialect.Delete(&sql.SQLComponent)

	var (
		res dbsql.Result
		err error
	)

	if sql.tx != nil {
		res, err = sql.diver.ExecWithTx(sql.tx, sql.Statement, sql.Args...)
	} else {
		res, err = sql.diver.ExecWithConnection(sql.conn, sql.Statement, sql.Args...)
	}

	if err != nil {
		return err
	}

	if affectRow, _ := res.RowsAffected(); affectRow < 1 {
		return errors.New("no affect row")
	}

	return nil
}

// 返回所有符合查詢的結果
func (sql *SQL) All() ([]map[string]interface{}, error) {
	//最後清空sql資訊
	defer RecycleSQL(sql)

	sql.dialect.Select(&sql.SQLComponent)

	if sql.tx != nil {
		return sql.diver.QueryWithTx(sql.tx, sql.Statement, sql.Args...)
	}
	return sql.diver.QueryWithConnection(sql.conn, sql.Statement, sql.Args...)
}

// 回傳第一筆符合的資料
func (sql *SQL) First() (map[string]interface{}, error) {
	// 執行結束後清空sql資訊
	defer RecycleSQL(sql)

	//尋找資料
	sql.dialect.Select(&sql.SQLComponent)

	var (
		res []map[string]interface{}
		err error
	)

	//假設有tx在tx中執行查詢，反之一般資料庫執行
	if sql.tx != nil {
		res, err = sql.diver.QueryWithTx(sql.tx, sql.Statement, sql.Args...)
	} else {
		res, err = sql.diver.QueryWithConnection(sql.conn, sql.Statement, sql.Args...)
	}

	if err != nil {
		return nil, err
	}

	if len(res) < 1 {
		return nil, errors.New("out of index")
	}
	return res[0], nil
}

// sql 語法 where = ...，回傳 SQl struct
// 將值設置至SQL.Wheres、SQL.Args
func (sql *SQL) Where(field string, operation string, arg interface{}) *SQL {
	sql.Wheres = append(sql.Wheres, dialect.Where{
		Field:     field,     // 欄位名稱
		Operation: operation, // 符號
		Qmark:     "?",
	})
	sql.Args = append(sql.Args, arg)
	return sql
}

// 將參數raw、args設置至SQL(struct)
func (sql *SQL) WhereRaw(raw string, args ...interface{}) *SQL {
	sql.WhereRaws = raw
	sql.Args = append(sql.Args, args...)
	return sql
}

// where多個數值，ex where id IN (1,2,3,4);
func (sql *SQL) WhereIn(field string, arg []interface{}) *SQL {
	if len(arg) == 0 {
		panic("wrong parameter")
	}
	sql.Wheres = append(sql.Wheres, dialect.Where{
		Field:     field,
		Operation: "in",
		Qmark:     "(" + strings.Repeat("?,", len(arg)-1) + "?)",
	})
	sql.Args = append(sql.Args, arg...)
	return sql
}

// LeftJoin
func (sql *SQL) LeftJoin(table string, fieldA string, operation string, fieldB string) *SQL {
	sql.Leftjoins = append(sql.Leftjoins, dialect.Join{
		FieldA:    fieldA,
		FieldB:    fieldB,
		Table:     table,
		Operation: operation,
	})
	return sql
}

func (sql *SQL) wrap(field string) string {
	if sql.diver.Name() == "mssql" {
		return fmt.Sprintf(`[%s]`, field)
	}
	return sql.diver.GetDelimiter() + field + sql.diver.GetDelimiter()
}

// OrderBy set order fields.
func (sql *SQL) OrderBy(fields ...string) *SQL {
	if len(fields) == 0 {
		panic("wrong order field")
	}
	for i := 0; i < len(fields); i++ {
		if i == len(fields)-2 {
			sql.Order += " " + sql.wrap(fields[i]) + " " + fields[i+1]
			return sql
		}
		sql.Order += " " + sql.wrap(fields[i]) + " and "
	}
	return sql
}

// 將SQL(struct)資訊清除
func (sql *SQL) clean() {
	sql.Functions = make([]string, 0)
	sql.Group = ""
	sql.Values = make(map[string]interface{})
	sql.Fields = make([]string, 0)
	sql.TableName = ""
	sql.Wheres = make([]dialect.Where, 0)
	sql.Leftjoins = make([]dialect.Join, 0)
	sql.Args = make([]interface{}, 0)
	sql.Order = ""
	sql.Offset = ""
	sql.Limit = ""
	sql.WhereRaws = ""
	sql.UpdateRaws = make([]dialect.RawUpdate, 0)
	sql.Statement = ""
}

// 清空的sql 資訊放入SQLPool中
func RecycleSQL(sql *SQL) {
	// sql資訊清除
	sql.clean()

	sql.conn = ""
	sql.diver = nil
	sql.tx = nil
	sql.dialect = nil

	//清空的sql 資訊放入SQLPool中
	SQLPool.Put(sql)
}
