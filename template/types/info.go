package types

import (
	"goadminapi/context"
	"goadminapi/modules/db"
	"goadminapi/template/types/form"
	"goadminapi/template/types/table"
	"html/template"

	"goadminapi/plugins/admin/modules/parameter"
)

type FieldList []Field

// Field is the table field.
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

// 過濾表單欄位
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