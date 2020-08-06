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

func (t Type) IsMultiFile() bool {
	return t == Multifile
}

// 是否設置值為範圍
func (t Type) IsRange() bool {
	return t == DatetimeRange || t == NumberRange
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