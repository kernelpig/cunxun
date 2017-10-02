package handler

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/gavv/httpexpect"
	"github.com/stretchr/testify/assert"

	"wangqingang/cunxun/checkcode"
	"wangqingang/cunxun/common"
	"wangqingang/cunxun/db"
	"wangqingang/cunxun/error"
	"wangqingang/cunxun/model"
	"wangqingang/cunxun/script"
	"wangqingang/cunxun/test"
)

func testUserSignup(t *testing.T, e *httpexpect.Expect, request *UserSignupRequest) string {
	assert := assert.New(t)
	resp := e.POST("/api/u/signup").WithJSON(request).
		Expect().Status(http.StatusOK)

	// TODO: update错误不应该降级
	object := &UserSignupResponse{}
	err := json.Unmarshal([]byte(resp.Body().Raw()), object)
	assert.Nil(err)

	return object.UserId
}

func testUserSignupHandler(t *testing.T, e *httpexpect.Expect) {
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
}

func testUserLogin(t *testing.T, e *httpexpect.Expect, request *UserLoginRequest) string {
	resp := e.POST("/api/u/login").
		WithJSON(request).
		Expect().Status(http.StatusOK)

	respObj := resp.JSON().Object()
	return respObj.Value("user_token").String().NotEmpty().Raw()
}

func testUserLoginHandler(t *testing.T, e *httpexpect.Expect) {
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
	testUserLogin(t, e, loginRequest)
}

func testUserLogout(t *testing.T, e *httpexpect.Expect, token string) {
	resp := e.POST("/api/u/logout").
		WithHeader(common.AuthHeaderKey, token).
		Expect().Status(http.StatusOK)
	resp.JSON().Object().
		Value("code").Number().Equal(error.OK)
}

func testUserGetAvatar(t *testing.T, e *httpexpect.Expect, userId int) {
	e.GET("/api/u/{user_id}/avatar").
		WithPath("user_id", userId).
		Expect().Status(http.StatusOK).ContentType("image/png")
}

func testUserGetAvatarHandler(t *testing.T, e *httpexpect.Expect) {
	test.InitTestCaseEnv(t)
	testUserGetAvatar(t, e, 0)
}

func testUserLogoutHandler(t *testing.T, e *httpexpect.Expect) {
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
	token := testUserLogin(t, e, loginRequest)
	testUserLogout(t, e, token)
}

func testUserSignupHandler_UserAlreadyExist(t *testing.T, e *httpexpect.Expect) {
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

	user, err := model.CreateUser(db.Mysql, &model.User{
		Phone:          sendRequest.Phone,
		NickName:       test.GenRandString(),
		HashedPassword: test.GenRandString(),
		PasswordLevel:  test.GenRandInt(5),
		RegisterSource: test.TestWebSource,
		Avatar:         test.GenRandString(),
	})
	assert.Nil(err)
	assert.NotNil(user)

	signupRequest := &UserSignupRequest{
		Phone:      sendRequest.Phone,
		NickName:   test.GenRandString(),
		Source:     sendRequest.Source,
		Password:   test.GenFakePassword(),
		VerifyCode: code,
	}
	resp := e.POST("/api/u/signup").WithJSON(signupRequest).
		Expect().Status(http.StatusBadRequest)

	userAlreadyExistCode := error.Code{
		ServiceIndex:   error.SCunxun,
		InterfaceIndex: error.IUserSignup,
		SubModuleIndex: error.MUserErr,
		SubErrorIndex:  error.UserAlreadyExist,
	}

	resp.JSON().Object().Value("code").Number().Equal(userAlreadyExistCode.C())
}

func testUserGetInfo(t *testing.T, e *httpexpect.Expect, xToken string, userId string) {
	resp := e.GET("/api/u/{user_id}").
		WithPath("user_id", userId).
		WithHeader(common.AuthHeaderKey, xToken).
		Expect().Status(http.StatusOK)

	resp.JSON().Object().
		Value("code").Number().Equal(error.OK)
}

func testUserGetInfoHandler(t *testing.T, e *httpexpect.Expect) {
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
	userId := testUserSignup(t, e, signupRequest)

	captchaId = testCaptchaCreate(t, e)
	captchaValue = testDebugGetCaptchaValue(t, e, captchaId)

	loginRequest := &UserLoginRequest{
		Phone:        sendRequest.Phone,
		Source:       sendRequest.Source,
		Password:     signupRequest.Password,
		CaptchaId:    captchaId,
		CaptchaValue: captchaValue,
	}
	token := testUserLogin(t, e, loginRequest)

	testUserGetInfo(t, e, token, userId)
}

func testUserGetList(t *testing.T, e *httpexpect.Expect, xToken string, pageNum, pageSize int) []*User {
	assert := assert.New(t)

	resp := e.GET("/api/u").
		WithQuery("order_by", model.OrderByCreateDate).
		WithQuery("page_num", pageNum).
		WithQuery("page_size", pageSize).
		WithHeader(common.AuthHeaderKey, xToken).
		Expect().Status(http.StatusOK)

	respObj := resp.JSON().Object()
	respObj.Value("code").Number().Equal(error.OK)

	var result struct {
		Code int     `json:"code"`
		List []*User `json:"list"`
	}
	err := json.Unmarshal([]byte(resp.Body().Raw()), &result)
	assert.Nil(err)

	return result.List
}

func testSuperAdminLogin(t *testing.T, e *httpexpect.Expect) string {
	assert := assert.New(t)
	_, err := script.CreateSuperAdmin()
	assert.Nil(err)

	captchaId := testCaptchaCreate(t, e)
	captchaValue := testDebugGetCaptchaValue(t, e, captchaId)

	loginRequest := &UserLoginRequest{
		Phone:        "86 " + common.Config.User.SuperAdminPhone,
		Source:       test.TestWebSource,
		Password:     common.Config.User.SuperAdminPassword,
		CaptchaId:    captchaId,
		CaptchaValue: captchaValue,
	}
	return testUserLogin(t, e, loginRequest)
}

func testUserGetListHandler(t *testing.T, e *httpexpect.Expect) {
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

	xSuperToken := testSuperAdminLogin(t, e)

	list := testUserGetList(t, e, xSuperToken, 1, 20)
	assert.Equal(2, len(list))
}

func testUserCreate(t *testing.T, e *httpexpect.Expect, xToken string, request *UserCreateRequest) string {
	assert := assert.New(t)
	resp := e.POST("/api/u/").
		WithHeader(common.AuthHeaderKey, xToken).
		WithJSON(request).Expect()

	object := &UserSignupResponse{}
	err := json.Unmarshal([]byte(resp.Body().Raw()), object)
	assert.Nil(err)
	assert.Equal(error.OK, object.Code)

	return object.UserId
}

func testUserCreateHandler(t *testing.T, e *httpexpect.Expect) {
	test.InitTestCaseEnv(t)

	xSuperToken := testSuperAdminLogin(t, e)

	createRequest := &UserCreateRequest{
		Phone:    test.GenFakePhone(),
		NickName: test.GenRandString(),
		Password: test.GenFakePassword(),
		Role:     model.UserRoleNormal,
	}
	testUserCreate(t, e, xSuperToken, createRequest)

	captchaId := testCaptchaCreate(t, e)
	captchaValue := testDebugGetCaptchaValue(t, e, captchaId)

	loginRequest := &UserLoginRequest{
		Phone:        createRequest.Phone,
		Source:       common.WebSource,
		Password:     createRequest.Password,
		CaptchaId:    captchaId,
		CaptchaValue: captchaValue,
	}
	testUserLogin(t, e, loginRequest)
}

func testUserUpdate(t *testing.T, e *httpexpect.Expect, xToken string, userId string, request *UserUpdateRequest) {
	resp := e.PUT("/api/u/{user_id}").
		WithPath("user_id", userId).
		WithHeader(common.AuthHeaderKey, xToken).
		WithJSON(request).Expect()

	respObj := resp.Status(http.StatusOK).JSON().Object()
	respObj.Value("code").Number().Equal(error.OK)
}

func testUserUpdateHandler(t *testing.T, e *httpexpect.Expect) {
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
	userId := testUserSignup(t, e, signupRequest)

	xSuperToken := testSuperAdminLogin(t, e)
	updateRequest := &UserUpdateRequest{
		Role: model.UserRoleNormal,
	}
	testUserUpdate(t, e, xSuperToken, userId, updateRequest)
}

func testUserDeleteById(t *testing.T, e *httpexpect.Expect, xToken string, userId string) {
	resp := e.DELETE("/api/u/{user_id}").
		WithPath("user_id", userId).
		WithHeader(common.AuthHeaderKey, xToken).
		Expect().Status(http.StatusOK)

	respObj := resp.JSON().Object()
	respObj.Value("code").Number().Equal(error.OK)
}

func testUserDeleteByIdHandler(t *testing.T, e *httpexpect.Expect) {
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
	userId := testUserSignup(t, e, signupRequest)

	captchaId = testCaptchaCreate(t, e)
	captchaValue = testDebugGetCaptchaValue(t, e, captchaId)

	loginRequest := &UserLoginRequest{
		Phone:        sendRequest.Phone,
		Source:       sendRequest.Source,
		Password:     signupRequest.Password,
		CaptchaId:    captchaId,
		CaptchaValue: captchaValue,
	}
	testUserLogin(t, e, loginRequest)

	xSuperToken := testSuperAdminLogin(t, e)
	testUserDeleteById(t, e, xSuperToken, userId)
}
