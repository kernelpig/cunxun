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

func SQLQueryRows(db sqlExec, selects []interface{}, wheres map[string]interface{}) (int64, error) {
	var f []string
	tableName, f := utils.StructGetFieldName(selects[0], columnTagKey)

	var w []string
	var q []interface{}

	for key, value := range wheres {
		w = append(w, fmt.Sprintf("%s = ?", key))
		q = append(q, value)
	}

	_SQL := fmt.Sprintf("SELECT %s FROM `%s` WHERE %s", strings.Join(f, ", "), tableName, strings.Join(w, " and "))
	rows, err := db.Query(_SQL, q...)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, nil
		}
		return 0, e.SP(e.MMysqlErr, e.MysqlSelectErr, err)
	}
	defer rows.Close()

	var rowsAffected int64
	for rows.Next() {
		var s []interface{}
		_, selectsMap := utils.Struct2MapWithAddr(selects[rowsAffected], columnTagKey)

		for _, value := range selectsMap {
			s = append(s, value)
		}
		if err := rows.Scan(s...); err != nil {
			return 0, e.SP(e.MMysqlErr, e.MysqlRowScanErr, err)
		}
		rowsAffected++
	}

	if err := rows.Err(); err != nil {
		return 0, e.SP(e.MMysqlErr, e.MysqlRowScanErr, err)
	}

	return rowsAffected, nil
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
