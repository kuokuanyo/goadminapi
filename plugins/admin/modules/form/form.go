package form

type Values map[string][]string

// Get 透過參數key判斷Values[key]長度是否大於0，如果大於零回傳Values[key][0]，反之回傳""
func (f Values) Get(key string) string {
	if len(f[key]) > 0 {
		return f[key][0]
	}
	return ""
}

// Add 將參數key、value加入Values(map[string][]string)中
func (f Values) Add(key string, value string) {
	f[key] = []string{value}
}

// 透過參數key刪除Values(map[string][]string)[key]
func (f Values) Delete(key string) {
	delete(f, key)
}

// ToMap 將Values(struct)的值都加入map[string]string
func (f Values) ToMap() map[string]string {
	var m = make(map[string]string)
	for key, v := range f {
		if len(v) > 0 {
			m[key] = v[0]
		}
	}
	return m
}

// IsSingleUpdatePost check the param if from an single update post request type or not.
func (f Values) IsSingleUpdatePost() bool {
	return f.Get("__is_single_update") == "1"
}

// 刪除__post_type與__is_single_update的鍵與值後回傳map[string][]string
func (f Values) RemoveRemark() Values {
	f.Delete("__post_type")
	f.Delete("__is_single_update")
	return f
}
