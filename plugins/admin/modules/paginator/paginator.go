package paginator

import (
	"goadminapi/plugins/admin/modules/parameter"
	template2 "goadminapi/template"
	"goadminapi/template/components"
	"goadminapi/template/types"
	"html/template"
	"math"
	"strconv"
)

// Config 分頁設置
type Config struct {
	Size         int
	Param        parameter.Parameters
	PageSizeList []string
}

// Get 設置分頁資訊
func Get(cfg Config) types.PaginatorAttribute {
	// PaginatorAttribute(struct) 分頁的struct
	paginator := template2.Default().Paginator().(*components.PaginatorAttribute)

	// 判斷總共頁數
	// cfg.Size為總共資料，cfg.Param.PageSizeInt為設置一頁的資料數
	totalPage := int(math.Ceil(float64(cfg.Size) / float64(cfg.Param.PageSizeInt)))

	// cfg.Param.PageInt為預設第n頁
	// 如果設為第一頁則previousClass=disabled
	if cfg.Param.PageInt == 1 {
		paginator.PreviousClass = "disabled"
		paginator.PreviousUrl = cfg.Param.URLPath
	} else {
		paginator.PreviousClass = ""
		// GetLastPageRouteParamStr 取得前一頁的url參數，例如在第二頁時回傳第一頁的url參數
		paginator.PreviousUrl = cfg.Param.URLPath + cfg.Param.GetLastPageRouteParamStr()
	}

	// 如果PageInt=totalPage(總頁數)則沒有下一頁，NextClass=disabled
	if cfg.Param.PageInt == totalPage {
		paginator.NextClass = "disabled"
		paginator.NextUrl = cfg.Param.URLPath
	} else {
		paginator.NextClass = ""
		paginator.NextUrl = cfg.Param.URLPath + cfg.Param.GetNextPageRouteParamStr()
	}

	// 取得第一頁的url
	paginator.Url = cfg.Param.URLPath + cfg.Param.GetRouteParamStrWithoutPageSize("1") + "&" + "__no_animation_" + "=true"
	paginator.CurPageEndIndex = strconv.Itoa((cfg.Param.PageInt) * cfg.Param.PageSizeInt)       // 該頁起始資料的index
	paginator.CurPageStartIndex = strconv.Itoa((cfg.Param.PageInt - 1) * cfg.Param.PageSizeInt) // 該頁結束資料的index
	paginator.Total = strconv.Itoa(cfg.Size)

	// 單頁資料顯示筆數
	if len(cfg.PageSizeList) == 0 {
		cfg.PageSizeList = []string{"10", "20", "50", "100"}
	}

	// 在每一個PageSizeList([]string)設值空HTML
	paginator.Option = make(map[string]template.HTML, len(cfg.PageSizeList))
	for i := 0; i < len(cfg.PageSizeList); i++ {
		paginator.Option[cfg.PageSizeList[i]] = template.HTML("")
	}
	// 在選擇單頁筆數的資料增加HTML
	paginator.Option[cfg.Param.PageSize] = template.HTML("selected")

	paginator.Pages = []map[string]string{}

	if totalPage < 10 {
		var pagesArr []map[string]string
		for i := 1; i < totalPage+1; i++ {
			if i == cfg.Param.PageInt {
				pagesArr = append(pagesArr, map[string]string{
					"page":    cfg.Param.Page,
					"active":  "active",
					"isSplit": "0",
					"url":     cfg.Param.URLNoAnimation(cfg.Param.Page),
				})
			} else {
				page := strconv.Itoa(i)
				pagesArr = append(pagesArr, map[string]string{
					"page":    page,
					"active":  "",
					"isSplit": "0",
					"url":     cfg.Param.URLNoAnimation(page),
				})
			}
		}
		paginator.Pages = pagesArr
	} else {
		var pagesArr []map[string]string
		if cfg.Param.PageInt < 6 {
			for i := 1; i < totalPage+1; i++ {

				if i == cfg.Param.PageInt {
					pagesArr = append(pagesArr, map[string]string{
						"page":    cfg.Param.Page,
						"active":  "active",
						"isSplit": "0",
						"url":     cfg.Param.URLNoAnimation(cfg.Param.Page),
					})
				} else {
					page := strconv.Itoa(i)
					pagesArr = append(pagesArr, map[string]string{
						"page":    page,
						"active":  "",
						"isSplit": "0",
						"url":     cfg.Param.URLNoAnimation(page),
					})
				}

				if i == 6 {
					pagesArr = append(pagesArr, map[string]string{
						"page":    "",
						"active":  "",
						"isSplit": "1",
						"url":     cfg.Param.URLNoAnimation("6"),
					})
					i = totalPage - 1
				}
			}
		} else if cfg.Param.PageInt < totalPage-4 {
			for i := 1; i < totalPage+1; i++ {

				if i == cfg.Param.PageInt {
					pagesArr = append(pagesArr, map[string]string{
						"page":    cfg.Param.Page,
						"active":  "active",
						"isSplit": "0",
						"url":     cfg.Param.URLNoAnimation(cfg.Param.Page),
					})
				} else {
					page := strconv.Itoa(i)
					pagesArr = append(pagesArr, map[string]string{
						"page":    page,
						"active":  "",
						"isSplit": "0",
						"url":     cfg.Param.URLNoAnimation(page),
					})
				}

				if i == 2 {
					pagesArr = append(pagesArr, map[string]string{
						"page":    "",
						"active":  "",
						"isSplit": "1",
						"url":     cfg.Param.URLNoAnimation("2"),
					})
					if cfg.Param.PageInt < 7 {
						i = 5
					} else {
						i = cfg.Param.PageInt - 2
					}
				}

				if cfg.Param.PageInt < 7 {
					if i == cfg.Param.PageInt+5 {
						pagesArr = append(pagesArr, map[string]string{
							"page":    "",
							"active":  "",
							"isSplit": "1",
							"url":     cfg.Param.URLNoAnimation(strconv.Itoa(i)),
						})
						i = totalPage - 1
					}
				} else {
					if i == cfg.Param.PageInt+3 {
						pagesArr = append(pagesArr, map[string]string{
							"page":    "",
							"active":  "",
							"isSplit": "1",
							"url":     cfg.Param.URLNoAnimation(strconv.Itoa(i)),
						})
						i = totalPage - 1
					}
				}
			}
		} else {
			for i := 1; i < totalPage+1; i++ {

				if i == cfg.Param.PageInt {
					pagesArr = append(pagesArr, map[string]string{
						"page":    cfg.Param.Page,
						"active":  "active",
						"isSplit": "0",
						"url":     cfg.Param.URLNoAnimation(cfg.Param.Page),
					})
				} else {
					page := strconv.Itoa(i)
					pagesArr = append(pagesArr, map[string]string{
						"page":    page,
						"active":  "",
						"isSplit": "0",
						"url":     cfg.Param.URLNoAnimation(page),
					})
				}

				if i == 2 {
					pagesArr = append(pagesArr, map[string]string{
						"page":    "",
						"active":  "",
						"isSplit": "1",
						"url":     cfg.Param.URLNoAnimation("2"),
					})
					i = totalPage - 4
				}
			}
		}
		paginator.Pages = pagesArr
	}
	// SetPageSizeList將顯示資料比數選項設置至(struct)
	return paginator.SetPageSizeList(cfg.PageSizeList)
}
