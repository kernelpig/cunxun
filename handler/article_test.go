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

func testArticleCreate(t *testing.T, e *httpexpect.Expect, xToken string, request *ArticleCreateRequest) int {
	resp := e.POST("/article/").
		WithHeader(common.AuthHeaderKey, xToken).
		WithJSON(request).
		Expect().Status(http.StatusOK)

	respObj := resp.JSON().Object()
	respObj.Value("code").Number().Equal(error.OK)
	return int(respObj.Value("article_id").Number().Raw())
}

func testArticleCreateHandler(t *testing.T, e *httpexpect.Expect) {
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

	createColumnRequest := &ColumnCreateRequest{
		Name: test.GenRandString(),
	}
	columnID := testColumnCreate(t, e, xToken, createColumnRequest)

	createArticleRequest := &ArticleCreateRequest{
		ColumnId: columnID,
		Title:    test.GenRandString(),
		Content:  test.GenRandString() + test.GenRandString(),
	}
	testArticleCreate(t, e, xToken, createArticleRequest)
}

func testArticleGetList(t *testing.T, e *httpexpect.Expect, columnID, pageNum, pageSize int) []*model.Article {
	assert := assert.New(t)

	resp := e.GET("/article/").
		WithQuery("column_id", columnID).
		WithQuery("order_by", "comment_count").
		WithQuery("page_num", pageNum).
		WithQuery("page_size", pageSize).
		Expect().Status(http.StatusOK)

	respObj := resp.JSON().Object()
	respObj.Value("code").Number().Equal(error.OK)

	var result struct {
		Code int              `json:"code"`
		List []*model.Article `json:"list"`
	}
	err := json.Unmarshal([]byte(resp.Body().Raw()), &result)
	assert.Nil(err)
	return result.List
}

func testArticleGetListHandler(t *testing.T, e *httpexpect.Expect) {
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

	createColumnRequest := &ColumnCreateRequest{
		Name: test.GenRandString(),
	}
	columnID := testColumnCreate(t, e, xToken, createColumnRequest)

	for i := 0; i < 5; i++ {
		createArticleRequest := &ArticleCreateRequest{
			ColumnId: columnID,
			Title:    test.GenRandString(),
			Content:  test.GenRandString() + test.GenRandString(),
		}
		testArticleCreate(t, e, xToken, createArticleRequest)
	}

	list := testArticleGetList(t, e, columnID, 1, 10)
	assert.Equal(5, len(list))
}

func testArticleGet(t *testing.T, e *httpexpect.Expect, articleID int) *model.Article {
	assert := assert.New(t)

	resp := e.GET("/article/{article_id}").
		WithPath("article_id", articleID).
		Expect().Status(http.StatusOK)

	respObject := resp.JSON().Object()
	respObject.Value("code").Number().Equal(0)
	var result struct {
		Code int            `json:"code"`
		Item *model.Article `json:"item"`
	}

	err := json.Unmarshal([]byte(resp.Body().Raw()), &result)
	assert.Nil(err)
	return result.Item
}

func testArticleGetHandler(t *testing.T, e *httpexpect.Expect) {
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

	createColumnRequest := &ColumnCreateRequest{
		Name: test.GenRandString(),
	}
	columnID := testColumnCreate(t, e, xToken, createColumnRequest)

	createArticleRequest := &ArticleCreateRequest{
		ColumnId: columnID,
		Title:    test.GenRandString(),
		Content:  test.GenRandString() + test.GenRandString(),
	}
	articleID := testArticleCreate(t, e, xToken, createArticleRequest)
	article := testArticleGet(t, e, articleID)
	assert.Equal(articleID, article.ID)
}
