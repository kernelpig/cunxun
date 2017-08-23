package handler

import (
	"net/http/httptest"
	"testing"

	"github.com/gavv/httpexpect"

	"wangqingang/cunxun/test"
)

func init() {
	test.TestInit()
}

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
	testCreateCaptchaHandler(t, e)
	testGetCaptchaImageHandler(t, e)
}

// debug接口测试
func testDebugHandler(t *testing.T, e *httpexpect.Expect) {
	testDebugPingHandler(t, e)
	testDebugGetCaptchaValueHandler(t, e)
}

// 内部接口测试
func testInternelHandler(t *testing.T, e *httpexpect.Expect) {

}

// 异常测试
func testExceptions(t *testing.T, e *httpexpect.Expect) {

}
