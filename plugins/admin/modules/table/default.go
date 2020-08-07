package table

import (
	"encoding/json"
	"fmt"
	"goadminapi/modules/db/dialect"
	"goadminapi/plugins/admin/modules"
	"goadminapi/plugins/admin/modules/form"
	"goadminapi/plugins/admin/modules/paginator"
	"goadminapi/plugins/admin/modules/parameter"
	"html/template"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"

	"goadminapi/modules/db"
	"goadminapi/template/types"
)

type GetDataFromURLRes struct {
	Data []map[string]interface{}
	Size int
}

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

// NewDefaultTable 將參數值設置至預設DefaultTable(struct)
func NewDefaultTable(cfgs ...Config) Table {
	var cfg Config
	if len(cfgs) > 0 && cfgs[0].PrimaryKey.Name != "" {
		cfg = cfgs[0]
	} else {
		cfg = DefaultConfig()
	}
	return &DefaultTable{
		BaseTable: &BaseTable{
			Info:           types.NewInfoPanel(cfg.PrimaryKey.Name), // 預設InfoPanel(struct)
			Form:           types.NewFormPanel(),                    // 預設FormPanel(struct)
			Detail:         types.NewInfoPanel(cfg.PrimaryKey.Name),
			CanAdd:         cfg.CanAdd,
			Editable:       cfg.Editable,
			Deletable:      cfg.Deletable,
			Exportable:     cfg.Exportable,
			PrimaryKey:     cfg.PrimaryKey,
			OnlyNewForm:    cfg.OnlyNewForm,
			OnlyUpdateForm: cfg.OnlyUpdateForm,
			OnlyDetail:     cfg.OnlyDetail,
			OnlyInfo:       cfg.OnlyInfo,
		},
		connectionDriver: cfg.Driver,
		connection:       cfg.Connection,
		sourceURL:        cfg.SourceURL,
		getDataFun:       cfg.GetDataFun,
	}
}

//-----------------------------table(interface)的方法--------------------------------

// GetData 透過參數處理後取得前端介面顯示資料，將值設置至PanelInfo(struct)
// PanelInfo裡的資訊有主題、描述名稱、可以篩選條件的欄位、選擇顯示的欄位、分頁、[]TheadItem(欄位資訊)等資訊
func (tb *DefaultTable) GetData(params parameter.Parameters) (PanelInfo, error) {
	var (
		data      []map[string]interface{}
		size      int
		beginTime = time.Now()
	)

	// -------一般用戶、角色、權限介面都不會執行---------
	if tb.Info.QueryFilterFn != nil {
		// db透過參數取得匹配的Service(interface)，接著將參數轉換為Connection(interface)
		ids, stop := tb.Info.QueryFilterFn(params, tb.db())
		if stop {
			return tb.GetDataWithIds(params.WithPKs(ids...))
		}
	}

	if tb.getDataFun != nil { // 一般都nil
		data, size = tb.getDataFun(params)
	} else if tb.sourceURL != "" { // 一般都為空
		data, size = tb.getDataFromURL(params)
	} else if tb.Info.GetDataFn != nil { // 一般都為nil;
		data, size = tb.Info.GetDataFn(params)
	} else if params.IsAll() {
		// 透過參數處理後取得前端介面顯示資料，將值設置至PanelInfo(struct)
		// 設置PanelInfo裡的主題、描述名稱、選擇顯示的資料、[]TheadItem(欄位資訊)
		return tb.getAllDataFromDatabase(params)
	} else {
		// ----------------一般執行此條件後直接return-----------------
		// 透過參數處理後取得前端介面顯示資料，將值設置至PanelInfo(struct)
		// 設置PanelInfo裡的主題、描述名稱、可以篩選條件的欄位、選擇顯示的資料、分頁、[]TheadItem(欄位資訊)
		return tb.getDataFromDatabase(params)
	}

	infoList := make(types.InfoList, 0)
	for i := 0; i < len(data); i++ {
		infoList = append(infoList, tb.getTempModelData(data[i], params, []string{}))
	}

	thead, _, _, _, _, filterForm := tb.getTheadAndFilterForm(params, []string{})
	endTime := time.Now()
	extraInfo := ""

	if !tb.Info.IsHideQueryInfo {
		extraInfo = fmt.Sprintf("<b>" + "query time" + ": </b>" +
			fmt.Sprintf("%.3fms", endTime.Sub(beginTime).Seconds()*1000))
	}
	return PanelInfo{
		Thead:    thead,
		InfoList: infoList,
		Paginator: paginator.Get(paginator.Config{
			Size:         size,
			Param:        params,
			PageSizeList: tb.Info.GetPageSizeList(),
		}).SetExtraInfo(template.HTML(extraInfo)),
		Title:          tb.Info.Title,
		FilterFormData: filterForm,
		Description:    tb.Info.Description,
	}, nil
}

// GetDataWithIds 透過參數(選擇取得特定id資料)處理後取得前端介面顯示資料，將值設置至PanelInfo(struct)
// PanelInfo裡的資訊有主題、描述名稱、可以篩選條件的欄位、選擇顯示的欄位、分頁、[]TheadItem(欄位資訊)等資訊
func (tb *DefaultTable) GetDataWithIds(params parameter.Parameters) (PanelInfo, error) {
	var (
		data      []map[string]interface{}
		size      int
		beginTime = time.Now()
	)

	if tb.getDataFun != nil {
		data, size = tb.getDataFun(params)
	} else if tb.sourceURL != "" {
		data, size = tb.getDataFromURL(params)
	} else if tb.Info.GetDataFn != nil {
		data, size = tb.Info.GetDataFn(params)
	} else {
		// 透過參數處理後取得前端介面顯示資料，將值設置至PanelInfo(struct)
		// 設置PanelInfo裡的主題、描述名稱、可以篩選條件的欄位、選擇顯示的資料、分頁、[]TheadItem(欄位資訊)
		// ----------大部分取得資料都執行此動作後直接return-------------------
		return tb.getDataFromDatabase(params)
	}

	infoList := make([]map[string]types.InfoItem, 0)

	for i := 0; i < len(data); i++ {
		infoList = append(infoList, tb.getTempModelData(data[i], params, []string{}))
	}

	thead, _, _, _, _, filterForm := tb.getTheadAndFilterForm(params, []string{})

	endTime := time.Now()

	return PanelInfo{
		Thead:    thead,
		InfoList: infoList,
		Paginator: paginator.Get(paginator.Config{
			Size:         size,
			Param:        params,
			PageSizeList: tb.Info.GetPageSizeList(),
		}).
			SetExtraInfo(template.HTML(fmt.Sprintf("<b>" + "query time" + ": </b>" +
				fmt.Sprintf("%.3fms", endTime.Sub(beginTime).Seconds()*1000)))),
		Title:          tb.Info.Title,
		FilterFormData: filterForm,
		Description:    tb.Info.Description,
	}, nil
}

// InsertData insert data.
func (tb *DefaultTable) InsertData(dataList form.Values) error {
	var (
		id     = int64(0)
		err    error
		errMsg = ""
	)

	dataList.Add("__post_type", "1")

	// -------------只有新增權限會執行(權限有設置PostHook)----------------
	if tb.Form.PostHook != nil {
		defer func() {
			dataList.Add("__post_type", "1")
			dataList.Add(tb.GetPrimaryKey().Name, strconv.Itoa(int(id)))
			dataList.Add("__post_result", errMsg)
			go func() {
				defer func() {
					if err := recover(); err != nil {
						panic(err)
					}
				}()
				err := tb.Form.PostHook(dataList)
				if err != nil {
					if err.Error() == "no affect row" {
						err = nil
					}
					panic(err)
				}
			}()
		}()
	}

	// -------------只有新增權限會執行(權限有設置Validator)----------------
	if tb.Form.Validator != nil {
		if err := tb.Form.Validator(dataList); err != nil {
			errMsg = "post error: " + err.Error()
			return err
		}
	}

	// -------都沒有設置PreProcessFn，不會執行---------
	// if tb.Form.PreProcessFn != nil {
	// 	dataList = tb.Form.PreProcessFn(dataList)
	// }

	// 用戶及角色頁面會執行新增資料的動作，直接return結果
	// --------------新增權限頁面不會執行------------------
	if tb.Form.InsertFn != nil {
		dataList.Delete("__post_type")
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

// UpdateData update data.
func (tb *DefaultTable) UpdateData(dataList form.Values) error {
	var (
		errMsg = ""
		err    error
	)

	dataList.Add("__post_type", "0")

	// -------------只有權限編輯介面會執行--------------
	if tb.Form.PostHook != nil {
		defer func() {
			dataList.Add("__post_type", "0")
			dataList.Add("__post_result", errMsg)
			go func() {
				defer func() {
					if err := recover(); err != nil {
						panic(err)
					}
				}()

				err := tb.Form.PostHook(dataList)
				if err != nil {
					panic(err)
				}
			}()
		}()
	}

	// ----------只有權限編輯介面會執行------------
	if tb.Form.Validator != nil {
		if err := tb.Form.Validator(dataList); err != nil {
			errMsg = "post error: " + err.Error()
			return err
		}
	}

	// -------都沒有設置PreProcessFn，不會執行---------
	// if tb.Form.PreProcessFn != nil {
	// 	dataList = tb.Form.PreProcessFn(dataList)
	// }

	// ------------用戶、角色介面有設置更新函式，直接執行並return--------
	if tb.Form.UpdateFn != nil {
		dataList.Delete("__post_type")
		// 更新資料
		err = tb.Form.UpdateFn(dataList)
		if err != nil {
			errMsg = "post error: " + err.Error()
		}
		return err
	}

	// ------------權限會執行更新--------------
	_, err = tb.sql().Table(tb.Form.Table).
		Where(tb.PrimaryKey.Name, "=", dataList.Get(tb.PrimaryKey.Name)).
		Update(tb.getInjectValueFromFormValue(dataList, types.PostTypeUpdate))
	if db.CheckError(err, db.UPDATE) {
		if err != nil {
			errMsg = "post error: " + err.Error()
		}
		return err
	}
	return nil
}

// DeleteData delete data.
func (tb *DefaultTable) DeleteData(id string) error {
	var (
		idArr = strings.Split(id, ",")
		err   error
	)

	// 目前沒有設置DeleteHook、DeleteHookWithRes
	// if tb.Info.DeleteHook != nil {
	// 	defer func() {
	// 		go func() {
	// 			defer func() {
	// 				if recoverErr := recover(); recoverErr != nil {
	// 					panic(recoverErr)
	// 				}
	// 			}()

	// 			if hookErr := tb.Info.DeleteHook(idArr); hookErr != nil {
	// 				panic(hookErr)
	// 			}
	// 		}()
	// 	}()
	// }
	// if tb.Info.DeleteHookWithRes != nil {
	// 	defer func() {
	// 		go func() {
	// 			defer func() {
	// 				if recoverErr := recover(); recoverErr != nil {
	// 					panic(recoverErr)
	// 				}
	// 			}()

	// 			if hookErr := tb.Info.DeleteHookWithRes(idArr, err); hookErr != nil {
	// 				panic(hookErr)
	// 			}
	// 		}()
	// 	}()
	// }
	// 用戶、角色、權限都是設置DeleteFn，沒有設置PreDeleteFn
	// if tb.Info.PreDeleteFn != nil {
	// 	if err = tb.Info.PreDeleteFn(idArr); err != nil {
	// 		return err
	// 	}
	// }

	// -------------用戶、角色、權限介面執行(都有在info設置刪除函式)後直接return------------------
	if tb.Info.DeleteFn != nil {
		err = tb.Info.DeleteFn(idArr)
		return err
	}

	// if len(idArr) == 0 || tb.Info.Table == "" {
	// 	err = errors.New("delete error: wrong parameter")
	// 	return err
	// }

	err = tb.delete(tb.Info.Table, tb.PrimaryKey.Name, idArr)
	return err
}

//-----------------------------table(interface)的方法--------------------------------

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

// delete delete data
func (tb *DefaultTable) delete(table, key string, values []string) error {
	var vals = make([]interface{}, len(values))
	for i, v := range values {
		vals[i] = v
	}

	return tb.sql().Table(table).
		WhereIn(key, vals).
		Delete()
}

// 透過參數取得匹配的Service(interface)，接著將參數轉換為Connection(interface)回傳並回傳
func (tb *DefaultTable) db() db.Connection {
	if tb.connectionDriver != "" && tb.getDataFromDB() {
		return db.GetConnectionFromService(services.Get(tb.connectionDriver))
	}
	return nil
}

// 取得分隔符號
func (tb *DefaultTable) delimiter() string {
	if tb.getDataFromDB() {
		return tb.db().GetDelimiter()
	}
	return ""
}

// 取得[]TheadItem(欄位資訊)、欄位名稱、joinFields(ex:group_concat(goadmin_roles.`name`...)、join語法(left join....)、合併的資料表、可篩選過濾的欄位
func (tb *DefaultTable) getTheadAndFilterForm(params parameter.Parameters, columns Columns) (types.Thead,
	string, string, string, []string, []types.FormField) {
	// TableInfo(struct) 資料表資訊
	return tb.Info.FieldList.GetTheadAndFilterForm(types.TableInfo{
		Table:      tb.Info.Table,       // ex: users
		Delimiter:  tb.delimiter(),      // ex:'
		Driver:     tb.connectionDriver, // ex: mysql
		PrimaryKey: tb.PrimaryKey.Name,  // ex: id
	}, params, columns, func() *db.SQL {
		return tb.sql()
	})
}

// getTempModelData 取得顯示在頁面上的資料(只取得選擇要顯示的欄位)
// 不管有沒有選擇要顯示ID欄位，都會將ID數值加入map[string]types.InfoItem
func (tb *DefaultTable) getTempModelData(res map[string]interface{}, params parameter.Parameters, columns Columns) map[string]types.InfoItem {
	var tempModelData = make(map[string]types.InfoItem)
	headField := ""

	// GetValueFromDatabaseType 取得id(主鍵)的值
	primaryKeyValue := db.GetValueFromDatabaseType(tb.PrimaryKey.Type, res[tb.PrimaryKey.Name], len(columns) == 0)

	// DefaultTable.Info.FieldList為介面所有欄位資訊
	for _, field := range tb.Info.FieldList {
		headField = field.Field // 欄位名稱

		// 用戶介面的角色欄位會執行(因為需要join其他表)
		if field.Joins.Valid() {
			// ex: roles_join_name
			headField = field.Joins.Last().Table + "_join_" + field.Field
		}

		if field.Hide {
			continue
		}
		// params.Columns為頁面中所選擇顯示的欄位
		// 如果沒有選擇要顯示該欄位則直接continue
		if !modules.InArrayWithoutEmpty(params.Columns, headField) {
			continue
		}

		typeName := field.TypeName
		// 將有join其他資料表的欄位(例如用戶頁面的角色欄位)類別設為VARCHAR
		if field.Joins.Valid() {
			typeName = db.Varchar
		}

		// GetValueFromDatabaseType 判斷條件後取得Value(string)
		combineValue := db.GetValueFromDatabaseType(typeName, res[headField], len(columns) == 0).String()

		var value interface{}
		// 取得欄位的值，角色欄位會取得ex: <span class="label label-succ.....HTML語法
		if len(columns) == 0 || modules.InArray(columns, headField) || field.Joins.Valid() {
			value = field.ToDisplay(types.FieldModel{
				ID:    primaryKeyValue.String(),
				Value: combineValue,
				Row:   res,
			})
		} else {
			value = field.ToDisplay(types.FieldModel{
				ID:    primaryKeyValue.String(),
				Value: "",
				Row:   res,
			})
		}

		if valueStr, ok := value.(string); ok {
			tempModelData[headField] = types.InfoItem{
				Content: template.HTML(valueStr),
				Value:   combineValue,
			}
		} else {
			// -----角色欄位執行------
			tempModelData[headField] = types.InfoItem{
				Content: value.(template.HTML),
				Value:   combineValue,
			}
		}
	}

	// 不管有沒有顯示id(主鍵)欄位，最後都會加上id欄位資訊
	primaryKeyField := tb.Info.FieldList.GetFieldByFieldName(tb.PrimaryKey.Name)
	value := primaryKeyField.ToDisplay(types.FieldModel{
		ID:    primaryKeyValue.String(),
		Value: primaryKeyValue.String(),
		Row:   res,
	})

	if valueStr, ok := value.(string); ok {
		tempModelData[tb.PrimaryKey.Name] = types.InfoItem{
			Content: template.HTML(valueStr),
			Value:   primaryKeyValue.String(),
		}
	} else {
		tempModelData[tb.PrimaryKey.Name] = types.InfoItem{
			Content: value.(template.HTML),
			Value:   primaryKeyValue.String(),
		}
	}
	return tempModelData
}

// 透過參數處理後取得前端介面顯示資料，將值設置至PanelInfo(struct)
// 設置PanelInfo裡的主題、描述名稱、選擇顯示的資料、[]TheadItem(欄位資訊)
func (tb *DefaultTable) getAllDataFromDatabase(params parameter.Parameters) (PanelInfo, error) {
	var (
		connection     = tb.db()
		queryStatement = "select %s from %s %s %s %s order by " + modules.Delimiter(connection.GetDelimiter(), "%s") + " %s"
	)

	// columns, _ := tb.getColumns(tb.Info.Table)
	columns, _ := tb.getColumns(tb.Info.Table)

	thead, fields, joins := tb.Info.FieldList.GetThead(types.TableInfo{
		Table:      tb.Info.Table,
		Delimiter:  tb.db().GetDelimiter(),
		Driver:     tb.connectionDriver,
		PrimaryKey: tb.PrimaryKey.Name,
	}, params, columns)

	fields += tb.Info.Table + "." + modules.FilterField(tb.PrimaryKey.Name, connection.GetDelimiter())

	groupBy := ""
	if joins != "" {
		groupBy = " GROUP BY " + tb.Info.Table + "." + modules.Delimiter(connection.GetDelimiter(), tb.PrimaryKey.Name)
	}

	var (
		wheres    = ""
		whereArgs = make([]interface{}, 0)
		existKeys = make([]string, 0)
	)

	// Statement 取得wheres語法、where數值、existKeys([]string)
	wheres, whereArgs, existKeys = params.Statement(wheres, tb.Info.Table, connection.GetDelimiter(), whereArgs, columns, existKeys,
		tb.Info.FieldList.GetFieldFilterProcessValue)
	// Statement 處理wheres語法及where值後回傳
	wheres, whereArgs = tb.Info.Wheres.Statement(wheres, connection.GetDelimiter(), whereArgs, existKeys, columns)
	wheres, whereArgs = tb.Info.WhereRaws.Statement(wheres, whereArgs)
	if wheres != "" {
		wheres = " where " + wheres
	}

	if !modules.InArray(columns, params.SortField) {
		params.SortField = tb.PrimaryKey.Name
	}

	queryCmd := fmt.Sprintf(queryStatement, fields, tb.Info.Table, joins, wheres, groupBy, params.SortField, params.SortType)
	res, err := connection.QueryWithConnection(tb.connection, queryCmd, whereArgs...)
	if err != nil {
		return PanelInfo{}, err
	}

	infoList := make([]map[string]types.InfoItem, 0)
	for i := 0; i < len(res); i++ {
		// getTempModelData 取得顯示在頁面上的資料(只取得選擇要顯示的欄位)
		// 不管有沒有選擇要顯示ID欄位，都會將數值加入map[string]types.InfoItem
		infoList = append(infoList, tb.getTempModelData(res[i], params, columns))
	}

	return PanelInfo{
		InfoList:    infoList,
		Thead:       thead,
		Title:       tb.Info.Title,
		Description: tb.Info.Description,
	}, nil
}

// 透過參數處理後取得前端介面顯示資料，將值設置至PanelInfo(struct)
// 設置PanelInfo裡的主題、描述名稱、可以篩選條件的欄位、選擇顯示的資料、分頁、[]TheadItem(欄位資訊)
func (tb *DefaultTable) getDataFromDatabase(params parameter.Parameters) (PanelInfo, error) {
	var (
		connection = tb.db()
		// Delimiter使用該資料庫引擎的符號
		placeholder    = modules.Delimiter(connection.GetDelimiter(), "%s") // ex: '%s'(mysql)
		queryStatement string
		countStatement string
		// 透過參數__pk尋找Parameters.Fields[__pk]是否存在，如果存在則回傳第一個value值(string)並且用","拆解成[]string
		ids = params.PKs()                                                                           // ex:[]
		pk  = tb.Info.Table + "." + modules.Delimiter(connection.GetDelimiter(), tb.PrimaryKey.Name) // ex: users.`id`
	)

	// 判斷是否資料庫引擎為postgresql
	if connection.Name() == "postgresql" {
		placeholder = "%s"
	}
	beginTime := time.Now()

	// 判斷是否挑選特定id資料
	if len(ids) > 0 {
		countExtra := ""
		if connection.Name() == "mssql" {
			countExtra = "as [size]"
		}
		queryStatement = "select %s from " + placeholder + " %s where " + pk + " in (%s) %s ORDER BY %s." + placeholder + " %s"
		countStatement = "select count(*) " + countExtra + " from " + placeholder + " %s where " + pk + " in (%s)"
	} else {
		if connection.Name() == "mssql" {
			queryStatement = "SELECT * FROM (SELECT ROW_NUMBER() OVER (ORDER BY %s." + placeholder + " %s) as ROWNUMBER_, %s from " +
				placeholder + "%s %s %s ) as TMP_ WHERE TMP_.ROWNUMBER_ > ? AND TMP_.ROWNUMBER_ <= ?"
			countStatement = "select count(*) as [size] from " + placeholder + " %s %s"
		} else {
			queryStatement = "select %s from " + placeholder + "%s %s %s order by %s." + placeholder + " %s LIMIT ? OFFSET ?"
			countStatement = "select count(*) from " + placeholder + " %s %s"
		}
	}

	// getColumns(取得資料表欄位)將欄位名稱加入columns([]string)
	columns, _ := tb.getColumns(tb.Info.Table)

	// 取得[]TheadItem(欄位資訊)、欄位名稱、joinFields(ex:group_concat(goadmin_roles.`name`...)、join語法(left join....)、合併的資料表、可篩選過濾的欄位
	thead, fields, joinFields, joins, joinTables, filterForm := tb.getTheadAndFilterForm(params, columns)

	fields += pk
	allFields := fields
	groupFields := fields

	if joinFields != "" {
		allFields += "," + joinFields[:len(joinFields)-1]
		// -------mssql引擎才會執行下面語法處理------------
		if connection.Name() == "mssql" {
			for _, field := range tb.Info.FieldList {
				if field.TypeName == db.Text || field.TypeName == db.Longtext {
					f := modules.Delimiter(connection.GetDelimiter(), field.Field)
					headField := tb.Info.Table + "." + f
					allFields = strings.Replace(allFields, headField, "CAST("+headField+" AS NVARCHAR(MAX)) as "+f, -1)
					groupFields = strings.Replace(groupFields, headField, "CAST("+headField+" AS NVARCHAR(MAX))", -1)
				}
			}
		}
	}

	// ------------一般不會執行，都有設置排序欄位------------------
	if !modules.InArray(columns, params.SortField) {
		params.SortField = tb.PrimaryKey.Name
	}

	var (
		wheres    = ""
		whereArgs = make([]interface{}, 0)
		args      = make([]interface{}, 0)
		existKeys = make([]string, 0)
	)

	// -------一般介面的ids為空[]-------------------
	if len(ids) > 0 {
		for _, value := range ids {
			if value != "" {
				wheres += "?,"
				args = append(args, value)
			}
		}
		wheres = wheres[:len(wheres)-1]
	} else {
		// Statement 取得wheres語法、where數值、existKeys([]string)
		// --------一般取得空值，只有existKeys=[__is_all]----------------
		wheres, whereArgs, existKeys = params.Statement(wheres, tb.Info.Table, connection.GetDelimiter(), whereArgs, columns, existKeys,
			tb.Info.FieldList.GetFieldFilterProcessValue)

		// Statement 處理wheres語法及where值後回傳
		// ----一般DefaultTable.Info.Wheres為空，回傳的值不變-------
		wheres, whereArgs = tb.Info.Wheres.Statement(wheres, connection.GetDelimiter(), whereArgs, existKeys, columns)

		// ----一般DefaultTable.Info.WhereRaws為空，回傳的值不變-------
		wheres, whereArgs = tb.Info.WhereRaws.Statement(wheres, whereArgs)

		if wheres != "" {
			wheres = " where " + wheres
		}

		if connection.Name() == "mssql" {
			args = append(whereArgs, (params.PageInt-1)*params.PageSizeInt, params.PageInt*params.PageSizeInt)
		} else {
			args = append(whereArgs, params.PageSizeInt, (params.PageInt-1)*params.PageSizeInt)
		}
	}

	groupBy := ""
	if len(joinTables) > 0 {
		if connection.Name() == "mssql" {
			groupBy = " GROUP BY " + groupFields
		} else {
			groupBy = " GROUP BY " + pk
		}
	}

	queryCmd := ""
	if connection.Name() == "mssql" && len(ids) == 0 {
		queryCmd = fmt.Sprintf(queryStatement, tb.Info.Table, params.SortField, params.SortType,
			allFields, tb.Info.Table, joins, wheres, groupBy)
	} else {
		queryCmd = fmt.Sprintf(queryStatement, allFields, tb.Info.Table, joins, wheres, groupBy,
			tb.Info.Table, params.SortField, params.SortType)
	}

	res, err := connection.QueryWithConnection(tb.connection, queryCmd, args...)
	if err != nil {
		return PanelInfo{}, err
	}

	infoList := make([]map[string]types.InfoItem, 0)
	// 將取得資料處理
	for i := 0; i < len(res); i++ {
		// getTempModelData 取得顯示在頁面上的資料(只取得選擇要顯示的欄位)
		// 不管有沒有選擇要顯示ID欄位，都會將數值加入map[string]types.InfoItem
		infoList = append(infoList, tb.getTempModelData(res[i], params, columns))
	}

	var size int

	// ------一般會len(ids)=0，會執行------------
	if len(ids) == 0 {
		// 計算資料筆數語法
		countCmd := fmt.Sprintf(countStatement, tb.Info.Table, joins, wheres)
		// ex: total: [map[count(*):4]](4筆)
		total, err := connection.QueryWithConnection(tb.connection, countCmd, whereArgs...)
		if err != nil {
			return PanelInfo{}, err
		}

		if tb.connectionDriver == "postgresql" {
			size = int(total[0]["count"].(int64))
		} else if tb.connectionDriver == "mssql" {
			size = int(total[0]["size"].(int64))
		} else {
			size = int(total[0]["count(*)"].(int64)) // ex:4(4筆符合)
		}
	}

	endTime := time.Now()

	return PanelInfo{
		Thead:    thead,
		InfoList: infoList, // 顯示在介面上的所有資料
		// GetPaginator 設置分頁資訊
		Paginator: tb.GetPaginator(size, params,
			template.HTML(fmt.Sprintf("<b>"+"query time"+": </b>"+
				fmt.Sprintf("%.3fms", endTime.Sub(beginTime).Seconds()*1000)))),
		Title:          tb.Info.Title,       // 左上角主題
		FilterFormData: filterForm,          // 可以篩選條件的欄位
		Description:    tb.Info.Description, //主題旁的描述
	}, nil
}

// getDataFromURL(從url中取得data)
func (tb *DefaultTable) getDataFromURL(params parameter.Parameters) ([]map[string]interface{}, int) {

	u := ""
	if strings.Contains(tb.sourceURL, "?") {
		u = tb.sourceURL + "&" + params.Join()
	} else {
		u = tb.sourceURL + "?" + params.Join()
	}
	// 透過參數__pk尋找Parameters.Fields[__pk]是否存在，如果存在則回傳第一個value值(string)並且用","拆解成[]string
	res, err := http.Get(u + "&pk=" + strings.Join(params.PKs(), ","))

	if err != nil {
		return []map[string]interface{}{}, 0
	}

	defer func() {
		_ = res.Body.Close()
	}()

	body, err := ioutil.ReadAll(res.Body)

	if err != nil {
		return []map[string]interface{}{}, 0
	}

	var data GetDataFromURLRes

	err = json.Unmarshal(body, &data)

	if err != nil {
		return []map[string]interface{}{}, 0
	}

	return data.Data, data.Size
}
