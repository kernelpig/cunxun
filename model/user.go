package model

import (
	e "wangqingang/cunxun/error"
)

// User 对应于数据库user表中的一行
type User struct {
	ID             int    `column:"id"`
	Phone          string `column:"phone"`
	NickName       string `column:"nickname"`
	HashedPassword string `column:"hashed_password"`
	PasswordLevel  int    `column:"password_level"`
	RegisterSource string `column:"register_source"`
	Avatar         string `column:"avatar"`
}

func GetUserByPhone(db sqlExec, phone string) (*User, error) {
	u := &User{}
	isFound, err := SQLQueryRow(db, u, map[string]interface{}{
		"phone": phone,
	})
	if err != nil || !isFound {
		return nil, e.SE(e.MUserErr, e.UserGetErr, err)
	}
	return u, nil
}

func GetUserByID(db sqlExec, userId int) (*User, error) {
	u := &User{}
	isFound, err := SQLQueryRow(db, u, map[string]interface{}{
		"id": userId,
	})
	if err != nil || !isFound {
		return nil, e.SE(e.MUserErr, e.UserGetErr, err)
	}
	return u, nil
}

func CreateUser(db sqlExec, user *User) (*User, error) {
	id, err := SQLInsert(db, user)
	if err != nil {
		return nil, e.SE(e.MUserErr, e.UserCreateErr, err)
	}
	user.ID = int(id)
	return user, nil
}

//// 根据id更新用户信息
//func UpdateUserById(db sqlExecutor, id string, fieldValues map[string]interface{}) (int64, error) {
//	_SQLTemp := "update user set %s where id = ?"
//
//	sqlArgs := make([]interface{}, 0)
//	sqlFields := make([]string, 0)
//	for key, value := range fieldValues {
//		sqlFields = append(sqlFields, fmt.Sprintf("`%s` = ?", key))
//		sqlArgs = append(sqlArgs, value)
//	}
//	sqlArgs = append(sqlArgs, id)
//	_SQL := fmt.Sprintf(_SQLTemp, strings.Join(sqlFields, ","))
//	result, err := db.Exec(_SQL, sqlArgs...)
//	if err != nil {
//		return 0, err
//	}
//
//	rowCnt, err := result.RowsAffected()
//	if err != nil {
//		return 0, err
//	}
//
//	return rowCnt, nil
//
//}
//
//// 根据id删除用户
//func DeleteUserById(db sqlExecutor, id string) error {
//	_SQL := "delete from user where id = ?"
//
//	_, err := db.Exec(_SQL, id)
//	if err != nil {
//		return err
//	}
//
//	return nil
//}
