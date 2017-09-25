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

	items, err := GetColumnList(db.Mysql, map[string]interface{}{})
	assert.Nil(err)
	assert.NotNil(items)
	assert.Equal(10, len(items))
}

func TestUpdateColumnById(t *testing.T) {
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
	count, err := UpdateColumnById(db.Mysql, c.ID, &Column{Name: newName})
	assert.Nil(err)
	assert.NotZero(count)

	c, err = GetColumnByID(db.Mysql, c.ID)
	assert.Nil(err)
	assert.NotNil(c)
	assert.Equal(newName, c.Name)
}

func TestDeleteColumnById(t *testing.T) {
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

	count, err := DeleteColumnById(db.Mysql, c.ID)
	assert.Nil(err)
	assert.Equal(int64(1), count)

	c, err = GetColumnByID(db.Mysql, c.ID)
	assert.Nil(err)
	assert.Nil(c)
}
