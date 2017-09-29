package script

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"wangqingang/cunxun/test"
)

const (
	testPubKeyPath = "./pub_key"
	testPriKeyPath = "./pri_key"
)

func TestCreateTokenKey(t *testing.T) {
	test.InitTestCaseEnv(t)
	assert := assert.New(t)

	err := CreateTokenKey(testPubKeyPath, testPriKeyPath)
	assert.Nil(err)
	t.Logf("%v", err)
}
