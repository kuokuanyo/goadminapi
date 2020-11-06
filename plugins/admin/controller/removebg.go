package controller

import (
	b64 "encoding/base64"
	"encoding/json"
	"fmt"
	"image/png"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"goadminapi/context"
	"goadminapi/modules/db"
	"goadminapi/plugins/admin/models"
	"goadminapi/plugins/admin/modules"
	"goadminapi/plugins/admin/modules/response"
	"os"

	"github.com/nfnt/resize"
)

// RemoveBg 去背處理
func (h *Handler) RemoveBg(ctx *context.Context) {
	image := ctx.FormValue("image_file_b64")
	userid := ctx.FormValue("userid")
	filename := modules.Uuid()

	// 請求去背api
	body := requestremoveAPI(ctx, image)
	if body == "" {
		response.BadRequest(ctx, "圖片無法辨識，請重新上傳圖片!")
		return
	}

	// 處理原圖並寫入檔案
	original(ctx, filename, body)

	// 處理縮圖並寫入檔案
	thumbnail(ctx, filename, body)

	_, err := models.RemoveBg().SetConn(h.conn).New(userid, filename)
	if db.CheckError(err, db.INSERT) {
		response.BadRequest(ctx, err.Error())
		return
	}

	response.OkWithData(ctx, map[string]interface{}{
		"filename":  filename,
		"original":  "https://hilive.com.tw/original/" + filename + ".png",
		"thumbnail": "https://hilive.com.tw/thumbnail/" + filename + ".png",
	})
	return
}

// 請求去背API
func requestremoveAPI(ctx *context.Context, image string) (content string) {
	client := &http.Client{}
	data := url.Values{}
	data.Set("image_file_b64", image)
	// data.Set("size", "full")
	if data["image_file_b64"][0] == "" {
		response.BadRequest(ctx, "Please provide the source image in the image_file_b64 parameter.")
		return
	}

	// request
	req, err := http.NewRequest("POST", "https://api.remove.bg/v1.0/removebg", strings.NewReader(data.Encode()))
	if err != nil {
		response.BadRequest(ctx, err.Error())
		return
	}
	// add header
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("X-API-Key", "7B4j9zSjGMJ1jXgzgUTjQrh5")
	resp, err := client.Do(req)
	if err != nil {
		response.BadRequest(ctx, err.Error())
		return
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		response.BadRequest(ctx, err.Error())
		return
	}
	defer resp.Body.Close()

	var f interface{}
	json.Unmarshal(body, &f)
	if strings.Contains(fmt.Sprintf("%v", f), "errors") {
		response.BadRequest(ctx, "圖片無法辨識，請重新上傳圖片!")
		return ""
	}

	// 上傳圖片參數必須為base64
	content = b64.StdEncoding.EncodeToString(body)
	return content
}

// 請求上傳圖片API
// func requestUpoladAPI(ctx *context.Context, content string) (result map[string]interface{}) {
// 	var f interface{}
// 	client := &http.Client{}
// 	var endpotin = oauth2.Endpoint{
// 		AuthURL:  "https://api.imgur.com/oauth2/authorize",
// 		TokenURL: "https://api.imgur.com/oauth2/token",
// 	}

// 	var OauthConfig = &oauth2.Config{
// 		ClientID:     "736681296be0679",
// 		ClientSecret: "8bcd89d811bc35f8e888cd39b4a3b51d4c1b2ca3",
// 		RedirectURL:  "https://www.getpostman.com/oauth2/callback",
// 		Scopes:       []string{},
// 		Endpoint:     endpotin,
// 	}

// 	// 請求需要登入帳號密碼
// 	token, err := OauthConfig.PasswordCredentialsToken(oauth2.NoContext, "a167829435@gmail.com", "asdf4440")
// 	if err != nil {
// 		response.BadRequest(ctx, err.Error())
// 		return
// 	}

// 	data := url.Values{}
// 	// 參數必須為base64
// 	data.Set("image", content)
// 	data.Set("type", "base64")

// 	req, err := http.NewRequest("POST", "https://api.imgur.com/3/upload", strings.NewReader(data.Encode()))
// 	if err != nil {
// 		response.BadRequest(ctx, err.Error())
// 		return
// 	}

// 	// add header
// 	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
// 	req.Header.Add("Client-ID", "736681296be0679")
// 	req.Header.Set("Authorization", "Bearer "+token.AccessToken)

// 	resp, err := client.Do(req)
// 	if err != nil {
// 		response.BadRequest(ctx, err.Error())
// 		return
// 	}
// 	defer resp.Body.Close()

// 	body, err := ioutil.ReadAll(resp.Body)
// 	if err != nil {
// 		response.BadRequest(ctx, err.Error())
// 		return
// 	}

// 	json.Unmarshal(body, &f)
// 	result = f.(map[string]interface{})["data"].(map[string]interface{})
// 	return result
// }

// 處理原圖並寫入檔案
func original(ctx *context.Context, file, content string) {
	img, err := png.Decode(b64.NewDecoder(b64.StdEncoding, strings.NewReader(content)))
	if err != nil {
		response.BadRequest(ctx, "圖片辨識出現錯誤，無法處理")
		return
	}
	//寫入檔案
	out, err := os.Create("./themes/adminlte/resource/assets/dist/original/" + file + ".png")
	if err != nil {
		response.BadRequest(ctx, err.Error())
		return
	}
	defer out.Close()
	// write new image to file
	png.Encode(out, img)
}

// 處理縮圖並寫入檔案
func thumbnail(ctx *context.Context, file, content string) {
	img, err := png.Decode(b64.NewDecoder(b64.StdEncoding, strings.NewReader(content)))
	if err != nil {
		response.BadRequest(ctx, "圖片辨識出現錯誤，無法處理")
		return
	}
	// 縮圖
	m := resize.Thumbnail(200, 200, img, resize.Lanczos2)

	//寫入檔案
	out, err := os.Create("./themes/adminlte/resource/assets/dist/thumbnail/" + file + ".png")
	if err != nil {
		response.BadRequest(ctx, err.Error())
		return
	}
	defer out.Close()
	// write new image to file
	png.Encode(out, m)
}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}
