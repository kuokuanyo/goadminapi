package language

// import (
// 	"strings"

// 	"goadminapi/modules/config"

// 	"golang.org/x/text/language"
// )

// type LangSet map[string]string
// type LangMap map[string]LangSet

// var (
// 	CN    = language.Chinese.String()
// 	Langs = [...]string{CN}
// )

// var Lang = LangMap{
// 	language.Chinese.String(): cn,
// 	"cn":                      cn,
// }

// func FixedLanguageKey(key string) string {
// 	if key == "cn" {
// 		return CN
// 	}
// 	return key
// }

// func JoinScopes(scopes []string) string {
// 	j := ""
// 	for _, scope := range scopes {
// 		j += scope + "."
// 	}
// 	return j
// }

// // GetWithScope return the value of given scopes.
// func GetWithScope(value string, scopes ...string) string {
// 	if config.GetLanguage() == "" {
// 		return value
// 	}

// 	if locale, ok := Lang[config.GetLanguage()][JoinScopes(scopes)+strings.ToLower(value)]; ok {
// 		return locale
// 	}

// 	return value
// }
