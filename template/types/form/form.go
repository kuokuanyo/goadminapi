package form

import (
	"html/template"
)

type Type uint8

const (
	Default Type = iota
	Text
	SelectSingle
	Select
	IconPicker
	SelectBox
	File
	Multifile
	Password
	RichText
	Datetime
	DatetimeRange
	Radio
	Checkbox
	CheckboxStacked
	CheckboxSingle
	Email
	Date
	DateRange
	Url
	Ip
	Color
	Array
	Currency
	Rate
	Number
	Table
	NumberRange
	TextArea
	Custom
	Switch
	Code
	Slider
)

// 判斷參數t是否存在於AllType([]Type)裡，如果存在則回傳t，不存在回傳參數def
func CheckType(t, def Type) Type {
	for _, item := range AllType {
		if t == item {
			return t
		}
	}
	return def
}

var AllType = []Type{Default, Text, Array, SelectSingle, Select, IconPicker, SelectBox, File, Multifile, Password,
	RichText, Datetime, DatetimeRange, Checkbox, CheckboxStacked, Radio, Table, Email, Url, Ip, Color, Currency, Number, NumberRange,
	TextArea, Custom, Switch, Code, Rate, Slider, Date, DateRange, CheckboxSingle}

type Layout uint8

const (
	LayoutDefault Layout = iota
	LayoutTwoCol
	LayoutThreeCol
	LayoutFourCol
	LayoutFiveCol
	LayoutSixCol
	LayoutFlow
	LayoutTab
)

func (l Layout) Default() bool {
	return l == LayoutDefault
}

func (l Layout) Flow() bool {
	return l == LayoutFlow
}

func (t Type) Name() string {
	switch t {
	case Default:
		return "Default"
	case Text:
		return "Text"
	case SelectSingle:
		return "SelectSingle"
	case Select:
		return "Select"
	case IconPicker:
		return "IconPicker"
	case SelectBox:
		return "SelectBox"
	case File:
		return "File"
	case Table:
		return "Table"
	case Multifile:
		return "Multifile"
	case Password:
		return "Password"
	case RichText:
		return "RichText"
	case Rate:
		return "Rate"
	case Checkbox:
		return "Checkbox"
	case CheckboxStacked:
		return "CheckboxStacked"
	case CheckboxSingle:
		return "CheckboxSingle"
	case Date:
		return "Date"
	case DateRange:
		return "DateRange"
	case Datetime:
		return "Datetime"
	case DatetimeRange:
		return "DatetimeRange"
	case Radio:
		return "Radio"
	case Slider:
		return "Slider"
	case Array:
		return "Array"
	case Email:
		return "Email"
	case Url:
		return "Url"
	case Ip:
		return "Ip"
	case Color:
		return "Color"
	case Currency:
		return "Currency"
	case Number:
		return "Number"
	case NumberRange:
		return "NumberRange"
	case TextArea:
		return "TextArea"
	case Custom:
		return "Custom"
	case Switch:
		return "Switch"
	case Code:
		return "Code"
	default:
		panic("wrong form type")
	}
}

func (t Type) String() string {
	switch t {
	case Default:
		return "default"
	case Text:
		return "text"
	case SelectSingle:
		return "select_single"
	case Select:
		return "select"
	case IconPicker:
		return "iconpicker"
	case SelectBox:
		return "selectbox"
	case File:
		return "file"
	case Table:
		return "table"
	case Multifile:
		return "multi_file"
	case Password:
		return "password"
	case RichText:
		return "richtext"
	case Rate:
		return "rate"
	case Checkbox:
		return "checkbox"
	case CheckboxStacked:
		return "checkbox_stacked"
	case CheckboxSingle:
		return "checkbox_single"
	case Date:
		return "datetime"
	case DateRange:
		return "datetime_range"
	case Datetime:
		return "datetime"
	case DatetimeRange:
		return "datetime_range"
	case Radio:
		return "radio"
	case Slider:
		return "slider"
	case Array:
		return "array"
	case Email:
		return "email"
	case Url:
		return "url"
	case Ip:
		return "ip"
	case Color:
		return "color"
	case Currency:
		return "currency"
	case Number:
		return "number"
	case NumberRange:
		return "number_range"
	case TextArea:
		return "textarea"
	case Custom:
		return "custom"
	case Switch:
		return "switch"
	case Code:
		return "code"
	default:
		panic("wrong form type")
	}
}

func (l Layout) String() string {
	switch l {
	case LayoutDefault:
		return "LayoutDefault"
	case LayoutTwoCol:
		return "LayoutTwoCol"
	case LayoutThreeCol:
		return "LayoutThreeCol"
	case LayoutFourCol:
		return "LayoutFourCol"
	case LayoutFiveCol:
		return "LayoutFiveCol"
	case LayoutSixCol:
		return "LayoutSixCol"
	case LayoutFlow:
		return "LayoutFlow"
	case LayoutTab:
		return "LayoutTab"
	default:
		return "LayoutDefault"
	}
}

func GetLayoutFromString(s string) Layout {
	switch s {
	case "LayoutDefault":
		return LayoutDefault
	case "LayoutTwoCol":
		return LayoutTwoCol
	case "LayoutThreeCol":
		return LayoutThreeCol
	case "LayoutFourCol":
		return LayoutFourCol
	case "LayoutFiveCol":
		return LayoutFiveCol
	case "LayoutSixCol":
		return LayoutSixCol
	case "LayoutFlow":
		return LayoutFlow
	case "LayoutTab":
		return LayoutTab
	default:
		return LayoutDefault
	}
}

// 取得日期時間選項
func getDateTimeRangeOptions(f Type) (map[string]interface{}, map[string]interface{}) {
	format := "YYYY-MM-DD HH:mm:ss"
	if f == DateRange {
		format = "YYYY-MM-DD"
	}
	m := map[string]interface{}{
		"format": format,
	}
	m1 := map[string]interface{}{
		"format":     format,
		"useCurrent": false,
	}
	return m, m1
}

// 取得日期時間選項
func getDateTimeOptions(f Type) map[string]interface{} {
	format := "YYYY-MM-DD HH:mm:ss"
	if f == Date {
		format = "YYYY-MM-DD"
	}
	m := map[string]interface{}{
		"format":           format,
		"allowInputToggle": true,
	}
	return m
}

// GetDefaultOptions 設置表單欄位選項
func (t Type) GetDefaultOptions(field string) (map[string]interface{}, map[string]interface{}, template.JS) {
	switch t {
	case File, Multifile:
		return map[string]interface{}{
			"overwriteInitial":     true,
			"initialPreviewAsData": true,
			"browseLabel":          "瀏覽",
			"showRemove":           false,
			"previewClass":         "preview-" + field,
			"showUpload":           false,
			"allowedFileTypes":     []string{"image"},
		}, nil, ""
	case Slider:
		return map[string]interface{}{
			"type":     "single",
			"prettify": false,
			"hasGrid":  true,
			"max":      100,
			"min":      1,
			"step":     1,
			"postfix":  "",
		}, nil, ""
	case DatetimeRange:
		op1, op2 := getDateTimeRangeOptions(DatetimeRange)
		return op1, op2, ""
	case Datetime:
		return getDateTimeOptions(Datetime), nil, ""
	case Date:
		return getDateTimeOptions(Date), nil, ""
	case DateRange:
		op1, op2 := getDateTimeRangeOptions(DateRange)
		return op1, op2, ""
	case Code:
		return nil, nil, `
	theme = "monokai";
	font_size = 14;
	language = "html";
	options = {useWorker: false};
`
	}
	return nil, nil, ""
}

func (t Type) IsArray() bool {
	return t == Array
}

func (t Type) IsTable() bool {
	return t == Table
}

// 判斷t(unit8)是否符合條件
func (t Type) IsSelect() bool {
	return t == Select || t == SelectSingle || t == SelectBox || t == Radio || t == Switch ||
		t == Checkbox || t == CheckboxStacked || t == CheckboxSingle
}

// 判斷t(unit8)是否符合條件
func (t Type) IsSingleSelect() bool {
	return t == SelectSingle || t == Radio || t == Switch || t == CheckboxSingle
}

// 判斷t(unit8)是否符合條件，是否有多個選擇
func (t Type) IsMultiSelect() bool {
	return t == Select || t == SelectBox || t == Checkbox || t == CheckboxStacked
}

func (l Layout) Col() int {
	if l == LayoutTwoCol {
		return 2
	}
	if l == LayoutThreeCol {
		return 3
	}
	if l == LayoutFourCol {
		return 4
	}
	if l == LayoutFiveCol {
		return 5
	}
	if l == LayoutSixCol {
		return 6
	}
	return 0
}

func (t Type) IsCode() bool {
	return t == Code
}

func (t Type) FixOptions(m map[string]interface{}) map[string]interface{} {
	switch t {
	case Slider:
		if _, ok := m["type"]; !ok {
			m["type"] = "single"
		}
		if _, ok := m["prettify"]; !ok {
			m["prettify"] = false
		}
		if _, ok := m["hasGrid"]; !ok {
			m["hasGrid"] = true
		}
		return m
	}
	return m
}

func (t Type) IsFile() bool {
	// File = 6
	// Multifile = 7
	return t == File || t == Multifile
}

func (t Type) IsMultiFile() bool {
	return t == Multifile
}

// 是否設置值為範圍
func (t Type) IsRange() bool {
	return t == DatetimeRange || t == NumberRange
}

// 是否為自定義
func (t Type) IsCustom() bool {
	return t == Custom
}

// 判斷條件後設置[]template.HTML
func (t Type) SelectedLabel() []template.HTML {
	if t == Select || t == SelectSingle || t == SelectBox {
		return []template.HTML{"selected", ""}
	}
	if t == Radio || t == Switch || t == Checkbox || t == CheckboxStacked || t == CheckboxSingle {
		return []template.HTML{"checked", ""}
	}
	return []template.HTML{"", ""}
}
