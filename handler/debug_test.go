package handler

import (
	"net/http"
	"testing"

	"github.com/gavv/httpexpect"

	"wangqingang/cunxun/checkcode"
	"wangqingang/cunxun/common"
	"wangqingang/cunxun/test"
)

func testDebugPing(t *testing.T, e *httpexpect.Expect) {
	e.GET("/debug/ping").
		Expect().
		Status(http.StatusOK)
}

func testDebugPingHandler(t *testing.T, e *httpexpect.Expect) {
	test.InitTestCaseEnv(t)
	testDebugPing(t, e)
}

func testDebugGetCaptchaValue(t *testing.T, e *httpexpect.Expect, captchaId string) string {
	resp := e.GET("/debug/captcha/{captcha_id}").
		WithPath("captcha_id", captchaId).
		Expect().Status(http.StatusOK).
		JSON().Object()

	resp.Value("code").Number().Equal(common.OK)

	return resp.Value("captcha_value").String().NotEmpty().Raw()
}

func testDebugCaptchaGetValueHandler(t *testing.T, e *httpexpect.Expect) {
	test.InitTestCaseEnv(t)

	captchaID := testCaptchaCreate(t, e)

	testDebugGetCaptchaValue(t, e, captchaID)
}

func testDebugCheckcodeGetValue(t *testing.T, e *httpexpect.Expect, key *checkcode.CheckCodeKey) string {
	resp := e.GET("/debug/checkcode/").
		WithQuery("phone", key.Phone).
		WithQuery("purpose", key.Purpose).
		WithQuery("source", key.Source).
		Expect().Status(http.StatusOK).
		JSON().Object()

	resp.Value("code").Number().Equal(common.OK)

	return resp.Value("checkcode").String().NotEmpty().Raw()
}

func testDebugCheckcodeGetValueHandler(t *testing.T, e *httpexpect.Expect) {
	test.InitTestCaseEnv(t)

	captchaID := testCaptchaCreate(t, e)
	captchaValue := testDebugGetCaptchaValue(t, e, captchaID)

	sendRequest := &CheckcodeSendRequest{
		Phone:        test.GenFakePhone(),
		Purpose:      test.TestSignupPurpose,
		Source:       test.TestWebSource,
		CaptchaId:    captchaID,
		CaptchaValue: captchaValue,
	}
	testCheckcodeSend(t, e, sendRequest)
}
