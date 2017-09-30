package model

import (
	"time"

	e "wangqingang/cunxun/error"
)

const (
	CarpoolingDisable = 0
	CarpoolingEnable  = 1
)

type Carpooling struct {
	ID          int       `json:"id" column:"id"`
	CreatedAt   time.Time `json:"created_at" column:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" column:"updated_at"`
	CreaterUid  int       `json:"creater_uid" column:"creater_uid"`
	UpdaterUid  int       `json:"updater_uid" column:"updater_uid"`
	FromCity    string    `json:"from_city" column:"from_city"`
	ToCity      string    `json:"to_city" column:"to_city"`
	DepartTIme  time.Time `json:"depart_time" column:"depart_time"`
	PeopleCount int       `json:"people_count" column:"people_count"`
	Status      int       `json:"status" column:"status"`
	Remark      string    `json:"remark" column:"remark"`
}

type CarpoolingDetailView struct {
	ID          int       `json:"id" column:"id"`
	CreatedAt   time.Time `json:"created_at" column:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" column:"updated_at"`
	CreaterUid  int       `json:"creater_uid" column:"creater_uid"`
	UpdaterUid  int       `json:"updater_uid" column:"updater_uid"`
	FromCity    string    `json:"from_city" column:"from_city"`
	ToCity      string    `json:"to_city" column:"to_city"`
	DepartTIme  time.Time `json:"depart_time" column:"depart_time"`
	PeopleCount int       `json:"people_count" column:"people_count"`
	Status      int       `json:"status" column:"status"`
	Remark      string    `json:"remark" column:"remark"`
	Nickname    string    `json:"nickname" column:"nickname"`
}

func GetCarpoolingByID(db sqlExec, CarpoolingID int) (*CarpoolingDetailView, error) {
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

func CreateCarpooling(db sqlExec, Carpooling *Carpooling) (*Carpooling, error) {
	id, err := SQLInsert(db, Carpooling)
	if err != nil {
		if isDBDuplicateErr(err) {
			return nil, e.SP(e.MCarpoolingErr, e.CarpoolingAlreadyExist, err)
		}
		return nil, e.SP(e.MCarpoolingErr, e.CarpoolingCreateErr, err)
	}
	Carpooling.ID = int(id)
	return Carpooling, nil
}
