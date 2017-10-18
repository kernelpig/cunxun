package model

import (
	"database/sql"
	"fmt"
	"strings"

	"wangqingang/cunxun/common"
	e "wangqingang/cunxun/error"
	"wangqingang/cunxun/utils"
)

// orm处理常量定义
const (
	columnTagKey = "column"
	pageNumStart = 1
	pageSize     = 20
)

// 排序方式枚举定义,禁止用户输入DB字段映射,防止SQL注入
const (
	OrderByCreateDate   = "create_date"
	OrderByUpdateDate   = "update_date"
	OrderByCommentCount = "comment_count"
	OrderByPriority     = "priority"
	OrderByIgnore       = ""
	OrderById           = "id"
)

// 参数与DB字段映射表, 禁止用户输入DB字段映射,防止SQL注入
var OrderByMap map[string]string

func init() {
	OrderByMap = make(map[string]string)
	OrderByMap[OrderByIgnore] = "created_at desc"
	OrderByMap[OrderByCreateDate] = "created_at desc"
	OrderByMap[OrderByUpdateDate] = "updated_at desc"
	OrderByMap[OrderByCommentCount] = "comment_count desc"
	OrderByMap[OrderByPriority] = "priority desc"
	OrderByMap[OrderById] = "id asc"
}

type sqlExec interface {
	Exec(query string, args ...interface{}) (sql.Result, error)
	QueryRow(query string, args ...interface{}) *sql.Row
	Query(query string, args ...interface{}) (*sql.Rows, error)
}

func isMysqlDuplicateErr(err error) bool {
	return strings.Contains(err.Error(), "1062")
}

func isDBDuplicateErr(err error) bool {
	if messageErr, ok := err.(e.Message); ok {
		return messageErr.Code.IsSubError(e.MMysqlErr, e.MysqlDuplicateErr)
	}
	return false
}

func dumpSQL(_SQL string) string {
	if !common.Config.ReleaseMode {
		fmt.Println("_SQL: ", _SQL)
	}
	return _SQL
}

// 获取order by字符串
func getOrderByStr(orderByKeys []string) (string, error) {
	orderByValues := make([]string, 0)
	for _, k := range orderByKeys {
		if v, ok := OrderByMap[k]; ok {
			orderByValues = append(orderByValues, v)
		} else {
			return "", fmt.Errorf("invalid order by key: %s", k)
		}
	}
	// 增加时间次要排序
	orderByValues = append(orderByValues, OrderByMap[OrderByCreateDate])
	return strings.Join(orderByValues, ", "), nil
}

// pageSize<=len(selects), pageNum待获取的页数数据, 数据页码从1开始, 遇到任何错误都返回处理完成, 不再处理后续页面
// 默认自带created_at字段次要排序，不需要再加了
func SQLQueryRows(db sqlExec, selects *[]interface{}, wheres map[string]interface{}, orderByKeys []string, pageSize, pageNum int) (bool, error) {
	var f []string
	var _SQL string
	var rows *sql.Rows
	var err error
	var orderBy string

	// 如果有排序则必须符合orderMap
	orderByValues, err := getOrderByStr(orderByKeys)
	if err != nil {
		return true, e.SD(e.MParamsErr, e.ParamsInvalidOrderBy, err.Error())
	}
	orderBy = fmt.Sprintf("order by %s", orderByValues)

	// 数据页码从1开始
	if pageNum < pageNumStart {
		return true, e.S(e.MMysqlErr, e.MysqlInvalidPageNum)
	} else if len(*selects) < pageSize {
		return true, e.S(e.MMysqlErr, e.MysqlNoEnoughModelBuf)
	} else if len(*selects) > pageSize {
		*selects = (*selects)[:pageSize]
	}

	// 计算limit分页偏移
	pageOffset := pageSize * (pageNum - 1)

	// 获取表名及所有字段
	tableName, f := utils.StructGetFieldName((*selects)[0], columnTagKey)

	// 不带有where查询条件
	if wheres == nil || len(wheres) == 0 {
		_SQL = fmt.Sprintf("SELECT %s FROM `%s` %s limit %d, %d",
			strings.Join(f, ", "), tableName, orderBy, pageOffset, pageSize)
		rows, err = db.Query(dumpSQL(_SQL))
	} else {
		// 带有where查询条件
		var w []string
		var q []interface{}

		for key, value := range wheres {
			// 为类型的默认值则不增加此查询条件
			if !utils.IsTypeDefaultValue(value) {
				w = append(w, fmt.Sprintf("%s = ?", key))
				q = append(q, value)
			}
		}
		// 过滤掉类型默认值后无查询条件
		if len(w) == 0 || len(q) == 0 {
			_SQL = fmt.Sprintf("SELECT %s FROM `%s` %s limit %d, %d",
				strings.Join(f, ", "), tableName, orderBy, pageOffset, pageSize)
			rows, err = db.Query(dumpSQL(_SQL))
		} else {
			_SQL = fmt.Sprintf("SELECT %s FROM `%s` WHERE %s %s limit %d, %d",
				strings.Join(f, ", "), tableName, strings.Join(w, " and "), orderBy, pageOffset, pageSize)
			rows, err = db.Query(dumpSQL(_SQL), q...)
		}
	}

	if err != nil {
		if err == sql.ErrNoRows {
			return true, nil
		}
		return true, e.SP(e.MMysqlErr, e.MysqlSelectErr, err)
	}
	defer rows.Close()

	var rowsAffected int64
	for rows.Next() {
		// 获取当前的结构中所有字段地址, 按照字段排序
		_, fieldAddr := utils.StructGetFieldAddr((*selects)[rowsAffected], columnTagKey)
		if err := rows.Scan(fieldAddr...); err != nil {
			return true, e.SP(e.MMysqlErr, e.MysqlRowScanErr, err)
		}
		rowsAffected++
	}

	if err := rows.Err(); err != nil {
		return true, e.SP(e.MMysqlErr, e.MysqlRowScanErr, err)
	}

	// 收缩buf, 去掉没有用到的数据
	*selects = (*selects)[:rowsAffected]
	if rowsAffected < int64(pageSize) {
		return true, nil
	}
	return false, nil
}

func SQLQueryRow(db sqlExec, selects interface{}, wheres map[string]interface{}) (bool, error) {
	// 没有查询条件
	if wheres == nil || len(wheres) == 0 {
		return false, nil
	}

	var f []string
	var s []interface{}
	tableName, selectsMap := utils.Struct2MapWithAddr(selects, columnTagKey)
	for key, value := range selectsMap {
		f = append(f, key)
		s = append(s, value)
	}

	var w []string
	var q []interface{}
	for key, value := range wheres {
		// 过滤掉类型默认值查询条件
		if !utils.IsTypeDefaultValue(value) {
			w = append(w, fmt.Sprintf("%s = ?", key))
			q = append(q, value)
		}
	}
	// 过滤掉类型默认值后无查询条件
	if len(w) == 0 || len(q) == 0 {
		return false, nil
	}

	_SQL := fmt.Sprintf("SELECT %s FROM `%s` WHERE %s", strings.Join(f, ", "), tableName, strings.Join(w, " and "))
	err := db.QueryRow(dumpSQL(_SQL), q...).Scan(s...)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, e.SP(e.MMysqlErr, e.MysqlSelectErr, err)
	}

	return true, nil
}

// nullValue直接使用对应model类型的空结构, 只是为了获得表名
func SQLDelete(db sqlExec, nullValue interface{}, wheres map[string]interface{}) (int64, error) {
	tableName, _ := utils.Struct2MapWithValue(nullValue, columnTagKey, true)

	var w []string
	var q []interface{}

	// 没有查询条件
	if wheres == nil || len(wheres) == 0 {
		return 0, nil
	}
	for key, value := range wheres {
		if !utils.IsTypeDefaultValue(value) {
			// 过滤掉类型默认值筛选条件
			w = append(w, fmt.Sprintf("%s = ?", key))
			q = append(q, value)
		}
	}
	// 过滤掉类型默认值无筛选条件
	if len(w) == 0 || len(q) == 0 {
		return 0, nil
	}

	_SQL := fmt.Sprintf("DELETE FROM `%s` WHERE %s", tableName, strings.Join(w, " and "))
	sqlResult, err := db.Exec(dumpSQL(_SQL), q...)
	if err != nil {
		return 0, e.SP(e.MMysqlErr, e.MysqlDeleteErr, err)
	}

	rowAffected, err := sqlResult.RowsAffected()
	if err != nil {
		return 0, e.SP(e.MMysqlErr, e.MysqlRowAffectErr, err)
	}
	return rowAffected, nil
}

// updates使用新的对象, 不要使用携带多余字段值的对象, 防止误修改
func SQLUpdate(db sqlExec, updates interface{}, wheres map[string]interface{}) (int64, error) {
	// 无筛选条件
	if wheres == nil || len(wheres) == 0 {
		return 0, nil
	}

	var u []string
	var q []interface{}
	tableName, updatesMap := utils.Struct2MapWithValue(updates, columnTagKey, true)
	for key, value := range updatesMap {
		u = append(u, fmt.Sprintf("%s = ?", key))
		q = append(q, value)
	}

	var w []string
	for key, value := range wheres {
		if !utils.IsTypeDefaultValue(value) {
			// 过滤掉类型默认值条件
			w = append(w, fmt.Sprintf("%s = ?", key))
			q = append(q, value)
		}
	}
	// 筛选后无条件
	if len(w) == 0 || len(q) == 0 {
		return 0, nil
	}

	_SQL := fmt.Sprintf("UPDATE `%s` SET %s WHERE %s", tableName, strings.Join(u, ", "), strings.Join(w, " and "))
	sqlResult, err := db.Exec(dumpSQL(_SQL), q...)
	if err != nil {
		if isMysqlDuplicateErr(err) {
			return 0, e.SP(e.MMysqlErr, e.MysqlDuplicateErr, err)
		}
		return 0, e.SP(e.MMysqlErr, e.MysqlUpdateErr, err)
	}

	rowAffected, err := sqlResult.RowsAffected()
	if err != nil {
		return 0, e.SP(e.MMysqlErr, e.MysqlRowAffectErr, err)
	}
	return rowAffected, nil
}

func SQLInsert(db sqlExec, inserts interface{}) (int64, error) {
	var q []interface{}

	var u []string
	tableName, insertssMap := utils.Struct2MapWithValue(inserts, columnTagKey, true)

	for key, value := range insertssMap {
		u = append(u, fmt.Sprintf("%s = ?", key))
		q = append(q, value)
	}

	_SQL := fmt.Sprintf("INSERT INTO `%s` SET %s", tableName, strings.Join(u, ", "))
	sqlResult, err := db.Exec(dumpSQL(_SQL), q...)

	if err != nil {
		if isMysqlDuplicateErr(err) {
			return 0, e.SP(e.MMysqlErr, e.MysqlDuplicateErr, err)
		}
		return 0, e.SP(e.MMysqlErr, e.MysqlInsertErr, err)
	}

	lastInsertId, err := sqlResult.LastInsertId()
	if err != nil {
		return 0, e.SP(e.MMysqlErr, e.MysqlLastInsertErr, err)
	}
	return lastInsertId, nil
}
