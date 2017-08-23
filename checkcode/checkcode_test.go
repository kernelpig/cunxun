package checkcode

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"time"
	"wangqingang/cunxun/common"
	"wangqingang/cunxun/test"
)

const (
	testWebSource     = common.WebSource
	testSignupPurpose = common.SignupPurpose
)

func TestCheckCodeKey_CreateCheckCode(t *testing.T) {
	test.InitTestCaseEnv(t)
	assert := assert.New(t)

	key := CheckCodeKey{
		Source:  testWebSource,
		Purpose: testSignupPurpose,
		Phone:   test.GenFakePhone(),
	}

	code := test.GenFakeCheckcode()
	checkcode, _ := key.CreateCheckCode(code)

	checkcode, _ = key.GetVerify()
	assert.NotNil(checkcode)

	isEqual, _ := checkcode.Check(code)
	assert.True(isEqual)

	checkcode.Clean()
	checkcode, _ = key.GetVerify()
	assert.Nil(checkcode)
}

func TestCheckCode_Timeout(t *testing.T) {
	test.InitTestCaseEnv(t)
	assert := assert.New(t)

	key := CheckCodeKey{
		Source:  testWebSource,
		Purpose: testSignupPurpose,
		Phone:   test.GenFakePhone(),
	}

	// 设置1秒测试老化
	originTTL := common.Config.Verify.TTL.Duration
	common.Config.Verify.TTL.Duration = time.Duration(1) * time.Second

	// 创建TTL为1秒
	code := test.GenFakeCheckcode()
	key.CreateCheckCode(code)

	// 恢复正常的TTL
	common.Config.Verify.TTL.Duration = originTTL

	time.Sleep(time.Duration(2) * time.Second)
	checkcode, _ := key.GetVerify()
	assert.Nil(checkcode)
}
