package handler

import (
	"net/http"
	"testing"

	"github.com/gavv/httpexpect"

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

func testDebugGetCaptchaValue(t *testing.T, e *httpexpect.Expect) string {
	id := testCreateCaptcha(t, e)
	resp := e.GET("/debug/captcha/{captcha_id}").
		WithPath("captcha_id", id).
		Expect().Status(http.StatusOK).
		JSON().Object()
	resp.Value("code").Number().Equal(common.OK)
	return resp.Value("captcha_value").String().NotEmpty().Raw()
}

func testDebugGetCaptchaValueHandler(t *testing.T, e *httpexpect.Expect) {
	test.InitTestCaseEnv(t)
	testDebugGetCaptchaValue(t, e)
}
