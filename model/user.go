package model

// User 对应于数据库user表中的一行
type User struct {
	ID             int    `column:"id"`
	Phone          string `column:"phone"`
	Name           string `column:"name"`
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
		return nil, err
	}
	return u, nil
}

//
//func GetUserById(db sqlExecutor, id string) (*User, error) {
//	_SQL := "select id, phone, hashed_password, password_level, name, register_source from user where id = ?"
//
//	user := &User{}
//	err := db.QueryRow(_SQL, id).Scan(
//		&user.ID,
//		&user.Phone,
//		&user.HashedPassword,
//		&user.PasswordLevel,
//		&user.Name,
//		&user.RegisterSource,
//	)
//	if err == sql.ErrNoRows {
//		return nil, nil
//	}
//	if err != nil {
//		return nil, err
//	}
//
//	return user, nil
//}

func CreateUser(db sqlExec, user *User) (*User, error) {
	id, err := SQLInsert(db, user)
	if err != nil {
		return nil, err
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