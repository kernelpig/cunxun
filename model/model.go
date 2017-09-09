package model

import (
	"database/sql"
	"fmt"
	"strings"

	e "wangqingang/cunxun/error"
	"wangqingang/cunxun/utils"
)

const (
	columnTagKey = "column"
	pageNumStart = 1
	pageSize     = 20
)

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

// pageSize<=len(selects), pageNum待获取的页数数据, 数据页码从1开始, 遇到任何错误都返回处理完成, 不再处理后续页面
func SQLQueryRows(db sqlExec, selects *[]interface{}, wheres map[string]interface{}, pageSize, pageNum int) (bool, error) {
	var f []string
	var _SQL string
	var rows *sql.Rows
	var err error

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
		_SQL = fmt.Sprintf("SELECT %s FROM `%s` limit %d, %d",
			strings.Join(f, ", "), tableName, pageOffset, pageSize)
		rows, err = db.Query(_SQL)
	} else {
		// 带有where查询条件
		var w []string
		var q []interface{}

		for key, value := range wheres {
			w = append(w, fmt.Sprintf("%s = ?", key))
			q = append(q, value)
		}
		_SQL = fmt.Sprintf("SELECT %s FROM `%s` WHERE %s limit %d, %d",
			strings.Join(f, ", "), tableName, strings.Join(w, " and "), pageOffset, pageSize)
		rows, err = db.Query(_SQL, q...)
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
	if rowsAffected <= int64(pageSize) {
		return true, nil
	}
	return false, nil
}

func SQLQueryRow(db sqlExec, selects interface{}, wheres map[string]interface{}) (bool, error) {
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
		w = append(w, fmt.Sprintf("%s = ?", key))
		q = append(q, value)
	}

	_SQL := fmt.Sprintf("SELECT %s FROM `%s` WHERE %s", strings.Join(f, ", "), tableName, strings.Join(w, " and "))
	err := db.QueryRow(_SQL, q...).Scan(s...)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, e.SP(e.MMysqlErr, e.MysqlSelectErr, err)
	}

	return true, nil
}

func SQLUpdate(db sqlExec, updates interface{}, wheres map[string]interface{}) (int64, error) {
	var q []interface{}

	var u []string
	tableName, updatesMap := utils.Struct2MapWithValue(updates, columnTagKey, true)

	for key, value := range updatesMap {
		u = append(u, fmt.Sprintf("%s = ?", key))
		q = append(q, value)
	}

	var w []string
	for key, value := range wheres {
		w = append(w, fmt.Sprintf("%s = ?", key))
		q = append(q, value)
	}

	_SQL := fmt.Sprintf("UPDATE `%s` SET %s WHERE %s", tableName, strings.Join(u, ", "), strings.Join(w, " and "))
	sqlResult, err := db.Exec(_SQL, q...)
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
	sqlResult, err := db.Exec(_SQL, q...)

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
