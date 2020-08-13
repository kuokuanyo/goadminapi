package parameter

import (
	"goadminapi/plugins/admin/modules"
	"net/url"
	"strconv"
	"strings"
)

var keys = []string{"__page", "__pageSize", "__sort", "__columns", "__prefix", "_pjax", "__no_animation_"}

var operators = map[string]string{
	"like": "like",
	"gr":   ">",
	"gq":   ">=",
	"eq":   "=",
	"ne":   "!=",
	"le":   "<",
	"lq":   "<=",
	"free": "free",
}

type Parameters struct {
	Page        string
	PageInt     int
	PageSize    string
	PageSizeInt int
	SortField   string
	Columns     []string // 選擇顯示的欄位
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

	sortField := getDefault(values, "__sort", primaryKey)          // __sort
	sortType := getDefault(values, "__sort_type", defaultSortType) // __sort_type

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

// 取得wheres語法、where數值、existKeys([]string)
func (param Parameters) Statement(wheres, table, delimiter string, whereArgs []interface{}, columns, existKeys []string,
	filterProcess func(string, string, string) string) (string, []interface{}, []string) {
	var multiKey = make(map[string]uint8)

	// 處理param.Fields，ex: map[__is_all:[false]]
	for key, value := range param.Fields {
		keyIndexSuffix := ""

		keyArr := strings.Split(key, "__index__")
		// -----一般下面兩個條件式不會執行---------
		if len(keyArr) > 1 {
			key = keyArr[0]
			keyIndexSuffix = "__index__" + keyArr[1]
		}
		if keyIndexSuffix != "" {
			multiKey[key] = 0
		} else if _, exist := multiKey[key]; !exist && modules.InArray(existKeys, key) {
			continue
		}

		// 取的運算式符號
		var op string
		if strings.Contains(key, "_end") {
			key = strings.Replace(key, "_end", "", -1)
			op = "<="
		} else if strings.Contains(key, "_start") {
			key = strings.Replace(key, "_start", "", -1)
			op = ">="
		} else if len(value) > 1 {
			op = "in"
		} else if !strings.Contains(key, "__operator__") {
			// -------一般執行此條件-----
			op = operators[param.GetFieldOperator(key, keyIndexSuffix)] // op = '='(用戶頁面)
		}

		// --------一般__is_all都不在columns，因此執行else --------------
		if modules.InArray(columns, key) {
			if op == "in" {
				qmark := ""
				for range value {
					qmark += "?,"
				}
				wheres += table + "." + modules.FilterField(key, delimiter) + " " + op + " (" + qmark[:len(qmark)-1] + ") and "
			} else {
				wheres += table + "." + modules.FilterField(key, delimiter) + " " + op + " ? and "
			}
			if op == "like" && !strings.Contains(value[0], "%") {
				whereArgs = append(whereArgs, "%"+filterProcess(key, value[0], keyIndexSuffix)+"%")
			} else {
				for _, v := range value {
					whereArgs = append(whereArgs, filterProcess(key, v, keyIndexSuffix))
				}
			}
		} else {
			keys := strings.Split(key, "_join_") // ex: [__is_all]

			// -----一般下面條件式不會執行，len(keys)=0---------
			if len(keys) > 1 {
				val := filterProcess(key, value[0], keyIndexSuffix)
				if op == "in" {
					qmark := ""
					for range value {
						qmark += "?,"
					}
					wheres += keys[0] + "." + modules.FilterField(keys[1], delimiter) + " " + op + " (" + qmark[:len(qmark)-1] + ") and "
				} else {
					wheres += keys[0] + "." + modules.FilterField(keys[1], delimiter) + " " + op + " ? and "
				}
				if op == "like" && !strings.Contains(val, "%") {
					whereArgs = append(whereArgs, "%"+val+"%")
				} else {
					for _, v := range value {
						whereArgs = append(whereArgs, filterProcess(key, v, keyIndexSuffix))
					}
				}
			}
		}
		existKeys = append(existKeys, key)
	}
	if len(wheres) > 3 {
		wheres = wheres[:len(wheres)-4]
	}

	return wheres, whereArgs, existKeys
}

// 處理url後(?...)的部分(頁面設置、排序方式....等)
func (param Parameters) GetRouteParamStr() string {
	p := param.GetFixedParamStr()
	p.Add("__page", param.Page)
	return "?" + p.Encode()
}

// PK透過參數__pk尋找Parameters.Fields[__pk]是否存在，如果存在則回傳第一個value值(string)並且用","拆解成[]string，回傳第一個數值
func (param Parameters) PK() string {
	return param.PKs()[0]
}

// 將參數(多個string)結合並設置至Parameters.Fields["__pk"]
func (param Parameters) WithPKs(id ...string) Parameters {
	param.Fields["__pk"] = []string{strings.Join(id, ",")}
	return param
}

// 刪除Parameters.Fields[參數]
func (param Parameters) DeleteField(field string) Parameters {
	delete(param.Fields, field)
	return param
}

// DeletePK 刪除Parameters.Fields[__pk]
func (param Parameters) DeletePK() Parameters {
	// PrimaryKey = __pk
	delete(param.Fields, "__pk")
	return param
}

// DeleteEditPk 刪除Parameters.Fields[__edit_pk]後回傳
func (param Parameters) DeleteEditPk() Parameters {
	delete(param.Fields, "__edit_pk")
	return param
}

// DeleteDetailPk 刪除Parameters.Fields[__detail_pk]
func (param Parameters) DeleteDetailPk() Parameters {
	delete(param.Fields, "__detail_pk")
	return param
}

func (param Parameters) Join() string {
	p := param.GetFixedParamStr()
	p.Add("__page", param.Page)
	return p.Encode()
}

// 取得欄位設置數值
func (param Parameters) GetFieldValue(field string) string {
	value, ok := param.Fields[field]
	if ok && len(value) > 0 {
		return value[0]
	}
	return ""
}

// 取得欄位設置數值(多個值([]string))
func (param Parameters) GetFieldValues(field string) []string {
	return param.Fields[field]
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

// 設置第n(參數)頁至Parameters.Page
func (param Parameters) SetPage(page string) Parameters {
	param.Page = page
	return param
}

// 取得url參數(沒有PageSize)
func (param Parameters) GetRouteParamStrWithoutPageSize(page string) string {
	p := url.Values{}
	p.Add("__sort", param.SortField)
	p.Add("__page", page)
	p.Add("__sort_type", param.SortType)
	if len(param.Columns) > 0 {
		p.Add("__columns", strings.Join(param.Columns, ","))
	}
	for key, value := range param.Fields {
		p[key] = value
	}
	return "?" + p.Encode()
}

// 取得第n(參數)頁的url參數
func (param Parameters) URLNoAnimation(page string) string {
	return param.URLPath + param.SetPage(page).GetRouteParamStr() + "&" + "__no_animation_" + "=true"
}

// 將__pageSize、__no_animation_...等資訊加入url.Values(map[string][]string)後編碼回傳(沒有sort參數)
func (param Parameters) GetFixedParamStrWithoutSort() string {
	p := url.Values{}
	p.Add("__pageSize", param.PageSize)
	for key, value := range param.Fields {
		p[key] = value
	}
	// NoAnimationKey = __go_admin_no_animation_
	p.Add("__no_animation_", "true")
	if len(param.Columns) > 0 {
		// Columns = __columns
		p.Add("__columns", strings.Join(param.Columns, ","))
	}
	return "&" + p.Encode()
}

// GetLastPageRouteParamStr 取得前一頁的url參數，例如在第二頁時回傳第一頁的url參數
func (param Parameters) GetLastPageRouteParamStr() string {
	p := param.GetFixedParamStr()
	p.Add("__page", strconv.Itoa(param.PageInt-1))
	return "?" + p.Encode()
}

// GetNextPageRouteParamStr 取得下一頁的url參數，例如在第二頁時回傳第三頁的url參數
func (param Parameters) GetNextPageRouteParamStr() string {
	p := param.GetFixedParamStr()
	p.Add("__page", strconv.Itoa(param.PageInt+1))
	return "?" + p.Encode()
}

// 如果過濾值為範圍，設值起始值
func (param Parameters) GetFilterFieldValueStart(field string) string {
	return param.GetFieldValue(field + "_start")
}

// 如果過濾值為範圍，設值結束值
func (param Parameters) GetFilterFieldValueEnd(field string) string {
	return param.GetFieldValue(field + "_end")
}

// 將[]string利用"__separator__"字串join
func (param Parameters) GetFieldValuesStr(field string) string {
	return strings.Join(param.Fields[field], "__separator__")
}

// GetFieldOperator 取得運算子(ex: =.<.>...)
func (param Parameters) GetFieldOperator(field, suffix string) string {
	op := param.GetFieldValue(field + "__operator__" + suffix)
	if op == "" {
		return "eq"
	}
	return op
}

// 判斷__is_all是否 = true
func (param Parameters) IsAll() bool {
	return param.GetFieldValue("__is_all") == "true"
}

// WithIsAll 添加Parameters.Fields["__is_all"]
func (param Parameters) WithIsAll(isAll bool) Parameters {
	if isAll {
		param.Fields["__is_all"] = []string{"true"}
	} else {
		param.Fields["__is_all"] = []string{"false"}
	}
	return param
}

// DeleteIsAll 刪除Parameters.Fields["__is_all"]
func (param Parameters) DeleteIsAll() Parameters {
	delete(param.Fields, "__is_all")
	return param
}
