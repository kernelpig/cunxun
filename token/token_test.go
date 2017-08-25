package token

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"wangqingang/cunxun/common"
	"wangqingang/cunxun/test"
)

const (
	TestAccountId       = 1234
	TestTokenString     = "adcfbgfdfsfdsfsadfdfdsfda"
	TestTokenSource     = common.WebSource
	TestTokenTimeoutTTL = 1
)

var testTokenKey = &TokenKey{
	UserId: TestAccountId,
	Source: TestTokenSource,
}

func TestCreateToken(t *testing.T) {
	test.InitTestCaseEnv(t)
	assert := assert.New(t)

	token, err := testTokenKey.CreateToken(TestTokenString, common.Config.Token.AccessTokenTTL.D())
	assert.Nil(err)
	assert.NotNil(token)

	token, err = testTokenKey.GetToken()
	assert.Nil(err)
	assert.NotNil(token)
	assert.Equal(TestTokenString, token.Token)

	token.Clean()
	token, err = testTokenKey.GetToken()
	assert.Nil(err)
	assert.Nil(token)
}

func TestTokenTimeout(t *testing.T) {
	test.InitTestCaseEnv(t)
	assert := assert.New(t)

	token, err := testTokenKey.CreateToken(TestTokenString, time.Duration(TestTokenTimeoutTTL)*time.Second)
	assert.Nil(err)
	assert.NotNil(token)
	assert.Equal(TestTokenString, token.Token)

	time.Sleep(time.Duration(2) * time.Second)
	token, err = testTokenKey.GetToken()
	assert.Nil(err)
	assert.Nil(token)
}

func TestTokenCreateAndStore(t *testing.T) {
	test.InitTestCaseEnv(t)
	assert := assert.New(t)

	token, err := TokenCreateAndStore(TestAccountId, test.TestWebSource)
	assert.Nil(err)
	assert.NotEmpty(token)
}
