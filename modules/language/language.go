package language

import (
	"html/template"
	"strings"

	"goadminapi/modules/config"

	"golang.org/x/text/language"
)

type LangSet map[string]string
type LangMap map[string]LangSet

var (
	CN    = language.Chinese.String()
	Langs = [...]string{CN}
)

var Lang = LangMap{
	language.Chinese.String(): cn,
	"cn":                      cn,
}

// 判斷globalCfg.Language是否為空，如果不為空則依照設定的語言處理參數，最後回傳
func GetWithScope(value string, scopes ...string) string {
	if config.GetLanguage() == "" {
		return value
	}

	if locale, ok := Lang[config.GetLanguage()][JoinScopes(scopes)+strings.ToLower(value)]; ok {
		return locale
	}

	return value
}

func Get(value string) string {
	return GetWithScope(value)
}

// 判斷globalCfg.Language是否為空，接著處理參數並回傳HTML
func GetFromHtml(value template.HTML, scopes ...string) template.HTML {
	if config.GetLanguage() == "" {
		return value
	}

	if locale, ok := Lang[config.GetLanguage()][JoinScopes(scopes)+strings.ToLower(string(value))]; ok {
		return template.HTML(locale)
	}

	return value
}

func FixedLanguageKey(key string) string {
	if key == "cn" {
		return CN
	}
	return key
}

func JoinScopes(scopes []string) string {
	j := ""
	for _, scope := range scopes {
		j += scope + "."
	}
	return j
}
