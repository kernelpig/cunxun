package middleware

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"wangqingang/cunxun/test"
	"wangqingang/cunxun/token"
)

const (
	testTimeoutTTL = 1
	testUserRole   = 2
)

func TestCheckAccessToken(t *testing.T) {
	test.InitTestCaseEnv(t)
	assert := assert.New(t)

	userId := test.GenRandInt(100)
	userToken, err := token.TokenCreateAndStore(userId, testUserRole, test.TestWebSource, time.Duration(1)*time.Minute)
	assert.Nil(err)
	assert.NotEmpty(userToken)

	payload, err := CheckAccessToken(userToken)
	assert.Nil(err)
	assert.NotNil(payload)
}
