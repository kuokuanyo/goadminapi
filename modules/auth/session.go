package auth

import (
	"encoding/json"
	"goadminapi/context"
	"goadminapi/modules/db"
	"goadminapi/modules/db/dialect"
	"net/http"
	"strconv"
	"time"

	"goadminapi/modules/config"

	"goadminapi/plugins/admin/modules"
)

// 使用資料庫當作持久性的驅動程式
// DBDriver也是PersistenceDriver(interface)
// 紀錄goadmin_session資料表(觀察用戶登入狀態)
type DBDriver struct {
	conn      db.Connection
	tableName string
}

// Session contains info of session.
type Session struct {
	Expires time.Duration //cookie存在時間
	Cookie  string
	Values  map[string]interface{}
	Driver  PersistenceDriver
	Sid     string
	Context *context.Context
}

// Config wraps the Session info.
type Config struct {
	Expires time.Duration
	Cookie  string
}

// 儲存與獲取session資訊
type PersistenceDriver interface {
	Load(string) (map[string]interface{}, error)
	Update(sid string, values map[string]interface{}) error
}

func (driver *DBDriver) table() *db.SQL {
	return db.Table(driver.tableName).WithDriver(driver.conn)
}

// DBDriver也是PersistenceDriver(interface)
// 尋找資料表中符合參數(sid)的user資料，將資料表values欄位值(ex:{"user_id":1})JSON解碼並回傳values
func (driver *DBDriver) Load(sid string) (map[string]interface{}, error) {
	// table取得sql(struct)
	// 取得user(透過sid尋找符合資料)
	sesModel, err := driver.table().Where("sid", "=", sid).First()

	// 檢查錯誤
	if db.CheckError(err, db.QUERY) {
		return nil, err
	}

	//如果沒有找到符合的資料
	if sesModel == nil {
		return map[string]interface{}{}, nil
	}

	var values map[string]interface{}
	err = json.Unmarshal([]byte(sesModel["values"].(string)), &values)

	//回傳goadmin_session表中的values欄位
	return values, err
}

// 刪除超過時間的cookie(session)
func (driver *DBDriver) deleteOverdueSession() {

	defer func() {
		if err := recover(); err != nil {
			panic(err)
		}
	}()

	var (
		duration = strconv.Itoa(config.GetSessionLifeTime() + 1000)
		// GetDatabases將globalCfg.Databases的driver值設置至DatabaseList(map[string]Database)
		// GetDefault取得預設資料庫DatabaseList["default"]的值
		driverName = config.GetDatabases().GetDefault().Driver
		raw        = ``
	)

	if "postgresql" == driverName {
		raw = `extract(epoch from now()) - ` + duration + ` > extract(epoch from created_at)`
	} else if "mysql" == driverName {
		raw = `unix_timestamp(created_at) < unix_timestamp() - ` + duration
	} else if "sqlite" == driverName {
		raw = `strftime('%s', created_at) < strftime('%s', 'now') - ` + duration
	} else if "mssql" == driverName {
		raw = `DATEDIFF(second, [created_at], GETDATE()) > ` + duration
	}

	if raw != "" {
		// WhereRaw將參數raw設置至SQL(struct)
		// Delete刪除資料
		_ = driver.table().WhereRaw(raw).Delete()
	}
}

// 刪除逾時的cookie，尋找符合參數sid的資料，如果沒有符合的資料則將sesValue(處理過後的參數values)與sid加入資料表，如有符合的資料則是更新sesValue的值
func (driver *DBDriver) Update(sid string, values map[string]interface{}) error {
	// deleteOverdueSession刪除超過時間的cookie(session)
	go driver.deleteOverdueSession()

	if sid != "" {
		if len(values) == 0 {
			// 刪除資料
			err := driver.table().Where("sid", "=", sid).Delete()
			if db.CheckError(err, db.DELETE) {
				return err
			}
		}
		valuesByte, err := json.Marshal(values)
		if err != nil {
			return err
		}
		sesValue := string(valuesByte)
		// 尋找符合的資料
		sesModel, _ := driver.table().Where("sid", "=", sid).First()
		if sesModel == nil {
			// GetNoLimitLoginIP回傳globalCfg.NoLimitLoginIP
			if !config.GetNoLimitLoginIP() {
				// 刪除資料
				err = driver.table().Where("values", "=", sesValue).Delete()
				if db.CheckError(err, db.DELETE) {
					return err
				}
			}
			// 將sesValue與sid加入資料表
			_, err := driver.table().Insert(dialect.H{
				"values": sesValue,
				"sid":    sid,
			})
			if db.CheckError(err, db.INSERT) {
				return err
			}
		} else {
			// 更新資料表Update(dialect.H{"values": sesValue,})的值
			_, err := driver.table().
				Where("sid", "=", sid).
				Update(dialect.H{
					"values": sesValue,
				})
			if db.CheckError(err, db.UPDATE) {
				return err
			}
		}
	}
	return nil
}

// 將參數(conn)設置並回傳DBDriver(struct)
func newDBDriver(conn db.Connection) *DBDriver {
	return &DBDriver{
		conn: conn,
		// 資料表名(取得目前有哪些用戶登入)
		tableName: "session",
	}
}

// 設置Session(struct)資訊並取得cookie及設置cookie值
func InitSession(ctx *context.Context, conn db.Connection) (*Session, error) {

	sessions := new(Session)
	// 更新Session(struct)的Expires(時間)與Cookie
	sessions.UpdateConfig(Config{
		Expires: time.Second * time.Duration(config.GetSessionLifeTime()),
		Cookie:  "session",
	})

	// UseDriver透過參數(newDBDriver(conn))設置Session.Driver
	// newDBDriver透過參數(conn)回傳DBDriver(struct)
	sessions.UseDriver(newDBDriver(conn))
	sessions.Values = make(map[string]interface{})

	// 取得cookie並設置值，接著設定Session(struct)資訊，將參數ctx設置至Session.Context
	return sessions.StartCtx(ctx)
}

// 透過參數(driver)設置Session.Driver
func (ses *Session) UseDriver(driver PersistenceDriver) {
	ses.Driver = driver
}

// 更新Session(struct)的Expires(時間)與Cookie
func (ses *Session) UpdateConfig(config Config) {
	ses.Expires = config.Expires
	ses.Cookie = config.Cookie
}

// 尋找資料表中符合參數(sesKey)的user資料，回傳user_id
func GetSessionByKey(sesKey, key string, conn db.Connection) (interface{}, error) {
	// newDBDriver將參數(conn)設置並回傳DBDriver(struct)
	// 尋找資料表中符合參數(sesKey)的user資料，將資料表values欄位值(ex:{"user_id":1})JSON解碼並回傳values
	m, err := newDBDriver(conn).Load(sesKey)
	return m[key], err
}

// 取得cookie並設置值，接著設定Session(struct)資訊，將參數ctx設置至Session.Context
func (ses *Session) StartCtx(ctx *context.Context) (*Session, error) {
	if cookie, err := ctx.Request.Cookie(ses.Cookie); err == nil && cookie.Value != "" {
		ses.Sid = cookie.Value
		// 尋找資料表中符合參數(cookie.Value)的user資料，將資料的values欄位值JSON解碼並回傳values
		valueFromDriver, err := ses.Driver.Load(cookie.Value)
		if err != nil {
			return nil, err
		}
		if len(valueFromDriver) > 0 {
			ses.Values = valueFromDriver
		}
	} else {
		//給Session.Sid一組uuid
		ses.Sid = modules.Uuid()
	}
	ses.Context = ctx
	return ses, nil
}

// 藉由參數(key)取得Session.Values[key]
func (ses *Session) Get(key string) interface{} {
	return ses.Values[key]
}

// 將參數key、value加入Session.Values後檢查是否有符合Session.Sid的資料，判斷插入或是更新資料
// 最後設置cookie(struct)並儲存在response header Set-Cookie中
func (ses *Session) Add(key string, value interface{}) error {
	ses.Values[key] = value

	// 刪除逾時的cookie，尋找符合參數sid的資料
	// 如果沒有符合的資料則將sesValue(處理過後的參數values)與sid加入資料表
	// 如有符合的資料則是更新sesValue的值
	if err := ses.Driver.Update(ses.Sid, ses.Values); err != nil {
		return err
	}
	cookie := http.Cookie{
		Name:  ses.Cookie,
		Value: ses.Sid,
		// 回傳globalCfg.SessionLifeTime
		MaxAge: config.GetSessionLifeTime(),
		// cookie存在時間
		Expires:  time.Now().Add(ses.Expires),
		HttpOnly: true,
		Path:     "/",
	}
	if config.GetDomain() != "" {
		cookie.Domain = config.GetDomain()
	}

	// 設置cookie(struct)在response header Set-Cookie中
	ses.Context.SetCookie(&cookie)
	return nil
}
