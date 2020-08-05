package modules

import uuid "github.com/satori/go.uuid"

// 判斷第二個參數(string)是否存在[]string(第一個參數中)
func InArray(arr []string, str string) bool {
	for _, v := range arr {
		if v == str {
			return true
		}
	}
	return false
}

// 如果第一個參數(source)為空則回傳第二個參數(def)，否則回傳source
func SetDefault(source, def string) string {
	if source == "" {
		return def
	}
	return source
}

// 判斷第二個參數符號(分隔符)，如果為[則回傳[field(第一個參數)]，否則回傳ex: 'field'
func FilterField(filed, delimiter string) string {
	if delimiter == "[" {
		return "[" + filed + "]"
	}
	return delimiter + filed + delimiter
}

// RemoveBlankFromArray 將參數中不為空的參數加入[]string
func RemoveBlankFromArray(s []string) []string {
	var r []string
	for _, str := range s {
		if str != "" {
			r = append(r, str)
		}
	}
	return r
}

// 判斷參數del後回傳del+s(參數)+del或[s(參數)]
func Delimiter(del, s string) string {
	if del == "[" {
		return "[" + s + "]"
	}
	return del + s + del
}

// 判斷條件，true return a，false return b
func AorB(condition bool, a, b string) string {
	if condition {
		return a
	}
	return b
}

func Uuid() string {
	uid, _ := uuid.NewV4()
	rst := uid.String()
	return rst
}

