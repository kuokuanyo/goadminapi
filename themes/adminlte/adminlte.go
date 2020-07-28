package adminlte

import (
	adminTemplate "goadminapi/template"
	"strings"

	"goadminapi/template/components"
	"goadminapi/template/types"
	"goadminapi/themes/common"
	"goadminapi/themes/resource"

	"github.com/gobuffalo/packr/v2"
)

type Theme struct {
	ThemeName string
	components.Base
	*common.BaseTheme
}

var Adminlte = Theme{
	ThemeName: "adminlte",
	Base: components.Base{
		Attribute: types.Attribute{
			TemplateList: TemplateList,
		},
	},
	BaseTheme: &common.BaseTheme{
		AssetPaths:   resource.AssetPaths,
		TemplateList: TemplateList,
	},
}

func init() {
	adminTemplate.Add("adminlte", &Adminlte)
}

// --------------------template(interface)所有方法
func (t *Theme) Name() string {
	return t.ThemeName
}

func (t *Theme) GetAsset(path string) ([]byte, error) {
	path = strings.Replace(path, "/assets/dist", "", -1)
	box := packr.New("adminlte", "./resource/assets/dist")
	return box.Find(path)
}

