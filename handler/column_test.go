package handler

import (
	"net/http"
	"testing"

	"github.com/gavv/httpexpect"

	"wangqingang/cunxun/checkcode"
	"wangqingang/cunxun/common"
	"wangqingang/cunxun/error"
	"wangqingang/cunxun/test"
)

func testColumnCreate(t *testing.T, e *httpexpect.Expect, xToken string, request *ColumnCreateRequest) {
	resp := e.POST("/column/").
		WithHeader(common.AuthHeaderKey, xToken).
		WithJSON(request).
		Expect().Status(http.StatusOK)

	respObj := resp.JSON().Object()
	respObj.Value("code").Number().Equal(error.OK)
}

func testColumnCreateHandler(t *testing.T, e *httpexpect.Expect) {
	test.InitTestCaseEnv(t)

	captchaId := testCaptchaCreate(t, e)
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

	signupRequest := &UserSignupRequest{
		Phone:      sendRequest.Phone,
		NickName:   test.GenRandString(),
		Source:     sendRequest.Source,
		Password:   test.GenFakePassword(),
		VerifyCode: code,
	}
	testUserSignup(t, e, signupRequest)

	captchaId = testCaptchaCreate(t, e)
	captchaValue = testDebugGetCaptchaValue(t, e, captchaId)

	loginRequest := &UserLoginRequest{
		Phone:        sendRequest.Phone,
		Source:       sendRequest.Source,
		Password:     signupRequest.Password,
		CaptchaId:    captchaId,
		CaptchaValue: captchaValue,
	}
	xToken := testUserLogin(t, e, loginRequest)

	createRequest := &ColumnCreateRequest{
		Name: test.GenRandString(),
	}
	testColumnCreate(t, e, xToken, createRequest)
}
