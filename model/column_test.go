package model

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"wangqingang/cunxun/db"
	"wangqingang/cunxun/test"
)

func TestCreateColumn(t *testing.T) {
	test.InitTestCaseEnv(t)
	assert := assert.New(t)

	c := &Column{
		Name:       test.GenRandString(),
		CreaterUid: test.GenRandInt(5),
	}

	c, err := CreateColumn(db.Mysql, c)
	assert.Nil(err)
	assert.NotNil(c)

	c, err = GetColumnByID(db.Mysql, c.ID)
	assert.Nil(err)
	assert.NotNil(c)
}
