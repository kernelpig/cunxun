package handler

import (
	"net/http"
	"testing"
	"encoding/json"

	"github.com/gavv/httpexpect"

	"github.com/stretchr/testify/assert"
	"wangqingang/cunxun/error"
	"wangqingang/cunxun/test"
)

func testCaptchaCreate(t *testing.T, e *httpexpect.Expect) string {
	assert := assert.New(t)
	resp := e.POST("/api/captcha").
		Expect().Status(http.StatusOK)

	object := &CreateResponse{}
	err := json.Unmarshal([]byte(resp.Body().Raw()), object)
	assert.Nil(err)
	assert.Equal(error.OK, object.Code)

	return object.Id
}

func testCaptchaGetImage(t *testing.T, e *httpexpect.Expect, id string) {
	e.GET("/api/captcha/{captcha_id}").
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
