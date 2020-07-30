package table

import (
	"database/sql"
	"goadminapi/context"
	"goadminapi/modules/config"
	"goadminapi/modules/db"
	"goadminapi/template/types"
	"goadminapi/template/types/form"
	"strings"

	"goadminapi/template"
)

type SystemTable struct {
	conn db.Connection
	c    *config.Config
}

// 將參數設置至SystemTable(struct)後回傳
func NewSystemTable(conn db.Connection, c *config.Config) *SystemTable {
	return &SystemTable{conn: conn, c: c}
}

// 設置success至LabelAttribute.Type
func label() types.LabelAttribute {
	return template.Get(config.GetTheme()).Label().SetType("success")
}

// 將[]string轉換成[]interface{}
func interfaces(arr []string) []interface{} {
	var iarr = make([]interface{}, len(arr))

	for key, v := range arr {
		iarr[key] = v
	}

	return iarr
}

// connection 設置SQL(struct)
func (s *SystemTable) connection() *db.SQL {
	return db.WithDriver(s.conn)
}

// 添加欄位資訊
func (s *SystemTable) GetManagerTable(ctx *context.Context) (managerTable Table) {
	// NewDefaultTable 將參數值設置至預設DefaultTable(struct)
	// 預設Config(struct)，driver設為參數，主鍵id
	// *******managerTable為DefaultTable(struct)，包含table(interface)所有方法***********
	managerTable = NewDefaultTable(DefaultConfigWithDriver(config.GetDatabases().GetDefault().Driver))

	// GetInfo 將參數值設置至base.Info(InfoPanel(struct)).primaryKey
	// AddXssJsFilter添加func(value FieldModel) interface{}至參數i.processChains([]FieldFilterFn)
	// HideFilterArea InfoPanel.IsHideFilterArea = true
	info := managerTable.GetInfo().AddXssJsFilter().HideFilterArea()

	// AddField 添加欄位資訊至InfoPanel.FieldList
	// FieldSortable 設置為可排序
	// FieldFilterable 設置為可篩選並添加篩選的表單欄位資訊至FilterFormFields
	// FieldJoin 添加join其他資料表資訊
	// FieldDisplay 將參數添加至InfoPanel.FieldList[].Display
	// ****************用戶名、暱稱、角色皆可以篩選**********************
	info.AddField("ID", "id", "INT").FieldSortable()
	info.AddField("用戶名稱", "username", db.Varchar).FieldFilterable()
	info.AddField("暱稱", "name", db.Varchar).FieldFilterable()
	// 用戶角色需先與role_user資料表join後再與roles資料表join
	info.AddField("角色", "name", db.Varchar).
		FieldJoin(types.Join{
			Table:     "role_users",
			JoinField: "user_id",
			Field:     "id",
		}).
		FieldJoin(types.Join{
			Table:     "roles",
			JoinField: "id",
			Field:     "role_id",
			BaseTable: "role_users",
		}).
		FieldDisplay(func(model types.FieldModel) interface{} {
			labels := template.HTML("")
			// 設置success至LabelAttribute.Type
			labelTpl := label().SetType("success")

			labelValues := strings.Split(model.Value, types.JoinFieldValueDelimiter)
			for key, label := range labelValues {
				if key == len(labelValues)-1 {
					labels += labelTpl.SetContent(template.HTML(label)).GetContent()
				} else {
					labels += labelTpl.SetContent(template.HTML(label)).GetContent() + "<br><br>"
				}
			}

			if labels == template.HTML("") {
				return "沒有角色"
			}

			return labels
		}).FieldFilterable()
	info.AddField("建立時間", "created_at", db.Timestamp)
	info.AddField("更新時間", "updated_at", db.Timestamp)

	info.SetTable("users").SetTitle("用戶管理").SetDescription("用戶管理").
		// 設置刪除函式，除了刪除users資料表，其他關連表資料也必須刪除
		SetDeleteFn(func(idArr []string) error {

			var ids = interfaces(idArr)

			_, txErr := s.connection().WithTransaction(func(tx *sql.Tx) (e error, i map[string]interface{}) {
				// 必須要刪除role_users、user_permissions、users資料表的資料
				deleteUserRoleErr := s.connection().WithTx(tx).
					Table("role_users").
					WhereIn("user_id", ids).
					Delete()
				if db.CheckError(deleteUserRoleErr, db.DELETE) {
					return deleteUserRoleErr, nil
				}
				deleteUserPermissionErr := s.connection().WithTx(tx).
					Table("user_permissions").
					WhereIn("user_id", ids).
					Delete()
				if db.CheckError(deleteUserPermissionErr, db.DELETE) {
					return deleteUserPermissionErr, nil
				}
				deleteUserErr := s.connection().WithTx(tx).
					Table("users").
					WhereIn("id", ids).
					Delete()
				if db.CheckError(deleteUserErr, db.DELETE) {
					return deleteUserErr, nil
				}
				return nil, nil
			})

			return txErr
		})

	// GetForm 將參數值設置至BaseTable.Form(FormPanel(struct)).primaryKey
	// AddXssJsFilter添加func(value FieldModel) interface{}至參數i.processChains([]FieldFilterFn)
	formList := managerTable.GetForm().AddXssJsFilter()

	formList.AddField("ID", "id", "INT", form.Default).FieldNotAllowEdit().FieldNotAllowAdd()
}
