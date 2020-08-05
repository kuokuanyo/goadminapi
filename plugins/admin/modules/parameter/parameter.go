package parameter

import (
	"goadminapi/plugins/admin/modules"
	"net/url"
	"strconv"
	"strings"
)

var keys = []string{"__page", "__pageSize", "__sort", "__columns", "__prefix", "_pjax", "__no_animation_"}

type Parameters struct {
	Page        string
	PageInt     int
	PageSize    string
	PageSizeInt int
	SortField   string
	Columns     []string
	SortType    string
	Animation   bool
	URLPath     string
	Fields      map[string][]string
}

// 設置(頁數及頁數Size)至Parameters(struct)
func BaseParam() Parameters {
	return Parameters{Page: "1", PageSize: "10", Fields: make(map[string][]string)}
}

// 透過參數key取得url中的值(value)，判斷是否為空，如果是空值回傳第三個參數def，如果不為空則回傳value
func getDefault(values url.Values, key, def string) string {
	value := values.Get(key)
	if value == "" {
		return def
	}
	return value
}

// 將頁面size、資料排列方式、選擇欄位...等資訊後設置至Parameters(struct)
func GetParam(u *url.URL, defaultPageSize int, p ...string) Parameters {
	// Query從url取得設定參數
	// ex: map[__columns:[id,username,name,goadmin_roles_goadmin_join_name,created_at,updated_at] __page:[1] __pageSize:[10]  __sort:[id] __sort_type:[desc] ...]
	values := u.Query()

	// 設定主鍵及排列方式
	primaryKey := "id"
	defaultSortType := "desc"
	if len(p) > 0 {
		primaryKey = p[0]
		defaultSortType = p[1]
	}

	// getDefault透過參數key取得url中的值(value)，判斷是否為空，如果是空值回傳第三個參數def，如果不為空則回傳value
	page := getDefault(values, "__page", "1")                                   // __page
	pageSize := getDefault(values, "__pageSize", strconv.Itoa(defaultPageSize)) // __pageSize
	pageInt, _ := strconv.Atoi(page)
	pageSizeInt, _ := strconv.Atoi(pageSize)

	sortField := getDefault(values, "__sort", primaryKey)                       // __sort
	sortType := getDefault(values, "__sort_type", defaultSortType)              // __sort_type

	// 如果有設定顯示欄位(則回傳欄位名稱至columnsArr，如果沒有設定則回傳空[])
	columns := getDefault(values, "__columns", "") // ex: id,username,name...
	columnsArr := make([]string, 0)
	if columns != "" {
		columns, _ = url.QueryUnescape(columns)
		columnsArr = strings.Split(columns, ",")
	}

	// 判斷是否有動畫參數
	animation := true
	if values.Get("__no_animation_") == "true" {
		animation = false
	}

	// fields為keys之外的鍵，ex:map[__edit_pk:[4]...]
	fields := make(map[string][]string)
	for key, value := range values {
		// keys = []string{"__page", "__pageSize", "__sort", "__columns", "__prefix", "_pjax", "__no_animation_"}
		if !modules.InArray(keys, key) && len(value) > 0 && value[0] != "" {
			if key == "__sort_type" {
				if value[0] != "desc" && value[0] != "asc" {
					fields[key] = []string{"desc"}
				}
			} else {
				if strings.Contains(key, "__operator__") &&
					values.Get(strings.Replace(key, "__operator__", "", -1)) == "" {
					continue
				}
				fields[strings.Replace(key, "[]", "", -1)] = value
			}
		}
	}

	return Parameters{
		Page:        page,
		PageSize:    pageSize,
		PageSizeInt: pageSizeInt,
		PageInt:     pageInt,
		URLPath:     u.Path,
		SortField:   sortField,
		SortType:    sortType,
		Fields:      fields,
		Animation:   animation,
		Columns:     columnsArr,
	}
}

// 將頁面size、資料排列方式、選擇欄位...等資訊後設置至Parameters(struct)
func GetParamFromURL(urlStr string, defaultPageSize int, defaultSortType, primaryKey string) Parameters {
	// 解析url
	// ex: /admin/info/manager?__page=1&__pageSize=10&__sort=id&__sort_type=desc
	u, err := url.Parse(urlStr)
	if err != nil {
		return BaseParam()
	}

	return GetParam(u, defaultPageSize, primaryKey, defaultSortType)
}

// 將Parameters(struct)的鍵與值加入至url.Values(map[string][]string)
func (param Parameters) GetFixedParamStr() url.Values {
	p := url.Values{}
	p.Add("__sort", param.SortField)
	p.Add("__pageSize", param.PageSize)
	p.Add("__sort_type", param.SortType)
	if len(param.Columns) > 0 {
		p.Add("__columns", strings.Join(param.Columns, ","))
	}
	for key, value := range param.Fields {
		p[key] = value
	}
	return p
}


// 處理url後(?...)的部分(頁面設置、排序方式....等)
func (param Parameters) GetRouteParamStr() string {
	p := param.GetFixedParamStr()
	p.Add("__page", param.Page)
	return "?" + p.Encode()
}

// 將參數(多個string)結合並設置至Parameters.Fields["__pk"]後回傳
func (param Parameters) WithPKs(id ...string) Parameters {
	param.Fields["__pk"] = []string{strings.Join(id, ",")}
	return param
}

func (param Parameters) Join() string {
	p := param.GetFixedParamStr()
	p.Add("__page", param.Page)
	return p.Encode()
}

// 透過參數field尋找Parameters.Fields[field]是否存在，如果存在則回傳第一個value值(string)，不存在則回傳""
func (param Parameters) GetFieldValue(field string) string {
	value, ok := param.Fields[field]
	if ok && len(value) > 0 {
		return value[0]
	}
	return ""
}

// 透過參數__pk尋找Parameters.Fields[__pk]是否存在，如果存在則回傳第一個value值(string)並且用","拆解成[]string
func (param Parameters) PKs() []string {
	// 透過參數__pk尋找Parameters.Fields[__pk]是否存在，如果存在則回傳第一個value值(string)
	// PrimaryKey = PrimaryKey
	pk := param.GetFieldValue("__pk")
	if pk == "" {
		return []string{}
	}
	return strings.Split(param.GetFieldValue("__pk"), ",")
}