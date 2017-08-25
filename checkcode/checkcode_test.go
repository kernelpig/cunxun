package checkcode

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"wangqingang/cunxun/test"
)

const (
	testTimeoutTTL = 1
)

func TestCheckCodeKey_CreateCheckCode(t *testing.T) {
	test.InitTestCaseEnv(t)
	assert := assert.New(t)

	key := CheckCodeKey{
		Source:  test.TestWebSource,
		Purpose: test.TestSignupPurpose,
		Phone:   test.GenFakePhone(),
	}

	checkcode, _ := key.CreateCheckCode(time.Duration(testTimeoutTTL) * time.Second)

	checkcode, _ = key.GetCheckcode()
	assert.NotNil(checkcode)

	isEqual, _ := checkcode.Check("invalid code")
	assert.False(isEqual)

	isEqual, _ = checkcode.Check(checkcode.Code)
	assert.True(isEqual)

	checkcode.Clean()
	checkcode, _ = key.GetCheckcode()
	assert.Nil(checkcode)
}

func TestCheckCode_Timeout(t *testing.T) {
	test.InitTestCaseEnv(t)
	assert := assert.New(t)

	key := CheckCodeKey{
		Source:  test.TestWebSource,
		Purpose: test.TestSignupPurpose,
		Phone:   test.GenFakePhone(),
	}

	// 创建TTL为1秒
	key.CreateCheckCode(time.Duration(testTimeoutTTL) * time.Second)

	time.Sleep(time.Duration(2) * time.Second)
	checkcode, _ := key.GetCheckcode()
	assert.Nil(checkcode)
}
