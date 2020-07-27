package form

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

// 判斷t(unit8)是否符合條件
func (t Type) IsSelect() bool {
	// Select = 3
	// SelectSingle = 2
	// SelectBox = 5
	// Radio = 12
	// Switch = 2
	return t == Select || t == SelectSingle || t == SelectBox || t == Radio || t == Switch ||
		t == Checkbox || t == CheckboxStacked || t == CheckboxSingle
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