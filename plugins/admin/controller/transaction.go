package controller

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"goadminapi/context"
	"goadminapi/plugins/admin/modules/response"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/line/line-bot-sdk-go/linebot"
)

// Transaction 確認是否付款
func (h *Handler) Transaction(ctx *context.Context) {
	var userid string
	liff := ctx.Query("liff.state")
	if liff == "" {
		response.BadRequest(ctx, "Please provide parameters.")
		return
	}

	transaction := strings.Split(liff, "transactionId=")[1]
	orderID := strings.Split(strings.Split(liff, "orderId=")[1], "&")[0]
	if transaction == "" {
		response.BadRequest(ctx, "Please provide transactionId.")
		return
	}
	if orderID == "" {
		response.BadRequest(ctx, "Please provide orderId.")
		return
	}

	//請求付款確認
	result := requestConfirmAPI(ctx, transaction)

	// 開啟line機器人
	bot := openLineBot(ctx)

	if fmt.Sprintf("%v", result["returnMessage"]) == "Success." {
		// 開啟heroku資料庫
		herokuDb := openHerokuDB(ctx)
		defer herokuDb.Close()

		rows, err := herokuDb.Query("SELECT userid FROM transaction WHERE orderid=?", orderID)
		if err != nil {
			response.BadRequest(ctx, err.Error())
			return
		}
		for rows.Next() {
			rows.Scan(&userid)
		}

		pictureurl := "https://hilive.com.tw/original/" + orderID + ".png"
		// 傳送訊息
		if _, err = bot.PushMessage(
			userid,
			linebot.NewTextMessage("已完成付款，圖片傳送中(如需退費請至 https://hilive.com.tw/admin/refund?orderID="+orderID+" 申請退費流程):"),
			linebot.NewImageMessage(pictureurl, pictureurl),
			linebot.NewTextMessage("謝謝您使用本公司的去背功能!\n如想繼續使用去背功能，麻煩請上傳需要去背的圖片!").WithQuickReplies(linebot.NewQuickReplyItems(
				linebot.NewQuickReplyButton(
					"",
					linebot.NewCameraAction("立刻拍攝")),
				linebot.NewQuickReplyButton(
					"",
					linebot.NewCameraRollAction("相簿選擇")),
			))).Do(); err != nil {
			response.BadRequest(ctx, err.Error())
			return
		}

		stmt, _ := herokuDb.Prepare("UPDATE transaction set remark = ?, transactionid = ? where orderid = ?")
		stmt.Exec("paid", transaction, orderID)
		stmt, _ = herokuDb.Prepare("UPDATE remark set remark = ? where userid = ?")
		stmt.Exec("remove background", userid)

		ctx.HTML(200, succ)
	} else {
		if _, err := bot.PushMessage(
			userid,
			linebot.NewTextMessage("付款失敗，請重新付款!")).
			Do(); err != nil {
			response.BadRequest(ctx, err.Error())
			return
		}
		ctx.HTML(200, fail)
	}
}

// 請求付款確認API
func requestConfirmAPI(ctx *context.Context, transaction string) (result map[string]interface{}) {
	var f interface{}
	client := &http.Client{}
	jsbody := map[string]string{"amount": "15", "currency": "TWD"}
	js, _ := json.Marshal(jsbody)

	// 請求付款確認API
	req, err := http.NewRequest("POST", "https://sandbox-api-pay.line.me/v2/payments/"+transaction+"/confirm", bytes.NewBuffer([]byte(js)))
	if err != nil {
		response.BadRequest(ctx, err.Error())
		return
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("X-LINE-ChannelId", "1654985667")
	req.Header.Add("X-LINE-ChannelSecret", "a991a81db5cee82b12e2378871132a78")

	resp, err := client.Do(req)
	if err != nil {
		response.BadRequest(ctx, err.Error())
		return
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		response.BadRequest(ctx, err.Error())
		return
	}

	json.Unmarshal(body, &f)
	result = f.(map[string]interface{})
	return
}

// 開啟heroku資料庫
func openHerokuDB(ctx *context.Context) (herokuDb *sql.DB) {
	herokuDb, err := sql.Open("mysql", "be94ad46dfd2c5:0986ac8c@tcp(us-cdbr-east-02.cleardb.com:3306)/heroku_340b0d6567ec671")
	if err != nil {
		response.BadRequest(ctx, err.Error())
		return
	}
	if err := herokuDb.Ping(); err != nil {
		response.BadRequest(ctx, err.Error())
		return
	}
	return
}

// 開啟line機器人
func openLineBot(ctx *context.Context) (bot *linebot.Client) {
	bot, err := linebot.New("d8ea5d696639efe57504d9b226157e32",
		"8R/t1FtEvfOXNiyUkNUZ/TSc94IN5tsDUXGwWmXYtp2mTicYiA7Bcg31hUmojMQUKLXDdj9RaUD9X0t/YdnJLalEVjigeKmy/Qs0dHjkAyCTjrNOx406f8l7E3cIuAhpN8181kpRPJvxTp2KYhYXIQdB04t89/1O/w1cDnyilFU=")
	if err != nil {
		response.BadRequest(ctx, err.Error())
		return
	}
	return
}

// 取得處理後的結果
func getFinalData(ctx *context.Context, db *sql.DB, table, col, user string) string {
	var data string
	rows, err := db.Query("SELECT "+col+" FROM "+table+" WHERE userid= ? ", user)
	if err != nil {
		response.BadRequest(ctx, err.Error())
		return ""
	}
	for rows.Next() {
		rows.Scan(&data)
	}
	return data
}

// 必須先登入平台取得cookies
func getCookies(ctx *context.Context) []*http.Cookie {
	client := &http.Client{}
	data := url.Values{}
	var cookies []*http.Cookie
	data.Set("phone", "0932530813")
	data.Set("password", "admin")

	req, err := http.NewRequest("POST", "https://hilive.com.tw/admin/signin", strings.NewReader(data.Encode()))
	if err != nil {
		response.BadRequest(ctx, err.Error())
		return cookies
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	resp, err := client.Do(req)
	if err != nil {
		response.BadRequest(ctx, err.Error())
		return cookies
	}
	defer resp.Body.Close()

	cookies = resp.Cookies()
	return cookies
}

// 執行平台上的轉檔API
func requestConvert(ctx *context.Context, file string, cookies []*http.Cookie) (fileurl string) {
	var f interface{}
	client := &http.Client{}
	data := url.Values{}
	data.Set("file", file)
	req, err := http.NewRequest("POST", "https://hilive.com.tw/admin/convert-file", strings.NewReader(data.Encode()))
	if err != nil {
		response.BadRequest(ctx, err.Error())
		return
	}

	for _, cookie := range cookies {
		req.AddCookie(&http.Cookie{Name: cookie.Name, Value: cookie.Value, HttpOnly: cookie.HttpOnly})
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

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

	json.Unmarshal(body, &f)
	fileurl = fmt.Sprintf("%v", f.(map[string]interface{})["data"].(map[string]interface{})["file"])
	return
}

var succ = `<!DOCTYPE html>
<html><head>
	<title>已付款</title>
	<meta charset="utf-8">
	<meta http-equiv="X-UA-Compatible" content="IE=EDGE">
	<link rel="stylesheet" type="text/css" href="//scdn.line-apps.com/linepay/web/202007021542/css/pc_standby.min.css">
	<script src="https://d.line-scdn.net/liff/1.0/sdk.js"></script>
	</head>
	<body>
	<div class="wrap">

		<div class="table sandBox">
			  <div class="table-cell simulation">
				  <span class="blind">LINE Pay</span><span class="ico"></span><span class="txt">Simulation</span><span class="border"></span>
			  </div>
		</div>
		
		<h2 class="logo"><span class="blind">LINE Pay</span><span class="ico"></span></h2>
		<div class="message">
			<div class="welcome">
				<span class="name">付款完成並已傳送圖片至對話中，請關閉頁面接收完整圖!</span>
			</div>
			<input type="button" onclick="liff.closeWindow();" value="關閉視窗" style="width:120px;height:40px;border:2px blue none;background-color:rgb(7, 184, 7);">
		</div>
		<footer class="footer">© LINE Pay</footer>
	</div>
	</script>
	</body></html>`

var fail = `<!DOCTYPE html>
	<html><head>
		<title>付款失敗</title>
		<meta charset="utf-8">
		<meta http-equiv="X-UA-Compatible" content="IE=EDGE">
		<link rel="stylesheet" type="text/css" href="//scdn.line-apps.com/linepay/web/202007021542/css/pc_standby.min.css">
		<script src="https://d.line-scdn.net/liff/1.0/sdk.js"></script>
		</head>
		<body>
		<div class="wrap">
	
			<div class="table sandBox">
				  <div class="table-cell simulation">
					  <span class="blind">LINE Pay</span><span class="ico"></span><span class="txt">Simulation</span><span class="border"></span>
				  </div>
			</div>
			
			<h2 class="logo"><span class="blind">LINE Pay</span><span class="ico"></span></h2>
			<div class="message">
				<div class="welcome">
					<span class="name">付款發生錯誤，請關閉視窗並重新執行付款動作!</span>
				</div>
				<input type="button" onclick="liff.closeWindow();" value="關閉視窗" style="width:120px;height:40px;border:2px blue none;background-color:rgb(7, 184, 7);">
			</div>
			<footer class="footer">© LINE Pay</footer>
		</div>
		</script>
		</body></html>`
