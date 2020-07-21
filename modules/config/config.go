package config

import "html/template"

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
	// global urlŤeşó
	UrlPrefix string `json:"prefix,omitempty" yaml:"prefix,omitempty" ini:"prefix,omitempty"`

	// The theme name of template.
	Theme string `json:"theme,omitempty" yaml:"theme,omitempty" ini:"theme,omitempty"`

	// The path where files will be stored into.
	Store Store `json:"store,omitempty" yaml:"store,omitempty" ini:"store,omitempty"`

	// The title of web page.
	Title string `json:"title,omitempty" yaml:"title,omitempty" ini:"title,omitempty"`

	// Logo is the top text in the sidebar.
	// Śb°źĂäÄćŞşłťłĄŞşlogo(Ŕł¸ÓŹOĽŞ¤W¨¤)
	Logo template.HTML `json:"logo,omitempty" yaml:"logo,omitempty" ini:"logo,omitempty"`

	// Mini-logo is the top text in the sidebar when folding.
	// Śb°źĂäÄćłťłĄŞşlogo(ˇí°źĂäÄćŚŹţ¨ÓŽÉ)
	MiniLogo template.HTML `json:"mini_logo,omitempty" yaml:"mini_logo,omitempty" ini:"mini_logo,omitempty"`

	// The url redirect to after login.
	// ľn¤JŤážÉŚVŞşurl
	IndexUrl string `json:"index,omitempty" yaml:"index,omitempty" ini:"index,omitempty"`

	// Login page URL
	// ľn¤J­ś­ąŞşurl
	LoginUrl string `json:"login_url,omitempty" yaml:"login_url,omitempty" ini:"login_url,omitempty"`

	prefix string
}

// DatabaseList is a map of Database.
type DatabaseList map[string]Database

// ¤ĺĽóŔxŚsŞşŚě¸mĽH¤Îprefix
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
