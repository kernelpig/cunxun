package model

import (
	e "wangqingang/cunxun/error"
)

type Column struct {
	ID         int    `column:"id"`
	Name       string `column:"name"`
	CreaterUid int    `column:"creater_uid"`
}

func GetColumnByID(db sqlExec, columnID int) (*Column, error) {
	u := &Column{}
	isFound, err := SQLQueryRow(db, u, map[string]interface{}{
		"id": columnID,
	})
	if err != nil {
		return nil, e.SP(e.MColumnErr, e.ColumnGetErr, err)
	} else if !isFound {
		return nil, nil
	} else {
		return u, nil
	}
}

func CreateColumn(db sqlExec, column *Column) (*Column, error) {
	id, err := SQLInsert(db, column)
	if err != nil {
		if isDBDuplicateErr(err) {
			return nil, e.SP(e.MColumnErr, e.ColumnAlreadyExist, err)
		}
		return nil, e.SP(e.MColumnErr, e.ColumnCreateErr, err)
	}
	column.ID = int(id)
	return column, nil
}
