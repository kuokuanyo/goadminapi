package auth

import (
	"goadminapi/context"
	"goadminapi/modules/config"
	"goadminapi/modules/db"
	"goadminapi/modules/errors"
	"goadminapi/modules/logger"
	"goadminapi/modules/page"
	"goadminapi/plugins/admin/models"
	"goadminapi/template/types"
	"net/http"
	"net/url"

	template2 "goadminapi/template"
)

// MiddlewareCallback is type of callback function.
type MiddlewareCallback func(ctx *context.Context)

// Invoker contains the callback functions which are used
// in the route middleware.
type Invoker struct {
	prefix                 string
	authFailCallback       MiddlewareCallback //驗證失敗Callback
	permissionDenyCallback MiddlewareCallback //權限拒絕Callback
	conn                   db.Connection
}

// DefaultInvoker 設置並回傳Invoker(Struct)
func DefaultInvoker(conn db.Connection) *Invoker {
	return &Invoker{
		prefix: config.Prefix(),
		authFailCallback: func(ctx *context.Context) {
			if ctx.Request.URL.Path == config.Url(config.GetLoginUrl()) {
				return
			}
			if ctx.Request.URL.Path == config.Url("/logout") {
				ctx.Write(302, map[string]string{
					"Location": config.Url(config.GetLoginUrl()),
				}, ``)
				return
			}
			param := ""
			// url後加入參數
			if ref := ctx.Headers("Referer"); ref != "" {
				param = "?ref=" + url.QueryEscape(ref)
			}

			u := config.Url(config.GetLoginUrl() + param)
			_, err := ctx.Request.Cookie("session")
			referer := ctx.Headers("Referer")

			if (ctx.Headers("X-PJAX") == "" && ctx.Method() != "GET") ||
				err != nil ||
				referer == "" {
				ctx.Write(302, map[string]string{
					"Location": u,
				}, ``)
			} else {
				// 登入時間過長或同個IP登入
				msg := "login overdue, please login again"
				//添加HTML
				ctx.HTML(http.StatusOK, `<script>
	if (typeof(swal) === "function") {
		swal({
			type: "info",
			title: `+"login info"+`,
			text: "`+msg+`",
			showCancelButton: false,
			confirmButtonColor: "#3c8dbc",
			confirmButtonText: '`+"got it"+`',
        })
		setTimeout(function(){ location.href = "`+u+`"; }, 3000);
	} else {
		alert("`+msg+`")
		location.href = "`+u+`"
    }
</script>`)
			}
		},
		permissionDenyCallback: func(ctx *context.Context) {
			if ctx.Headers("X-PJAX") == "" && ctx.Method() != "GET" {
				ctx.JSON(http.StatusForbidden, map[string]interface{}{
					"code": http.StatusForbidden,
					"msg":  "permission denied",
				})
			} else {
				page.SetPageContent(ctx, Auth(ctx), func(ctx interface{}) (types.Panel, error) {
					return template2.WarningPanel(errors.PermissionDenied, template2.NoPermission403Page), nil
				}, conn)
			}
		},
		conn: conn,
	}
}

// 透過參數ctx取得UserModel，並且取得該user的role、權限與可用menu，最後檢查用戶權限
func Filter(ctx *context.Context, conn db.Connection) (models.UserModel, bool, bool) {
	var (
		id   float64
		ok   bool
		user = models.User("users")
	)

	// 設置Session(struct)資訊並取得cookie及設置cookie值
	ses, err := InitSession(ctx, conn)

	if err != nil {
		// 驗證失敗
		logger.Error("retrieve auth user failed", err)
		return user, false, false
	}

	// 藉由參數取得Session.Values[user_id]
	if id, ok = ses.Get("user_id").(float64); !ok {
		return user, false, false
	}

	// GetCurUserByID取得該id的角色、權限以即可訪問的菜單
	user, ok = GetCurUserByID(int64(id), conn)

	if !ok {
		return user, false, false
	}

	// CheckPermissions透過path、method、param檢查用戶權限
	return user, true, CheckPermissions(user, ctx.Request.URL.String(), ctx.Method(), ctx.PostForm())
}

// 透過參數ctx取得UserModel，並且取得該user的role、權限與可用menu，最後檢查用戶權限
func (invoker *Invoker) Middleware() context.Handler {
	return func(ctx *context.Context) {
		// 透過參數ctx取得UserModel，並且取得該user的role、權限與可用menu，最後檢查用戶權限
		user, authOk, permissionOk := Filter(ctx, invoker.conn)
		if authOk && permissionOk {
			ctx.SetUserValue("user", user)
			ctx.Next()
			return
		}

		if !authOk {
			invoker.authFailCallback(ctx)
			ctx.Abort()
			return
		}
		if !permissionOk {
			ctx.SetUserValue("user", user)
			invoker.permissionDenyCallback(ctx)
			ctx.Abort()
			return
		}
	}
}

// 建立Invoker(Struct)並透過參數ctx取得UserModel，並且取得該user的role、權限與可用menu，最後檢查用戶權限
func Middleware(conn db.Connection) context.Handler {
	return DefaultInvoker(conn).Middleware()
}

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
