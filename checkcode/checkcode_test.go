package checkcode

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"wangqingang/cunxun/common"
	"wangqingang/cunxun/test"
)

func TestCheckCodeKey_CreateCheckCode(t *testing.T) {
	test.InitTestCaseEnv(t)
	assert := assert.New(t)

	key := CheckCodeKey{
		Source:  test.TestWebSource,
		Purpose: test.TestSignupPurpose,
		Phone:   test.GenFakePhone(),
	}

	checkcode, _ := key.CreateCheckCode()

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

	// 设置1秒测试老化
	originTTL := common.Config.Checkcode.TTL.Duration
	common.Config.Checkcode.TTL.Duration = time.Duration(1) * time.Second

	// 创建TTL为1秒
	key.CreateCheckCode()

	// 恢复正常的TTL
	common.Config.Checkcode.TTL.Duration = originTTL

	time.Sleep(time.Duration(2) * time.Second)
	checkcode, _ := key.GetCheckcode()
	assert.Nil(checkcode)
}
