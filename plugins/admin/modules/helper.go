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

func Uuid() string {
	uid, _ := uuid.NewV4()
	rst := uid.String()
	return rst
}
