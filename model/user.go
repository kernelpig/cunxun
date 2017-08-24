package model

import "database/sql"

type User struct {
}

func GetUserByPhone(db *sql.DB, phone string) (*User, error) {
	return nil, nil
}
