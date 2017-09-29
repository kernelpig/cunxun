package script

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"wangqingang/cunxun/common"
	"wangqingang/cunxun/db"
	"wangqingang/cunxun/model"
	"wangqingang/cunxun/test"
)

func TestCreateSuperAdmin(t *testing.T) {
	assert := assert.New(t)
	test.InitTestCaseEnv(t)

	user, err := CreateSuperAdmin()
	assert.Nil(err)
	assert.NotNil(user)

	user, err = model.GetUserByPhone(db.Mysql, "86 "+common.Config.User.SuperAdminPhone)
	assert.Nil(err)
	assert.NotNil(user)
}

func TestCreateColumns(t *testing.T) {
	assert := assert.New(t)
	test.InitTestCaseEnv(t)

	user, err := CreateSuperAdmin()
	assert.Nil(err)
	assert.NotNil(user)

	columns, err := CreateColumns(user)
	assert.Nil(err)
	assert.NotNil(columns)
	assert.NotZero(len(columns))
}
