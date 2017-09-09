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

func GetAllColumn(db sqlExec) ([]*Column, error) {
	var list []*Column

	// 初始化缓冲区
	var modelBuf = make([]interface{}, 0)
	for i := 0; i < pageSize; i++ {
		modelBuf = append(modelBuf, &Column{})
	}

	// 每次只取pageSize个
	for pageNum := 1; true; pageNum++ {
		isOver, err := SQLQueryRows(db, &modelBuf, map[string]interface{}{}, pageSize, pageNum)
		if err != nil {
			return nil, e.SP(e.MColumnErr, e.ColumnGetErr, err)
		}
		for _, item := range modelBuf {
			if model, ok := item.(*Column); ok {
				list = append(list, model)
			}
		}
		if isOver {
			break
		}
	}
	return list, nil
}
