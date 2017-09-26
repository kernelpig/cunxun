package main

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"wangqingang/cunxun/common"
	"wangqingang/cunxun/db"
	"wangqingang/cunxun/model"
	"wangqingang/cunxun/test"
)

func TestUserCreateSuperAdmin(t *testing.T) {
	assert := assert.New(t)
	test.InitTestCaseEnv(t)

	err := UserCreateSuperAdmin(common.Config.User)
	assert.Nil(err)

	user, err := model.GetUserByPhone(db.Mysql, "86 "+common.Config.User.SuperAdminPhone)
	assert.Nil(err)
	assert.NotNil(user)
}
