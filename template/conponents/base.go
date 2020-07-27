package components

import (
	"goadminapi/modules/config"
	"goadminapi/modules/menu"
	"goadminapi/template/types"
	"goadminapi/template/types/form"
)

type Base struct {
	Attribute types.Attribute
}

// ----------------Base(struct)有Template(interface)所有方法-----------------
func (b Base) Col() types.ColAttribute {
	return &ColAttribute{
		Name:      "col",
		Size:      "col-md-2",
		Content:   "",
		Attribute: b.Attribute,
	}
}

func (b Base) Row() types.RowAttribute {
	return &RowAttribute{
		Name:      "row",
		Content:   "",
		Attribute: b.Attribute,
	}
}

func (b Base) Form() types.FormAttribute {
	return &FormAttribute{
		Name:         "form",
		Content:      []types.FormField{},
		Url:          "/",
		Method:       "post",
		HiddenFields: make(map[string]string),
		Layout:       form.LayoutDefault,
		Title:        "edit",
		Attribute:    b.Attribute,
		CdnUrl:       config.GetAssetUrl(),
		HeadWidth:    2,
		InputWidth:   8,
	}
}

func (b Base) Table() types.TableAttribute {
	return &TableAttribute{
		Name:      "table",
		Thead:     make(types.Thead, 0),
		InfoList:  make([]map[string]types.InfoItem, 0),
		Type:      "table",
		Style:     "hover",
		Layout:    "auto",
		Attribute: b.Attribute,
	}
}

func (b Base) DataTable() types.DataTableAttribute {
	return &DataTableAttribute{
		TableAttribute: *(b.Table().
			SetStyle("hover").
			SetName("data-table").
			SetType("data-table").(*TableAttribute)),
		Attribute: b.Attribute,
	}
}

func (b Base) TreeView() types.TreeViewAttribute {
	return &TreeViewAttribute{
		Name:      "treeview",
		Attribute: b.Attribute,
	}
}

func (b Base) Tree() types.TreeAttribute {
	return &TreeAttribute{
		Name:      "tree",
		Tree:      make([]menu.Item, 0),
		Attribute: b.Attribute,
	}
}

func (b Base) Tabs() types.TabsAttribute {
	return &TabsAttribute{
		Name:      "tabs",
		Attribute: b.Attribute,
	}
}

func (b Base) Alert() types.AlertAttribute {
	return &AlertAttribute{
		Name:      "alert",
		Attribute: b.Attribute,
	}
}

