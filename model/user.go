package model

import (
	"time"

	e "wangqingang/cunxun/error"
)

const (
	UserRoleNormal     = 0
	UserRoleAdmin      = 1
	UserRoleSuperAdmin = 2
)

// User 对应于数据库user表中的一行
type User struct {
	ID             int       `json:"id" column:"id"`
	Phone          string    `json:"phone" column:"phone"`
	NickName       string    `json:"nickname" column:"nickname"`
	HashedPassword string    `json:"hashed_password" column:"hashed_password"`
	PasswordLevel  int       `json:"password_level" column:"password_level"`
	RegisterSource string    `json:"register_source" column:"register_source"`
	Avatar         string    `json:"avatar" column:"avatar"`
	CreatedAt      time.Time `json:"created_at" column:"created_at"`
	UpdatedAt      time.Time `json:"updated_at" column:"updated_at"`
	Role           int       `json:"role" column:"role"`
}

func GetUserByPhone(db sqlExec, phone string) (*User, error) {
	u := &User{}
	isFound, err := SQLQueryRow(db, u, map[string]interface{}{
		"phone": phone,
	})
	if err != nil {
		return nil, e.SP(e.MUserErr, e.UserGetErr, err)
	} else if !isFound {
		return nil, nil
	} else {
		return u, nil
	}
}

func GetUserByID(db sqlExec, userId int) (*User, error) {
	u := &User{}
	isFound, err := SQLQueryRow(db, u, map[string]interface{}{
		"id": userId,
	})
	if err != nil {
		return nil, e.SP(e.MUserErr, e.UserGetErr, err)
	} else if !isFound {
		return nil, nil
	} else {
		return u, nil
	}
}

func CreateUser(db sqlExec, user *User) (*User, error) {
	id, err := SQLInsert(db, user)
	if err != nil {
		if isDBDuplicateErr(err) {
			return nil, e.SP(e.MUserErr, e.UserAlreadyExist, err)
		}
		return nil, e.SP(e.MUserErr, e.UserCreateErr, err)
	}
	user.ID = int(id)
	return user, nil
}

func GetUserList(db sqlExec, where map[string]interface{}) ([]*User, error) {
	var list []*User

	// 初始化缓冲区
	var modelBuf = make([]interface{}, 0)
	for i := 0; i < pageSize; i++ {
		modelBuf = append(modelBuf, &User{})
	}

	// 每次只取pageSize个
	for pageNum := 1; true; pageNum++ {
		isOver, err := SQLQueryRows(db, &modelBuf, where, OrderById, pageSize, pageNum)
		if err != nil {
			return nil, e.SP(e.MUserErr, e.UserGetListErr, err)
		}
		for _, item := range modelBuf {
			if model, ok := item.(*User); ok {
				list = append(list, model)
			}
		}
		if isOver {
			break
		}
	}
	return list, nil
}
