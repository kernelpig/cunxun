package handler

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/gavv/httpexpect"
	"github.com/stretchr/testify/assert"

	"wangqingang/cunxun/checkcode"
	"wangqingang/cunxun/common"
	"wangqingang/cunxun/error"
	"wangqingang/cunxun/model"
	"wangqingang/cunxun/test"
)

func testCommentCreate(t *testing.T, e *httpexpect.Expect, xToken string, request *CommentCreateRequest) int {
	resp := e.POST("/api/comment/").
		WithHeader(common.AuthHeaderKey, xToken).
		WithJSON(request).
		Expect().Status(http.StatusOK)

	respObj := resp.JSON().Object()
	respObj.Value("code").Number().Equal(error.OK)
	return int(respObj.Value("comment_id").Number().Raw())
}

func testCommentCreateHandler(t *testing.T, e *httpexpect.Expect) {
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

	checkcodeKey := &checkcode.CheckCodeKey{
		Phone:   sendRequest.Phone,
		Purpose: sendRequest.Purpose,
		Source:  sendRequest.Source,
	}
	code := testDebugCheckcodeGetValue(t, e, checkcodeKey)

	signupRequest := &UserSignupRequest{
		Phone:      sendRequest.Phone,
		NickName:   test.GenRandString(),
		Source:     sendRequest.Source,
		Password:   test.GenFakePassword(),
		VerifyCode: code,
	}
	testUserSignup(t, e, signupRequest)

	captchaId = testCaptchaCreate(t, e)
	captchaValue = testDebugGetCaptchaValue(t, e, captchaId)

	loginRequest := &UserLoginRequest{
		Phone:        sendRequest.Phone,
		Source:       sendRequest.Source,
		Password:     signupRequest.Password,
		CaptchaId:    captchaId,
		CaptchaValue: captchaValue,
	}
	xToken := testUserLogin(t, e, loginRequest)
	xSuperToken := testSuperAdminLogin(t, e)

	createColumnRequest := &ColumnCreateRequest{
		Name: test.GenRandString(),
	}
	columnID := testColumnCreate(t, e, xSuperToken, createColumnRequest)

	createArticleRequest := &ArticleCreateRequest{
		ColumnId: columnID,
		Title:    test.GenRandString(),
		Content:  test.GenRandString() + test.GenRandString(),
	}
	articleID := testArticleCreate(t, e, xToken, createArticleRequest)

	createCommentRequest := &CommentCreateRequest{
		RelateId: articleID,
		Content:  test.GenRandString() + test.GenRandString(),
	}
	testCommentCreate(t, e, xToken, createCommentRequest)
}

func testCommentGetList(t *testing.T, e *httpexpect.Expect, relateID, pageNum, pageSize int) []*model.Comment {
	assert := assert.New(t)

	resp := e.GET("/api/comment/").
		WithQuery("creater_uid", 0).
		WithQuery("relate_id", relateID).
		WithQuery("page_num", pageNum).
		WithQuery("page_size", pageSize).
		Expect().Status(http.StatusOK)

	respObj := resp.JSON().Object()
	respObj.Value("code").Number().Equal(error.OK)

	var result struct {
		Code int              `json:"code"`
		List []*model.Comment `json:"list"`
	}
	err := json.Unmarshal([]byte(resp.Body().Raw()), &result)
	assert.Nil(err)
	return result.List
}

func testCommentGetListHandler(t *testing.T, e *httpexpect.Expect) {
	test.InitTestCaseEnv(t)
	assert := assert.New(t)

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

	checkcodeKey := &checkcode.CheckCodeKey{
		Phone:   sendRequest.Phone,
		Purpose: sendRequest.Purpose,
		Source:  sendRequest.Source,
	}
	code := testDebugCheckcodeGetValue(t, e, checkcodeKey)

	signupRequest := &UserSignupRequest{
		Phone:      sendRequest.Phone,
		NickName:   test.GenRandString(),
		Source:     sendRequest.Source,
		Password:   test.GenFakePassword(),
		VerifyCode: code,
	}
	testUserSignup(t, e, signupRequest)

	captchaId = testCaptchaCreate(t, e)
	captchaValue = testDebugGetCaptchaValue(t, e, captchaId)

	loginRequest := &UserLoginRequest{
		Phone:        sendRequest.Phone,
		Source:       sendRequest.Source,
		Password:     signupRequest.Password,
		CaptchaId:    captchaId,
		CaptchaValue: captchaValue,
	}
	xToken := testUserLogin(t, e, loginRequest)
	xSuperToken := testSuperAdminLogin(t, e)

	createColumnRequest := &ColumnCreateRequest{
		Name: test.GenRandString(),
	}
	columnID := testColumnCreate(t, e, xSuperToken, createColumnRequest)

	createArticleRequest := &ArticleCreateRequest{
		ColumnId: columnID,
		Title:    test.GenRandString(),
		Content:  test.GenRandString() + test.GenRandString(),
	}
	articleID := testArticleCreate(t, e, xToken, createArticleRequest)

	for i := 0; i < 5; i++ {
		createCommentRequest := &CommentCreateRequest{
			RelateId: articleID,
			Content:  test.GenRandString() + test.GenRandString(),
		}
		testCommentCreate(t, e, xToken, createCommentRequest)
	}

	list := testCommentGetList(t, e, columnID, 1, 10)
	assert.Equal(5, len(list))
}

func testCommentGet(t *testing.T, e *httpexpect.Expect, commentID int) *model.Comment {
	assert := assert.New(t)

	resp := e.GET("/api/comment/{comment_id}").
		WithPath("comment_id", commentID).
		Expect().Status(http.StatusOK)

	respObject := resp.JSON().Object()
	respObject.Value("code").Number().Equal(0)
	var result struct {
		Code int            `json:"code"`
		Item *model.Comment `json:"item"`
	}

	err := json.Unmarshal([]byte(resp.Body().Raw()), &result)
	assert.Nil(err)
	return result.Item
}

func testCommentGetHandler(t *testing.T, e *httpexpect.Expect) {
	test.InitTestCaseEnv(t)
	assert := assert.New(t)

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

	checkcodeKey := &checkcode.CheckCodeKey{
		Phone:   sendRequest.Phone,
		Purpose: sendRequest.Purpose,
		Source:  sendRequest.Source,
	}
	code := testDebugCheckcodeGetValue(t, e, checkcodeKey)

	signupRequest := &UserSignupRequest{
		Phone:      sendRequest.Phone,
		NickName:   test.GenRandString(),
		Source:     sendRequest.Source,
		Password:   test.GenFakePassword(),
		VerifyCode: code,
	}
	testUserSignup(t, e, signupRequest)

	captchaId = testCaptchaCreate(t, e)
	captchaValue = testDebugGetCaptchaValue(t, e, captchaId)

	loginRequest := &UserLoginRequest{
		Phone:        sendRequest.Phone,
		Source:       sendRequest.Source,
		Password:     signupRequest.Password,
		CaptchaId:    captchaId,
		CaptchaValue: captchaValue,
	}
	xToken := testUserLogin(t, e, loginRequest)
	xSuperToken := testSuperAdminLogin(t, e)

	createColumnRequest := &ColumnCreateRequest{
		Name: test.GenRandString(),
	}
	columnID := testColumnCreate(t, e, xSuperToken, createColumnRequest)

	createArticleRequest := &ArticleCreateRequest{
		ColumnId: columnID,
		Title:    test.GenRandString(),
		Content:  test.GenRandString() + test.GenRandString(),
	}
	articleID := testArticleCreate(t, e, xToken, createArticleRequest)

	createCommentRequest := &CommentCreateRequest{
		RelateId: articleID,
		Content:  test.GenRandString() + test.GenRandString(),
	}
	commentID := testCommentCreate(t, e, xToken, createCommentRequest)
	comment := testCommentGet(t, e, commentID)
	assert.Equal(commentID, comment.ID)
}

func testCommentUpdateById(t *testing.T, e *httpexpect.Expect, xToken string, commentId int, request *CommentCreateRequest) {
	resp := e.PUT("/api/comment/{comment_id}").
		WithPath("comment_id", commentId).
		WithHeader(common.AuthHeaderKey, xToken).
		WithJSON(request).
		Expect().Status(http.StatusOK)

	respObj := resp.JSON().Object()
	respObj.Value("code").Number().Equal(error.OK)
}

func testCommentUpdateByIdHandler(t *testing.T, e *httpexpect.Expect) {
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

	checkcodeKey := &checkcode.CheckCodeKey{
		Phone:   sendRequest.Phone,
		Purpose: sendRequest.Purpose,
		Source:  sendRequest.Source,
	}
	code := testDebugCheckcodeGetValue(t, e, checkcodeKey)

	signupRequest := &UserSignupRequest{
		Phone:      sendRequest.Phone,
		NickName:   test.GenRandString(),
		Source:     sendRequest.Source,
		Password:   test.GenFakePassword(),
		VerifyCode: code,
	}
	testUserSignup(t, e, signupRequest)

	captchaId = testCaptchaCreate(t, e)
	captchaValue = testDebugGetCaptchaValue(t, e, captchaId)

	loginRequest := &UserLoginRequest{
		Phone:        sendRequest.Phone,
		Source:       sendRequest.Source,
		Password:     signupRequest.Password,
		CaptchaId:    captchaId,
		CaptchaValue: captchaValue,
	}
	xToken := testUserLogin(t, e, loginRequest)
	xSuperToken := testSuperAdminLogin(t, e)

	createColumnRequest := &ColumnCreateRequest{
		Name: test.GenRandString(),
	}
	columnID := testColumnCreate(t, e, xSuperToken, createColumnRequest)

	createArticleRequest := &ArticleCreateRequest{
		ColumnId: columnID,
		Title:    test.GenRandString(),
		Content:  test.GenRandString() + test.GenRandString(),
	}
	articleID := testArticleCreate(t, e, xToken, createArticleRequest)

	createCommentRequest := &CommentCreateRequest{
		RelateId: articleID,
		Content:  test.GenRandString() + test.GenRandString(),
	}
	commentId := testCommentCreate(t, e, xToken, createCommentRequest)

	updateCommentRequest := &CommentCreateRequest{
		RelateId: articleID,
		Content:  test.GenRandString() + test.GenRandString(),
	}
	testCommentUpdateById(t, e, xToken, commentId, updateCommentRequest)
}

func testCommentDeleteById(t *testing.T, e *httpexpect.Expect, xToken string, commentId int) {
	resp := e.DELETE("/api/comment/{comment_id}").
		WithPath("comment_id", commentId).
		WithHeader(common.AuthHeaderKey, xToken).
		Expect().Status(http.StatusOK)

	respObj := resp.JSON().Object()
	respObj.Value("code").Number().Equal(error.OK)
}

func testCommentDeleteByIdHandler(t *testing.T, e *httpexpect.Expect) {
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

	checkcodeKey := &checkcode.CheckCodeKey{
		Phone:   sendRequest.Phone,
		Purpose: sendRequest.Purpose,
		Source:  sendRequest.Source,
	}
	code := testDebugCheckcodeGetValue(t, e, checkcodeKey)

	signupRequest := &UserSignupRequest{
		Phone:      sendRequest.Phone,
		NickName:   test.GenRandString(),
		Source:     sendRequest.Source,
		Password:   test.GenFakePassword(),
		VerifyCode: code,
	}
	testUserSignup(t, e, signupRequest)

	captchaId = testCaptchaCreate(t, e)
	captchaValue = testDebugGetCaptchaValue(t, e, captchaId)

	loginRequest := &UserLoginRequest{
		Phone:        sendRequest.Phone,
		Source:       sendRequest.Source,
		Password:     signupRequest.Password,
		CaptchaId:    captchaId,
		CaptchaValue: captchaValue,
	}
	xToken := testUserLogin(t, e, loginRequest)
	xSuperToken := testSuperAdminLogin(t, e)

	createColumnRequest := &ColumnCreateRequest{
		Name: test.GenRandString(),
	}
	columnID := testColumnCreate(t, e, xSuperToken, createColumnRequest)

	createArticleRequest := &ArticleCreateRequest{
		ColumnId: columnID,
		Title:    test.GenRandString(),
		Content:  test.GenRandString() + test.GenRandString(),
	}
	articleID := testArticleCreate(t, e, xToken, createArticleRequest)

	createCommentRequest := &CommentCreateRequest{
		RelateId: articleID,
		Content:  test.GenRandString() + test.GenRandString(),
	}
	commentId := testCommentCreate(t, e, xToken, createCommentRequest)

	testCommentDeleteById(t, e, xToken, commentId)
}
