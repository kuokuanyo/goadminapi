package types

import (
	"html/template"

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

	Editable    bool `json:"editable"`
	NotAllowAdd bool `json:"not_allow_add"`
	Must        bool `json:"must"`
	Hide        bool `json:"hide"`

	Width int `json:"width"`

	InputWidth int `json:"input_width"`
	HeadWidth  int `json:"head_width"`

	Joins Joins `json:"-"`

	Divider      bool   `json:"divider"`
	DividerTitle string `json:"divider_title"`

	HelpMsg template.HTML `json:"help_msg"`

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
	FieldList         FormFields
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

// 判斷FormFields[i].Field是否存在參數field，存在則回傳FormFields[i](FormField)
func (f FormFields) FindByFieldName(field string) *FormField {
	for i := 0; i < len(f); i++ {
		if f[i].Field == field {
			return &f[i]
		}
	}
	return nil
}
