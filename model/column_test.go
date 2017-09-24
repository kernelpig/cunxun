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

func TestGetColumnList(t *testing.T) {
	test.InitTestCaseEnv(t)
	assert := assert.New(t)

	var cs []*Column
	for i := 0; i < 10; i++ {
		c := &Column{
			Name:       test.GenRandString(),
			CreaterUid: test.GenRandInt(5),
		}
		_, err := CreateColumn(db.Mysql, c)
		assert.Nil(err)
		cs = append(cs, c)
	}

	items, err := GetColumnList(db.Mysql)
	assert.Nil(err)
	assert.NotNil(items)
	assert.Equal(10, len(items))
}

func TestUpdateColumnList(t *testing.T) {
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

	newName := test.GenRandString()
	count, err := UpdateColumnList(db.Mysql, map[string]interface{}{"id": c.ID}, &Column{Name: newName})
	assert.Nil(err)
	assert.NotZero(count)

	c, err = GetColumnByID(db.Mysql, c.ID)
	assert.Nil(err)
	assert.NotNil(c)
	assert.Equal(newName, c.Name)
}
