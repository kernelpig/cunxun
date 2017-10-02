package handler

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/gavv/httpexpect"

	"wangqingang/cunxun/checkcode"
	"wangqingang/cunxun/error"
	"wangqingang/cunxun/test"
)

func testCheckcodeSend(t *testing.T, e *httpexpect.Expect, request *CheckcodeSendRequest) {
	resp := e.POST("/api/checkcode/send").
		WithJSON(request).
		Expect().Status(http.StatusOK)

	fmt.Println(resp.Body().Raw())
	respObj := resp.JSON().Object()
	respObj.Value("code").Number().Equal(error.OK)
}

func testCheckcodeSendHandler(t *testing.T, e *httpexpect.Expect) {
	test.InitTestCaseEnv(t)

	captchaId := testCaptchaCreate(t, e)
	captchaValue := testDebugGetCaptchaValue(t, e, captchaId)

	request := &CheckcodeSendRequest{
		Phone:        test.GenFakePhone(),
		Purpose:      test.TestSignupPurpose,
		Source:       test.TestWebSource,
		CaptchaId:    captchaId,
		CaptchaValue: captchaValue,
	}

	testCheckcodeSend(t, e, request)
}

func testCheckcodeVerify(t *testing.T, e *httpexpect.Expect, request *CheckVerifyCodeRequest) {
	resp := e.POST("/api/checkcode/check").
		WithJSON(request).
		Expect().Status(http.StatusOK)

	fmt.Println(resp.Body().Raw())
	resp.JSON().Object().Value("code").Number().Equal(error.OK)
}

func testCheckcodeVerifyHandler(t *testing.T, e *httpexpect.Expect) {
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

	code := testDebugCheckcodeGetValue(t, e, &checkcode.CheckCodeKey{
		Phone:   sendRequest.Phone,
		Purpose: sendRequest.Purpose,
		Source:  sendRequest.Source,
	})

	checkRequest := &CheckVerifyCodeRequest{
		Phone:      sendRequest.Phone,
		Purpose:    sendRequest.Purpose,
		Source:     sendRequest.Source,
		VerifyCode: code,
	}
	testCheckcodeVerify(t, e, checkRequest)
}
