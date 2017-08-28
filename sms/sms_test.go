package sms

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"wangqingang/cunxun/common"
)

const (
	testAliAccessId     = "YOUR_ALIYUN_ACCESS_ID"
	testAliAccessSecret = "YOUR_ALIYUN_ACCESS_SECRET"
	testPurpse          = common.SignupPurpose
	testPhone           = "18653193004"
	testCheckcode       = "123456"
)

func TestSendCheckcode(t *testing.T) {
	assert := assert.New(t)
	config := &common.SmsConfig{
		AliAccessId:     testAliAccessId,
		AliAccessSecret: testAliAccessSecret,
	}
	response, err := SendCheckcode(config, testPhone, testPurpse, testCheckcode)
	assert.Nil(err)
	assert.NotEmpty(response)
	t.Log(response)
}
