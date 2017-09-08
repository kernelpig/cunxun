package handler

import (
	"net/http/httptest"
	"testing"

	"github.com/gavv/httpexpect"
)

// 主测试函数
func TestHandlers(t *testing.T) {
	server := httptest.NewServer(ServerEngine())
	defer server.Close()

	e := httpexpect.New(t, server.URL)

	testBaseHandler(t, e)
	testDebugHandler(t, e)
	testInternelHandler(t, e)
	testExceptions(t, e)
}

// 接口基础功能测试
func testBaseHandler(t *testing.T, e *httpexpect.Expect) {
	testCaptchaCreateHandler(t, e)
	testCaptchaGetImageHandler(t, e)

	testCheckcodeSendHandler(t, e)
	testCheckcodeVerifyHandler(t, e)

	testUserSignupHandler(t, e)
	testUserLoginHandler(t, e)
	testUserLogoutHandler(t, e)

	testColumnCreateHandler(t, e)
}

// debug接口测试
func testDebugHandler(t *testing.T, e *httpexpect.Expect) {
	testDebugPingHandler(t, e)
	testDebugCaptchaGetValueHandler(t, e)
	testDebugCheckcodeGetValueHandler(t, e)
}

// 内部接口测试
func testInternelHandler(t *testing.T, e *httpexpect.Expect) {

}

// 异常测试
func testExceptions(t *testing.T, e *httpexpect.Expect) {
	testUserSignupHandler_UserAlreadyExist(t, e)
}
