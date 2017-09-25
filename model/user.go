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
