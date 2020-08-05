package table

import (
	"goadminapi/context"
	"goadminapi/modules/db"
	"goadminapi/modules/service"
	"goadminapi/plugins/admin/modules/form"
	"goadminapi/plugins/admin/modules/parameter"
	"sync"
	"sync/atomic"

	"goadminapi/template/types"
)

var (
	services service.List
	count    uint32
	lock     sync.Mutex
)

type BaseTable struct {
	Info           *types.InfoPanel
	Form           *types.FormPanel
	Detail         *types.InfoPanel
	CanAdd         bool
	Editable       bool
	Deletable      bool
	Exportable     bool
	OnlyInfo       bool
	OnlyDetail     bool
	OnlyNewForm    bool
	OnlyUpdateForm bool
	PrimaryKey     PrimaryKey
}

type Generator func(ctx *context.Context) Table
type GeneratorList map[string]Generator

type PrimaryKey struct {
	Type db.DatabaseType
	Name string
}

type PanelInfo struct {
	Thead          types.Thead              `json:"thead"`
	InfoList       types.InfoList           `json:"info_list"`
	FilterFormData types.FormFields         `json:"filter_form_data"`
	Paginator      types.PaginatorAttribute `json:"-"`
	Title          string                   `json:"title"`
	Description    string                   `json:"description"`
}

type FormInfo struct {
	FieldList         types.FormFields        `json:"field_list"`
	GroupFieldList    types.GroupFormFields   `json:"group_field_list"`
	GroupFieldHeaders types.GroupFieldHeaders `json:"group_field_headers"`
	Title             string                  `json:"title"`
	Description       string                  `json:"description"`
}

type Table interface {
	GetInfo() *types.InfoPanel
	GetDetail() *types.InfoPanel
	// GetDetailFromInfo() *types.InfoPanel
	GetForm() *types.FormPanel

	// GetCanAdd() bool
	GetEditable() bool
	GetDeletable() bool
	// GetExportable() bool

	GetPrimaryKey() PrimaryKey

	GetData(params parameter.Parameters) (PanelInfo, error)
	GetDataWithIds(params parameter.Parameters) (PanelInfo, error)
	// GetDataWithId(params parameter.Parameters) (FormInfo, error)
	UpdateData(dataList form.Values) error
	InsertData(dataList form.Values) error
	DeleteData(pk string) error

	// GetNewForm() FormInfo

	// GetOnlyInfo() bool
	GetOnlyDetail() bool
	GetOnlyNewForm() bool
	GetOnlyUpdateForm() bool

	// Copy() Table
}

// 透過參數key及gen(function)添加至GeneratorList(map[string]Generator)
func (g GeneratorList) Add(key string, gen Generator) {
	g[key] = gen
}

// 透過參數list判斷GeneratorList已經有該key、value，如果不存在則加入該鍵與值
func (g GeneratorList) Combine(list GeneratorList) GeneratorList {
	for key, gen := range list {
		if _, ok := g[key]; !ok {
			g[key] = gen
		}
	}
	return g
}

// 透過參數gens判斷GeneratorList已經有該key、value，如果不存在則加入該鍵與值
func (g GeneratorList) CombineAll(gens []GeneratorList) GeneratorList {
	for _, list := range gens {
		for key, gen := range list {
			if _, ok := g[key]; !ok {
				g[key] = gen
			}
		}
	}
	return g
}

// 將參數(srv)設置給services(map[string]Service)
func SetServices(srv service.List) {
	lock.Lock()
	defer lock.Unlock()

	if atomic.LoadUint32(&count) != 0 {
		panic("can not initialize twice")
	}
	services = srv
}

// ------------------------table(interface)的方法---------------------

// GetInfo 將參數值設置至base.Info(InfoPanel(struct)).primaryKey
func (base *BaseTable) GetInfo() *types.InfoPanel {
	return base.Info.SetPrimaryKey(base.PrimaryKey.Name, base.PrimaryKey.Type)
}

// GetPrimaryKey return BaseTable.PrimaryKey(struct)
func (base *BaseTable) GetPrimaryKey() PrimaryKey { return base.PrimaryKey }

// GetForm 將參數值設置至BaseTable.Form(FormPanel(struct)).primaryKey
func (base *BaseTable) GetForm() *types.FormPanel {
	return base.Form.SetPrimaryKey(base.PrimaryKey.Name, base.PrimaryKey.Type)
}

// GetDetail 將參數值設置至base.Detail(InfoPanel(struct)).primaryKey中後回傳InfoPanel(struct)
func (base *BaseTable) GetDetail() *types.InfoPanel {
	return base.Detail.SetPrimaryKey(base.PrimaryKey.Name, base.PrimaryKey.Type)
}

// 回傳BaseTable.Editable(是否可以編輯)
func (base *BaseTable) GetEditable() bool { return base.Editable }

// 回傳BaseTable.Deletable(是否可以刪除)
func (base *BaseTable) GetDeletable() bool { return base.Deletable }

// 回傳是否只有更新表單功能
func (base *BaseTable) GetOnlyUpdateForm() bool { return base.OnlyUpdateForm }

// 回傳是否只有新增表單功能
func (base *BaseTable) GetOnlyNewForm() bool { return base.OnlyNewForm }

// 回傳是否只有取得細節的權限
func (base *BaseTable) GetOnlyDetail() bool { return base.OnlyDetail }
