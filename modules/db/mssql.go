package db

import (
	"database/sql"
	"fmt"
	"goadminapi/modules/config"
	"regexp"
	"strconv"
	"strings"
)

// Mssql is a Connection of mssql.
type Mssql struct {
	Base
}

func GetMssqlDB() *Mssql {
	return &Mssql{
		Base: Base{
			DbList: make(map[string]*sql.DB),
		},
	}
}

func (db *Mssql) GetDelimiter() string {
	return "["
}

// Name implements the method Connection.Name.
func (db *Mssql) Name() string {
	return "mssql"
}

func replaceStringFunc(pattern, src string, rpl func(s string) string) (string, error) {

	r, err := regexp.Compile(pattern)
	if err != nil {
		return "", err
	}

	bytes := r.ReplaceAllFunc([]byte(src), func(bytes []byte) []byte {
		return []byte(rpl(string(bytes)))
	})

	return string(bytes), nil
}

func replace(pattern string, replace, src []byte) ([]byte, error) {

	r, err := regexp.Compile(pattern)
	if err != nil {
		return nil, err
	}

	return r.ReplaceAll(src, replace), nil
}

func replaceString(pattern, rep, src string) (string, error) {
	r, e := replace(pattern, []byte(rep), []byte(src))
	return string(r), e
}

func matchAllString(pattern string, src string) ([][]string, error) {
	r, err := regexp.Compile(pattern)
	if err != nil {
		return nil, err
	}
	return r.FindAllStringSubmatch(src, -1), nil
}

func isMatch(pattern string, src []byte) bool {
	r, err := regexp.Compile(pattern)
	if err != nil {
		return false
	}
	return r.Match(src)
}

func isMatchString(pattern string, src string) bool {
	return isMatch(pattern, []byte(src))
}

func matchString(pattern string, src string) ([]string, error) {
	r, err := regexp.Compile(pattern)
	if err != nil {
		return nil, err
	}
	return r.FindStringSubmatch(src), nil
}

func (db *Mssql) handleSqlBeforeExec(query string) string {
	index := 0
	str, _ := replaceStringFunc("\\?", query, func(s string) string {
		index++
		return fmt.Sprintf("@p%d", index)
	})

	str, _ = replaceString("\"", "", str)

	return db.parseSql(str)
}

func (db *Mssql) parseSql(sql string) string {

	patten := `^\s*(?i)(SELECT)|(LIMIT\s*(\d+)\s*,\s*(\d+))`
	if isMatchString(patten, sql) == false {
		return sql
	}

	res, err := matchAllString(patten, sql)
	if err != nil {
		return ""
	}

	index := 0
	keyword := strings.TrimSpace(res[index][0])
	keyword = strings.ToUpper(keyword)

	index++
	switch keyword {
	case "SELECT":
		if len(res) < 2 || (strings.HasPrefix(res[index][0], "LIMIT") == false && strings.HasPrefix(res[index][0], "limit") == false) {
			break
		}

		if isMatchString("((?i)SELECT)(.+)((?i)LIMIT)", sql) == false {
			break
		}

		selectStr := ""
		orderbyStr := ""
		haveOrderby := isMatchString("((?i)SELECT)(.+)((?i)ORDER BY)", sql)
		if haveOrderby {
			queryExpr, _ := matchString("((?i)SELECT)(.+)((?i)ORDER BY)", sql)

			if len(queryExpr) != 4 || strings.EqualFold(queryExpr[1], "SELECT") == false || strings.EqualFold(queryExpr[3], "ORDER BY") == false {
				break
			}
			selectStr = queryExpr[2]

			orderbyExpr, _ := matchString("((?i)ORDER BY)(.+)((?i)LIMIT)", sql)
			if len(orderbyExpr) != 4 || strings.EqualFold(orderbyExpr[1], "ORDER BY") == false || strings.EqualFold(orderbyExpr[3], "LIMIT") == false {
				break
			}
			orderbyStr = orderbyExpr[2]
		} else {
			queryExpr, _ := matchString("((?i)SELECT)(.+)((?i)LIMIT)", sql)
			if len(queryExpr) != 4 || strings.EqualFold(queryExpr[1], "SELECT") == false || strings.EqualFold(queryExpr[3], "LIMIT") == false {
				break
			}
			selectStr = queryExpr[2]
		}

		first, limit := 0, 0
		for i := 1; i < len(res[index]); i++ {
			if len(strings.TrimSpace(res[index][i])) == 0 {
				continue
			}

			if strings.HasPrefix(res[index][i], "LIMIT") || strings.HasPrefix(res[index][i], "limit") {
				first, _ = strconv.Atoi(res[index][i+1])
				limit, _ = strconv.Atoi(res[index][i+2])
				break
			}
		}

		if haveOrderby {
			sql = fmt.Sprintf("SELECT * FROM (SELECT ROW_NUMBER() OVER (ORDER BY %s) as ROWNUMBER_, %s   ) as TMP_ WHERE TMP_.ROWNUMBER_ > %d AND TMP_.ROWNUMBER_ <= %d", orderbyStr, selectStr, first, limit)
		} else {
			if first == 0 {
				first = limit
			} else {
				first = limit - first
			}
			sql = fmt.Sprintf("SELECT * FROM (SELECT TOP %d * FROM (SELECT TOP %d %s) as TMP1_ ) as TMP2_ ", first, limit, selectStr)
		}
	default:
	}
	return sql
}

// InitDB implements the method Connection.InitDB.
func (db *Mssql) InitDB(cfglist map[string]config.Database) Connection {
	db.Once.Do(func() {
		for conn, cfg := range cfglist {

			if cfg.Dsn == "" {
				cfg.Dsn = fmt.Sprintf("user id=%s;password=%s;server=%s;port=%s;database=%s;"+cfg.ParamStr(),
					cfg.User, cfg.Pwd, cfg.Host, cfg.Port, cfg.Name)
			}

			sqlDB, err := sql.Open("sqlserver", cfg.Dsn)

			if sqlDB == nil {
				panic("invalid connection")
			}

			if err != nil {
				_ = sqlDB.Close()
				panic(err.Error())
			} else {
				sqlDB.SetMaxIdleConns(cfg.MaxIdleCon)
				sqlDB.SetMaxOpenConns(cfg.MaxOpenCon)

				db.DbList[conn] = sqlDB
			}

			if err := sqlDB.Ping(); err != nil {
				panic(err)
			}
		}
	})
	return db
}

// -------------connection(interface)的所有方法--------------------------
// Query implements the method Connection.Query.
func (db *Mssql) Query(query string, args ...interface{}) ([]map[string]interface{}, error) {
	query = db.handleSqlBeforeExec(query)
	return CommonQuery(db.DbList["default"], query, args...)
}

// Exec implements the method Connection.Exec.
func (db *Mssql) Exec(query string, args ...interface{}) (sql.Result, error) {
	query = db.handleSqlBeforeExec(query)
	return CommonExec(db.DbList["default"], query, args...)
}

// QueryWithConnection implements the method Connection.QueryWithConnection.
func (db *Mssql) QueryWithConnection(con string, query string, args ...interface{}) ([]map[string]interface{}, error) {
	query = db.handleSqlBeforeExec(query)
	return CommonQuery(db.DbList[con], query, args...)
}

// ExecWithConnection implements the method Connection.ExecWithConnection.
func (db *Mssql) ExecWithConnection(con string, query string, args ...interface{}) (sql.Result, error) {
	query = db.handleSqlBeforeExec(query)
	return CommonExec(db.DbList[con], query, args...)
}

// QueryWithTx is query method within the transaction.
func (db *Mssql) QueryWithTx(tx *sql.Tx, query string, args ...interface{}) ([]map[string]interface{}, error) {
	query = db.handleSqlBeforeExec(query)
	return CommonQueryWithTx(tx, query, args...)
}

// ExecWithTx is exec method within the transaction.
func (db *Mssql) ExecWithTx(tx *sql.Tx, query string, args ...interface{}) (sql.Result, error) {
	query = db.handleSqlBeforeExec(query)
	return CommonExecWithTx(tx, query, args...)
}

// BeginTxAndConnection starts a transaction with level LevelDefault and connection.
func (db *Mssql) BeginTxAndConnection(conn string) *sql.Tx {
	return CommonBeginTxWithLevel(db.DbList[conn], sql.LevelDefault)
}

// -------------connection(interface)的所有方法--------------------------
