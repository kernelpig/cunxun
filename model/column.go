package model

import (
	"time"

	e "wangqingang/cunxun/error"
	"wangqingang/cunxun/id"
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
	ID         uint64    `json:"id" column:"id"`
	Name       string    `json:"name" column:"name"`
	CreaterUid uint64    `json:"creater_uid" column:"creater_uid"`
	CreatedAt  time.Time `json:"created_at" column:"created_at"`
	UpdatedAt  time.Time `json:"updated_at" column:"updated_at"`
}

type ColumnDetailView struct {
	ID         uint64    `json:"id" column:"id"`
	Name       string    `json:"name" column:"name"`
	CreaterUid uint64    `json:"creater_uid" column:"creater_uid"`
	CreatedAt  time.Time `json:"created_at" column:"created_at"`
	UpdatedAt  time.Time `json:"updated_at" column:"updated_at"`
	Nickname   string    `json:"nickname" column:"nickname"`
	// TODO: 下面这个字段干啥的? 暂时没有发现用处, 确认一下前端是否用到
	ColumnCount int `json:"column_count" column:"column_count"`
}

func GetColumnByID(db sqlExec, columnID uint64) (*Column, error) {
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
	// 未设置ID使用自动生成的id, 主要是考虑到人工设置特殊的ID场景
	if column.ID == 0 {
		id, err := id.Generate()
		if err != nil {
			return nil, err
		}
		column.ID = id
	}
	_, err := SQLInsert(db, column)
	if err != nil {
		if isDBDuplicateErr(err) {
			return nil, e.SP(e.MColumnErr, e.ColumnAlreadyExist, err)
		}
		return nil, e.SP(e.MColumnErr, e.ColumnCreateErr, err)
	}
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

func UpdateColumnById(db sqlExec, columnId uint64, valueWillSet *Column) (int64, error) {
	return UpdateColumnList(db, map[string]interface{}{"id": columnId}, valueWillSet)
}

func UpdateColumnByIdOfSelf(db sqlExec, columnId, userId uint64, valueWillSet *Column) (int64, error) {
	return UpdateColumnList(db, map[string]interface{}{"id": columnId, "creater_uid": userId}, valueWillSet)
}

func DeleteColumnById(db sqlExec, columnId uint64) (int64, error) {
	return DeleteColumnList(db, map[string]interface{}{"id": columnId})
}

func DeleteColumnByIdOfSelf(db sqlExec, columnId, userId uint64) (int64, error) {
	return DeleteColumnList(db, map[string]interface{}{"id": columnId, "creater_uid": userId})
}

func GetColumnList(db sqlExec, where map[string]interface{}) ([]*ColumnDetailView, error) {
	var list []*ColumnDetailView

	// 初始化缓冲区
	var modelBuf = make([]interface{}, 0)
	for i := 0; i < pageSize; i++ {
		modelBuf = append(modelBuf, &ColumnDetailView{})
	}

	// 每次只取pageSize个
	for pageNum := 1; true; pageNum++ {
		isOver, err := SQLQueryRows(db, &modelBuf, where, []string{OrderById}, pageSize, pageNum)
		if err != nil {
			return nil, e.SP(e.MColumnErr, e.ColumnGetAllErr, err)
		}
		for _, item := range modelBuf {
			if model, ok := item.(*ColumnDetailView); ok {
				list = append(list, model)
			}
		}
		if isOver {
			break
		}
	}
	return list, nil
}
