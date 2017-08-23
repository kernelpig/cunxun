package handler

import (
	"net/http"
	"testing"

	"github.com/gavv/httpexpect"
)

func testDebugPing(t *testing.T, e *httpexpect.Expect) {
	e.GET("/debug/ping").
		Expect().
		Status(http.StatusOK)
}

func testDebugPingHandler(t *testing.T, e *httpexpect.Expect) {
	initTestCaseEnv(t)
	testDebugPing(t, e)
}

func testDebugGetCaptchaImage(t *testing.T, e *httpexpect.Expect) {
	id := testCreateCaptcha(t, e)
	e.GET("/debug/captcha/{id}").
		WithPath("id", id).
		Expect().Status(http.StatusOK)
}
