package model

import (
	"time"

	e "wangqingang/cunxun/error"
)

// 固定类别ID
const (
	ColumnIdNews       = 1
	ColumnIdBar        = 2
	ColumnIdRental     = 3
	ColumnIdCarpooling = 4
)

// 固定类别名称
const (
	ColumnNameNews       = "资讯"
	ColumnNameBar        = "贴吧"
	ColumnNameRental     = "租房"
	ColumnNameCarpooling = "拼车"
)

// 初始化栏目列表
var ColumnsOfOrigin map[int]string

func init() {
	ColumnsOfOrigin = make(map[int]string)
	ColumnsOfOrigin[ColumnIdNews] = ColumnNameNews
	ColumnsOfOrigin[ColumnIdBar] = ColumnNameBar
	ColumnsOfOrigin[ColumnIdRental] = ColumnNameRental
	ColumnsOfOrigin[ColumnIdCarpooling] = ColumnNameCarpooling
}

type Column struct {
	ID         int       `json:"id" column:"id"`
	Name       string    `json:"name" column:"name"`
	CreaterUid int       `json:"creater_uid" column:"creater_uid"`
	CreatedAt  time.Time `json:"created_at" column:"created_at"`
	UpdatedAt  time.Time `json:"updated_at" column:"updated_at"`
}

type ColumnListView struct {
	ID          int       `json:"id" column:"id"`
	Name        string    `json:"name" column:"name"`
	CreaterUid  int       `json:"creater_uid" column:"creater_uid"`
	CreatedAt   time.Time `json:"created_at" column:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" column:"updated_at"`
	Nickname    string    `json:"nickname" column:"nickname"`
	ColumnCount int       `json:"column_count" column:"column_count"`
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

func UpdateColumnList(db sqlExec, wheres map[string]interface{}, valueWillSet *Column) (int64, error) {
	count, err := SQLUpdate(db, valueWillSet, wheres)
	if err != nil {
		return 0, e.SP(e.MColumnErr, e.ColumnUpdateErr, err)
	} else {
		return count, nil
	}
}

func DeleteColumnList(db sqlExec, wheres map[string]interface{}) (int64, error) {
	count, err := SQLDelete(db, &Column{}, wheres)
	if err != nil {
		return 0, e.SP(e.MColumnErr, e.ColumnDeleteErr, err)
	} else {
		return count, nil
	}
}

func UpdateColumnById(db sqlExec, columnId int, valueWillSet *Column) (int64, error) {
	return UpdateColumnList(db, map[string]interface{}{"id": columnId}, valueWillSet)
}

func UpdateColumnByIdOfSelf(db sqlExec, columnId, userId int, valueWillSet *Column) (int64, error) {
	return UpdateColumnList(db, map[string]interface{}{"id": columnId, "creater_uid": userId}, valueWillSet)
}

func DeleteColumnById(db sqlExec, columnId int) (int64, error) {
	return DeleteColumnList(db, map[string]interface{}{"id": columnId})
}

func DeleteColumnByIdOfSelf(db sqlExec, columnId, userId int) (int64, error) {
	return DeleteColumnList(db, map[string]interface{}{"id": columnId, "creater_uid": userId})
}

func GetColumnList(db sqlExec, where map[string]interface{}) ([]*ColumnListView, error) {
	var list []*ColumnListView

	// 初始化缓冲区
	var modelBuf = make([]interface{}, 0)
	for i := 0; i < pageSize; i++ {
		modelBuf = append(modelBuf, &ColumnListView{})
	}

	// 每次只取pageSize个
	for pageNum := 1; true; pageNum++ {
		isOver, err := SQLQueryRows(db, &modelBuf, where, OrderById, pageSize, pageNum)
		if err != nil {
			return nil, e.SP(e.MColumnErr, e.ColumnGetAllErr, err)
		}
		for _, item := range modelBuf {
			if model, ok := item.(*ColumnListView); ok {
				list = append(list, model)
			}
		}
		if isOver {
			break
		}
	}
	return list, nil
}
