package types

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"strings"

	"goadminapi/modules/config"
	"goadminapi/modules/db"
	"goadminapi/modules/utils"

	"goadminapi/context"
	"goadminapi/plugins/admin/modules"
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

	Editable    bool `json:"editable"`      // 允許編輯
	NotAllowAdd bool `json:"not_allow_add"` // 不允許增加
	Must        bool `json:"must"`          // 該欄位必填
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

func (fo FieldOptions) Copy() FieldOptions {
	newOptions := make(FieldOptions, len(fo))
	copy(newOptions, fo)
	return newOptions
}

// updateValue 將FormField值處理後(設置選項...等資訊)回傳
func (f *FormField) updateValue(id, val string, res map[string]interface{}, typ PostType, sql *db.SQL) *FormField {
	m := FieldModel{
		ID:       id,
		Value:    val,
		Row:      res,
		PostType: typ,
	}

	// ----------一般都為false----------------------
	if f.isBelongToATable() {
		if f.FormType.IsSelect() {
			if len(f.OptionsArr) == 0 && f.OptionArrInitFn != nil {
				f.OptionsArr = f.OptionArrInitFn(m)
				for i := 0; i < len(f.OptionsArr); i++ {
					f.OptionsArr[i] = f.OptionsArr[i].SetSelectedLabel(f.FormType.SelectedLabel())
				}
			} else {
				f.setOptionsFromSQL(sql)
				if f.FormType.IsSingleSelect() {
					values := f.ToDisplayStringArray(m)
					f.OptionsArr = make([]FieldOptions, len(values))
					for k, value := range values {
						f.OptionsArr[k] = f.Options.Copy().SetSelected(value, f.FormType.SelectedLabel())
					}
				} else {
					values := f.ToDisplayStringArrayArray(m)
					f.OptionsArr = make([]FieldOptions, len(values))
					for k, value := range values {
						f.OptionsArr[k] = f.Options.Copy().SetSelected(value, f.FormType.SelectedLabel())
					}
				}
			}
		} else {
			f.ValueArr = f.ToDisplayStringArray(m)
		}
	} else {
		if f.FormType.IsSelect() { // 判斷表般欄位是否是用select的(ex:角色、權限欄位)
			if len(f.Options) == 0 && f.OptionInitFn != nil {
				// SetSelectedLabel 設置值至FieldOptions[k].SelectedLabel
				f.Options = f.OptionInitFn(m).SetSelectedLabel(f.FormType.SelectedLabel())
			} else { // ---------一般執行此條件------------------
				// setOptionsFromSQL 從資料庫中設置值(欄位名稱及數值)至FormField.Options
				f.setOptionsFromSQL(sql)
				// SetSelected 判斷條件後將參數加入FieldOptions[k].SelectedLabel
				f.Options.SetSelected(f.ToDisplay(m), f.FormType.SelectedLabel())
			}
		} else if f.FormType.IsArray() {
			f.ValueArr = f.ToDisplayStringArray(m)
		} else { // ----一般表單欄位執行此條件---------
			f.Value = f.ToDisplayHTML(m) // -------一般為空值----------
			if f.FormType.IsFile() {     // ex:頭像欄位
				if f.Value != template.HTML("") {
					f.Value2 = config.GetStore().URL(string(f.Value))
				}
			}
		}
	}
	return f
}

// UpdateDefaultValue 將FormField值處理後(設置選項...等資訊)回傳
func (f *FormField) UpdateDefaultValue(sql *db.SQL) *FormField {
	// -------一般都沒有設置FormField.Default，為空值---------------------
	f.Value = f.Default // template.HTML
	return f.updateValue("", string(f.Value), make(map[string]interface{}), PostTypeCreate, sql)
}

func (f *FormField) fillCustom(src string) string {
	t := template.New("custom")
	t, err := t.Parse(src)
	if err != nil {
		panic(err)
	}
	buf := new(bytes.Buffer)
	err = t.Execute(buf, f)
	if err != nil {
		panic(err)
	}
	return buf.String()
}

// FillCustomContent 填寫至定義內容
func (f *FormField) FillCustomContent() *FormField {
	// TODO: optimize
	if f.CustomContent != "" {
		f.CustomContent = template.HTML(f.fillCustom(string(f.CustomContent)))
	}
	if f.CustomJs != "" {
		f.CustomJs = template.JS(f.fillCustom(string(f.CustomJs)))
	}
	if f.CustomCss != "" {
		f.CustomCss = template.CSS(f.fillCustom(string(f.CustomCss)))
	}
	return f
}

// FillCustomContent 判斷是否為自定義，處理FormField(struct)
func (f FormFields) FillCustomContent() FormFields {
	for i := range f {
		if f[i].FormType.IsCustom() { // 表單欄位判斷是否為自定義，一般都不是
			f[i] = *(f[i]).FillCustomContent()
		}
	}
	return f
}

func (f FormFields) RemoveNotShow() FormFields {
	ff := f
	for i := 0; i < len(ff); {
		if ff[i].FatherFormType == form2.Table {
			ff = append(ff[:i], ff[i+1:]...)
		} else { // ------一般都執行此條件-----------
			i++
		}
	}
	return ff
}

// FieldsWithValue 設置選項、預設值...等資訊至FormFields(帶有預設值)
func (f *FormPanel) FieldsWithValue(pk, id string, columns []string, res map[string]interface{}, sql func() *db.SQL) FormFields {
	var (
		list  = make(FormFields, 0)
		hasPK = false
	)

	// field為表單上所有欄位的資訊
	for _, field := range f.FieldList {
		// rowValue為該欄位的值
		rowValue := field.GetRawValue(columns, res[field.Field])

		if field.FatherField != "" {
			f.FieldList.FindTableField(field.Field, field.FatherField).UpdateValue(id, rowValue, res, sql())
		} else if field.FormType.IsTable() {
			list = append(list, field)
		} else {
			// ------------一般都執行此條件----------------
			// UpdateValue 將FormField值處理後(設置選項、預設值...等資訊)回傳
			list = append(list, *(field.UpdateValue(id, rowValue, res, sql())))
		}

		if field.Field == pk {
			hasPK = true
		}
	}

	// hasPK判斷是否有primary key
	if !hasPK {
		list = list.Add(FormField{
			Head:       pk,
			FieldClass: pk,
			Field:      pk,
			Value:      template.HTML(id),
			FormType:   form2.Default,
			Hide:       true,
		})
	}

	// FillCustomContent 判斷是否為自定義，處理FormField(struct)
	return list.FillCustomContent()
}

// Add add FormField to FormFields([]FormField)
func (f FormFields) Add(field FormField) FormFields {
	return append(f, field)
}

// FieldsWithDefaultValue 將表單欄位處理後(設置選項...等資訊)加入FormFields([]FormField)中
func (f *FormPanel) FieldsWithDefaultValue(sql ...func() *db.SQL) FormFields {
	var list = make(FormFields, 0)

	// FormPanel.FieldList 為表單的欄位資訊
	// 下面迴圈將表單欄位處理後(設置選項...等資訊)加入list
	for _, v := range f.FieldList {
		// 判斷欄位是否允許添加，例如ID欄位無法手動增加
		if v.allowAdd() {
			v.Editable = true
			if v.FatherField != "" {
				if len(sql) > 0 {
					f.FieldList.FindTableField(v.Field, v.FatherField).UpdateDefaultValue(sql[0]())
				} else {
					f.FieldList.FindTableField(v.Field, v.FatherField).UpdateDefaultValue(nil)
				}
			} else if v.FormType.IsTable() {
				list = append(list, v)
			} else {
				if len(sql) > 0 {
					// ---------一般都執行此函式-------------
					// UpdateDefaultValue 將FormField值處理後(設置選項...等資訊)回傳
					list = append(list, *(v.UpdateDefaultValue(sql[0]())))
				} else {

					list = append(list, *(v.UpdateDefaultValue(nil)))
				}
			}
		}
	}

	return list.FillCustomContent().RemoveNotShow()
}

func (f *FormField) isNotBelongToATable() bool {
	return f.FatherField == "" && !f.FatherFormType.IsTable()
}

// GroupField 先判斷條件後處理FormField，最後將FormField與TabHeader加入至groupFormList與groupHeaders
func (f *FormPanel) GroupField(sql ...func() *db.SQL) ([]FormFields, []string) {
	var (
		groupFormList = make([]FormFields, 0)
		groupHeaders  = make([]string, 0)
	)

	// FormPanel.TabGroups [][]string
	// 判斷條件
	for index, group := range f.TabGroups {
		list := make(FormFields, 0)
		for _, fieldName := range group {
			field := f.FieldList.FindByFieldName(fieldName)
			if field != nil && field.isNotBelongToATable() && field.allowAdd() {
				field.Editable = true
				if field.FormType.IsTable() {
					for z := 0; z < len(field.TableFields); z++ {
						if len(sql) > 0 {
							field.TableFields[z] = *(field.TableFields[z].UpdateDefaultValue(sql[0]()))
						} else {
							field.TableFields[z] = *(field.TableFields[z].UpdateDefaultValue(nil))
						}
					}
					list = append(list, *field)
				} else {
					if len(sql) > 0 {
						// 在template\types\form.go
						// UpdateDefaultValue首先對FieldOptions([]FieldOption)執行迴圈，判斷條件後將參數(html)設置至FieldOptions[k].SelectedLabel後回傳
						// 最後判斷條件後將參數f.FormType.SelectedLabel()([]template.HTML)加入FieldOptions[k].SelectedLabel，回傳FormField
						// FillCustomContent(填寫自定義內容)對FormFields([]FormField)執行迴圈，判斷條件後設置FormField，最後回傳FormFields([]FormField)
						list = append(list, *(field.UpdateDefaultValue(sql[0]())))
					} else {
						list = append(list, *(field.UpdateDefaultValue(nil)))
					}
				}
			}
		}
		groupFormList = append(groupFormList, list.FillCustomContent())
		// TabHeaders []string
		groupHeaders = append(groupHeaders, f.TabHeaders[index])
	}

	return groupFormList, groupHeaders
}

// GetRawValue為取得該欄位的值
func (f *FormField) GetRawValue(columns []string, v interface{}) string {
	isJSON := len(columns) == 0
	return modules.AorB(isJSON || modules.InArray(columns, f.Field),
		db.GetValueFromDatabaseType(f.TypeName, v, isJSON).String(), "")
}

// UpdateValue 將FormField值處理後(設置選項、預設值...等資訊)回傳
func (f *FormField) UpdateValue(id, val string, res map[string]interface{}, sql *db.SQL) *FormField {
	return f.updateValue(id, val, res, PostTypeUpdate, sql)
}

func (f *FormPanel) GroupFieldWithValue(pk, id string, columns []string, res map[string]interface{}, sql func() *db.SQL) ([]FormFields, []string) {
	var (
		groupFormList = make([]FormFields, 0)
		groupHeaders  = make([]string, 0)
		hasPK         = false
	)

	if len(f.TabGroups) > 0 {
		for index, group := range f.TabGroups {
			list := make(FormFields, 0)
			for _, fieldName := range group {
				field := f.FieldList.FindByFieldName(fieldName)
				if field != nil && field.isNotBelongToATable() {
					if field.FormType.IsTable() {
						for z := 0; z < len(field.TableFields); z++ {
							rowValue := field.TableFields[z].GetRawValue(columns, res[field.TableFields[z].Field])
							if field.TableFields[z].Field == pk {
								hasPK = true
							}
							field.TableFields[z] = *(field.TableFields[z].UpdateValue(id, rowValue, res, sql()))
						}
						list = append(list, *field)
					} else {
						if field.Field == pk {
							hasPK = true
						}
						rowValue := field.GetRawValue(columns, res[field.Field])
						list = append(list, *(field.UpdateValue(id, rowValue, res, sql())))
					}
				}
			}

			groupFormList = append(groupFormList, list.FillCustomContent())
			groupHeaders = append(groupHeaders, f.TabHeaders[index])
		}

		if len(groupFormList) > 0 && !hasPK {
			groupFormList[len(groupFormList)-1] = groupFormList[len(groupFormList)-1].Add(FormField{
				Head:       pk,
				FieldClass: pk,
				Field:      pk,
				Value:      template.HTML(id),
				Hide:       true,
			})
		}
	}

	return groupFormList, groupHeaders
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

// FindByFieldName 判斷FormFields[i].Field是否存在參數field，存在則回傳FormFields[i](FormField)
func (f FormFields) FindByFieldName(field string) *FormField {
	for i := 0; i < len(f); i++ {
		if f[i].Field == field {
			return &f[i]
		}
	}
	return nil
}

// allowAdd 允許增加
func (f *FormField) allowAdd() bool {
	return !f.NotAllowAdd
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

func (f *FormField) isBelongToATable() bool {
	return f.FatherField != "" && f.FatherFormType.IsTable()
}

// AddXssJsFilter 添加func(value FieldModel) interface{}至參數i.processChains([]FieldFilterFn)
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

// setOptionsFromSQL 從資料庫中設置值(欄位名稱及數值)至FormField.Options
func (f *FormField) setOptionsFromSQL(sql *db.SQL) {
	if sql != nil && f.OptionTable.Table != "" && len(f.Options) == 0 {
		// Select將參數設置至SQL(struct).Fields並且設置SQL(struct).Functions
		sql.Table(f.OptionTable.Table).Select(f.OptionTable.ValueField, f.OptionTable.TextField)

		if f.OptionTable.QueryProcessFn != nil {
			f.OptionTable.QueryProcessFn(sql)
		}

		// 返回所有符合查詢的結果
		queryRes, err := sql.All()
		if err == nil {
			for _, item := range queryRes {
				f.Options = append(f.Options, FieldOption{
					Value: fmt.Sprintf("%v", item[f.OptionTable.ValueField]), // ex: id
					Text:  fmt.Sprintf("%v", item[f.OptionTable.TextField]),  // ex:slug的值
				})
			}
		}

		if f.OptionTable.ProcessFn != nil {
			f.Options = f.OptionTable.ProcessFn(f.Options)
		}
	}
}

func (f FormFields) FindTableField(field, father string) *FormField {
	// FindByFieldName 判斷FormFields[i].Field是否存在參數field，存在則回傳FormFields[i](FormField)
	ff := f.FindByFieldName(father)
	return ff.TableFields.FindByFieldName(field)
}

// SetSelectedLabel 設置值至FieldOptions[k].SelectedLabel
func (fo FieldOptions) SetSelectedLabel(labels []template.HTML) FieldOptions {
	for k := range fo {
		if fo[k].Selected {
			fo[k].SelectedLabel = labels[0]
		} else {
			fo[k].SelectedLabel = labels[1]
		}
	}
	return fo
}

// SetSelected 判斷條件後將參數加入FieldOptions[k].SelectedLabel
func (fo FieldOptions) SetSelected(val interface{}, labels []template.HTML) FieldOptions {

	if valArr, ok := val.([]string); ok {
		for k := range fo {
			text := fo[k].Text // ex:slug的值 or http_method的值
			if text == "" {
				text = string(fo[k].TextHTML)
			}
			fo[k].Selected = utils.InArray(valArr, fo[k].Value) || utils.InArray(valArr, text) // -----一般為false-----
			if fo[k].Selected {
				fo[k].SelectedLabel = labels[0]
			} else { // -----一般執行此條件--------
				// -----一般都為空值--------
				fo[k].SelectedLabel = labels[1]
			}
		}
	} else {
		for k := range fo {
			text := fo[k].Text
			if text == "" {
				text = string(fo[k].TextHTML)
			}
			fo[k].Selected = fo[k].Value == val || text == val
			if fo[k].Selected {
				fo[k].SelectedLabel = labels[0]
			} else {
				fo[k].SelectedLabel = labels[1]
			}
		}
	}

	return fo
}
