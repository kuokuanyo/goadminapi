package models

import (
	"goadminapi/modules/db"
	"goadminapi/modules/db/dialect"
)

type UserModel struct {
	Base          `json:"-"`
	Id            int64             `json:"id"`
	Name          string            `json:"name"`
	UserName      string            `json:"user_name"`
	Password      string            `json:"password"`
	Avatar        string            `json:"avatar"`
	RememberToken string            `json:"remember_token"`
	Permissions   []PermissionModel `json:"permissions"`
	MenuIds       []int64           `json:"menu_ids"`
	Roles         []RoleModel       `json:"role"`
	Level         string            `json:"level"`
	LevelName     string            `json:"level_name"`
	CreatedAt     string            `json:"created_at"`
	UpdatedAt     string            `json:"updated_at"`
}

// 設定資料表名給UserModel(struct)
func User(tablename string) UserModel {
	return UserModel{Base: Base{TableName: tablename}}
}

// 將參數設置給UserModel.conn後回傳
func (t UserModel) SetConn(con db.Connection) UserModel {
	t.Conn = con
	return t
}

// 透過參數username尋找符合的資料並設置至UserModel
func (t UserModel) FindByUserName(username interface{}) UserModel {
	// Table藉由給定的table回傳sql(struct)
	// sql 語法 where = ...，回傳 SQl struct
	// First回傳第一筆符合的資料
	item, _ := t.Table(t.TableName).Where("username", "=", username).First()
	// 將item資訊設置至UserModel後回傳
	return t.MapToModel(item)
}

// 查詢role藉由user
func (t UserModel) WithRoles() UserModel {
	roleModel, _ := t.Table("role_users").
		LeftJoin("roles", "roles.id", "=", "role_users.role_id").
		Where("user_id", "=", t.Id).
		Select("roles.id", "roles.name", "roles.slug",
			"roles.created_at", "roles.updated_at").
		All()

	for _, role := range roleModel {
		//Role取得初始化role model
		//MapToModel將role設置至rolemodel中
		t.Roles = append(t.Roles, Role().MapToModel(role))
	}

	if len(t.Roles) > 0 {
		//設置slug、name(roles表)至user model
		t.Level = t.Roles[0].Slug
		t.LevelName = t.Roles[0].Name
	}

	return t
}

// 取得該用戶的role_id
func (t UserModel) GetAllRoleId() []interface{} {

	var ids = make([]interface{}, len(t.Roles))

	for key, role := range t.Roles {
		ids[key] = role.Id
	}

	return ids
}

// 查詢user的permission
func (t UserModel) WithPermissions() UserModel {

	var permissions = make([]map[string]interface{}, 0)

	//可能會有多個role id(可以設定多個role)
	roleIds := t.GetAllRoleId()

	//----------------------------------------------------------------------------------------------
	// permission會依照user_id以及role_id取得不同的權限，因此需要做下列兩次判斷
	//----------------------------------------------------------------------------------------------
	if len(roleIds) > 0 {
		//查詢role_id的permission
		permissions, _ = t.Table("role_permissions").
			LeftJoin("permissions", "permissions.id", "=", "role_permissions.permission_id").
			WhereIn("role_id", roleIds).
			Select("permissions.http_method", "permissions.http_path",
				"permissions.id", "permissions.name", "permissions.slug",
				"permissions.created_at", "permissions.updated_at").
			All()
	}

	// 跟上面role的方式一樣(藉由user_id取得permission)
	// 可能有多個permission
	userPermissions, _ := t.Table("user_permissions").
		LeftJoin("permissions", "permissions.id", "=", "user_permissions.permission_id").
		Where("user_id", "=", t.Id).
		Select("permissions.http_method", "permissions.http_path",
			"permissions.id", "permissions.name", "permissions.slug",
			"permissions.created_at", "permissions.updated_at").
		All()

	permissions = append(permissions, userPermissions...)

	// 加入權限
	for i := 0; i < len(permissions); i++ {
		exist := false
		//如果裡面已經有相同的權限加入，就停止迴圈
		for j := 0; j < len(t.Permissions); j++ {
			if t.Permissions[j].Id == permissions[i]["id"] {
				exist = true
				break
			}
		}

		//Permission()為初始化Permission model
		//MapToModel 將role設置至Permission model
		if exist {
			continue
		}
		t.Permissions = append(t.Permissions, Permission().MapToModel(permissions[i]))
	}

	return t
}

// 取得可用menu
func (t UserModel) WithMenus() UserModel {

	var menuIdsModel []map[string]interface{}

	// 判斷是否為超級管理員
	if t.IsSuperAdmin() {
		menuIdsModel, _ = t.Table("goadmin_role_menu").
			LeftJoin("goadmin_menu", "goadmin_menu.id", "=", "goadmin_role_menu.menu_id").
			Select("menu_id", "parent_id").
			All()
	} else {
		// 取得menuIdsModel藉由role_id
		rolesId := t.GetAllRoleId()
		if len(rolesId) > 0 {
			menuIdsModel, _ = t.Table("goadmin_role_menu").
				LeftJoin("goadmin_menu", "goadmin_menu.id", "=", "goadmin_role_menu.menu_id").
				WhereIn("goadmin_role_menu.role_id", rolesId).
				Select("menu_id", "parent_id").
				All()
		}
	}

	var menuIds []int64

	// 將menu_id加入menuIds中
	for _, mid := range menuIdsModel {
		if parentId, ok := mid["parent_id"].(int64); ok && parentId != 0 {
			for _, mid2 := range menuIdsModel {
				if mid2["menu_id"].(int64) == mid["parent_id"].(int64) {
					menuIds = append(menuIds, mid["menu_id"].(int64))
					break
				}
			}
		} else {
			menuIds = append(menuIds, mid["menu_id"].(int64))
		}
	}
	t.MenuIds = menuIds
	return t
}

// 將參數password設置至UserModel.UserModel並且更新dialect.H{"password": password,}
func (t UserModel) UpdatePwd(password string) UserModel {

	_, _ = t.Table(t.TableName).
		Where("id", "=", t.Id).
		Update(dialect.H{
			"password": password,
		})

	t.Password = password
	return t
}

// 判斷是否為超級管理員
func (t UserModel) IsSuperAdmin() bool {
	for _, per := range t.Permissions {
		if len(per.HttpPath) > 0 && per.HttpPath[0] == "*" && per.HttpMethod[0] == "" {
			return true
		}
	}
	return false
}

// 判斷是否為空
func (t UserModel) IsEmpty() bool {
	return t.Id == int64(0)
}

// 將取得的值(參數m)設置usermodel從map中
func (t UserModel) MapToModel(m map[string]interface{}) UserModel {
	t.Id, _ = m["id"].(int64)
	t.Name, _ = m["name"].(string)
	t.UserName, _ = m["username"].(string)
	t.Password, _ = m["password"].(string)
	t.Avatar, _ = m["avatar"].(string)
	t.RememberToken, _ = m["remember_token"].(string)
	t.CreatedAt, _ = m["created_at"].(string)
	t.UpdatedAt, _ = m["updated_at"].(string)
	return t
}
