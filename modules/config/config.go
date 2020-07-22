package config

import "html/template"

var (
	globalCfg = new(Config)
)

type Config struct {
	// An map supports multi database connection. The first
	// element of Databases is the default connection. See the
	// file connection.go.
	Databases DatabaseList `json:"database,omitempty" yaml:"database,omitempty" ini:"database,omitempty"`

	// The cookie domain used in the auth modules. see
	// the session.go.
	Domain string `json:"domain,omitempty" yaml:"domain,omitempty" ini:"domain,omitempty"`

	// Used to set as the localize language which show in the
	// interface.
	Language string `json:"language,omitempty" yaml:"language,omitempty" ini:"language,omitempty"`

	// The global url prefix.
	UrlPrefix string `json:"prefix,omitempty" yaml:"prefix,omitempty" ini:"prefix,omitempty"`

	// The theme name of template.
	Theme string `json:"theme,omitempty" yaml:"theme,omitempty" ini:"theme,omitempty"`

	// The path where files will be stored into.
	Store Store `json:"store,omitempty" yaml:"store,omitempty" ini:"store,omitempty"`

	// The title of web page.
	Title string `json:"title,omitempty" yaml:"title,omitempty" ini:"title,omitempty"`

	// Logo is the top text in the sidebar.
	Logo template.HTML `json:"logo,omitempty" yaml:"logo,omitempty" ini:"logo,omitempty"`

	// Mini-logo is the top text in the sidebar when folding.
	MiniLogo template.HTML `json:"mini_logo,omitempty" yaml:"mini_logo,omitempty" ini:"mini_logo,omitempty"`

	// The url redirect to after login.
	IndexUrl string `json:"index,omitempty" yaml:"index,omitempty" ini:"index,omitempty"`

	// Login page URL
	LoginUrl string `json:"login_url,omitempty" yaml:"login_url,omitempty" ini:"login_url,omitempty"`

	prefix string

	// Session valid time duration,units are seconds. Default 7200.
	SessionLifeTime int `json:"session_life_time,omitempty" yaml:"session_life_time,omitempty" ini:"session_life_time,omitempty"`

	// Limit login with different IPs
	NoLimitLoginIP bool `json:"no_limit_login_ip,omitempty" yaml:"no_limit_login_ip,omitempty" ini:"no_limit_login_ip,omitempty"`

}

// DatabaseList is a map of Database.
type DatabaseList map[string]Database

type Store struct {
	Path   string `json:"path,omitempty" yaml:"path,omitempty" ini:"path,omitempty"`
	Prefix string `json:"prefix,omitempty" yaml:"prefix,omitempty" ini:"prefix,omitempty"`
}

type Database struct {
	Host       string            `json:"host,omitempty" yaml:"host,omitempty" ini:"host,omitempty"`
	Port       string            `json:"port,omitempty" yaml:"port,omitempty" ini:"port,omitempty"`
	User       string            `json:"user,omitempty" yaml:"user,omitempty" ini:"user,omitempty"`
	Pwd        string            `json:"pwd,omitempty" yaml:"pwd,omitempty" ini:"pwd,omitempty"`
	Name       string            `json:"name,omitempty" yaml:"name,omitempty" ini:"name,omitempty"`
	MaxIdleCon int               `json:"max_idle_con,omitempty" yaml:"max_idle_con,omitempty" ini:"max_idle_con,omitempty"`
	MaxOpenCon int               `json:"max_open_con,omitempty" yaml:"max_open_con,omitempty" ini:"max_open_con,omitempty"`
	Driver     string            `json:"driver,omitempty" yaml:"driver,omitempty" ini:"driver,omitempty"`
	File       string            `json:"file,omitempty" yaml:"file,omitempty" ini:"file,omitempty"`
	Dsn        string            `json:"dsn,omitempty" yaml:"dsn,omitempty" ini:"dsn,omitempty"`
	Params     map[string]string `json:"params,omitempty" yaml:"params,omitempty" ini:"params,omitempty"`
}

// 取得預設資料庫DatabaseList["default"]的值
func (d DatabaseList) GetDefault() Database {
	return d["default"]
}

// 將globalCfg.Databases[key]的driver值設置至DatabaseList(map[string]Database).Database.Driver
func GetDatabases() DatabaseList {
	var list = make(DatabaseList, len(globalCfg.Databases))
	for key := range globalCfg.Databases {
		list[key] = Database{
			Driver: globalCfg.Databases[key].Driver,
		}
	}
	return list
}

// 取得Config.IndexUrl
func (c *Config) Index() string {
	if c.IndexUrl == "" {
		return "/"
	}
	if c.IndexUrl[0] != '/' {
		return "/" + c.IndexUrl
	}
	return c.IndexUrl
}

// 取得Config.prefix
func (c *Config) Prefix() string {
	return c.prefix
}

// 處理Config.IndexUrl(登入後導向的url)後回傳
func (c *Config) GetIndexURL() string {
	// 取得Config.IndexUrl(登入後導向的url)
	index := c.Index()
	if index == "/" {
		return c.Prefix()
	}

	return c.Prefix() + index
}

func (d Database) ParamStr() string {
	p := ""
	if d.Params == nil {
		d.Params = make(map[string]string)
	}
	if d.Driver == "mysql" || d.Driver == "sqlite" {
		if d.Driver == "mysql" {
			if _, ok := d.Params["charset"]; !ok {
				d.Params["charset"] = "utf8mb4"
			}
		}
		if len(d.Params) > 0 {
			p = "?"
			for k, v := range d.Params {
				p += k + "=" + v + "&"
			}
			p = p[:len(p)-1]
		}
	}
	// if d.Driver == "mssql" {
	// 	if _, ok := d.Params["encrypt"]; !ok {
	// 		d.Params["encrypt"] = "disable"
	// 	}
	// 	for k, v := range d.Params {
	// 		p += k + "=" + v + ";"
	// 	}
	// 	p = p[:len(p)-1]
	// }
	// if d.Driver == "postgresql" {
	// 	if _, ok := d.Params["sslmode"]; !ok {
	// 		d.Params["sslmode"] = "disable"
	// 	}
	// 	p = " "
	// 	for k, v := range d.Params {
	// 		p += k + "=" + v + " "
	// 	}
	// 	p = p[:len(p)-1]
	// }
	return p
}

func GetDomain() string {
	return globalCfg.Domain
}

func GetNoLimitLoginIP() bool {
	return globalCfg.NoLimitLoginIP
}

func GetSessionLifeTime() int {
	return globalCfg.SessionLifeTime
}