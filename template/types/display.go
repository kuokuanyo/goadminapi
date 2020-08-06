package types

import (
	"fmt"
	"goadminapi/modules/config"
	"goadminapi/template/types/form"
	"html/template"
	"strings"
)

// globalDisplayProcessChains 類別為[]FieldFilterFn，FieldFilterFn類別為func(value FieldModel) interface{}
var globalDisplayProcessChains = make(DisplayProcessFnChains, 0)

type FieldDisplay struct {
	Display              FieldFilterFn
	DisplayProcessChains DisplayProcessFnChains
}

type DisplayProcessFnChains []FieldFilterFn

// 將參數f(func(value FieldModel) interface{})加入DisplayProcessFnChains([]FieldFilterFn)
func (d DisplayProcessFnChains) Add(f FieldFilterFn) DisplayProcessFnChains {
	return append(d, f)
}

// 添加func(value FieldModel) interface{}至參數chains([]FieldFilterFn)
func addXssJsFilter(chains DisplayProcessFnChains) DisplayProcessFnChains {
	chains = chains.Add(func(value FieldModel) interface{} {
		replacer := strings.NewReplacer("<script>", "&lt;script&gt;", "</script>", "&lt;/script&gt;")
		return replacer.Replace(value.Value)
	})
	return chains
}

// 複製[]FieldFilterFn後回傳
func (d DisplayProcessFnChains) Copy() DisplayProcessFnChains {
	if len(d) == 0 {
		return make(DisplayProcessFnChains, 0)
	} else {
		var newDisplayProcessFnChains = make(DisplayProcessFnChains, len(d))
		copy(newDisplayProcessFnChains, d)
		return newDisplayProcessFnChains
	}
}

// 如果參數長度大於0則回傳參數，否則複製全域變數globalDisplayProcessChains([]FieldFilterFn)後回傳
func chooseDisplayProcessChains(internal DisplayProcessFnChains) DisplayProcessFnChains {
	if len(internal) > 0 {
		return internal
	}
	return globalDisplayProcessChains.Copy()
}

// setDefaultDisplayFnOfFormType 設置表單類型函式
func setDefaultDisplayFnOfFormType(f *FormPanel, typ form.Type) {
	if typ.IsMultiFile() {
		f.FieldList[f.curFieldListIndex].Display = func(value FieldModel) interface{} {
			if value.Value == "" {
				return ""
			}
			arr := strings.Split(value.Value, ",")
			res := "["
			for i, item := range arr {
				if i == len(arr)-1 {
					res += "'" + config.GetStore().URL(item) + "']"
				} else {
					res += "'" + config.GetStore().URL(item) + "',"
				}
			}
			return res
		}
	}
	if typ.IsSelect() {
		f.FieldList[f.curFieldListIndex].Display = func(value FieldModel) interface{} {
			return strings.Split(value.Value, ",")
		}
	}
}

// IsNotSelectRes判斷參數類別，如果為HTML、[]string、[][]string則回傳false
func (f FieldDisplay) IsNotSelectRes(v interface{}) bool {
	switch v.(type) {
	case template.HTML:
		return false
	case []string:
		return false
	case [][]string:
		return false
	default:
		return true
	}
}

// 判斷條件後回傳數值(interface{})
func (f FieldDisplay) ToDisplay(value FieldModel) interface{} {
	// FieldDisplay.Display(func(value FieldModel) interface{})
	val := f.Display(value)
	// IsNotSelectRes判斷參數類別，如果為HTML、[]string、[][]string則回傳false
	if len(f.DisplayProcessChains) > 0 && f.IsNotSelectRes(val) {
		valStr := fmt.Sprintf("%v", val)
		for _, process := range f.DisplayProcessChains {
			valStr = fmt.Sprintf("%v", process(FieldModel{
				Row:   value.Row,
				Value: valStr,
				ID:    value.ID,
			}))
		}
		return valStr
	}

	return val
}
