package auth

import (
	"goadminapi/context"
	"goadminapi/modules/db"
	"goadminapi/modules/service"
	"goadminapi/plugins/admin/models"
	"sync"

	"golang.org/x/crypto/bcrypt"
)

// CSRFToken is type of a csrf token list.
type CSRFToken []string

type TokenService struct {
	tokens CSRFToken //[]string
	lock   sync.Mutex
}

// 為設置Service(interface)方法
func (s *TokenService) Name() string {
	return "token_csrf_helper"
}

// ***************初始化*******************
// 將token_csrf_helper加入services(map[string]Generator)
func init() {
	service.Register("token_csrf_helper", func() (service.Service, error) {
		return &TokenService{
			tokens: make(CSRFToken, 0),
		}, nil
	})
}

// 透過參數ctx回傳目前登入的用戶(Context.UserValue["user"])並轉換成UserModel
func Auth(ctx *context.Context) models.UserModel {
	return ctx.User().(models.UserModel)
}

// 檢查user密碼是否正確之後取得user的role、permission及可用menu，最後更新資料表(goadmin_users)的密碼值(加密)
func Check(password string, username string, conn db.Connection) (user models.UserModel, ok bool) {

	// User設置UserModel.Base.TableName(struct)並回傳設置UserModel(struct)
	// SetConn將參數conn(db.Connection)設置至UserModel.conn(UserModel.Base.Conn)
	user = models.User("users").SetConn(conn).FindByUserName(username)
	// 判斷user是否為空
	if user.IsEmpty() {

		ok = false
	} else {
		// 檢查密碼
		if comparePassword(password, user.Password) {
			ok = true
			//取得user的role、permission及可用menu
			user = user.WithRoles().WithPermissions().WithMenus()
			// EncodePassword將參數pwd加密
			// UpdatePwd更新資料表密碼(加密)
			user.UpdatePwd(EncodePassword([]byte(password)))
		} else {
			ok = false
		}
	}
	return
}

// 檢查密碼是否相符
func comparePassword(comPwd, pwdHash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(pwdHash), []byte(comPwd))
	return err == nil
}

// 將參數pwd加密
func EncodePassword(pwd []byte) string {
	hash, err := bcrypt.GenerateFromPassword(pwd, bcrypt.DefaultCost)
	if err != nil {
		return ""
	}
	return string(hash[:])
}

// 設置cookie(struct)並儲存在response header Set-Cookie中
func SetCookie(ctx *context.Context, user models.UserModel, conn db.Connection) error {
	// 設置Session(struct)資訊並取得cookie及設置cookie值
	ses, err := InitSession(ctx, conn)

	if err != nil {
		return err
	}

	// Add將參數"user_id"、user.Id加入Session.Values後檢查是否有符合Session.Sid的資料，判斷插入或是更新資料
	// 最後設置cookie(struct)並儲存在response header Set-Cookie中
	return ses.Add("user_id", user.Id)
}

// 將參數s轉換成TokenService(struct)類別後回傳
func GetTokenService(s interface{}) *TokenService {
	if srv, ok := s.(*TokenService); ok {
		return srv
	}
	panic("wrong service")
}

// 如果符合參數，將在TokenService.tokens([]string)裡將符合的toCheckToken從[]string拿出
func (s *TokenService) CheckToken(toCheckToken string) bool {
	for i := 0; i < len(s.tokens); i++ {
		if (s.tokens)[i] == toCheckToken {
			s.tokens = append((s.tokens)[:i], (s.tokens)[i+1:]...)
			return true
		}
	}
	return false
}
