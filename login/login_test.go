package login

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"wangqingang/cunxun/test"
)

const (
	testTimeoutTTL = 1
)

func TestLogin_CreateLogin(t *testing.T) {
	test.InitTestCaseEnv(t)
	assert := assert.New(t)

	key := LoginKey{
		Source:  test.TestWebSource,
		Purpose: test.TestSignupPurpose,
		Phone:   test.GenFakePhone(),
	}

	login, _ := key.CreateLogin(time.Duration(1) * time.Second)
	login, _ = key.GetLogin()
	assert.NotNil(login)

	login.Clean()
	login, _ = key.GetLogin()
	assert.Nil(login)
}

func TestLogin_Timeout(t *testing.T) {
	test.InitTestCaseEnv(t)
	assert := assert.New(t)

	key := LoginKey{
		Source:  test.TestWebSource,
		Purpose: test.TestSignupPurpose,
		Phone:   test.GenFakePhone(),
	}

	// 创建TTL为1秒
	key.CreateLogin(time.Duration(testTimeoutTTL) * time.Second)
	login, _ := key.GetLogin()
	assert.NotNil(login)

	time.Sleep(time.Duration(2) * time.Second)
	login, _ = key.GetLogin()
	assert.Nil(login)
}
