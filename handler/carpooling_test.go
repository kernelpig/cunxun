package handler

import (
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"github.com/gavv/httpexpect"
	"github.com/stretchr/testify/assert"

	"wangqingang/cunxun/checkcode"
	"wangqingang/cunxun/common"
	"wangqingang/cunxun/error"
	"wangqingang/cunxun/model"
	"wangqingang/cunxun/test"
)

func testCarpoolingCreate(t *testing.T, e *httpexpect.Expect, xToken string, request *CarpoolingCreateRequest) int {
	resp := e.POST("/api/carpooling/").
		WithHeader(common.AuthHeaderKey, xToken).
		WithJSON(request).
		Expect().Status(http.StatusOK)

	respObj := resp.JSON().Object()
	respObj.Value("code").Number().Equal(error.OK)
	return int(respObj.Value("carpooling_id").Number().Raw())
}

func testCarpoolingCreateHandler(t *testing.T, e *httpexpect.Expect) {
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

	createCarpoolingRequest := &CarpoolingCreateRequest{
		FromCity:    test.GenRandString(),
		ToCity:      test.GenRandString(),
		DepartTIme:  time.Now(),
		PeopleCount: test.GenRandInt(5),
		Remark:      test.GenRandString() + test.GenRandString(),
	}
	testCarpoolingCreate(t, e, xToken, createCarpoolingRequest)
}

func testCarpoolingGetList(t *testing.T, e *httpexpect.Expect, pageNum, pageSize int) []*model.CarpoolingDetailView {
	assert := assert.New(t)

	resp := e.GET("/api/carpooling/").
		WithQuery("creater_uid", 0).
		WithQuery("order_by", model.OrderByCreateDate).
		WithQuery("page_num", pageNum).
		WithQuery("page_size", pageSize).
		Expect().Status(http.StatusOK)

	respObj := resp.JSON().Object()
	respObj.Value("code").Number().Equal(error.OK)

	var result struct {
		Code int                           `json:"code"`
		List []*model.CarpoolingDetailView `json:"list"`
	}
	err := json.Unmarshal([]byte(resp.Body().Raw()), &result)
	assert.Nil(err)
	return result.List
}

func testCarpoolingGetListHandler(t *testing.T, e *httpexpect.Expect) {
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

	for i := 0; i < 5; i++ {
		createCarpoolingRequest := &CarpoolingCreateRequest{
			FromCity:    test.GenRandString(),
			ToCity:      test.GenRandString(),
			DepartTIme:  time.Now(),
			PeopleCount: test.GenRandInt(5),
			Remark:      test.GenRandString() + test.GenRandString(),
		}
		testCarpoolingCreate(t, e, xToken, createCarpoolingRequest)
	}

	list := testCarpoolingGetList(t, e, 1, 10)
	assert.Equal(5, len(list))
}

func testCarpoolingUpdateById(t *testing.T, e *httpexpect.Expect, xToken string, carpoolingId int, request *CarpoolingUpdateRequest) {
	resp := e.PUT("/api/carpooling/{carpooling_id}").
		WithPath("carpooling_id", carpoolingId).
		WithHeader(common.AuthHeaderKey, xToken).
		WithJSON(request).
		Expect().Status(http.StatusOK)

	respObj := resp.JSON().Object()
	respObj.Value("code").Number().Equal(error.OK)
}

func testCarpoolingUpdateByIdHandler(t *testing.T, e *httpexpect.Expect) {
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

	createCarpoolingRequest := &CarpoolingCreateRequest{
		FromCity:    test.GenRandString(),
		ToCity:      test.GenRandString(),
		DepartTIme:  time.Now(),
		PeopleCount: test.GenRandInt(5),
		Remark:      test.GenRandString() + test.GenRandString(),
	}
	carpoolingId := testCarpoolingCreate(t, e, xToken, createCarpoolingRequest)

	updateCarpoolingRequest := &CarpoolingUpdateRequest{
		FromCity:    test.GenRandString(),
		ToCity:      test.GenRandString(),
		DepartTIme:  time.Now(),
		PeopleCount: test.GenRandInt(5),
		Remark:      test.GenRandString() + test.GenRandString(),
	}
	testCarpoolingUpdateById(t, e, xToken, carpoolingId, updateCarpoolingRequest)
}
