package token

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"wangqingang/cunxun/test"
)

const (
	testUserId          = 1234
	testUserRole        = 2
	testTokenString     = "adcfbgfdfsfdsfsadfdfdsfda"
	testTokenTimeoutTTL = 1
)

var testTokenKey = &TokenKey{
	UserId: testUserId,
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
	assert.NotNil(err)
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
	assert.NotNil(err)
	assert.Nil(token)
}

func TestTokenCreateAndStore(t *testing.T) {
	test.InitTestCaseEnv(t)
	assert := assert.New(t)

	token, err := TokenCreateAndStore(testUserId, testUserRole, test.TestWebSource, time.Duration(1)*time.Minute)
	assert.Nil(err)
	assert.NotEmpty(token)
	TokenClean(testUserId, test.TestWebSource)
}
