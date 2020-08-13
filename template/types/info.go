package types

import (
	"encoding/json"
	"goadminapi/context"
	"goadminapi/modules/db"
	"goadminapi/modules/utils"
	"goadminapi/template/types/form"
	"goadminapi/template/types/table"
	"html/template"
	"strconv"
	"strings"

	"goadminapi/plugins/admin/modules"
	"goadminapi/plugins/admin/modules/parameter"
)

var DefaultPageSizeList = []int{10, 20, 30, 50, 100} // 單頁顯示資料筆數選項
var JoinFieldValueDelimiter = utils.Uuid(8)

const DefaultPageSize = 10

type FieldList []Field

// TableInfo 資料表資訊
type TableInfo struct {
	Table      string
	PrimaryKey string
	Delimiter  string
	Driver     string
}

// Field 資料表資訊
type Field struct {
	Head     string
	Field    string
	TypeName db.DatabaseType

	Joins Joins

	Width      int
	Sortable   bool
	EditAble   bool
	Fixed      bool
	Filterable bool
	Hide       bool

	EditType    table.Type
	EditOptions FieldOptions

	FilterFormFields []FilterFormField

	FieldDisplay
}

type PostType uint8

const (
	PostTypeCreate = iota
	PostTypeUpdate
)

type FieldModel struct {
	// The primaryKey of the table.
	ID string
	// The value of the single query result.
	Value string
	// The current row data.
	Row map[string]interface{}
	// Post type
	PostType PostType
}

type FieldFilterFn func(value FieldModel) interface{}

// 過濾表單的欄位資訊
type FilterFormField struct {
	Type        form.Type
	Options     FieldOptions
	OptionTable OptionTable
	Width       int
	Operator    FilterOperator
	OptionExt   template.JS
	Head        string
	Placeholder string
	HelpMsg     template.HTML
	ProcessFn   func(string) string
}

type FilterType struct {
	FormType    form.Type
	Operator    FilterOperator
	Head        string
	Placeholder string
	NoHead      bool
	Width       int
	HelpMsg     template.HTML
	Options     FieldOptions
	Process     func(string) string
	OptionExt   map[string]interface{}
}

// join表資訊
type Join struct {
	Table     string
	Field     string
	JoinField string
	BaseTable string
}
type Joins []Join

type TabGroups [][]string
type TabHeaders []string
type Sort uint8
type primaryKey struct {
	Type db.DatabaseType
	Name string
}

type Where struct {
	Join     string
	Field    string
	Operator string
	Arg      interface{}
}
type Wheres []Where

type WhereRaw struct {
	Raw  string
	Args []interface{}
}

type Callbacks []context.Node

type DeleteFn func(ids []string) error
type DeleteFnWithRes func(ids []string, res error) error

type GetDataFn func(param parameter.Parameters) ([]map[string]interface{}, int)
type QueryFilterFn func(param parameter.Parameters, conn db.Connection) (ids []string, stopQuery bool)

type ContentWrapper func(content template.HTML) template.HTML

// InfoPanel
type InfoPanel struct {
	FieldList         FieldList
	curFieldListIndex int

	Table       string
	Title       string
	Description string

	// Warn: may be deprecated future.
	TabGroups  TabGroups
	TabHeaders TabHeaders

	Sort      Sort
	SortField string

	PageSizeList    []int
	DefaultPageSize int

	ExportType int

	primaryKey primaryKey

	IsHideNewButton    bool
	IsHideExportButton bool
	IsHideEditButton   bool
	IsHideDeleteButton bool
	IsHideDetailButton bool
	IsHideFilterButton bool
	IsHideRowSelector  bool
	IsHidePagination   bool
	IsHideFilterArea   bool
	IsHideQueryInfo    bool
	FilterFormLayout   form.Layout

	FilterFormHeadWidth  int
	FilterFormInputWidth int

	Wheres    Wheres
	WhereRaws WhereRaw

	Callbacks Callbacks

	Buttons Buttons

	TableLayout string

	DeleteHook  DeleteFn
	PreDeleteFn DeleteFn
	DeleteFn    DeleteFn

	DeleteHookWithRes DeleteFnWithRes

	GetDataFn GetDataFn

	processChains DisplayProcessFnChains

	ActionButtons Buttons

	DisplayGeneratorRecords map[string]struct{}

	QueryFilterFn QueryFilterFn

	Wrapper ContentWrapper

	// column operation buttons
	Action     template.HTML
	HeaderHtml template.HTML
	FooterHtml template.HTML
}

type PostFieldFilterFn func(value PostFieldModel) interface{}
type FieldModelValue []string
type PostFieldModel struct {
	ID    string
	Value FieldModelValue
	Row   map[string]string
	// Post type
	PostType PostType
}

type InfoList []map[string]InfoItem
type InfoItem struct {
	Content template.HTML `json:"content"`
	Value   string        `json:"value"`
}

type DefaultAction struct {
	Attr   template.HTML
	JS     template.JS
	Ext    template.HTML
	Footer template.HTML
}

// 預設InfoPanel(struct)
func NewInfoPanel(pk string) *InfoPanel {
	return &InfoPanel{
		curFieldListIndex:       -1,
		PageSizeList:            DefaultPageSizeList, // []int{10, 20, 30, 50, 100}
		DefaultPageSize:         DefaultPageSize,     // 10
		processChains:           make(DisplayProcessFnChains, 0),
		Buttons:                 make(Buttons, 0),
		Callbacks:               make(Callbacks, 0),
		DisplayGeneratorRecords: make(map[string]struct{}),
		Wheres:                  make([]Where, 0),
		WhereRaws:               WhereRaw{},
		SortField:               pk,
		TableLayout:             "auto",
		FilterFormInputWidth:    10,
		FilterFormHeadWidth:     2,
	}
}

// SetPrimaryKey 將參數設置至InfoPanel(struct).primaryKey中並回傳
func (i *InfoPanel) SetPrimaryKey(name string, typ db.DatabaseType) *InfoPanel {
	i.primaryKey = primaryKey{Name: name, Type: typ}
	return i
}

//  AddField 添加欄位資訊至InfoPanel.FieldList
func (i *InfoPanel) AddField(head, field string, typeName db.DatabaseType) *InfoPanel {
	i.FieldList = append(i.FieldList, Field{
		Head:     head,
		Field:    field,
		TypeName: typeName,
		Sortable: false,
		Joins:    make(Joins, 0),
		EditAble: false,
		EditType: table.Text,
		FieldDisplay: FieldDisplay{
			Display: func(value FieldModel) interface{} {
				return value.Value
			},
			// chooseDisplayProcessChains 如果參數長度大於0則回傳參數
			// 否則複製全域變數globalDisplayProcessChains([]FieldFilterFn)後回傳
			DisplayProcessChains: chooseDisplayProcessChains(i.processChains),
		},
	})
	i.curFieldListIndex++
	return i
}

// FieldJoin 添加join其他資料表資訊
func (i *InfoPanel) FieldJoin(join Join) *InfoPanel {
	i.FieldList[i.curFieldListIndex].Joins = append(i.FieldList[i.curFieldListIndex].Joins, join)
	return i
}

// 設置為可篩選並添加篩選的表單欄位資訊至FilterFormFields
func (i *InfoPanel) FieldFilterable(filterType ...FilterType) *InfoPanel {
	i.FieldList[i.curFieldListIndex].Filterable = true

	// 如沒設置參數則添加一個過濾的表單欄位資訊至FilterFormFields
	if len(filterType) == 0 {
		i.FieldList[i.curFieldListIndex].FilterFormFields = append(i.FieldList[i.curFieldListIndex].FilterFormFields,
			FilterFormField{
				Type:        form.Text,
				Head:        i.FieldList[i.curFieldListIndex].Head,
				Placeholder: "輸入" + " " + i.FieldList[i.curFieldListIndex].Head,
			})
	}

	for _, filter := range filterType {
		var ff FilterFormField
		ff.Operator = filter.Operator
		if filter.FormType == form.Default {
			ff.Type = form.Text
		} else {
			ff.Type = filter.FormType
		}
		ff.Head = modules.AorB(!filter.NoHead && filter.Head == "",
			i.FieldList[i.curFieldListIndex].Head, filter.Head)
		ff.Width = filter.Width
		ff.HelpMsg = filter.HelpMsg
		ff.ProcessFn = filter.Process
		ff.Placeholder = modules.AorB(filter.Placeholder == "", "輸入"+" "+ff.Head, filter.Placeholder)
		ff.Options = filter.Options
		if len(filter.OptionExt) > 0 {
			s, _ := json.Marshal(filter.OptionExt)
			ff.OptionExt = template.JS(s)
		}
		i.FieldList[i.curFieldListIndex].FilterFormFields = append(i.FieldList[i.curFieldListIndex].FilterFormFields, ff)
	}
	return i
}

// GetFilterFormFields 將Field(struct)處理後取得可篩選的欄位資訊
func (f Field) GetFilterFormFields(params parameter.Parameters, headField string, sql ...*db.SQL) []FormField {
	var (
		filterForm               = make([]FormField, 0)
		value, value2, keySuffix string
	)

	// 處理可以篩選條件的欄位
	for index, filter := range f.FilterFormFields {
		// 一般index = 0，不會執行
		if index > 0 {
			keySuffix = "__index__" + strconv.Itoa(index)
		}

		if filter.Type.IsRange() { // 是否設定篩選範圍
			value = params.GetFilterFieldValueStart(headField)
			value2 = params.GetFilterFieldValueEnd(headField)
		} else if filter.Type.IsMultiSelect() { // 是否篩選多個條件
			value = params.GetFieldValuesStr(headField)
		} else {
			if filter.Operator == FilterOperatorFree {
				// GetFieldOperator 取得運算子(ex: =.<.>...)
				value2 = GetOperatorFromValue(params.GetFieldOperator(headField, keySuffix)).String()
			}
			// 一般篩選欄位都執行GetFieldValue函式
			// GetFieldValue 取得欄位設置數值
			value = params.GetFieldValue(headField + keySuffix)
		}

		var (
			optionExt1 = filter.OptionExt
			optionExt2 template.JS
		)
		if filter.OptionExt == template.JS("") {
			// ------------一般可篩選欄位只會執行GetDefaultOptions，結果都為空--------------
			// GetDefaultOptions 設置表單欄位選項
			op1, op2, js := filter.Type.GetDefaultOptions(headField + keySuffix)
			if op1 != nil {
				s, _ := json.Marshal(op1)
				optionExt1 = template.JS(string(s))
			}
			if op2 != nil {
				s, _ := json.Marshal(op2)
				optionExt2 = template.JS(string(s))
			}
			if js != template.JS("") {
				optionExt1 = js
			}
		}

		field := &FormField{
			Field:       headField + keySuffix,
			FieldClass:  headField + keySuffix,
			Head:        filter.Head,
			TypeName:    f.TypeName,
			HelpMsg:     filter.HelpMsg,
			FormType:    filter.Type,
			Editable:    true,
			Width:       filter.Width,
			Placeholder: filter.Placeholder,
			Value:       template.HTML(value),
			Value2:      value2,
			Options:     filter.Options,
			OptionExt:   optionExt1,
			OptionExt2:  optionExt2,
			OptionTable: filter.OptionTable,
			Label:       filter.Operator.Label(),
		}

		// 判斷條件後，從資料庫取得資料設置選項
		// ----------一般不會執行------------------
		field.setOptionsFromSQL(sql[0])

		// ----------一般下面兩個條件不會執行------------------
		if filter.Type.IsSingleSelect() {
			// SetSelected 判斷條件後將參數加入FieldOptions[k].SelectedLabel
			field.Options = field.Options.SetSelected(params.GetFieldValue(f.Field), filter.Type.SelectedLabel())
		}
		if filter.Type.IsMultiSelect() {
			field.Options = field.Options.SetSelected(params.GetFieldValues(f.Field), filter.Type.SelectedLabel())
		}

		filterForm = append(filterForm, *field)

		// ----------一般下面條件不會執行------------------
		if filter.Operator.AddOrNot() {
			filterForm = append(filterForm, FormField{
				Field:      headField + "__operator__" + keySuffix,
				FieldClass: headField + "__operator__" + keySuffix,
				Head:       f.Head,
				TypeName:   f.TypeName,
				Value:      template.HTML(filter.Operator.Value()),
				FormType:   filter.Type,
				Hide:       true,
			})
		}

	}
	return filterForm
}

// 取得[]TheadItem(欄位資訊)、欄位名稱、join語法(left join....)
func (f FieldList) GetThead(info TableInfo, params parameter.Parameters, columns []string) (Thead, string, string) {
	var (
		thead      = make(Thead, 0)
		fields     = ""
		joins      = ""
		joinTables = make([]string, 0)
	)

	for _, field := range f {
		// ------ID以及跟其他表關聯的欄位不會執行--------
		if field.Field != info.PrimaryKey && modules.InArray(columns, field.Field) &&
			!field.Joins.Valid() {
			// ex: users.`username`,users.`name`,users.`created_at`,users.`updated_at`,
			fields += info.Table + "." + modules.FilterField(field.Field, info.Delimiter) + ","
		}

		headField := field.Field

		if field.Joins.Valid() {
			headField = field.Joins.Last().Table + "_join_" + field.Field
			for _, join := range field.Joins {
				if !modules.InArray(joinTables, join.Table) {
					joinTables = append(joinTables, join.Table)
					if join.BaseTable == "" {
						join.BaseTable = info.Table
					}
					// ex: joins =  left join `role_users` on role_users.`user_id` = users.`id` left join....
					joins += " left join " + modules.FilterField(join.Table, info.Delimiter) + " on " +
						join.Table + "." + modules.FilterField(join.JoinField, info.Delimiter) + " = " +
						join.BaseTable + "." + modules.FilterField(join.Field, info.Delimiter)
				}
			}
		}
		// 檢查欄位是否隱藏
		if field.Hide {
			continue
		}

		thead = append(thead, TheadItem{
			Head:     field.Head,
			Sortable: field.Sortable,
			Field:    headField,
			// params.Columns為顯示的欄位
			Hide:       !modules.InArrayWithoutEmpty(params.Columns, headField), // 是否隱藏欄位
			Editable:   field.EditAble,
			EditType:   field.EditType.String(),
			EditOption: field.EditOptions,
			Width:      strconv.Itoa(field.Width) + "px",
		})
	}
	return thead, fields, joins
}

// 取得[]TheadItem(欄位資訊)、欄位名稱、joinFields(ex:group_concat(goadmin_roles.`name`...)、join語法(left join....)、合併的資料表、可篩選過濾的欄位
func (f FieldList) GetTheadAndFilterForm(info TableInfo, params parameter.Parameters, columns []string,
	sql ...func() *db.SQL) (Thead, string, string, string, []string, []FormField) {
	var (
		thead      = make(Thead, 0)
		fields     = ""                   // 欄位
		joinFields = ""                   // ex: group_concat(roles.`name` separator 'CkN694kH') as roles_join_name,
		joins      = ""                   // join資料表語法，ex: left join `role_users` on role_users.`user_id` = users.`id` left join....
		joinTables = make([]string, 0)    // ex:{roles role_id id role_users}(用戶頁面)
		filterForm = make([]FormField, 0) // 可以篩選過濾的欄位
	)

	// field為介面顯示的欄位
	for _, field := range f {
		// ------ID以及跟其他表關聯的欄位不會執行--------
		if field.Field != info.PrimaryKey && modules.InArray(columns, field.Field) &&
			// Valid對joins([]join(struct))執行迴圈，假設Join的Table、Field、JoinField都不為空，回傳true
			!field.Joins.Valid() {
			// ex: users.`username`,users.`name`,users.`created_at`,users.`updated_at`,
			fields += info.Table + "." + modules.FilterField(field.Field, info.Delimiter) + ","
		}

		headField := field.Field

		// -------------編輯介面(用戶的roles欄位(需要join其他表)會執行)-------------
		// 處理join語法
		// ex: [{role_users id user_id } {roles role_id id role_users}]
		if field.Joins.Valid() {
			// ex:roles_join_name
			headField = field.Joins.Last().Table + "_join_" + field.Field

			// GetAggregationExpression取得資料庫引擎的Aggregation表達式，將參數值加入表達式
			// ex: group_concat(roles.`name` separator 'CkN694kH') as roles_join_name,
			joinFields += db.GetAggregationExpression(info.Driver, field.Joins.Last().Table+"."+
				modules.FilterField(field.Field, info.Delimiter), headField, JoinFieldValueDelimiter) + ","

			for _, join := range field.Joins {
				if !modules.InArray(joinTables, join.Table) {
					joinTables = append(joinTables, join.Table)
					if join.BaseTable == "" {
						join.BaseTable = info.Table
					}
					// ex: joins =  left join `role_users` on role_users.`user_id` = users.`id` left join....
					joins += " left join " + modules.FilterField(join.Table, info.Delimiter) + " on " +
						join.Table + "." + modules.FilterField(join.JoinField, info.Delimiter) + " = " +
						join.BaseTable + "." + modules.FilterField(join.Field, info.Delimiter)

				}
			}
		}

		// 取得可篩選的欄位資訊，例如用戶頁面的用戶名、暱稱、角色
		if field.Filterable {
			if len(sql) > 0 {
				// GetFilterFormFields 將Field(struct)處理後取得可篩選的欄位資訊
				filterForm = append(filterForm, field.GetFilterFormFields(params, headField, sql[0]())...)
			} else {
				filterForm = append(filterForm, field.GetFilterFormFields(params, headField)...)
			}
		}

		// 檢查欄位是否隱藏
		if field.Hide {
			continue
		}

		// 將值添加至[]TheadItem
		thead = append(thead, TheadItem{
			Head:     field.Head,
			Sortable: field.Sortable, // 是否可以排序
			Field:    headField,
			// params.Columns為顯示的欄位
			Hide:       !modules.InArrayWithoutEmpty(params.Columns, headField), // 是否隱藏欄位
			Editable:   field.EditAble,
			EditType:   field.EditType.String(),
			EditOption: field.EditOptions,
			Width:      strconv.Itoa(field.Width) + "px",
		})
	}

	return thead, fields, joinFields, joins, joinTables, filterForm
}

// 透過參數取得欄位資訊
func (f FieldList) GetFieldByFieldName(name string) Field {
	for _, field := range f {
		if field.Field == name {
			return field
		}
		if JoinField(field.Joins.Last().Table, field.Field) == name {
			return field
		}
	}
	return Field{}
}

// 取得過濾欄位的處裡值
func (f FieldList) GetFieldFilterProcessValue(key, value, keyIndex string) string {

	field := f.GetFieldByFieldName(key)
	index := 0
	if keyIndex != "" {
		index, _ = strconv.Atoi(keyIndex)
	}
	if field.FilterFormFields[index].ProcessFn != nil {
		value = field.FilterFormFields[index].ProcessFn(value)
	}
	return value
}

// Statement 處理wheres語法及where值後回傳
func (whs Wheres) Statement(wheres, delimiter string, whereArgs []interface{}, existKeys, columns []string) (string, []interface{}) {
	pwheres := ""

	for k, wh := range whs {
		whFieldArr := strings.Split(wh.Field, ".")
		whField := ""
		whTable := ""

		if len(whFieldArr) > 1 {
			whField = whFieldArr[1]
			whTable = whFieldArr[0]
		} else {
			whField = whFieldArr[0]
		}

		if modules.InArray(existKeys, whField) {
			continue
		}

		// TODO: support like operation and join table
		if modules.InArray(columns, whField) {

			joinMark := ""
			if k != len(whs)-1 {
				joinMark = whs[k+1].Join
			}

			if whTable != "" {
				pwheres += whTable + "." + modules.FilterField(whField, delimiter) + " " + wh.Operator + " ? " + joinMark + " "
			} else {
				pwheres += modules.FilterField(whField, delimiter) + " " + wh.Operator + " ? " + joinMark + " "
			}
			whereArgs = append(whereArgs, wh.Arg)
		}
	}
	if wheres != "" && pwheres != "" {
		wheres += " and "
	}
	return wheres + pwheres, whereArgs
}

// Statement 處理wheres語法及where值後回傳
func (wh WhereRaw) Statement(wheres string, whereArgs []interface{}) (string, []interface{}) {

	if wh.Raw == "" {
		return wheres, whereArgs
	}

	if wheres != "" {
		if wh.check() != 0 {
			wheres += wh.Raw + " "
		} else {
			wheres += " and " + wh.Raw + " "
		}

		whereArgs = append(whereArgs, wh.Args...)
	} else {
		wheres += wh.Raw[wh.check():] + " "
		whereArgs = append(whereArgs, wh.Args...)
	}

	return wheres, whereArgs
}

// SetTable 設置資料表
func (i *InfoPanel) SetTable(table string) *InfoPanel {
	i.Table = table
	return i
}

// SetTitle 設置主題名稱
func (i *InfoPanel) SetTitle(title string) *InfoPanel {
	i.Title = title
	return i
}

// SetDescription 設置描述
func (i *InfoPanel) SetDescription(desc string) *InfoPanel {
	i.Description = desc
	return i
}

// GetPageSizeList 取得單頁顯示資料筆數選項
func (i *InfoPanel) GetPageSizeList() []string {
	var pageSizeList = make([]string, len(i.PageSizeList))
	for j := 0; j < len(i.PageSizeList); j++ {
		pageSizeList[j] = strconv.Itoa(i.PageSizeList[j])
	}
	return pageSizeList
}

// FieldSortable 設置為可排序
func (i *InfoPanel) FieldSortable() *InfoPanel {
	i.FieldList[i.curFieldListIndex].Sortable = true
	return i
}

// SetDeleteFn 設置刪除函式
func (i *InfoPanel) SetDeleteFn(fn DeleteFn) *InfoPanel {
	i.DeleteFn = fn
	return i
}

// FieldDisplay 將參數添加至InfoPanel.FieldList[].Display
func (i *InfoPanel) FieldDisplay(filter FieldFilterFn) *InfoPanel {
	i.FieldList[i.curFieldListIndex].Display = filter
	return i
}

// HideFilterArea(隱藏篩選區域) InfoPanel.IsHideFilterArea = true
func (i *InfoPanel) HideFilterArea() *InfoPanel {
	i.IsHideFilterArea = true
	return i
}

// 添加func(value FieldModel) interface{}至參數i.processChains([]FieldFilterFn)
func (i *InfoPanel) AddXssJsFilter() *InfoPanel {
	i.processChains = addXssJsFilter(i.processChains)
	return i
}

// 添加func(param parameter.Parameters) ([]map[string]interface{}, int)至參數i.GetDataFn
func (i *InfoPanel) SetGetDataFn(fn GetDataFn) *InfoPanel {
	i.GetDataFn = fn
	return i
}

// 判斷資料是升冪或降冪
func (i *InfoPanel) GetSort() string {
	switch i.Sort {
	case 1:
		return "asc"
	default:
		return "desc"
	}
}

// GroupBy get []InfoList
func (i InfoList) GroupBy(groups TabGroups) []InfoList {
	var res = make([]InfoList, len(groups))

	for key, value := range groups {
		var newInfoList = make(InfoList, len(i))

		for index, info := range i {
			var newRow = make(map[string]InfoItem)
			for mk, m := range info {
				if modules.InArray(value, mk) {
					newRow[mk] = m
				}
			}
			newInfoList[index] = newRow
		}

		res[key] = newInfoList
	}

	return res
}

// return table_join_field
func JoinField(table, field string) string {
	return table + "_join_" + field
}

// 假設Join的Table、Field、JoinField都不為空，回傳true
func (j Join) Valid() bool {
	return j.Table != "" && j.Field != "" && j.JoinField != ""
}

// 對joins([]join(struct))執行迴圈，假設Join的Table、Field、JoinField都不為空，回傳true
func (j Joins) Valid() bool {
	for i := 0; i < len(j); i++ {
		// 假設Join的Table、Field、JoinField都不為空，回傳true
		if j[i].Valid() {
			return true
		}
	}
	return false
}

// Last 判斷Joins([]Join)長度，如果大於0回傳Joins[len(j)-1](最後一個數值)
func (j Joins) Last() Join {
	if len(j) > 0 {
		return j[len(j)-1]
	}
	return Join{}
}

// *****************FieldModelValue([]string)的方法*******************

// return FieldModelValue[0]
func (r FieldModelValue) First() string {
	return r[0]
}

// return FieldModelValue[0]
func (r FieldModelValue) Value() string {
	return r.First()
}

func (wh WhereRaw) check() int {
	index := 0
	for i := 0; i < len(wh.Raw); i++ {
		if wh.Raw[i] == ' ' {
			continue
		} else {
			if wh.Raw[i] == 'a' {
				if len(wh.Raw) < i+3 {
					break
				} else {
					if wh.Raw[i+1] == 'n' && wh.Raw[i+2] == 'd' {
						index = i + 3
					}
				}
			} else if wh.Raw[i] == 'o' {
				if len(wh.Raw) < i+2 {
					break
				} else {
					if wh.Raw[i+1] == 'r' {
						index = i + 2
					}
				}
			} else {
				break
			}
		}
	}
	return index
}

// *****************FieldModelValue([]string)的方法*******************

// 設置DefaultAction(struct)
func NewDefaultAction(attr, ext, footer template.HTML, js template.JS) *DefaultAction {
	return &DefaultAction{Attr: attr, Ext: ext, Footer: footer, JS: js}
}

// *****************Action(interface)的所有方法*******************

// SetBtnId no return
func (def *DefaultAction) SetBtnId(btnId string)        {}
// SetBtnData no return 
func (def *DefaultAction) SetBtnData(data interface{})  {}
// Js get DefaultAction.JS
func (def *DefaultAction) Js() template.JS              { return def.JS }
// BtnAttribute get DefaultAction.Attr
func (def *DefaultAction) BtnAttribute() template.HTML  { return def.Attr }
// BtnClass return ""
func (def *DefaultAction) BtnClass() template.HTML      { return "" }
// ExtContent get DefaultAction.Ext
func (def *DefaultAction) ExtContent() template.HTML    { return def.Ext }
// FooterContent get DefaultAction.Footer
func (def *DefaultAction) FooterContent() template.HTML { return def.Footer }
// GetCallbacks get context.Node{}
func (def *DefaultAction) GetCallbacks() context.Node   { return context.Node{} }

// *****************Action(interface)的所有方法*******************

// 判斷TabGroups([][]string)是否長度大於0
func (t TabGroups) Valid() bool {
	return len(t) > 0
}