package config

import (
	"fmt"
	"goadminapi/modules/utils"
	"html/template"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
)

var (
	globalCfg  = new(Config)
	declare    sync.Once
	updateLock sync.Mutex
	lock       sync.Mutex
	count      uint32
)

// 頁面動畫
type PageAnimation struct {
	Type     string  `json:"type,omitempty" yaml:"type,omitempty" ini:"type,omitempty"`
	Duration float32 `json:"duration,omitempty" yaml:"duration,omitempty" ini:"duration,omitempty"`
	Delay    float32 `json:"delay,omitempty" yaml:"delay,omitempty" ini:"delay,omitempty"`
}

type Config struct {
	// An map supports multi database connection. The first
	// element of Databases is the default connection. See the
	// file connection.go.
	Databases DatabaseList `json:"database,omitempty" yaml:"database,omitempty" ini:"database,omitempty"`

	// The cookie domain used in the auth modules. see
	// the session.go.
	Domain string `json:"domain,omitempty" yaml:"domain,omitempty" ini:"domain,omitempty"`

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

	// Assets visit link.
	AssetUrl string `json:"asset_url,omitempty" yaml:"asset_url,omitempty" ini:"asset_url,omitempty"`

	// File upload engine,default "local"
	FileUploadEngine FileUploadEngine `json:"file_upload_engine,omitempty" yaml:"file_upload_engine,omitempty" ini:"file_upload_engine,omitempty"`

	// Session valid time duration,units are seconds. Default 7200.
	SessionLifeTime int `json:"session_life_time,omitempty" yaml:"session_life_time,omitempty" ini:"session_life_time,omitempty"`

	// Env is the environment,which maybe local,test,prod.
	Env string `json:"env,omitempty" yaml:"env,omitempty" ini:"env,omitempty"`

	// Page animation
	Animation PageAnimation `json:"animation,omitempty" yaml:"animation,omitempty" ini:"animation,omitempty"`

	// Limit login with different IPs
	NoLimitLoginIP bool `json:"no_limit_login_ip,omitempty" yaml:"no_limit_login_ip,omitempty" ini:"no_limit_login_ip,omitempty"`

	// Custom html in the tag head.
	CustomHeadHtml template.HTML `json:"custom_head_html,omitempty" yaml:"custom_head_html,omitempty" ini:"custom_head_html,omitempty"`

	// Custom html after body.
	CustomFootHtml template.HTML `json:"custom_foot_html,omitempty" yaml:"custom_foot_html,omitempty" ini:"custom_foot_html,omitempty"`

	// Footer Info html
	FooterInfo template.HTML `json:"footer_info,omitempty" yaml:"footer_info,omitempty" ini:"footer_info,omitempty"`

	// Login page title
	LoginTitle string `json:"login_title,omitempty" yaml:"login_title,omitempty" ini:"login_title,omitempty"`

	// Login page logo
	LoginLogo template.HTML `json:"login_logo,omitempty" yaml:"login_logo,omitempty" ini:"login_logo,omitempty"`

	// Auth user table
	AuthUserTable string `json:"auth_user_table,omitempty" yaml:"auth_user_table,omitempty" ini:"auth_user_table,omitempty"`

	// Debug mode
	Debug bool `json:"debug,omitempty" yaml:"debug,omitempty" ini:"debug,omitempty"`

	// Color scheme.
	ColorScheme string `json:"color_scheme,omitempty" yaml:"color_scheme,omitempty" ini:"color_scheme,omitempty"`

	Custom404HTML template.HTML `json:"custom_404_html,omitempty" yaml:"custom_404_html,omitempty" ini:"custom_404_html,omitempty"`

	Custom403HTML template.HTML `json:"custom_403_html,omitempty" yaml:"custom_403_html,omitempty" ini:"custom_403_html,omitempty"`

	Custom500HTML template.HTML `json:"custom_500_html,omitempty" yaml:"custom_500_html,omitempty" ini:"custom_500_html,omitempty"`

	ExcludeThemeComponents []string `json:"exclude_theme_components,omitempty" yaml:"exclude_theme_components,omitempty" ini:"exclude_theme_components,omitempty"`
}

// DatabaseList is a map of Database.
type DatabaseList map[string]Database

// 文件上傳引擎
type FileUploadEngine struct {
	Name   string                 `json:"name,omitempty" yaml:"name,omitempty" ini:"name,omitempty"`
	Config map[string]interface{} `json:"config,omitempty" yaml:"config,omitempty" ini:"config,omitempty"`
}

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

// 設置Config(struct)title、theme、登入url、前綴url...資訊，如果參數cfg(struct)有些數值為空值，設置預設值
// 最後回傳globalCfg
func Set(cfg Config) *Config {

	lock.Lock()
	defer lock.Unlock()

	// 不能設置config兩次
	if atomic.LoadUint32(&count) != 0 {
		panic("can not set config twice")
	}
	atomic.StoreUint32(&count, 1)

	cfg = SetDefault(cfg)

	//global url前綴
	if cfg.UrlPrefix == "" {
		cfg.prefix = "/"
	} else if cfg.UrlPrefix[0] != '/' {
		cfg.prefix = "/" + cfg.UrlPrefix
	} else {
		cfg.prefix = cfg.UrlPrefix
	}

	// 紀錄器(cfg)初始化
	// initLogger(cfg)

	// if cfg.SqlLog {
	// 	// 將logger(struct).sqlLogOpen設為true
	// 	logger.OpenSQLLog()
	// }

	if cfg.Debug {
		declare.Do(func() {
			fmt.Println(`GoAdmin is now running.
Running in "debug" mode. Switch to "release" mode in production.`)
		})
	}

	globalCfg = &cfg

	return globalCfg
}

// 如果參數cfg(struct)有些數值為空值，設置預設值
func SetDefault(cfg Config) Config {
	// SetDefault假設第一個參數 = 第二個參數回傳第三個參數，沒有的話回傳第一個參數
	cfg.Title = utils.SetDefault(cfg.Title, "", "Orange")
	cfg.LoginTitle = utils.SetDefault(cfg.LoginTitle, "", "Orange")
	cfg.Logo = template.HTML(utils.SetDefault(string(cfg.Logo), "", "<b>Go</b>Orange"))
	cfg.MiniLogo = template.HTML(utils.SetDefault(string(cfg.MiniLogo), "", "<b>O</b>O"))
	cfg.Theme = utils.SetDefault(cfg.Theme, "", "adminlte")
	cfg.IndexUrl = utils.SetDefault(cfg.IndexUrl, "", "/info/manager")
	cfg.LoginUrl = utils.SetDefault(cfg.LoginUrl, "", "/login")
	cfg.AuthUserTable = utils.SetDefault(cfg.AuthUserTable, "", "users")
	if cfg.Theme == "adminlte" {
		cfg.ColorScheme = utils.SetDefault(cfg.ColorScheme, "", "skin-black")
	}
	// cfg.FileUploadEngine.Name = utils.SetDefault(cfg.FileUploadEngine.Name, "", "local")
	// cfg.Env = utils.SetDefault(cfg.Env, "", EnvProd)
	if cfg.SessionLifeTime == 0 {
		// default two hours
		cfg.SessionLifeTime = 7200
	}
	return cfg
}

// 取得預設資料庫DatabaseList["default"]的值
func (d DatabaseList) GetDefault() Database {
	return d["default"]
}

// 將參數key、db設置至DatabaseList(map[string]Database)
func (d DatabaseList) Add(key string, db Database) {
	d[key] = db
}

// 將資料庫依照資料庫引擎分組(ex:mysql一組mssql一組)
func (d DatabaseList) GroupByDriver() map[string]DatabaseList {
	drivers := make(map[string]DatabaseList)
	for key, item := range d {
		if driverList, ok := drivers[item.Driver]; ok {
			driverList.Add(key, item)
		} else {
			drivers[item.Driver] = make(DatabaseList)
			drivers[item.Driver].Add(key, item)
		}
	}
	return drivers
}

func (d DatabaseList) Copy() DatabaseList {
	var c = make(DatabaseList)
	for k, v := range d {
		c[k] = v
	}
	return c
}

// 將所有globalCfg.Databases[key]的driver值設置至DatabaseList(map[string]Database).Database.Driver後回傳
func GetDatabases() DatabaseList {
	var list = make(DatabaseList, len(globalCfg.Databases))
	for key := range globalCfg.Databases {
		list[key] = Database{
			Driver: globalCfg.Databases[key].Driver,
		}
	}
	return list
}

// 將參數s轉換成Service(struct)並回傳Service.C(Config struct)
func GetService(s interface{}) *Config {
	if srv, ok := s.(*Service); ok {
		return srv.C
	}
	panic("wrong service")
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

// 取得Config.prefix
func (c *Config) AssertPrefix() string {
	if c.prefix == "/" {
		return ""
	}
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

// 將參數suffix(後綴)與Config.prefix(前綴)處理後回傳
func (c *Config) Url(suffix string) string {
	if c.prefix == "/" {
		return suffix
	}
	if suffix == "/" {
		return c.prefix
	}
	return c.prefix + suffix
}

// 判斷Config.Env是否是"prod"
func (c *Config) IsProductionEnvironment() bool {
	return c.Env == "prod"
}

// 將url前綴處理後回傳
func (c *Config) PrefixFixSlash() string {
	if c.UrlPrefix == "/" {
		return ""
	}
	if c.UrlPrefix != "" && c.UrlPrefix[0] != '/' {
		return "/" + c.UrlPrefix
	}
	return c.UrlPrefix
}

// 處理URL
func (s Store) URL(suffix string) string {
	if len(suffix) > 4 && suffix[:4] == "http" {
		return suffix
	}
	if s.Prefix == "" {
		if suffix[0] == '/' {
			return suffix
		}
		return "/" + suffix
	}
	if s.Prefix[0] == '/' {
		if suffix[0] == '/' {
			return s.Prefix + suffix
		}
		return s.Prefix + "/" + suffix
	}
	if suffix[0] == '/' {
		if len(s.Prefix) > 4 && s.Prefix[:4] == "http" {
			return s.Prefix + suffix
		}
		return "/" + s.Prefix + suffix
	}
	if len(s.Prefix) > 4 && s.Prefix[:4] == "http" {
		return s.Prefix + "/" + suffix
	}
	return "/" + s.Prefix + "/" + suffix
}

// URLRemovePrefix將URL的前綴(ex:/admin)去除
func (c *Config) URLRemovePrefix(url string) string {
	if url == c.prefix {
		return "/"
	}
	if c.prefix == "/" {
		return url
	}
	return strings.Replace(url, c.prefix, "", 1)
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

// globalCfg.prefix
func Prefix() string {
	return globalCfg.prefix
}

// 處理globalCfg(Config struct).IndexUrl(登入後導向的url)後回傳
func GetIndexURL() string {
	return globalCfg.GetIndexURL()
}

func GetTitle() string {
	return globalCfg.Title
}

func GetLogo() template.HTML {
	return globalCfg.Logo
}

func GetMiniLogo() template.HTML {
	return globalCfg.MiniLogo
}

// 將globalCfg.suffix(後綴)與globalCfg.prefix(前綴)處理後回傳
func Url(suffix string) string {
	return globalCfg.Url(suffix)
}

// globalCfg(Config struct).prefix將URL的前綴去除
func URLRemovePrefix(url string) string {
	return globalCfg.URLRemovePrefix(url)
}

// globalCfg.LoginUrl
func GetLoginUrl() string {
	return globalCfg.LoginUrl
}

func AssertPrefix() string {
	return globalCfg.AssertPrefix()
}

// globalCfg.AssetUrl
func GetAssetUrl() string {
	return globalCfg.AssetUrl
}

// globalCfg.AuthUserTable
func GetAuthUserTable() string {
	return globalCfg.AuthUserTable
}

func GetStore() Store {
	return globalCfg.Store
}

// globalCfg.Debug
func GetDebug() bool {
	return globalCfg.Debug
}

func GetDomain() string {
	return globalCfg.Domain
}

func GetTheme() string {
	return globalCfg.Theme
}

// 將Config.Databases[key].Driver設置至Config.Databases[key]後回傳(迴圈)
func (c *Config) EraseSens() *Config {
	for key := range c.Databases {
		c.Databases[key] = Database{
			Driver: c.Databases[key].Driver,
		}
	}
	return c
}

// 複製globalCfg(Config struct)後將Config.Databases[key].Driver設置至Config.Databases[key]後回傳
func Get() *Config {
	return globalCfg.Copy().EraseSens()
}

func GetNoLimitLoginIP() bool {
	return globalCfg.NoLimitLoginIP
}

func GetAnimation() PageAnimation {
	return globalCfg.Animation
}

func GetSessionLifeTime() int {
	return globalCfg.SessionLifeTime
}

func GetColorScheme() string {
	return globalCfg.ColorScheme
}

// 判斷globalCfg(Config).Env是否是"prod"
func IsProductionEnvironment() bool {
	return globalCfg.IsProductionEnvironment()
}

// 排除主題元件
func GetExcludeThemeComponents() []string {
	return globalCfg.ExcludeThemeComponents
}

func GetCustomHeadHtml() template.HTML {
	return globalCfg.CustomHeadHtml
}

func GetCustomFootHtml() template.HTML {
	return globalCfg.CustomFootHtml
}

func GetFooterInfo() template.HTML {
	return globalCfg.FooterInfo
}

func GetCustom500HTML() template.HTML {
	return globalCfg.Custom500HTML
}

func GetCustom404HTML() template.HTML {
	return globalCfg.Custom404HTML
}

func GetCustom403HTML() template.HTML {
	return globalCfg.Custom403HTML
}

// 將DatabaseList(map[string]Database)JSON編碼
func (d DatabaseList) JSON() string {
	return utils.JSON(d)
}

// 將Store(struct)JSON編碼
func (s Store) JSON() string {
	if s.Path == "" && s.Prefix == "" {
		return ""
	}
	return utils.JSON(s)
}

type Service struct {
	C *Config
}

// 回傳config(string)
func (s *Service) Name() string {
	return "config"
}

// 將參數c設置並回傳Service(struct)
func SrvWithConfig(c *Config) *Service {
	return &Service{c}
}

// 將Config的值設置至map[string]string
func (c *Config) ToMap() map[string]string {
	var m = make(map[string]string, 0)
	m["databases"] = c.Databases.JSON()
	m["domain"] = c.Domain
	m["url_prefix"] = c.UrlPrefix
	m["theme"] = c.Theme
	m["store"] = c.Store.JSON()
	m["title"] = c.Title
	m["logo"] = string(c.Logo)
	m["mini_logo"] = string(c.MiniLogo)
	m["index_url"] = c.IndexUrl
	// m["site_off"] = strconv.FormatBool(c.SiteOff)
	m["login_url"] = c.LoginUrl
	m["debug"] = strconv.FormatBool(c.Debug)
	// m["env"] = c.Env

	// Logger config
	// ========================

	// m["info_log_path"] = c.InfoLogPath
	// m["error_log_path"] = c.ErrorLogPath
	// m["access_log_path"] = c.AccessLogPath
	// m["sql_log"] = strconv.FormatBool(c.SqlLog)
	// m["access_log_off"] = strconv.FormatBool(c.AccessLogOff)
	// m["info_log_off"] = strconv.FormatBool(c.InfoLogOff)
	// m["error_log_off"] = strconv.FormatBool(c.ErrorLogOff)
	// m["access_assets_log_off"] = strconv.FormatBool(c.AccessAssetsLogOff)

	// m["logger_rotate_max_size"] = strconv.Itoa(c.Logger.Rotate.MaxSize)
	// m["logger_rotate_max_backups"] = strconv.Itoa(c.Logger.Rotate.MaxBackups)
	// m["logger_rotate_max_age"] = strconv.Itoa(c.Logger.Rotate.MaxAge)
	// m["logger_rotate_compress"] = strconv.FormatBool(c.Logger.Rotate.Compress)

	// m["logger_encoder_time_key"] = c.Logger.Encoder.TimeKey
	// m["logger_encoder_level_key"] = c.Logger.Encoder.LevelKey
	// m["logger_encoder_name_key"] = c.Logger.Encoder.NameKey
	// m["logger_encoder_caller_key"] = c.Logger.Encoder.CallerKey
	// m["logger_encoder_message_key"] = c.Logger.Encoder.MessageKey
	// m["logger_encoder_stacktrace_key"] = c.Logger.Encoder.StacktraceKey
	// m["logger_encoder_level"] = c.Logger.Encoder.Level
	// m["logger_encoder_time"] = c.Logger.Encoder.Time
	// m["logger_encoder_duration"] = c.Logger.Encoder.Duration
	// m["logger_encoder_caller"] = c.Logger.Encoder.Caller
	// m["logger_encoder_encoding"] = c.Logger.Encoder.Encoding
	// m["logger_level"] = strconv.Itoa(int(c.Logger.Level))

	// m["color_scheme"] = c.ColorScheme
	m["session_life_time"] = strconv.Itoa(c.SessionLifeTime)
	// m["asset_url"] = c.AssetUrl
	// m["file_upload_engine"] = c.FileUploadEngine.JSON()
	// m["custom_head_html"] = string(c.CustomHeadHtml)
	// m["custom_foot_html"] = string(c.CustomFootHtml)
	// m["custom_404_html"] = string(c.Custom404HTML)
	// m["custom_403_html"] = string(c.Custom403HTML)
	// m["custom_500_html"] = string(c.Custom500HTML)
	// m["footer_info"] = string(c.FooterInfo)
	m["login_title"] = c.LoginTitle
	m["login_logo"] = string(c.LoginLogo)
	m["auth_user_table"] = c.AuthUserTable
	// if len(c.Extra) == 0 {
	// 	m["extra"] = ""
	// } else {
	// 	m["extra"] = utils.JSON(c.Extra)
	// }

	// m["animation_type"] = c.Animation.Type
	// m["animation_duration"] = fmt.Sprintf("%.2f", c.Animation.Duration)
	// m["animation_delay"] = fmt.Sprintf("%.2f", c.Animation.Delay)

	m["no_limit_login_ip"] = strconv.FormatBool(c.NoLimitLoginIP)

	// m["hide_config_center_entrance"] = strconv.FormatBool(c.HideConfigCenterEntrance)
	// m["hide_app_info_entrance"] = strconv.FormatBool(c.HideAppInfoEntrance)
	// m["hide_tool_entrance"] = strconv.FormatBool(c.HideToolEntrance)

	return m
}

// 將參數m(map[string]string)的值更新至Config(struct)
func (c *Config) Update(m map[string]string) error {
	updateLock.Lock()
	defer updateLock.Unlock()
	c.Domain = m["domain"]
	c.Theme = m["theme"]
	c.Title = m["title"]
	c.Logo = template.HTML(m["logo"])
	c.MiniLogo = template.HTML(m["mini_logo"])
	// c.Debug = utils.ParseBool(m["debug"])
	// c.Env = m["env"]
	// c.SiteOff = utils.ParseBool(m["site_off"])

	// c.AccessLogOff = utils.ParseBool(m["access_log_off"])
	// c.InfoLogOff = utils.ParseBool(m["info_log_off"])
	// c.ErrorLogOff = utils.ParseBool(m["error_log_off"])
	// c.AccessAssetsLogOff = utils.ParseBool(m["access_assets_log_off"])

	// if c.InfoLogPath != m["info_log_path"] {
	// 	c.InfoLogPath = m["info_log_path"]
	// }
	// if c.ErrorLogPath != m["error_log_path"] {
	// 	c.ErrorLogPath = m["error_log_path"]
	// }
	// if c.AccessLogPath != m["access_log_path"] {
	// 	c.AccessLogPath = m["access_log_path"]
	// }
	// c.SqlLog = utils.ParseBool(m["sql_log"])

	// c.Logger.Rotate.MaxSize, _ = strconv.Atoi(m["logger_rotate_max_size"])
	// c.Logger.Rotate.MaxBackups, _ = strconv.Atoi(m["logger_rotate_max_backups"])
	// c.Logger.Rotate.MaxAge, _ = strconv.Atoi(m["logger_rotate_max_age"])
	// c.Logger.Rotate.Compress = utils.ParseBool(m["logger_rotate_compress"])

	// c.Logger.Encoder.Encoding = m["logger_encoder_encoding"]
	// loggerLevel, _ := strconv.Atoi(m["logger_level"])
	// c.Logger.Level = int8(loggerLevel)

	// if c.Logger.Encoder.Encoding == "json" {
	// 	c.Logger.Encoder.TimeKey = m["logger_encoder_time_key"]
	// 	c.Logger.Encoder.LevelKey = m["logger_encoder_level_key"]
	// 	c.Logger.Encoder.NameKey = m["logger_encoder_name_key"]
	// 	c.Logger.Encoder.CallerKey = m["logger_encoder_caller_key"]
	// 	c.Logger.Encoder.MessageKey = m["logger_encoder_message_key"]
	// 	c.Logger.Encoder.StacktraceKey = m["logger_encoder_stacktrace_key"]
	// 	c.Logger.Encoder.Level = m["logger_encoder_level"]
	// 	c.Logger.Encoder.Time = m["logger_encoder_time"]
	// 	c.Logger.Encoder.Duration = m["logger_encoder_duration"]
	// 	c.Logger.Encoder.Caller = m["logger_encoder_caller"]
	// }

	// initLogger(*c)

	// if c.Theme == "adminlte" {
	// 	c.ColorScheme = m["color_scheme"]
	// }
	ses, _ := strconv.Atoi(m["session_life_time"])
	if ses != 0 {
		c.SessionLifeTime = ses
	}
	// c.CustomHeadHtml = template.HTML(m["custom_head_html"])
	// c.CustomFootHtml = template.HTML(m["custom_foot_html"])
	// c.Custom404HTML = template.HTML(m["custom_404_html"])
	// c.Custom403HTML = template.HTML(m["custom_403_html"])
	// c.Custom500HTML = template.HTML(m["custom_500_html"])
	// c.FooterInfo = template.HTML(m["footer_info"])
	c.LoginTitle = m["login_title"]
	// c.AssetUrl = m["asset_url"]
	c.LoginLogo = template.HTML(m["login_logo"])
	// c.NoLimitLoginIP = utils.ParseBool(m["no_limit_login_ip"])

	// c.HideConfigCenterEntrance = utils.ParseBool(m["hide_config_center_entrance"])
	// c.HideAppInfoEntrance = utils.ParseBool(m["hide_app_info_entrance"])
	// c.HideToolEntrance = utils.ParseBool(m["hide_tool_entrance"])

	// c.FileUploadEngine = GetFileUploadEngineFromJSON(m["file_upload_engine"])

	// c.Animation.Type = m["animation_type"]
	// c.Animation.Duration = utils.ParseFloat32(m["animation_duration"])
	// c.Animation.Delay = utils.ParseFloat32(m["animation_delay"])

	// if m["extra"] != "" {
	// 	var extra = make(map[string]interface{}, 0)
	// 	_ = json.Unmarshal([]byte(m["extra"]), &extra)
	// 	c.Extra = extra
	// }

	return nil
}

func (c *Config) Copy() *Config {
	return &Config{
		Databases: c.Databases.Copy(),
		Domain:    c.Domain,
		// Language:                      c.Language,
		UrlPrefix: c.UrlPrefix,
		Theme:     c.Theme,
		Store:     c.Store,
		Title:     c.Title,
		Logo:      c.Logo,
		MiniLogo:  c.MiniLogo,
		IndexUrl:  c.IndexUrl,
		LoginUrl:  c.LoginUrl,
		Debug:     c.Debug,
		Env:       c.Env,
		// InfoLogPath:                   c.InfoLogPath,
		// ErrorLogPath:                  c.ErrorLogPath,
		// AccessLogPath:                 c.AccessLogPath,
		// SqlLog:                        c.SqlLog,
		// AccessLogOff:                  c.AccessLogOff,
		// InfoLogOff:                    c.InfoLogOff,
		// ErrorLogOff:                   c.ErrorLogOff,
		ColorScheme:     c.ColorScheme,
		SessionLifeTime: c.SessionLifeTime,
		AssetUrl:        c.AssetUrl,
		// FileUploadEngine:              c.FileUploadEngine,
		CustomHeadHtml: c.CustomHeadHtml,
		CustomFootHtml: c.CustomFootHtml,
		FooterInfo:     c.FooterInfo,
		LoginTitle:     c.LoginTitle,
		LoginLogo:      c.LoginLogo,
		AuthUserTable:  c.AuthUserTable,
		// Extra:                         c.Extra,
		Animation:      c.Animation,
		NoLimitLoginIP: c.NoLimitLoginIP,
		// Logger:                        c.Logger,
		// SiteOff:                       c.SiteOff,
		// HideConfigCenterEntrance:      c.HideConfigCenterEntrance,
		// HideAppInfoEntrance:           c.HideAppInfoEntrance,
		// HideToolEntrance:              c.HideToolEntrance,
		Custom404HTML: c.Custom404HTML,
		Custom500HTML: c.Custom500HTML,
		// UpdateProcessFn:               c.UpdateProcessFn,
		// OpenAdminApi:                  c.OpenAdminApi,
		// HideVisitorUserCenterEntrance: c.HideVisitorUserCenterEntrance,
		ExcludeThemeComponents: c.ExcludeThemeComponents,
		prefix:                 c.prefix,
	}
}
