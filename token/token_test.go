package token

import (
	"testing"
	"time"

	"wangqingang/cunxun/common"
	"wangqingang/cunxun/db"
	"github.com/stretchr/testify/assert"
)

const (
	TestAccountId       = "1234"
	TestTokenString     = "adcfbgfdfsfdsfsadfdfdsfda"
	TestTokenSource     = common.WebSource
	TestTokenTimeoutTTL = 1
)

var testTokenKey = &TokenKey{
	AccountId: TestAccountId,
	Source:    TestTokenSource,
}

func init() {
	common.InitConfig("../conf/config.dev.toml")
	db.InitRedis(common.Config.Redis)
}

func TestCreateToken(t *testing.T) {
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
