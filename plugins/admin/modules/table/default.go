package table

import (
	"goadminapi/modules/db/dialect"
	"goadminapi/modules/logger"
	"goadminapi/plugins/admin/modules"
	"goadminapi/plugins/admin/modules/form"
	"goadminapi/plugins/admin/modules/parameter"
	"strconv"
	"strings"

	"goadminapi/modules/db"
	"goadminapi/template/types"
)

type Columns []string

type GetDataFun func(params parameter.Parameters) ([]map[string]interface{}, int)

type DefaultTable struct {
	*BaseTable
	connectionDriver string
	connection       string
	sourceURL        string
	getDataFun       GetDataFun
}

func (tb *DefaultTable) getDataFromDB() bool {
	return tb.sourceURL == "" && tb.getDataFun == nil && tb.Info.GetDataFn == nil && tb.Detail.GetDataFn == nil
}

// 將參數設置(connName、conn)並回傳sql(struct)
func (tb *DefaultTable) sql() *db.SQL {
	if tb.connectionDriver != "" && tb.getDataFromDB() {
		return db.WithDriverAndConnection(tb.connection, db.GetConnectionFromService(services.Get(tb.connectionDriver)))
	}
	return nil
}

//-----------------------------table(interface)的方法--------------------------------
// InsertData insert data.
func (tb *DefaultTable) InsertData(dataList form.Values) error {
	var (
		id     = int64(0)
		err    error
		errMsg = ""
	)

	// 將__post_type:1 加入至map[string][]string
	dataList.Add("__post_type", "1")

	// -------------只有新增權限會執行----------------
	if tb.Form.PostHook != nil {
		defer func() {
			dataList.Add("__post_type", "1")
			dataList.Add(tb.GetPrimaryKey().Name, strconv.Itoa(int(id)))
			dataList.Add("__post_result", errMsg)
			go func() {
				defer func() {
					if err := recover(); err != nil {
						logger.Error(err)
					}
				}()

				err := tb.Form.PostHook(dataList)
				if err != nil {
					logger.Error(err)
				}
			}()
		}()
	}
	// -------------只有新增權限會執行----------------
	if tb.Form.Validator != nil {
		if err := tb.Form.Validator(dataList); err != nil {
			errMsg = "post error: " + err.Error()
			return err
		}
	}
	if tb.Form.PreProcessFn != nil {
		dataList = tb.Form.PreProcessFn(dataList)
	}

	// 用戶及角色頁面會執行新增資料的動作，直接return結果
	// --------------新增權限頁面不會執行------------------
	if tb.Form.InsertFn != nil {
		dataList.Delete("_post_type")
		err = tb.Form.InsertFn(dataList)
		if err != nil {
			errMsg = "post error: " + err.Error()
		}
		return err
	}

	id, err = tb.sql().Table(tb.Form.Table).Insert(tb.getInjectValueFromFormValue(dataList, types.PostTypeCreate))
	if db.CheckError(err, db.INSERT) {
		errMsg = "post error: " + err.Error()
		return err
	}

	return nil
}

// getColumns  取得所有欄位
func (tb *DefaultTable) getColumns(table string) (Columns, bool) {
	// 取得所有欄位資訊
	columnsModel, _ := tb.sql().Table(table).ShowColumns()
	columns := make(Columns, len(columnsModel))

	// 判斷資料庫引擎類型，將值加入columns([]string)
	switch tb.connectionDriver {
	case "mysql":
		auto := false
		for key, model := range columnsModel {
			columns[key] = model["Field"].(string)
			if columns[key] == tb.PrimaryKey.Name {
				if v, ok := model["Extra"].(string); ok {
					if v == "auto_increment" {
						auto = true
					}
				}
			}
		}
		return columns, auto
		// case db.DriverPostgresql:
		// 	auto := false
		// 	for key, model := range columnsModel {
		// 		columns[key] = model["column_name"].(string)
		// 		if columns[key] == tb.PrimaryKey.Name {
		// 			if v, ok := model["column_default"].(string); ok {
		// 				if strings.Contains(v, "nextval") {
		// 					auto = true
		// 				}
		// 			}
		// 		}
		// 	}
		// 	return columns, auto
		// case db.DriverSqlite:
		// 	for key, model := range columnsModel {
		// 		columns[key] = string(model["name"].(string))
		// 	}

		// 	num, _ := tb.sql().Table("sqlite_sequence").
		// 		Where("name", "=", tb.GetForm().Table).Count()

		// 	return columns, num > 0
		// case db.DriverMssql:
		// 	for key, model := range columnsModel {
		// 		columns[key] = string(model["column_name"].(string))
		// 	}
		// 	return columns, true
		// 	}
	default:
		panic("wrong driver")
	}
}

// --------------新增權限頁面會執行----------------
func (tb *DefaultTable) getInjectValueFromFormValue(dataList form.Values, typ types.PostType) dialect.H {
	var (
		value        = make(dialect.H)
		exceptString = make([]string, 0)
		// columns為資料表所有欄位
		// auto判斷是否有自動遞增(主鍵)的欄位
		columns, auto = tb.getColumns(tb.Form.Table)
		fun           types.PostFieldFilterFn
	)
	if auto {
		exceptString = []string{tb.PrimaryKey.Name, "__previous_", "__method_", "__token_",
			"__iframe", "__iframe_id"}
	} else {
		exceptString = []string{"__previous_", "__method_", "__token_",
			"__iframe", "__iframe_id"}
	}

	// ---------權限頁面執行新建動作會執行----------
	if !dataList.IsSingleUpdatePost() {
		// field為頁面顯示的所有欄位資訊
		for _, field := range tb.Form.FieldList {
			// 該欄位是否有多個選擇(ex: 權限的http_method欄位)
			if field.FormType.IsMultiSelect() {
				if _, ok := dataList[field.Field+"[]"]; !ok {
					dataList[field.Field+"[]"] = []string{""}
				}
			}
		}
	}

	// 刪除__post_type與__is_single_update的鍵與值
	dataList = dataList.RemoveRemark()
	// datalist為multipart/form-data設定的數值
	for k, v := range dataList {
		// 將名稱裡有[]取代成""(ex:http_method[]變成http_method)
		k = strings.Replace(k, "[]", "", -1)

		if !modules.InArray(exceptString, k) {
			if modules.InArray(columns, k) {
				// RemoveBlankFromArray 將參數中不為空的參數加入[]string
				vv := modules.RemoveBlankFromArray(v)

				delimiter := ","
				// FindByFieldName判斷FormFields[i].Field是否存在參數k，存在則回傳FormFields[i](FormField)
				// 取得欄位資訊(只取得新增資料頁面的欄位資訊)
				field := tb.Form.FieldList.FindByFieldName(k)
				if field != nil {
					fun = field.PostFilterFn
					// SetDefault如果第一個參數(source)為空則回傳第二個參數(def)，否則回傳第一個參數
					delimiter = modules.SetDefault(field.DefaultOptionDelimiter, ",") // ex: ,
				}

				if fun != nil {
					// -------新增權限的http_method、http_path欄位執行---------
					value[k] = fun(types.PostFieldModel{
						ID:    dataList.Get(tb.PrimaryKey.Name),
						Value: vv,
						// ToMap 將Values(struct)的值都加入map[string]string
						Row:      dataList.ToMap(),
						PostType: typ,
					})
				} else {
					// --------------新增權限頁面的name、slug欄位執行-----------------
					if len(vv) > 1 {
						value[k] = strings.Join(vv, delimiter)
					} else if len(vv) > 0 {
						value[k] = vv[0]
					} else {
						value[k] = ""
					}
				}
			} else {
				// 取得欄位資訊(只取得新增資料頁面的欄位資訊)
				field := tb.Form.FieldList.FindByFieldName(k)
				if field != nil && field.PostFilterFn != nil {
					field.PostFilterFn(types.PostFieldModel{
						ID:       dataList.Get(tb.PrimaryKey.Name),
						Value:    modules.RemoveBlankFromArray(v),
						Row:      dataList.ToMap(),
						PostType: typ,
					})
				}
			}
		}
	}
	return value
}
