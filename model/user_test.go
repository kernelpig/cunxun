package model

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"wangqingang/cunxun/db"
	"wangqingang/cunxun/test"
)

func TestCreateUser(t *testing.T) {
	test.InitTestCaseEnv(t)
	assert := assert.New(t)

	phone := test.GenFakePhone()
	u, err := GetUserByPhone(db.Mysql, phone)
	assert.Nil(err)
	assert.Nil(u)

	u = &User{
		Phone:          phone,
		NickName:       test.GenRandString(),
		HashedPassword: test.GenRandString(),
		PasswordLevel:  test.GenRandInt(5),
		RegisterSource: test.TestWebSource,
		Avatar:         test.GenRandString(),
	}

	u, err = CreateUser(db.Mysql, u)
	assert.Nil(err)
	assert.NotNil(u)

	u, err = GetUserByPhone(db.Mysql, phone)
	assert.Nil(err)
	assert.NotNil(u)

	u, err = GetUserByID(db.Mysql, u.ID)
	assert.Nil(err)
	assert.NotNil(u)
}
