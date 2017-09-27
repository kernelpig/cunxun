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

func TestGetUserList(t *testing.T) {
	test.InitTestCaseEnv(t)
	assert := assert.New(t)

	var cs []*User
	for i := 0; i < 10; i++ {
		c := &User{
			Phone:          test.GenFakePhone(),
			NickName:       test.GenRandString(),
			HashedPassword: test.GenRandString(),
			PasswordLevel:  test.GenRandInt(5),
			RegisterSource: test.TestWebSource,
			Avatar:         test.GenRandString(),
		}
		_, err := CreateUser(db.Mysql, c)
		assert.Nil(err)
		cs = append(cs, c)
	}

	items, isOver, err := GetUserList(db.Mysql, map[string]interface{}{}, OrderByIgnore, 20, 1)
	assert.Nil(err)
	assert.NotNil(items)
	assert.True(isOver)
	assert.Equal(10, len(items))
}

func TestUpdateUserById(t *testing.T) {
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

	count, err := UpdateUserById(db.Mysql, u.ID, &User{Role: UserRoleAdmin})
	assert.Nil(err)
	assert.NotZero(count)

	uNew, err := GetUserByID(db.Mysql, u.ID)
	assert.Nil(err)
	assert.NotNil(uNew)
	assert.Equal(UserRoleAdmin, uNew.Role)
	assert.Equal(u.ID, uNew.ID)
}

func TestDeleteUserById(t *testing.T) {
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

	count, err := DeleteUserById(db.Mysql, u.ID)
	assert.Nil(err)
	assert.NotZero(count)

	u, err = GetUserByID(db.Mysql, u.ID)
	assert.Nil(err)
	assert.Nil(u)
}
