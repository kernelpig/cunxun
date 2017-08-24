package handler

import (
	"github.com/gavv/httpexpect"
	"net/http"
	"testing"
	"wangqingang/cunxun/checkcode"
	"wangqingang/cunxun/common"
	"wangqingang/cunxun/test"
)

func testUserSignup(t *testing.T, e *httpexpect.Expect, request *UserSignupRequest) {
	resp := e.POST("/u/signup").WithJSON(request).
		Expect().Status(http.StatusOK)

	respObj := resp.JSON().Object()
	respObj.Value("code").Number().Equal(common.OK)
}

func testUserSignupHandler(t *testing.T, e *httpexpect.Expect) {
	test.InitTestCaseEnv(t)

	captchaId := testCreateCaptcha(t, e)
	captchaValue := testDebugGetCaptchaValue(t, e, captchaId)

	sendRequest := &CheckcodeSendRequest{
		Phone:        test.GenFakePhone(),
		Purpose:      test.TestSignupPurpose,
		Source:       test.TestWebSource,
		CaptchaId:    captchaId,
		CaptchaValue: captchaValue,
	}
	testCheckcodeSend(t, e, sendRequest)

	checkcodeKey := &checkcode.CheckCodeKey{
		Phone:   sendRequest.Phone,
		Purpose: sendRequest.Purpose,
		Source:  sendRequest.Source,
	}
	code := testDebugCheckcodeGetValue(t, e, checkcodeKey)

	signinRequest := &UserSignupRequest{
		Phone:      sendRequest.Phone,
		Source:     sendRequest.Source,
		Password:   test.GenFakePassword(),
		VerifyCode: code,
	}
	testUserSignup(t, e, signinRequest)
}
