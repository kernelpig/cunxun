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

	testDebugHandler(t, e)
	testInternelHandler(t, e)
	testExceptions(t, e)
	testBaseHandler(t, e)
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
	testUserGetAvatarHandler(t, e)
	testUserGetInfoHandler(t, e)
	testUserGetListHandler(t, e)
	testUserCreateHandler(t, e)
	testUserUpdateHandler(t, e)
	testUserDeleteByIdHandler(t, e)

	testColumnCreateHandler(t, e)
	testColumnGetListHandler(t, e)
	testColumnUpdateByIdHandler(t, e)
	testColumnDeleteByIdHandler(t, e)

	testArticleCreateHandler(t, e)
	testArticleGetListHandler(t, e)
	testArticleGetHandler(t, e)
	testArticleUpdateByIdHandler(t, e)
	testArticleDeleteByIdHandler(t, e)

	testCarpoolingCreateHandler(t, e)
	testCarpoolingGetListHandler(t, e)
	testCarpoolingGetHandler(t, e)
	testCarpoolingUpdateByIdHandler(t, e)
	testCarpoolingDeleteByIdHandler(t, e)

	testCommentCreateHandler(t, e)
	testCommentGetListHandler(t, e)
	testCommentGetHandler(t, e)
	testCommentDeleteByIdHandler(t, e)
	testCommentUpdateByIdHandler(t, e)
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
