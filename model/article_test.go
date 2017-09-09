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

	c := &Article{
		ColumnId:   test.GenRandInt(5),
		Title:      test.GenRandString(),
		Content:    test.GenRandString(),
		CreaterUid: test.GenRandInt(5),
		UpdaterUid: test.GenRandInt(5),
	}

	c, err := CreateArticle(db.Mysql, c)
	assert.Nil(err)
	assert.NotNil(c)

	c, err = GetArticleByID(db.Mysql, c.ID)
	assert.Nil(err)
	assert.NotNil(c)
}
