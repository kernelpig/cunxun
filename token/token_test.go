package token

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"wangqingang/cunxun/test"
)

const (
	testAccountId       = 1234
	testTokenString     = "adcfbgfdfsfdsfsadfdfdsfda"
	testTokenTimeoutTTL = 1
)

var testTokenKey = &TokenKey{
	UserId: testAccountId,
	Source: test.TestWebSource,
}

func TestTokenCreate(t *testing.T) {
	test.InitTestCaseEnv(t)
	assert := assert.New(t)

	token, err := testTokenKey.CreateToken(testTokenString, time.Duration(testTokenTimeoutTTL)*time.Second)
	assert.Nil(err)
	assert.NotNil(token)

	token, err = testTokenKey.GetToken()
	assert.Nil(err)
	assert.NotNil(token)
	assert.Equal(testTokenString, token.Token)

	token.Clean()
	token, err = testTokenKey.GetToken()
	assert.Nil(err)
	assert.Nil(token)
}

func TestTokenTimeout(t *testing.T) {
	test.InitTestCaseEnv(t)
	assert := assert.New(t)

	token, err := testTokenKey.CreateToken(testTokenString, time.Duration(testTokenTimeoutTTL)*time.Second)
	assert.Nil(err)
	assert.NotNil(token)
	assert.Equal(testTokenString, token.Token)

	time.Sleep(time.Duration(2) * time.Second)
	token, err = testTokenKey.GetToken()
	assert.Nil(err)
	assert.Nil(token)
}

func TestTokenCreateAndStore(t *testing.T) {
	test.InitTestCaseEnv(t)
	assert := assert.New(t)

	token, err := TokenCreateAndStore(testAccountId, test.TestWebSource)
	assert.Nil(err)
	assert.NotEmpty(token)
}
