package handler

import (
	"net/http"
	"testing"

	"github.com/gavv/httpexpect"

	"wangqingang/cunxun/common"
	"wangqingang/cunxun/test"
)

func testCaptchaCreate(t *testing.T, e *httpexpect.Expect) string {
	resp := e.POST("/captcha").
		Expect().Status(http.StatusOK).
		JSON().Object()

	resp.Value("code").Number().Equal(common.OK)
	return resp.Value("captcha_id").String().Raw()
}

func testCaptchaGetImage(t *testing.T, e *httpexpect.Expect, id string) {
	e.GET("/captcha/{captcha_id}").
		WithPath("captcha_id", id).
		Expect().Status(http.StatusOK)
}

func testCaptchaGetImageHandler(t *testing.T, e *httpexpect.Expect) {
	test.InitTestCaseEnv(t)
	id := testCaptchaCreate(t, e)
	testCaptchaGetImage(t, e, id)
}

func testCaptchaCreateHandler(t *testing.T, e *httpexpect.Expect) {
	test.InitTestCaseEnv(t)
	testCaptchaCreate(t, e)
}
