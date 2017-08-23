package handler

import (
	"net/http"
	"testing"

	"github.com/gavv/httpexpect"

	"wangqingang/cunxun/common"
)

func testCreateCaptcha(t *testing.T, e *httpexpect.Expect) string {
	resp := e.POST("/captcha").
		Expect().Status(http.StatusOK).
		JSON().Object()

	resp.Value("code").Number().Equal(common.OK)
	return resp.Value("id").String().Raw()
}

func testGetCaptchaImage(t *testing.T, e *httpexpect.Expect, id string) {
	e.GET("/captcha/{id}").
		WithPath("id", id).
		Expect().Status(http.StatusOK)
}

func testGetCaptchaImageHandler(t *testing.T, e *httpexpect.Expect) {
	initTestCaseEnv(t)
	id := testCreateCaptcha(t, e)
	testGetCaptchaImage(t, e, id)
}

func testCreateCaptchaHandler(t *testing.T, e *httpexpect.Expect) {
	initTestCaseEnv(t)
	testCreateCaptcha(t, e)
}
