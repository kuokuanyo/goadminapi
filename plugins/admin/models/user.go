package models

import (
	"database/sql"
	"goadminapi/modules/config"
	"goadminapi/modules/db"
	"goadminapi/modules/db/dialect"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"
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

// 透過參數取得UserModel(struct)
func UserWithId(id string) UserModel {
	idInt, _ := strconv.Atoi(id)
	// GetAuthUserTable return globalCfg.AuthUserTable
	return UserModel{Base: Base{TableName: config.GetAuthUserTable()}, Id: int64(idInt)}
}

// 透過參數(id)取得UserModel(struct)，將值設置至UserModel
func (t UserModel) Find(id interface{}) UserModel {
	// 藉由id取的符合資料
	item, _ := t.Table(t.TableName).Find(id)
	// 將取得的值(參數item)設置usermodel從map中
	return t.MapToModel(item)
}

// 透過參數username尋找符合的資料並設置至UserModel
func (t UserModel) FindByUserName(username interface{}) UserModel {
	item, _ := t.Table(t.TableName).Where("username", "=", username).First()

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

// New create a user model.
func (t UserModel) New(username, password, name, avatar string) (UserModel, error) {

	id, err := t.WithTx(t.Tx).Table(t.TableName).Insert(dialect.H{
		"username": username,
		"password": password,
		"name":     name,
		"avatar":   avatar,
	})

	t.Id = id
	t.UserName = username
	t.Password = password
	t.Avatar = avatar
	t.Name = name

	return t, err
}

// CheckRoleId check the role of the user model.
func (t UserModel) CheckRoleId(roleId string) bool {
	checkRole, _ := t.Table("role_users").
		Where("role_id", "=", roleId).
		Where("user_id", "=", t.Id).
		First()
	return checkRole != nil
}

// CheckPermissionById check the permission of the user.
func (t UserModel) CheckPermissionById(permissionId string) bool {
	checkPermission, _ := t.Table("user_permissions").
		Where("permission_id", "=", permissionId).
		Where("user_id", "=", t.Id).
		First()
	return checkPermission != nil
}

// AddRole add a role of the user model.
func (t UserModel) AddRole(roleId string) (int64, error) {
	if roleId != "" {
		if !t.CheckRoleId(roleId) {
			return t.WithTx(t.Tx).Table("role_users").
				Insert(dialect.H{
					"role_id": roleId,
					"user_id": t.Id,
				})
		}
	}
	return 0, nil
}

// AddPermission add a permission of the user model.
func (t UserModel) AddPermission(permissionId string) (int64, error) {
	if permissionId != "" {
		if !t.CheckPermissionById(permissionId) {
			return t.WithTx(t.Tx).Table("user_permissions").
				Insert(dialect.H{
					"permission_id": permissionId,
					"user_id":       t.Id,
				})
		}
	}
	return 0, nil
}

// 藉由參數檢查權限，如果有權限回傳第一個參數(path)，反之回傳""
func (t UserModel) GetCheckPermissionByUrlMethod(path, method string) string {
	// 檢查權限(藉由url、method)
	if !t.CheckPermissionByUrlMethod(path, method, url.Values{}) {
		return ""
	}
	return path
}

// Update update data
func (t UserModel) Update(username, password, name, avatar string) (int64, error) {
	fieldValues := dialect.H{
		"username":   username,
		"name":       name,
		"avatar":     avatar,
		"updated_at": time.Now().Format("2006-01-02 15:04:05"),
	}

	if password != "" {
		fieldValues["password"] = password
	}

	return t.WithTx(t.Tx).Table(t.TableName).
		Where("id", "=", t.Id).
		Update(fieldValues)
}

// DeleteRoles delete roles by id
func (t UserModel) DeleteRoles() error {
	return t.Table("role_users").
		Where("user_id", "=", t.Id).
		Delete()
}

// DeletePermissions delete all the permissions of the user model.
func (t UserModel) DeletePermissions() error {
	return t.WithTx(t.Tx).Table("user_permissions").
		Where("user_id", "=", t.Id).
		Delete()
}

// WithTx 將參數設置至UserModel.Tx
func (t UserModel) WithTx(tx *sql.Tx) UserModel {
	t.Tx = tx
	return t
}

// 取得參數
func getParam(u string) (string, url.Values) {
	m := make(url.Values)
	urr := strings.Split(u, "?")
	if len(urr) > 1 {
		m, _ = url.ParseQuery(urr[1])
	}
	return urr[0], m
}

func checkParam(src, comp url.Values) bool {
	if len(comp) == 0 {
		return true
	}
	if len(src) == 0 {
		return false
	}
	for key, value := range comp {
		v, find := src[key]
		if !find {
			return false
		}
		if len(value) == 0 {
			continue
		}
		if len(v) == 0 {
			return false
		}
		for i := 0; i < len(v); i++ {
			if v[i] == value[i] {
				continue
			} else {
				return false
			}
		}
	}
	return true
}

func inMethodArr(arr []string, str string) bool {
	for i := 0; i < len(arr); i++ {
		if strings.EqualFold(arr[i], str) {
			return true
		}
	}
	return false
}

// 檢查權限(藉由url、method)
func (t UserModel) CheckPermissionByUrlMethod(path, method string, formParams url.Values) bool {
	// 檢查是否為超級管理員
	if t.IsSuperAdmin() {
		return true
	}
	// 登出檢查
	logoutCheck, _ := regexp.Compile(config.Url("/logout") + "(.*?)")
	if logoutCheck.MatchString(path) {
		return true
	}

	if path == "" {
		return false
	}
	if path != "/" && path[len(path)-1] == '/' {
		path = path[:len(path)-1]
	}

	path = strings.Replace(path, "_edit_pk", "id", -1)
	path = strings.Replace(path, "_detail_pk", "id", -1)

	// 取得路徑及參數
	path, params := getParam(path)
	for key, value := range formParams {
		if len(value) > 0 {
			params.Add(key, value[0])
		}
	}

	for _, v := range t.Permissions {
		if v.HttpMethod[0] == "" || inMethodArr(v.HttpMethod, method) {

			if v.HttpPath[0] == "*" {
				return true
			}

			for i := 0; i < len(v.HttpPath); i++ {

				matchPath := config.Url(strings.TrimSpace(v.HttpPath[i]))
				matchPath, matchParam := getParam(matchPath)
				if matchPath == path {
					if checkParam(params, matchParam) {
						return true
					}
				}

				reg, err := regexp.Compile(matchPath)
				if err != nil {
					continue
				}

				if reg.FindString(path) == path {
					if checkParam(params, matchParam) {
						return true
					}
				}
			}
		}
	}
	return false
}

// 取得可用menu
func (t UserModel) WithMenus() UserModel {

	var menuIdsModel []map[string]interface{}

	// 判斷是否為超級管理員
	if t.IsSuperAdmin() {
		menuIdsModel, _ = t.Table("role_menu").
			LeftJoin("menu", "menu.id", "=", "role_menu.menu_id").
			Select("menu_id", "parent_id").
			All()
	} else {
		// 取得menuIdsModel藉由role_id
		rolesId := t.GetAllRoleId()
		if len(rolesId) > 0 {
			menuIdsModel, _ = t.Table("role_menu").
				LeftJoin("menu", "menu.id", "=", "role_menu.menu_id").
				WhereIn("role_menu.role_id", rolesId).
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

// 檢查用戶是否有可訪問的menu
func (t UserModel) HasMenu() bool {
	return len(t.MenuIds) != 0 || t.IsSuperAdmin()
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

// 設置UserModel.Conn = nil後回傳UserModel
func (t UserModel) ReleaseConn() UserModel {
	t.Conn = nil
	return t
}
