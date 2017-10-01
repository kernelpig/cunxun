package model

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"wangqingang/cunxun/db"
	"wangqingang/cunxun/test"
)

func TestCreateComment(t *testing.T) {
	test.InitTestCaseEnv(t)
	assert := assert.New(t)

	c := &Comment{
		RelateId:   test.GenRandId(5),
		Content:    test.GenRandString(),
		CreaterUid: test.GenRandId(5),
	}

	c, err := CreateComment(db.Mysql, c)
	assert.Nil(err)
	assert.NotNil(c)

	c, err = GetCommentByID(db.Mysql, c.ID)
	assert.Nil(err)
	assert.NotNil(c)
}

func TestGetCommentList(t *testing.T) {
	test.InitTestCaseEnv(t)
	assert := assert.New(t)

	u := &User{
		Phone:          test.GenFakePhone(),
		NickName:       test.GenRandString(),
		HashedPassword: test.GenRandString(),
		PasswordLevel:  test.GenRandInt(5),
		RegisterSource: test.TestWebSource,
		Avatar:         test.GenRandString(),
	}

	u, err := CreateUser(db.Mysql, u)
	assert.Nil(err)
	assert.NotNil(u)

	for i := 0; i < 10; i++ {
		c := &Comment{
			RelateId:   1,
			Content:    test.GenRandString(),
			CreaterUid: u.ID,
		}
		_, err := CreateComment(db.Mysql, c)
		assert.Nil(err)
	}

	items, isOver, err := GetCommentList(db.Mysql, map[string]interface{}{}, 20, 1)
	assert.Nil(err)
	assert.NotNil(items)
	assert.True(isOver)
	assert.Equal(10, len(items))
}

func TestUpdateCommentById(t *testing.T) {
	test.InitTestCaseEnv(t)
	assert := assert.New(t)

	c := &Comment{
		RelateId:   test.GenRandId(5),
		Content:    test.GenRandString(),
		CreaterUid: test.GenRandId(5),
	}

	c, err := CreateComment(db.Mysql, c)
	assert.Nil(err)
	assert.NotNil(c)

	c, err = GetCommentByID(db.Mysql, c.ID)
	assert.Nil(err)
	assert.NotNil(c)

	newContent := test.GenRandString()
	count, err := UpdateCommentById(db.Mysql, c.ID, &Comment{Content: newContent})
	assert.Nil(err)
	assert.NotZero(count)

	c, err = GetCommentByID(db.Mysql, c.ID)
	assert.Nil(err)
	assert.NotNil(c)
	assert.Equal(newContent, c.Content)
}

func TestDeleteCommentById(t *testing.T) {
	test.InitTestCaseEnv(t)
	assert := assert.New(t)

	c := &Comment{
		RelateId:   test.GenRandId(5),
		Content:    test.GenRandString(),
		CreaterUid: test.GenRandId(5),
	}

	c, err := CreateComment(db.Mysql, c)
	assert.Nil(err)
	assert.NotNil(c)

	c, err = GetCommentByID(db.Mysql, c.ID)
	assert.Nil(err)
	assert.NotNil(c)

	count, err := DeleteCommentById(db.Mysql, c.ID)
	assert.Nil(err)
	assert.NotZero(count)

	c, err = GetCommentByID(db.Mysql, c.ID)
	assert.Nil(err)
	assert.Nil(c)
}
