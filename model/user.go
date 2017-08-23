package model

import (
	"database/sql"
	"fmt"
	"strings"
)

// User 对应于数据库user表中的一行
type User struct {
	ID             string
	Name           string
	Phone          string
	HashedPassword string
	RegisterSource string
	PasswordLevel  int
}

func GetUserByPhone(db sqlExecutor, phone string) (*User, error) {
	_SQL := "select id, phone, hashed_password, password_level, name, register_source from account where phone = ?"

	account := &User{}
	err := db.QueryRow(_SQL, phone).Scan(
		&account.ID,
		&account.Phone,
		&account.HashedPassword,
		&account.PasswordLevel,
		&account.Name,
		&account.RegisterSource,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return account, nil
}

func GetUserById(db sqlExecutor, id string) (*User, error) {
	_SQL := "select id, phone, hashed_password, password_level, name, register_source from account where id = ?"

	account := &User{}
	err := db.QueryRow(_SQL, id).Scan(
		&account.ID,
		&account.Phone,
		&account.HashedPassword,
		&account.PasswordLevel,
		&account.Name,
		&account.RegisterSource,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return account, nil
}

func CreateUser(db sqlExecutor, id string, account User) (string, error) {
	_SQL := "insert into account(id, phone, hashed_password, password_level, name, register_source) value(?,?,?,?,?,?)"
	_, err := db.Exec(_SQL, id, account.Phone, account.HashedPassword, account.PasswordLevel, account.Name, account.RegisterSource)
	if err != nil {
		return "", err
	}

	return id, nil
}

// 根据id更新用户信息
func UpdateUserById(db sqlExecutor, id string, fieldValues map[string]interface{}) (int64, error) {
	_SQLTemp := "update account set %s where id = ?"

	sqlArgs := make([]interface{}, 0)
	sqlFields := make([]string, 0)
	for key, value := range fieldValues {
		sqlFields = append(sqlFields, fmt.Sprintf("`%s` = ?", key))
		sqlArgs = append(sqlArgs, value)
	}
	sqlArgs = append(sqlArgs, id)
	_SQL := fmt.Sprintf(_SQLTemp, strings.Join(sqlFields, ","))
	result, err := db.Exec(_SQL, sqlArgs...)
	if err != nil {
		return 0, err
	}

	rowCnt, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}

	return rowCnt, nil

}

// 根据id删除用户
func DeleteUserById(db sqlExecutor, id string) error {
	_SQL := "delete from account where id = ?"

	_, err := db.Exec(_SQL, id)
	if err != nil {
		return err
	}

	return nil
}
