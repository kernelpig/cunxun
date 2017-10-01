package model

import (
	"time"

	e "wangqingang/cunxun/error"
	"wangqingang/cunxun/id"
)

const (
	CarpoolingDisable = 0
	CarpoolingEnable  = 1
)

type Carpooling struct {
	ID          uint64    `json:"id" column:"id"`
	CreatedAt   time.Time `json:"created_at" column:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" column:"updated_at"`
	CreaterUid  uint64    `json:"creater_uid" column:"creater_uid"`
	UpdaterUid  uint64    `json:"updater_uid" column:"updater_uid"`
	FromCity    string    `json:"from_city" column:"from_city"`
	ToCity      string    `json:"to_city" column:"to_city"`
	DepartTime  time.Time `json:"depart_time" column:"depart_time"`
	PeopleCount int       `json:"people_count" column:"people_count"`
	Contact     string    `json:"contact" column:"contact"`
	Status      int       `json:"status" column:"status"`
	Remark      string    `json:"remark" column:"remark"`
}

type CarpoolingDetailView struct {
	ID          uint64    `json:"id" column:"id"`
	CreatedAt   time.Time `json:"created_at" column:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" column:"updated_at"`
	CreaterUid  uint64    `json:"creater_uid" column:"creater_uid"`
	UpdaterUid  uint64    `json:"updater_uid" column:"updater_uid"`
	FromCity    string    `json:"from_city" column:"from_city"`
	ToCity      string    `json:"to_city" column:"to_city"`
	DepartTime  time.Time `json:"depart_time" column:"depart_time"`
	PeopleCount int       `json:"people_count" column:"people_count"`
	Contact     string    `json:"contact" column:"contact"`
	Status      int       `json:"status" column:"status"`
	Remark      string    `json:"remark" column:"remark"`
	Nickname    string    `json:"nickname" column:"nickname"`
}

func GetCarpoolingByID(db sqlExec, CarpoolingID uint64) (*CarpoolingDetailView, error) {
	u := &CarpoolingDetailView{}
	isFound, err := SQLQueryRow(db, u, map[string]interface{}{
		"id": CarpoolingID,
	})
	if err != nil {
		return nil, e.SP(e.MCarpoolingErr, e.CarpoolingGetErr, err)
	} else if !isFound {
		return nil, nil
	} else {
		return u, nil
	}
}

func CreateCarpooling(db sqlExec, carpooling *Carpooling) (*Carpooling, error) {
	// 未设置ID使用自动生成的id, 主要是考虑到人工设置特殊的ID场景
	if carpooling.ID == 0 {
		id, err := id.Generate()
		if err != nil {
			return nil, err
		}
		carpooling.ID = id
	}
	_, err := SQLInsert(db, carpooling)
	if err != nil {
		if isDBDuplicateErr(err) {
			return nil, e.SP(e.MCarpoolingErr, e.CarpoolingAlreadyExist, err)
		}
		return nil, e.SP(e.MCarpoolingErr, e.CarpoolingCreateErr, err)
	}
	return carpooling, nil
}

func GetCarpoolingList(db sqlExec, where map[string]interface{}, orderBy string, pageSize, pageNum int) ([]*CarpoolingDetailView, bool, error) {
	var list []*CarpoolingDetailView

	// 初始化缓冲区
	var modelBuf = make([]interface{}, 0)
	for i := 0; i < pageSize; i++ {
		modelBuf = append(modelBuf, &CarpoolingDetailView{})
	}

	// 每次只取pageSize个
	isOver, err := SQLQueryRows(db, &modelBuf, where, orderBy, pageSize, pageNum)
	if err != nil {
		return nil, true, e.SP(e.MCarpoolingErr, e.CarpoolingAlreadyExist, err)
	}
	for _, item := range modelBuf {
		if model, ok := item.(*CarpoolingDetailView); ok {
			list = append(list, model)
		}
	}
	return list, isOver, nil
}

func UpdateCarpoolingList(db sqlExec, wheres map[string]interface{}, valueWillSet *Carpooling) (int64, error) {
	count, err := SQLUpdate(db, valueWillSet, wheres)
	if err != nil {
		return 0, e.SP(e.MCarpoolingErr, e.CarpoolingGetErr, err)
	} else {
		return count, nil
	}
}

func UpdateCarpoolingById(db sqlExec, carpoolingId uint64, valueWillSet *Carpooling) (int64, error) {
	return UpdateCarpoolingList(db, map[string]interface{}{"id": carpoolingId}, valueWillSet)
}

func UpdateCarpoolingByIdOfSelf(db sqlExec, carpoolingId, userID uint64, valueWillSet *Carpooling) (int64, error) {
	return UpdateCarpoolingList(db, map[string]interface{}{"id": carpoolingId, "creater_uid": userID}, valueWillSet)
}

func DeleteCarpoolingList(db sqlExec, wheres map[string]interface{}) (int64, error) {
	count, err := SQLDelete(db, &Carpooling{}, wheres)
	if err != nil {
		return 0, e.SP(e.MCarpoolingErr, e.CarpoolingDeleteErr, err)
	} else {
		return count, nil
	}
}

func DeleteCarpoolingById(db sqlExec, id uint64) (int64, error) {
	return DeleteCarpoolingList(db, map[string]interface{}{"id": id})
}

func DeleteCarpoolingByIdOfSelf(db sqlExec, id, userId uint64) (int64, error) {
	return DeleteCarpoolingList(db, map[string]interface{}{"id": id, "creater_uid": userId})
}
