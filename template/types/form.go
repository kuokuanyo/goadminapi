package types

import (
	"encoding/json"
	"fmt"
	"html/template"
	"strings"

	"goadminapi/modules/db"

	"goadminapi/context"
	"goadminapi/plugins/admin/modules/form"
	form2 "goadminapi/template/types/form"
)

type FieldOption struct {
	Text          string            `json:"text"`
	Value         string            `json:"value"`
	TextHTML      template.HTML     `json:"-"`
	Selected      bool              `json:"-"`
	SelectedLabel template.HTML     `json:"-"`
	Extra         map[string]string `json:"-"`
}
type FieldOptions []FieldOption

type OptionTableQueryProcessFn func(sql *db.SQL) *db.SQL
type OptionProcessFn func(options FieldOptions) FieldOptions
type OptionTable struct {
	Table          string
	TextField      string
	ValueField     string
	QueryProcessFn OptionTableQueryProcessFn
	ProcessFn      OptionProcessFn
}

type OptionInitFn func(val FieldModel) FieldOptions
type OptionArrInitFn func(val FieldModel) []FieldOptions

type FormFields []FormField

// FormField is the form field with different options.
type FormField struct {
	Field          string          `json:"field"`
	FieldClass     string          `json:"field_class"`
	TypeName       db.DatabaseType `json:"type_name"`
	Head           string          `json:"head"`
	Foot           template.HTML   `json:"foot"`
	FormType       form2.Type      `json:"form_type"`
	FatherFormType form2.Type      `json:"father_form_type"`
	FatherField    string          `json:"father_field"`

	RowWidth int
	RowFlag  uint8

	Default                template.HTML  `json:"default"`
	DefaultArr             interface{}    `json:"default_arr"`
	Value                  template.HTML  `json:"value"`
	Value2                 string         `json:"value_2"`
	ValueArr               []string       `json:"value_arr"`
	Value2Arr              []string       `json:"value_2_arr"`
	Options                FieldOptions   `json:"options"`
	OptionsArr             []FieldOptions `json:"options_arr"`
	DefaultOptionDelimiter string         `json:"default_option_delimiter"`
	Label                  template.HTML  `json:"label"`
	HideLabel              bool           `json:"hide_label"`

	Placeholder string `json:"placeholder"`

	CustomContent template.HTML `json:"custom_content"`
	CustomJs      template.JS   `json:"custom_js"`
	CustomCss     template.CSS  `json:"custom_css"`

	Editable    bool `json:"editable"` // 允許編輯
	NotAllowAdd bool `json:"not_allow_add"` // 不允許增加
	Must        bool `json:"must"` // 該欄位必填
	Hide        bool `json:"hide"`

	Width int `json:"width"`

	InputWidth int `json:"input_width"`
	HeadWidth  int `json:"head_width"`

	Joins Joins `json:"-"` 

	Divider      bool   `json:"divider"`
	DividerTitle string `json:"divider_title"`

	HelpMsg template.HTML `json:"help_msg"` // 欄位提示訊息

	TableFields FormFields

	OptionExt       template.JS     `json:"option_ext"`
	OptionExt2      template.JS     `json:"option_ext_2"`
	OptionInitFn    OptionInitFn    `json:"-"`
	OptionArrInitFn OptionArrInitFn `json:"-"`
	OptionTable     OptionTable     `json:"-"`

	FieldDisplay `json:"-"`
	PostFilterFn PostFieldFilterFn `json:"-"`
}

type FormPostFn func(values form.Values) error
type FormPreProcessFn func(values form.Values) form.Values
type Responder func(ctx *context.Context)

// FormPanel
type FormPanel struct {
	FieldList         FormFields // 表單欄位資訊
	curFieldListIndex int

	// Warn: may be deprecated in the future.
	TabGroups  TabGroups
	TabHeaders TabHeaders

	Table       string
	Title       string
	Description string

	Validator    FormPostFn
	PostHook     FormPostFn
	PreProcessFn FormPreProcessFn

	Callbacks Callbacks

	primaryKey primaryKey

	UpdateFn FormPostFn
	InsertFn FormPostFn

	IsHideContinueEditCheckBox bool
	IsHideContinueNewCheckBox  bool
	IsHideResetButton          bool
	IsHideBackButton           bool

	Layout form2.Layout

	HTMLContent template.HTML

	Header template.HTML

	InputWidth int
	HeadWidth  int

	Ajax          bool
	AjaxSuccessJS template.JS
	AjaxErrorJS   template.JS

	Responder Responder

	Wrapper ContentWrapper

	processChains DisplayProcessFnChains

	HeaderHtml template.HTML
	FooterHtml template.HTML
}

type GroupFormFields []FormFields
type GroupFieldHeaders []string

// 預設FormPanel(struct)
func NewFormPanel() *FormPanel {
	return &FormPanel{
		curFieldListIndex: -1,
		Callbacks:         make(Callbacks, 0),
		Layout:            form2.LayoutDefault,
	}
}

func (f *FormPanel) SetTable(table string) *FormPanel {
	f.Table = table
	return f
}

func (f *FormPanel) SetTitle(title string) *FormPanel {
	f.Title = title
	return f
}

func (f *FormPanel) SetDescription(desc string) *FormPanel {
	f.Description = desc
	return f
}

//  AddField 添加表單欄位資訊至FormPanel.FieldList並處理不同表單欄位類型的選項
func (f *FormPanel) AddField(head, field string, filedType db.DatabaseType, formType form2.Type) *FormPanel {
	f.FieldList = append(f.FieldList, FormField{
		Head:        head,
		Field:       field,
		FieldClass:  field,
		TypeName:    filedType,
		Editable:    true,
		Hide:        false,
		TableFields: make(FormFields, 0),
		Placeholder: "輸入" + " " + head,
		FormType:    formType,
		FieldDisplay: FieldDisplay{
			Display: func(value FieldModel) interface{} {
				return value.Value
			},
			DisplayProcessChains: chooseDisplayProcessChains(f.processChains),
		},
	})
	f.curFieldListIndex++

	// GetDefaultOptions 不同表單欄位類型設置不同選項
	op1, op2, js := formType.GetDefaultOptions(field)
	f.FieldOptionExt(op1)
	f.FieldOptionExt2(op2)
	f.FieldOptionExtJS(js)

	// setDefaultDisplayFnOfFormType 設置表單類型函式
	setDefaultDisplayFnOfFormType(f, formType)
	return f
}

// 將參數name、type設置至FormPanel.primaryKey後回傳
func (f *FormPanel) SetPrimaryKey(name string, typ db.DatabaseType) *FormPanel {
	f.primaryKey = primaryKey{Name: name, Type: typ}
	return f
}

// 判斷FormFields[i].Field是否存在參數field，存在則回傳FormFields[i](FormField)
func (f FormFields) FindByFieldName(field string) *FormField {
	for i := 0; i < len(f); i++ {
		if f[i].Field == field {
			return &f[i]
		}
	}
	return nil
}

// FieldMust 該表單欄位必填
func (f *FormPanel) FieldMust() *FormPanel {
	f.FieldList[f.curFieldListIndex].Must = true
	return f
}

// FieldNotAllowAdd 該表單欄位不允許增加
func (f *FormPanel) FieldNotAllowAdd() *FormPanel {
	f.FieldList[f.curFieldListIndex].NotAllowAdd = true
	return f
}

// FieldNotAllowEdit 該表單欄位不能編輯
func (f *FormPanel) FieldNotAllowEdit() *FormPanel {
	f.FieldList[f.curFieldListIndex].Editable = false
	return f
}

// FieldHelpMsg 增加提示資訊
func (f *FormPanel) FieldHelpMsg(s template.HTML) *FormPanel {
	f.FieldList[f.curFieldListIndex].HelpMsg = s
	return f
}

// FieldOptions 欄位選項
func (f *FormPanel) FieldOptions(options FieldOptions) *FormPanel {
	f.FieldList[f.curFieldListIndex].Options = options
	return f
}

// FieldOptionsFromTable 設置表單欄位的選項，第二個參數為顯示的選項名稱
func (f *FormPanel) FieldOptionsFromTable(table, textFieldName, valueFieldName string, process ...OptionTableQueryProcessFn) *FormPanel {
	var fn OptionTableQueryProcessFn
	if len(process) > 0 {
		fn = process[0]
	}
	f.FieldList[f.curFieldListIndex].OptionTable = OptionTable{
		Table:          table,
		TextField:      textFieldName,
		ValueField:     valueFieldName,
		QueryProcessFn: fn,
	}

	return f
}

// FieldDisplay 將參數(函式)添加至FormPanel.FieldList[].Display
func (f *FormPanel) FieldDisplay(filter FieldFilterFn) *FormPanel {
	f.FieldList[f.curFieldListIndex].Display = filter
	return f
}

// FieldPostFilterFn 添加函式func(value PostFieldModel) interface{}
func (f *FormPanel) FieldPostFilterFn(post PostFieldFilterFn) *FormPanel {
	f.FieldList[f.curFieldListIndex].PostFilterFn = post
	return f
}

// SetInsertFn 設置新增函式
func (f *FormPanel) SetInsertFn(fn FormPostFn) *FormPanel {
	f.InsertFn = fn
	return f
}

// SetUpdateFn 設置更新函式
func (f *FormPanel) SetUpdateFn(fn FormPostFn) *FormPanel {
	f.UpdateFn = fn
	return f
}

// SetPostValidator 新增函式func(values form.Values) error至FormPanel.Validator
func (f *FormPanel) SetPostValidator(va FormPostFn) *FormPanel {
	f.Validator = va
	return f
}

// SetPostHook 新增函式func(values form.Values) error至FormPanel.PostHook
func (f *FormPanel) SetPostHook(fn FormPostFn) *FormPanel {
	f.PostHook = fn
	return f
}

// AddXssJsFilter添加func(value FieldModel) interface{}至參數i.processChains([]FieldFilterFn)
func (f *FormPanel) AddXssJsFilter() *FormPanel {
	f.processChains = addXssJsFilter(f.processChains)
	return f
}

// 設置FormPanel.FieldList[].OptionExt(選項)
func (f *FormPanel) FieldOptionExt(m map[string]interface{}) *FormPanel {
	if m == nil {
		return f
	}
	if f.FieldList[f.curFieldListIndex].FormType.IsCode() {
		f.FieldList[f.curFieldListIndex].OptionExt = template.JS(fmt.Sprintf(`
	theme = "%s";
	font_size = %s;
	language = "%s";
	options = %s;
`, m["theme"], m["font_size"], m["language"], m["options"]))
		return f
	}

	m = f.FieldList[f.curFieldListIndex].FormType.FixOptions(m)
	s, _ := json.Marshal(m)

	if f.FieldList[f.curFieldListIndex].OptionExt != template.JS("") {
		ss := string(f.FieldList[f.curFieldListIndex].OptionExt)
		ss = strings.Replace(ss, "}", "", strings.Count(ss, "}"))
		ss = strings.TrimRight(ss, " ")
		ss += ","
		f.FieldList[f.curFieldListIndex].OptionExt = template.JS(ss) + template.JS(strings.Replace(string(s), "{", "", 1))
	} else {
		f.FieldList[f.curFieldListIndex].OptionExt = template.JS(string(s))
	}

	return f
}

// 設置FormPanel.FieldList[].OptionExt2(選項)
func (f *FormPanel) FieldOptionExt2(m map[string]interface{}) *FormPanel {
	if m == nil {
		return f
	}

	m = f.FieldList[f.curFieldListIndex].FormType.FixOptions(m)
	s, _ := json.Marshal(m)

	if f.FieldList[f.curFieldListIndex].OptionExt2 != template.JS("") {
		ss := string(f.FieldList[f.curFieldListIndex].OptionExt2)
		ss = strings.Replace(ss, "}", "", strings.Count(ss, "}"))
		ss = strings.TrimRight(ss, " ")
		ss += ","
		f.FieldList[f.curFieldListIndex].OptionExt2 = template.JS(ss) + template.JS(strings.Replace(string(s), "{", "", 1))
	} else {
		f.FieldList[f.curFieldListIndex].OptionExt2 = template.JS(string(s))
	}

	return f
}

// 設置FormPanel.FieldList[].OptionExt(選項)
func (f *FormPanel) FieldOptionExtJS(js template.JS) *FormPanel {
	if js != template.JS("") {
		f.FieldList[f.curFieldListIndex].OptionExt = js
	}
	return f
}
