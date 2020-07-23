package auth

import (
	"goadminapi/modules/config"
	"goadminapi/modules/db"
	"goadminapi/plugins/admin/models"
	"net/url"
)

// 尋找資料表中符合參數(sesKey)的user資料，回傳user_id值，如果沒有則回傳-1
func GetUserID(sesKey string, conn db.Connection) int64 {
	// GetSessionByKey尋找資料表中符合參數(sesKey)的user資料，回傳user_id
	id, err := GetSessionByKey(sesKey, "user_id", conn)
	if err != nil {
		return -1
	}
	if idFloat64, ok := id.(float64); ok {
		return int64(idFloat64)
	}
	return -1
}

// 透過參數(id)取得role、permission以及可使用menu並回傳UserModel(struct)
func GetCurUserByID(id int64, conn db.Connection) (user models.UserModel, ok bool) {
	// 透過參數(id)取得UserModel(struct)，將值設置至UserModel
	user = models.User("users").SetConn(conn).Find(id)
	if user.IsEmpty() {
		ok = false
		return
	}
	// 判斷是否有頭像
	// GetStore回傳globalCfg.Store
	if user.Avatar == "" || config.GetStore().Prefix == "" {
		user.Avatar = ""
	} else {
		user.Avatar = config.GetStore().URL(user.Avatar)
	}
	// 取得角色、權限及可使用菜單
	user = user.WithRoles().WithPermissions().WithMenus()
	// 檢查用戶是否有可訪問的menu
	ok = user.HasMenu()
	return
}

// 透過參數sesKey(cookie)取得id並利用id取得該user的role、permission以及可用menu，最後回傳UserModel(struct)
func GetCurUser(sesKey string, conn db.Connection) (user models.UserModel, ok bool) {
	if sesKey == "" {
		ok = false
		return
	}
	// 取得user_id(在goadmin_session資料表values欄位)
	// 尋找資料表中符合參數(sesKey)的user資料，回傳user_id值，如果沒有則回傳-1
	id := GetUserID(sesKey, conn)
	if id == -1 {
		ok = false
		return
	}
	// GetCurUserByID取得參數id的role、permission以及可使用menu並回傳UserModel(struct)
	return GetCurUserByID(id, conn)
}

// 透過參數檢查用戶權限
func CheckPermissions(user models.UserModel, path, method string, param url.Values) bool {
	// CheckPermissionByUrlMethod在plugins\admin\models\user.go中
	return user.CheckPermissionByUrlMethod(path, method, param)
}
