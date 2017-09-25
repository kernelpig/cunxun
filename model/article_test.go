package model

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"wangqingang/cunxun/db"
	"wangqingang/cunxun/test"
)

func TestCreateArticle(t *testing.T) {
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

	c := &Article{
		ColumnId:   test.GenRandInt(5),
		Title:      test.GenRandString(),
		Content:    test.GenRandString(),
		CreaterUid: u.ID,
		UpdaterUid: u.ID,
	}

	c, err = CreateArticle(db.Mysql, c)
	assert.Nil(err)
	assert.NotNil(c)

	cd, err := GetArticleByID(db.Mysql, c.ID)
	assert.Nil(err)
	assert.NotNil(cd)
}

func TestGetArticleList(t *testing.T) {
	test.InitTestCaseEnv(t)
	assert := assert.New(t)

	var cs []*Article
	for i := 0; i < 10; i++ {
		c := &Article{
			ColumnId:   1,
			Title:      test.GenRandString(),
			Content:    test.GenRandString(),
			CreaterUid: test.GenRandInt(5),
			UpdaterUid: test.GenRandInt(5),
		}
		_, err := CreateArticle(db.Mysql, c)
		assert.Nil(err)
		cs = append(cs, c)
	}

	items, isOver, err := GetArticleList(db.Mysql, map[string]interface{}{}, OrderByIgnore, 20, 1)
	assert.Nil(err)
	assert.NotNil(items)
	assert.True(isOver)
	assert.Equal(10, len(items))
}

func TestUpdateArticleById(t *testing.T) {
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

	c := &Article{
		ColumnId:   test.GenRandInt(5),
		Title:      test.GenRandString(),
		Content:    test.GenRandString(),
		CreaterUid: u.ID,
		UpdaterUid: u.ID,
	}

	c, err = CreateArticle(db.Mysql, c)
	assert.Nil(err)
	assert.NotNil(c)

	cd, err := GetArticleByID(db.Mysql, c.ID)
	assert.Nil(err)
	assert.NotNil(cd)

	newTitle := test.GenRandString()
	count, err := UpdateArticleById(db.Mysql, c.ID, &Article{Title: newTitle})
	assert.Nil(err)
	assert.NotZero(count)

	cd, err = GetArticleByID(db.Mysql, c.ID)
	assert.Nil(err)
	assert.NotNil(cd)
	assert.Equal(newTitle, cd.Title)
}

func TestDeleteArticleById(t *testing.T) {
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

	c := &Article{
		ColumnId:   test.GenRandInt(5),
		Title:      test.GenRandString(),
		Content:    test.GenRandString(),
		CreaterUid: u.ID,
		UpdaterUid: u.ID,
	}

	c, err = CreateArticle(db.Mysql, c)
	assert.Nil(err)
	assert.NotNil(c)

	cd, err := GetArticleByID(db.Mysql, c.ID)
	assert.Nil(err)
	assert.NotNil(cd)

	count, err := DeleteArticleById(db.Mysql, c.ID)
	assert.Nil(err)
	assert.NotZero(count)

	cd, err = GetArticleByID(db.Mysql, c.ID)
	assert.Nil(err)
	assert.Nil(cd)
}
