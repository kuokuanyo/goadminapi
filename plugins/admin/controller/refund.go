package controller

import (
	"bytes"
	"database/sql"
	"fmt"
	"goadminapi/context"
	"goadminapi/plugins/admin/modules/response"
	"goadminapi/template"
	template2 "html/template"
	"net/http"

	"github.com/line/line-bot-sdk-go/linebot"
)

func (h *Handler) ShowRefund(ctx *context.Context) {
	tmpl, name := template.GetComp("refund").GetTemplate()
	buf := new(bytes.Buffer)

	// 將第三個參數data寫入buf(struct)後輸出HTML
	if err := tmpl.ExecuteTemplate(buf, name, struct {
		UrlPrefix string
		Logo      template2.HTML
		CdnUrl    string
	}{
		UrlPrefix: h.config.AssertPrefix(),
		Logo:      h.config.LoginLogo,
		CdnUrl:    h.config.AssetUrl,
	}); err == nil {
		ctx.HTML(http.StatusOK, buf.String())
	} else {
		ctx.HTML(http.StatusOK, "parse template error (；′⌒`)")
		panic(err)
	}
}

func (h *Handler) Refund(ctx *context.Context) {
	var userid, transactionid, username, email, phone, remark string
	orderID := ctx.FormValue("orderID")
	reason := ctx.FormValue("reason")

	herokuDb := openHerokuDB(ctx)
	userid, transactionid, remark = getUserFromTransaction(ctx, herokuDb, orderID)
	if userid == "" {
		response.BadRequest(ctx, "抱歉，無此訂單紀錄")
		return
	}
	if reason == "" {
		response.BadRequest(ctx, "退費原因不能空白")
		return
	}
	if remark == "refund" {
		response.BadRequest(ctx, "此訂單已退費，請勿重複申請退費")
		return
	}

	username, email, phone = getInformationfromUserid(ctx, herokuDb, userid)
	pictureurl := "https://hilive.com.tw/original/" + orderID + ".png"

	// 開啟line機器人
	bot := openLineBot(ctx)
	// 傳送訊息
	if _, err := bot.PushMessage(
		"Ue25fbf858c52d4c4b84c19fd809e077a",
		linebot.NewTextMessage(fmt.Sprintf("申請退費\n用戶名: %s\nTransactionID: %s\nOrderID: %s\n電子郵件: %s\n手機號碼: %s\n退費原因: %s",
			username, transactionid, orderID, email, phone, reason)),
		linebot.NewImageMessage(pictureurl, pictureurl).
			WithQuickReplies(linebot.NewQuickReplyItems(
				linebot.NewQuickReplyButton(
					"",
					linebot.NewMessageAction("允許退費", "允許退費，OrderID為"+orderID)),
			))).Do(); err != nil {
		response.BadRequest(ctx, err.Error())
		return
	}
	if _, err := bot.PushMessage(
		userid,
		linebot.NewTextMessage("已收到您的申請退費訊息，若完成退費後會傳送訊息通知，謝謝!")).Do(); err != nil {
		response.BadRequest(ctx, err.Error())
		return
	}
	response.OkWithData(ctx, map[string]interface{}{})
	return
}

func getUserFromTransaction(ctx *context.Context, db *sql.DB, orderID string) (string, string, string) {
	var userid, transactionid, remark string
	rows, err := db.Query("SELECT userid, transactionid, remark FROM transaction WHERE orderid=?", orderID)
	if err != nil {
		response.BadRequest(ctx, err.Error())
		return "", "", ""
	}
	for rows.Next() {
		rows.Scan(&userid, &transactionid, &remark)
	}
	return userid, transactionid, remark
}

func getInformationfromUserid(ctx *context.Context, db *sql.DB, userid string) (string, string, string) {
	var username, email, phone string
	rows, err := db.Query("SELECT username, email, phone FROM information WHERE userid=?", userid)
	if err != nil {
		response.BadRequest(ctx, err.Error())
		return "", "", ""
	}
	for rows.Next() {
		rows.Scan(&username, &email, &phone)
	}
	return username, email, phone
}
