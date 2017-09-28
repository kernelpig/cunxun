package oss

import (
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/assert"

	"wangqingang/cunxun/test"
)

func TestPutObjectByBytes(t *testing.T) {
	test.InitTestCaseEnv(t)
	assert := assert.New(t)

	err := InitOss()
	assert.Nil(err)

	originName := "avatar.png"
	fileName := test.GenRandString() + path.Ext(originName)

	fd, err := os.Open(originName)
	defer fd.Close()
	assert.Nil(err)

	url, err := PutImageByFile(fileName, fd)
	assert.Nil(err)
	t.Log(url)
	assert.NotEmpty(url)
}
