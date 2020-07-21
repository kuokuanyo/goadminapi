package models

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
func (t UserModel) SetConn(con Connection) UserModel {
	t.Conn = con
	return t
}