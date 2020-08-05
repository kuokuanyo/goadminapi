package types

import (
	"encoding/json"
	"goadminapi/context"
	"goadminapi/modules/db"
	"goadminapi/modules/utils"
	"goadminapi/template/types/form"
	"goadminapi/template/types/table"
	"html/template"

	"goadminapi/plugins/admin/modules"
	"goadminapi/plugins/admin/modules/parameter"
)

var DefaultPageSizeList = []int{10, 20, 30, 50, 100}
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

// 透過參數並且將欄位、join語法...等資訊處理後，回傳[]TheadItem、欄位名稱、joinFields(ex:group_concat(goadmin_roles.`name`...)、join語法(left join....)、合併的資料表、可篩選過濾的欄位
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

		// 可以篩選的欄位，例如用戶頁面的用戶名、暱稱、角色
		if field.Filterable {
			if len(sql) > 0 {
				// GetFilterFormFields透過參數Parameters(struct)及string處理後回傳[]FormField
				filterForm = append(filterForm, field.GetFilterFormFields(params, headField, sql[0]())...)
			} else {
				filterForm = append(filterForm, field.GetFilterFormFields(params, headField)...)
			}
		}

	}

	return thead, fields, joinFields, joins, joinTables, filterForm
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

// 判斷資料是升冪或降冪
func (i *InfoPanel) GetSort() string {
	switch i.Sort {
	case 1:
		return "asc"
	default:
		return "desc"
	}
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
