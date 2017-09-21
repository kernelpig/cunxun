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
		ArticleId:  test.GenRandInt(5),
		Content:    test.GenRandString(),
		CreaterUid: test.GenRandInt(5),
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

	var cs []*Comment
	for i := 0; i < 10; i++ {
		c := &Comment{
			ArticleId:  1,
			Content:    test.GenRandString(),
			CreaterUid: u.ID,
		}
		_, err := CreateComment(db.Mysql, c)
		assert.Nil(err)
		cs = append(cs, c)
	}

	items, isOver, err := GetCommentList(db.Mysql, map[string]interface{}{}, 20, 1)
	assert.Nil(err)
	assert.NotNil(items)
	assert.True(isOver)
	assert.Equal(10, len(items))
}