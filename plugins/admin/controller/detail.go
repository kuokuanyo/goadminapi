package controller

import (
	"fmt"
	"goadminapi/context"
	"goadminapi/modules/auth"
	"goadminapi/plugins/admin/modules"
	"goadminapi/plugins/admin/modules/parameter"
	"goadminapi/template"
	"goadminapi/template/types"
	"goadminapi/template/types/form"
)

func (h *Handler) ShowDetail(ctx *context.Context) {
	var (
		prefix = ctx.Query("__prefix")
		id     = ctx.Query("__detail_pk")

		// 先透過參數prefix取得Table(interface)，接著判斷條件後將[]context.Node加入至Handler.operations後回傳
		panel = h.table(prefix, ctx)

		// 透過參數ctx回傳目前登入的用戶(Context.UserValue["user"])並轉換成UserModel
		user     = auth.Auth(ctx)
		newPanel = panel.Copy()

		// GetDetail 取得table中設置的InfoPanel(struct)，在generators.go設置的資訊
		detail = panel.GetDetail()

		// GetInfo 取得table中設置的InfoPanel(struct)，在generators.go設置的資訊
		info = panel.GetInfo() // info為所有顯示的欄位資訊

		// GetForm 取得table中設置的FormPanel(struct)，在generators.go設置的資訊
		formModel = newPanel.GetForm()

		fieldList = make(types.FieldList, 0)
	)

	// fieldList為所有欄位資訊
	// -------只有用戶建立detail資訊-----------
	if len(detail.FieldList) == 0 {
		// ------------角色、權限執行----------------
		fieldList = info.FieldList
	} else {
		// ------------用戶執行----------------
		fieldList = detail.FieldList
	}

	formModel.FieldList = make([]types.FormField, len(fieldList))

	// 將欄位資訊添加給表單欄位資訊
	for i, field := range fieldList {
		formModel.FieldList[i] = types.FormField{
			Field:        field.Field,
			FieldClass:   field.Field,
			TypeName:     field.TypeName, // 欄位類型
			Head:         field.Head,     // 欄位名稱(中文)
			Hide:         field.Hide,     // 是否隱藏
			Joins:        field.Joins,    // join
			FormType:     form.Default,   // ex: default
			FieldDisplay: field.FieldDisplay,
		}
	}

	// ------------一般detail.Table都為空---------------------
	// 將資料表名稱設置至formModel.Table
	if detail.Table != "" {
		formModel.Table = detail.Table
	} else {
		// -------一般都執行此條件---------
		formModel.Table = info.Table
	}

	// 將頁面size、資料排列方式、選擇欄位...等資訊後設置至Parameters(struct)
	param := parameter.GetParam(ctx.Request.URL,
		info.DefaultPageSize,
		info.SortField,
		info.GetSort())

	// DeleteDetailPk 刪除Parameters.Fields[__detail_pk]
	paramStr := param.DeleteDetailPk().GetRouteParamStr()

	// 編輯介面
	editUrl := modules.AorEmpty(!info.IsHideEditButton, h.routePathWithPrefix("show_edit", prefix)+paramStr+
		"&"+"__edit_pk"+"="+ctx.Query("__detail_pk"))
	// 刪除
	deleteUrl := modules.AorEmpty(!info.IsHideDeleteButton, h.routePathWithPrefix("delete", prefix)+paramStr)
	// 資料顯示介面
	infoUrl := h.routePathWithPrefix("info", prefix) + paramStr

	// 檢查是否有編輯及刪除的權限
	editUrl = user.GetCheckPermissionByUrlMethod(editUrl, h.route("show_edit").Method())
	deleteUrl = user.GetCheckPermissionByUrlMethod(deleteUrl, h.route("delete").Method())

	deleteJs := ""
	// 刪除資料的js語法
	if deleteUrl != "" {
		deleteJs = fmt.Sprintf(`<script>
function DeletePost(id) {
	swal({
			title: '%s',
			type: "warning",
			showCancelButton: true,
			confirmButtonColor: "#DD6B55",
			confirmButtonText: '%s',
			closeOnConfirm: false,
			cancelButtonText: '%s',
		},
		function () {
			$.ajax({
				method: 'post',
				url: '%s',
				data: {
					id: id
				},
				success: function (data) {
					if (typeof (data) === "string") {
						data = JSON.parse(data);
					}
					if (data.code === 200) {
						location.href = '%s'
					} else {
						swal(data.msg, '', 'error');
					}
				}
			});
		});
}

$('.delete-btn').on('click', function (event) {
	DeletePost(%s)
});

</script>`, "確定要刪除嗎?", "確定", "取消", deleteUrl, infoUrl, id)
	}

	title := ""
	desc := ""

	isNotIframe := ctx.Query("__iframe") != "true" // ex: true
	if isNotIframe {
		title = detail.Title
		if title == "" {
			title = info.Title + "細節"
		}

		desc = detail.Description
		if desc == "" {
			desc = info.Description + "細節"
		}
	}

	// GetDataWithId 透過id取得資料，並且將選項、預設值...等資訊設置至FormFields(帶有預設值)
	formInfo, err := newPanel.GetDataWithId(param.WithPKs(id))
	if err != nil {
		h.HTML(ctx, user, types.Panel{
			Content:     aAlert().Warning(err.Error()),
			Description: template.HTML(desc),
			Title:       template.HTML(title),
		}, param.Animation)
		return
	}

	h.HTML(ctx, user, types.Panel{
		Content: detailContent(aForm().
			SetTitle(template.HTML(title)). 
			SetContent(formInfo.FieldList). // 將欄位資訊及數值設置至表單的content
			SetFooter(template.HTML(deleteJs)).
			SetHiddenFields(map[string]string{
				"__previous_": infoUrl,
			}).
			SetPrefix(h.config.PrefixFixSlash()), editUrl, deleteUrl, !isNotIframe),
		Description: template.HTML(desc),
		Title:       template.HTML(title),
	}, param.Animation)
}
