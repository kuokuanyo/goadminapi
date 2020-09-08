package table

import (
	"database/sql"
	"errors"
	"goadminapi/context"
	"goadminapi/html"
	"goadminapi/modules/collection"
	"goadminapi/modules/config"
	"goadminapi/modules/db"
	"goadminapi/modules/db/dialect"
	"goadminapi/template/types"
	"goadminapi/template/types/form"
	"strconv"
	"strings"
	"time"

	"goadminapi/template"

	"goadminapi/plugins/admin/models"
	form2 "goadminapi/plugins/admin/modules/form"
	tmpl "html/template"

	"golang.org/x/crypto/bcrypt"
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

func (s *SystemTable) table(table string) *db.SQL {
	return s.connection().Table(table)
}

// link 新增一個連結(HTML)
func link(url, content string) tmpl.HTML {
	return html.AEl().
		SetAttr("href", url).
		SetContent(template.HTML(content)).
		Get()
}

func EncodePassword(pwd []byte) string {
	hash, err := bcrypt.GenerateFromPassword(pwd, bcrypt.DefaultCost)
	if err != nil {
		return ""
	}
	return string(hash[:])
}

// GetManagerTable 新增用戶頁面、表單欄位、細節頁面欄位資訊與函式
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
	info.AddField("userID", "userid", db.Varchar).FieldFilterable()
	info.AddField("用戶名稱", "username", db.Varchar).FieldFilterable()
	info.AddField("電話", "phone", db.Varchar).FieldFilterable()
	info.AddField("信箱", "email", db.Varchar).FieldFilterable()
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
	info.AddField("用戶照片", "pictureURL", db.Varchar)
	info.AddField("建立時間", "created_at", db.Timestamp)
	info.AddField("更新時間", "updated_at", db.Timestamp)

	info.SetTable("users").SetTitle("User").SetDescription("User Manage").
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

	// AddField 添加表單欄位資訊至FormPanel.FieldList並處理不同表單欄位類型的選項
	// FieldNotAllowEdit 該表單欄位不能編輯
	// FieldNotAllowAdd 該表單欄位不允許增加
	// FieldHelpMsg 增加提示資訊
	// FieldMust 該表單欄位必填
	// FieldOptionsFromTable 從資料表設置表單欄位的選項，第二個參數為顯示的選項名稱
	// FieldDisplay 將參數(函式)添加至FormPanel.FieldList[].Display
	formList.AddField("ID", "id", "INT", form.Default).FieldNotAllowEdit().FieldNotAllowAdd()
	formList.AddField("使用者名稱", "username", db.Varchar, form.Text).FieldMust()
	formList.AddField("電話號碼", "phone", db.Varchar, form.Text).
		FieldHelpMsg(template.HTML("用來登入")).FieldMust()
	formList.AddField("信箱", "email", db.Varchar, form.Text).FieldMust()
	formList.AddField("角色", "role_id", db.Varchar, form.Select).
		FieldOptionsFromTable("roles", "slug", "id").
		FieldDisplay(func(model types.FieldModel) interface{} {
			var roles []string

			if model.ID == "" {
				return roles
			}
			roleModel, _ := s.table("role_users").Select("role_id").
				Where("user_id", "=", model.ID).All()
			for _, v := range roleModel {
				roles = append(roles, strconv.FormatInt(v["role_id"].(int64), 10))
			}
			return roles
		}).FieldHelpMsg(template.HTML("沒有對應的選項?") +
		link("/admin/info/roles/new", "立刻新增一個"))
	formList.AddField("權限", "permission_id", db.Varchar, form.Select).
		FieldOptionsFromTable("permissions", "slug", "id").
		FieldDisplay(func(model types.FieldModel) interface{} {
			var permissions []string

			if model.ID == "" {
				return permissions
			}
			permissionModel, _ := s.table("user_permissions").
				Select("permission_id").Where("user_id", "=", model.ID).All()
			for _, v := range permissionModel {
				permissions = append(permissions, strconv.FormatInt(v["permission_id"].(int64), 10))
			}
			return permissions
		}).FieldHelpMsg(template.HTML("沒有對應的選項?") +
		link("/admin/info/permission/new", "立刻新增一個"))
	formList.AddField("密碼", "password", db.Varchar, form.Password).
		FieldDisplay(func(value types.FieldModel) interface{} {
			return ""
		})
	formList.AddField("請確認密碼", "password_again", db.Varchar, form.Password).
		FieldDisplay(func(value types.FieldModel) interface{} {
			return ""
		})

	formList.SetTable("users").SetTitle("User").SetDescription("User Manage")
	// 設置更新函式，必須先刪除角色及權限後再新增
	formList.SetUpdateFn(func(values form2.Values) error {
		if values.IsEmpty("name", "username") {
			return errors.New("username and password can not be empty")
		}

		user := models.UserWithId(values.Get("id")).SetConn(s.conn)
		password := values.Get("password")
		if password != "" {
			if password != values.Get("password_again") {
				return errors.New("password does not match")
			}
			password = EncodePassword([]byte(values.Get("password")))
		}

		// WithTransaction 取得tx(struct)，會持續並行Rollback、commit
		_, txErr := s.connection().WithTransaction(func(tx *sql.Tx) (e error, i map[string]interface{}) {
			_, updateUserErr := user.WithTx(tx).Update(values.Get("username"), password, values.Get("name"), values.Get("avatar"))
			if db.CheckError(updateUserErr, db.UPDATE) {
				return updateUserErr, nil
			}
			delRoleErr := user.WithTx(tx).DeleteRoles()
			if db.CheckError(delRoleErr, db.DELETE) {
				return delRoleErr, nil
			}

			for i := 0; i < len(values["role_id[]"]); i++ {
				_, addRoleErr := user.WithTx(tx).AddRole(values["role_id[]"][i])
				if db.CheckError(addRoleErr, db.INSERT) {
					return addRoleErr, nil
				}
			}

			delPermissionErr := user.WithTx(tx).DeletePermissions()
			if db.CheckError(delPermissionErr, db.DELETE) {
				return delPermissionErr, nil
			}

			for i := 0; i < len(values["permission_id[]"]); i++ {
				_, addPermissionErr := user.WithTx(tx).AddPermission(values["permission_id[]"][i])
				if db.CheckError(addPermissionErr, db.INSERT) {
					return addPermissionErr, nil
				}
			}
			return nil, nil
		})
		return txErr
	})

	// 設置新增資料函式
	formList.SetInsertFn(func(values form2.Values) error {
		if values.IsEmpty("username", "password", "email") {
			return errors.New("username,password and email can not be empty")
		}

		password := values.Get("password")
		if password != values.Get("password_again") {
			return errors.New("password does not match")
		}

		_, txErr := s.connection().WithTransaction(func(tx *sql.Tx) (e error, i map[string]interface{}) {
			user, createUserErr := models.User("users").WithTx(tx).SetConn(s.conn).New(values.Get("userid"), values.Get("username"),
				values.Get("phone"), values.Get("email"), EncodePassword([]byte(values.Get("password"))), values.Get("pictureURL"))
			if db.CheckError(createUserErr, db.INSERT) {
				return createUserErr, nil
			}

			// 新增角色、權限
			for i := 0; i < len(values["role_id[]"]); i++ {
				_, addRoleErr := user.WithTx(tx).AddRole(values["role_id[]"][i])
				if db.CheckError(addRoleErr, db.INSERT) {
					return addRoleErr, nil
				}
			}
			for i := 0; i < len(values["permission_id[]"]); i++ {
				_, addPermissionErr := user.WithTx(tx).AddPermission(values["permission_id[]"][i])
				if db.CheckError(addPermissionErr, db.INSERT) {
					return addPermissionErr, nil
				}
			}
			return nil, nil
		})
		return txErr
	})

	// 處理細節函式
	detail := managerTable.GetDetail()
	detail.AddField("ID", "id", db.Int)
	detail.AddField("userID", "userid", db.Varchar)
	detail.AddField("使用者名稱", "username", db.Varchar)
	detail.AddField("電話號碼", "phone", db.Varchar)
	detail.AddField("信箱", "email", db.Varchar)
	detail.AddField("用戶照片", "pictureURL", db.Varchar)
	detail.AddField("角色", "roles", db.Varchar).
		FieldDisplay(func(model types.FieldModel) interface{} {
			labelModels, _ := s.table("role_users").
				Select("roles.name").
				LeftJoin("roles", "roles.id", "=", "role_users.role_id").
				Where("user_id", "=", model.ID).
				All()

			labels := template.HTML("")
			labelTpl := label().SetType("success")

			for key, label := range labelModels {
				if key == len(labelModels)-1 {
					labels += labelTpl.SetContent(template.HTML(label["name"].(string))).GetContent()
				} else {
					labels += labelTpl.SetContent(template.HTML(label["name"].(string))).GetContent() + "<br><br>"
				}
			}
			if labels == template.HTML("") {
				return "沒有角色"
			}
			return labels
		})
	detail.AddField("權限", "roles", db.Varchar).
		FieldDisplay(func(model types.FieldModel) interface{} {
			permissionModel, _ := s.table("user_permissions").
				Select("permissions.name").
				LeftJoin("permissions", "permissions.id", "=", "user_permissions.permission_id").
				Where("user_id", "=", model.ID).
				All()

			permissions := template.HTML("")
			permissionTpl := label().SetType("success")

			for key, label := range permissionModel {
				if key == len(permissionModel)-1 {
					permissions += permissionTpl.SetContent(template.HTML(label["name"].(string))).GetContent()
				} else {
					permissions += permissionTpl.SetContent(template.HTML(label["name"].(string))).GetContent() + "<br><br>"
				}
			}
			return permissions
		})
	detail.AddField("建立時間", "created_at", db.Timestamp)
	detail.AddField("更新時間", "updated_at", db.Timestamp)
	return
}

// GetRolesTable 新增角色頁面、表單欄位資訊與函式
func (s *SystemTable) GetRolesTable(ctx *context.Context) (roleTable Table) {
	// NewDefaultTable 將參數值設置至預設DefaultTable(struct)
	// *******roleTable為DefaultTable(struct)，包含table(interface)所有方法***********
	roleTable = NewDefaultTable(DefaultConfigWithDriver(config.GetDatabases().GetDefault().Driver))

	// GetInfo 將參數值設置至base.Info(InfoPanel(struct)).primaryKey
	info := roleTable.GetInfo().AddXssJsFilter().HideFilterArea()

	// 增加角色頁面欄位資訊與函式
	// *********用戶角色、標誌可以篩選*************
	info.AddField("ID", "id", db.Int).FieldSortable()
	info.AddField("角色", "name", db.Varchar).FieldFilterable()
	info.AddField("標誌", "slug", db.Varchar).FieldFilterable()
	info.AddField("建立時間", "created_at", db.Timestamp)
	info.AddField("更新時間", "updated_at", db.Timestamp)

	// 刪除也必須刪除其他關連表資料
	info.SetTable("roles").SetTitle("Roles").SetDescription("Roles Manage").
		SetDeleteFn(func(idArr []string) error {
			var ids = interfaces(idArr)
			_, txErr := s.connection().WithTransaction(func(tx *sql.Tx) (e error, i map[string]interface{}) {
				deleteRoleUserErr := s.connection().WithTx(tx).
					Table("role_users").
					WhereIn("role_id", ids).
					Delete()
				if db.CheckError(deleteRoleUserErr, db.DELETE) {
					return deleteRoleUserErr, nil
				}

				deleteRoleMenuErr := s.connection().WithTx(tx).
					Table("role_menu").
					WhereIn("role_id", ids).
					Delete()
				if db.CheckError(deleteRoleMenuErr, db.DELETE) {
					return deleteRoleMenuErr, nil
				}

				deleteRolePermissionErr := s.connection().WithTx(tx).
					Table("role_permissions").
					WhereIn("role_id", ids).
					Delete()
				if db.CheckError(deleteRolePermissionErr, db.DELETE) {
					return deleteRolePermissionErr, nil
				}

				deleteRolesErr := s.connection().WithTx(tx).
					Table("roles").
					WhereIn("id", ids).
					Delete()
				if db.CheckError(deleteRolesErr, db.DELETE) {
					return deleteRolesErr, nil
				}
				return nil, nil
			})
			return txErr
		})

	// 處理表單欄位資訊與函式(更新、新增)
	formList := roleTable.GetForm().AddXssJsFilter()
	formList.AddField("ID", "id", db.Int, form.Default).FieldNotAllowEdit().FieldNotAllowAdd()
	formList.AddField("角色", "name", db.Varchar, form.Text).FieldHelpMsg(template.HTML("不能重複")).FieldMust()
	formList.AddField("標誌", "slug", db.Varchar, form.Text).FieldHelpMsg(template.HTML("不能重複")).FieldMust()
	formList.AddField("權限", "permission_id", db.Varchar, form.SelectBox).
		// 從permissions取得選項
		FieldOptionsFromTable("permissions", "name", "id").
		FieldDisplay(func(model types.FieldModel) interface{} {
			var permissions = make([]string, 0)
			if model.ID == "" {
				return permissions
			}

			perModel, _ := s.table("role_permissions").
				Select("permission_id").
				Where("role_id", "=", model.ID).
				All()
			for _, v := range perModel {
				permissions = append(permissions, strconv.FormatInt(v["permission_id"].(int64), 10))
			}
			return permissions
		}).FieldHelpMsg(template.HTML("沒有對應的選項?") +
		link("/admin/info/permission/new", "立刻新增一個"))
	formList.AddField("建立時間", "updated_at", db.Timestamp, form.Default).FieldNotAllowAdd()
	formList.AddField("更新時間", "created_at", db.Timestamp, form.Default).FieldNotAllowAdd()
	formList.SetTable("roles").SetTitle("Roles").SetDescription("Roles Manage")

	// 設置更新函式，必須先刪除所有相關資料表的權限，再新增設置的權限
	formList.SetUpdateFn(func(values form2.Values) error {
		if models.Role().SetConn(s.conn).IsSlugExist(values.Get("slug"), values.Get("id")) {
			return errors.New("slug exists")
		}

		role := models.RoleWithId(values.Get("id")).SetConn(s.conn)

		_, txErr := s.connection().WithTransaction(func(tx *sql.Tx) (e error, i map[string]interface{}) {
			_, updateRoleErr := role.WithTx(tx).Update(values.Get("name"), values.Get("slug"))
			if db.CheckError(updateRoleErr, db.UPDATE) {
				return updateRoleErr, nil
			}

			delPermissionErr := role.WithTx(tx).DeletePermissions()
			if db.CheckError(delPermissionErr, db.DELETE) {
				return delPermissionErr, nil
			}
			for i := 0; i < len(values["permission_id[]"]); i++ {
				_, addPermissionErr := role.WithTx(tx).AddPermission(values["permission_id[]"][i])
				if db.CheckError(addPermissionErr, db.INSERT) {
					return addPermissionErr, nil
				}
			}
			return nil, nil
		})
		return txErr
	})

	// 設置新增資料函式
	formList.SetInsertFn(func(values form2.Values) error {
		if models.Role().SetConn(s.conn).IsSlugExist(values.Get("slug"), "") {
			return errors.New("slug exists")
		}

		_, txErr := s.connection().WithTransaction(func(tx *sql.Tx) (e error, i map[string]interface{}) {
			role, createRoleErr := models.Role().WithTx(tx).SetConn(s.conn).New(values.Get("name"), values.Get("slug"))
			if db.CheckError(createRoleErr, db.INSERT) {
				return createRoleErr, nil
			}

			for i := 0; i < len(values["permission_id[]"]); i++ {
				_, addPermissionErr := role.WithTx(tx).AddPermission(values["permission_id[]"][i])
				if db.CheckError(addPermissionErr, db.INSERT) {
					return addPermissionErr, nil
				}
			}
			return nil, nil
		})
		return txErr
	})
	return
}

// GetPermissionTable 新增權限頁面、表單欄位欄位資訊與函式
func (s *SystemTable) GetPermissionTable(ctx *context.Context) (permissionTable Table) {
	// NewDefaultTable 將參數值設置至預設DefaultTable(struct)
	// *******permissionTable為DefaultTable(struct)，包含table(interface)所有方法***********
	permissionTable = NewDefaultTable(DefaultConfigWithDriver(config.GetDatabases().GetDefault().Driver))

	// GetInfo 將參數值設置至base.Info(InfoPanel(struct)).primaryKey
	info := permissionTable.GetInfo().AddXssJsFilter().HideFilterArea()

	// 增加權限頁面欄位資訊與函式
	// *********用戶名稱、標誌可以篩選*************
	info.AddField("ID", "id", db.Int).FieldSortable()
	info.AddField("權限", "name", db.Varchar).FieldFilterable()
	info.AddField("標誌", "slug", db.Varchar).FieldFilterable()
	info.AddField("方法", "http_method", db.Varchar).
		FieldDisplay(func(value types.FieldModel) interface{} {
			if value.Value == "" {
				return "All methods"
			}
			return value.Value
		})
	info.AddField("路徑", "http_path", db.Varchar).
		FieldDisplay(func(model types.FieldModel) interface{} {
			pathArr := strings.Split(model.Value, "\n")
			res := ""
			for i := 0; i < len(pathArr); i++ {
				if i == len(pathArr)-1 {
					res += string(label().SetContent(template.HTML(pathArr[i])).GetContent())
				} else {
					res += string(label().SetContent(template.HTML(pathArr[i])).GetContent()) + "<br><br>"
				}
			}
			return res
		})
	info.AddField("建立時間", "created_at", db.Timestamp)
	info.AddField("更新時間", "updated_at", db.Timestamp)

	// 刪除也必須刪除其他關連表資料
	info.SetTable("permissions").
		SetTitle("Permission").
		SetDescription("Permission Manage").
		SetDeleteFn(func(idArr []string) error {
			var ids = interfaces(idArr)

			_, txErr := s.connection().WithTransaction(func(tx *sql.Tx) (e error, i map[string]interface{}) {
				deleteRolePermissionErr := s.connection().WithTx(tx).
					Table("role_permissions").
					WhereIn("permission_id", ids).
					Delete()
				if db.CheckError(deleteRolePermissionErr, db.DELETE) {
					return deleteRolePermissionErr, nil
				}

				deleteUserPermissionErr := s.connection().WithTx(tx).
					Table("user_permissions").
					WhereIn("permission_id", ids).
					Delete()
				if db.CheckError(deleteUserPermissionErr, db.DELETE) {
					return deleteUserPermissionErr, nil
				}

				deletePermissionsErr := s.connection().WithTx(tx).
					Table("permissions").
					WhereIn("id", ids).
					Delete()
				if deletePermissionsErr != nil {
					return deletePermissionsErr, nil
				}

				return nil, nil
			})

			return txErr
		})

	// 處理表單欄位資訊與函式(更新、新增)
	formList := permissionTable.GetForm().AddXssJsFilter()
	formList.AddField("ID", "id", db.Int, form.Default).FieldNotAllowEdit().FieldNotAllowAdd()
	formList.AddField("權限", "name", db.Varchar, form.Text).FieldHelpMsg(template.HTML("不能重複")).FieldMust()
	formList.AddField("標誌", "slug", db.Varchar, form.Text).FieldHelpMsg(template.HTML("不能重複")).FieldMust()
	formList.AddField("方法", "http_method", db.Varchar, form.Select).
		FieldOptions(types.FieldOptions{
			{Value: "GET", Text: "GET"},
			{Value: "PUT", Text: "PUT"},
			{Value: "POST", Text: "POST"},
			{Value: "DELETE", Text: "DELETE"},
			{Value: "PATCH", Text: "PATCH"},
			{Value: "OPTIONS", Text: "OPTIONS"},
			{Value: "HEAD", Text: "HEAD"},
		}).
		FieldDisplay(func(model types.FieldModel) interface{} {
			return strings.Split(model.Value, ",")
		}).
		// FieldPostFilterFn 添加函式func(value PostFieldModel) interface{}
		FieldPostFilterFn(func(model types.PostFieldModel) interface{} {
			return strings.Join(model.Value, ",")
		}).
		FieldHelpMsg(template.HTML("如果為空代表所有方法"))
	formList.AddField("路徑", "http_path", db.Varchar, form.TextArea).
		FieldPostFilterFn(func(model types.PostFieldModel) interface{} {
			return strings.TrimSpace(model.Value.Value())
		}).
		FieldHelpMsg(template.HTML("路徑不包含全局前綴且必須一行設置一個路徑，換行輸入新路徑"))
	formList.AddField("建立時間", "updated_at", db.Timestamp, form.Default).FieldNotAllowAdd()
	formList.AddField("更新時間", "created_at", db.Timestamp, form.Default).FieldNotAllowAdd()

	formList.SetTable("permissions").
		SetTitle("Permission").
		SetDescription("Permission Manage").
		// SetPostValidator 新增函式func(values form.Values) error至FormPanel.Validator
		SetPostValidator(func(values form2.Values) error {
			if values.IsEmpty("slug", "http_path", "name") {
				return errors.New("slug or http_path or name should not be empty")
			}
			if models.Permission().SetConn(s.conn).IsSlugExist(values.Get("slug"), values.Get("id")) {
				return errors.New("slug exists")
			}
			return nil
			// SetPostHook 新增函式func(values form.Values) error至FormPanel.PostHook
		}).SetPostHook(func(values form2.Values) error {
		_, err := s.connection().Table("permissions").
			Where("id", "=", values.Get("id")).Update(dialect.H{
			"updated_at": time.Now().Format("2006-01-02 15:04:05"),
		})
		return err
	})
	return
}

// GetMenuTable 新增菜單頁面、表單欄位欄位資訊與函式
func (s *SystemTable) GetMenuTable(ctx *context.Context) (menuTable Table) {
	// NewDefaultTable 將參數值設置至預設DefaultTable(struct)
	// *******permissionTable為DefaultTable(struct)，包含table(interface)所有方法***********
	menuTable = NewDefaultTable(DefaultConfigWithDriver(config.GetDatabases().GetDefault().Driver))

	// 增加菜單頁面欄位資訊與函式
	info := menuTable.GetInfo().AddXssJsFilter().HideFilterArea()
	info.AddField("ID", "id", db.Int).FieldSortable()
	info.AddField("父級", "parent_id", db.Int)
	info.AddField("menu名", "title", db.Varchar)
	info.AddField("圖標", "icon", db.Varchar)
	info.AddField("路徑", "uri", db.Varchar)
	info.AddField("角色", "roles", db.Varchar)
	info.AddField("標頭", "header", db.Varchar)
	info.AddField("建立時間", "created_at", db.Timestamp)
	info.AddField("更新時間", "updated_at", db.Timestamp)

	// 設置刪除函式，與菜單有關聯的資料都必須刪除
	info.SetTable("menu").SetTitle("Menu").SetDescription("Menu Manage").
		SetDeleteFn(func(idArr []string) error {
			var ids = interfaces(idArr)
			_, txErr := s.connection().WithTransaction(func(tx *sql.Tx) (e error, i map[string]interface{}) {
				deleteRoleMenuErr := s.connection().WithTx(tx).
					Table("role_menu").
					WhereIn("menu_id", ids).
					Delete()
				if db.CheckError(deleteRoleMenuErr, db.DELETE) {
					return deleteRoleMenuErr, nil
				}

				deleteMenusErr := s.connection().WithTx(tx).
					Table("menu").
					WhereIn("id", ids).
					Delete()
				if db.CheckError(deleteMenusErr, db.DELETE) {
					return deleteMenusErr, nil
				}

				return nil, map[string]interface{}{}
			})
			return txErr
		})

	var parentIDOptions = types.FieldOptions{
		{
			Text:  "ROOT",
			Value: "0",
		},
	}

	// ***********處理父級選項******************
	// 先取出所有的父級(父級的parent_id = 0)
	allMenus, _ := s.connection().Table("menu").
		Where("parent_id", "=", 0).
		Select("id", "title").
		OrderBy("order", "asc").
		All()
		// 所有父級的id
	allMenuIDs := make([]interface{}, len(allMenus))

	if len(allMenuIDs) > 0 {
		for i := 0; i < len(allMenus); i++ {
			allMenuIDs[i] = allMenus[i]["id"]
		}
		// 取得父級下的menu
		secondLevelMenus, _ := s.connection().Table("menu").
			WhereIn("parent_id", allMenuIDs).
			Select("id", "title", "parent_id").
			All()

		// 轉換成[]map[string]interface{}
		secondLevelMenusCol := collection.Collection(secondLevelMenus)

		for i := 0; i < len(allMenus); i++ {
			// 新增父級的選項名稱
			parentIDOptions = append(parentIDOptions, types.FieldOption{
				TextHTML: "&nbsp;&nbsp;┝  " + template.HTML(allMenus[i]["title"].(string)),
				Value:    strconv.Itoa(int(allMenus[i]["id"].(int64))),
			})

			// 取得父級之下的menu
			col := secondLevelMenusCol.Where("parent_id", "=", allMenus[i]["id"].(int64))
			for i := 0; i < len(col); i++ {
				parentIDOptions = append(parentIDOptions, types.FieldOption{
					TextHTML: "&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;┝  " +
						template.HTML(col[i]["title"].(string)),
					Value: strconv.Itoa(int(col[i]["id"].(int64))),
				})
			}
		}
	}
	// ***********處理父級選項*********************

	// 處理表單欄位資訊與函式(更新、新增)
	formList := menuTable.GetForm().AddXssJsFilter()
	formList.AddField("ID", "id", db.Int, form.Default).FieldNotAllowEdit().FieldNotAllowAdd()
	formList.AddField("父級", "parent_id", db.Int, form.SelectSingle).
		FieldOptions(parentIDOptions).
		FieldDisplay(func(model types.FieldModel) interface{} {
			var menuItem []string
			if model.ID == "" {
				return menuItem
			}

			menuModel, _ := s.table("menu").Select("parent_id").Find(model.ID)
			menuItem = append(menuItem, strconv.FormatInt(menuModel["parent_id"].(int64), 10))
			return menuItem
		})
	formList.AddField("menu名", "title", db.Varchar, form.Text).FieldMust()
	formList.AddField("標頭", "header", db.Varchar, form.Text)
	formList.AddField("圖標", "icon", db.Varchar, form.IconPicker)
	formList.AddField("路徑", "uri", db.Varchar, form.Text)
	// 角色選項
	formList.AddField("角色", "roles", db.Int, form.Select).
		FieldOptionsFromTable("roles", "slug", "id").
		FieldDisplay(func(model types.FieldModel) interface{} {
			var roles []string
			if model.ID == "" {
				return roles
			}

			roleModel, _ := s.table("role_menu").
				Select("role_id").
				Where("menu_id", "=", model.ID).
				All()

			for _, v := range roleModel {
				roles = append(roles, strconv.FormatInt(v["role_id"].(int64), 10))
			}
			return roles
		})
	formList.AddField("更新時間", "updated_at", db.Timestamp, form.Default).FieldNotAllowAdd()
	formList.AddField("建立時間", "created_at", db.Timestamp, form.Default).FieldNotAllowAdd()

	formList.SetTable("menu").
		SetTitle("Menu").
		SetDescription("Menu Manage")

	return
}
