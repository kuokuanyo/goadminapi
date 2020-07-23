package utils

import "encoding/json"

// 假設第一個參數 = 第二個參數回傳第三個參數，沒有的話回傳第一個參數
func SetDefault(value, condition, def string) string {
	if value == condition {
		return def
	}
	return value
}

// 將參數a執行JSON編碼並回傳
func JSON(a interface{}) string {
	if a == nil {
		return ""
	}
	b, _ := json.Marshal(a)
	return string(b)
}
