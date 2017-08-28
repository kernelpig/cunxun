package handler

import (
	"net/http"

	"github.com/ahmetb/go-linq"
	"github.com/gin-gonic/gin"

	"fmt"
	"wangqingang/cunxun/captcha"
	"wangqingang/cunxun/common"
	"wangqingang/cunxun/db"
	e "wangqingang/cunxun/error"
	"wangqingang/cunxun/login"
	"wangqingang/cunxun/middleware"
	"wangqingang/cunxun/model"
	"wangqingang/cunxun/password"
	"wangqingang/cunxun/phone"
	"wangqingang/cunxun/token"
)

func UserSignupHandler(c *gin.Context) {
	var req UserSignupRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, e.IP(e.IUserSignup, e.MParamsErr, e.ParamsBindErr, err))
		return
	}

	detail := fmt.Sprintf("%+v", req)
	if !linq.From(common.SourceRange).Contains(req.Source) {
		c.JSON(http.StatusBadRequest, e.ID(e.IUserSignup, e.MParamsErr, e.ParamsInvalidSource, detail))
		return
	}

	if _, err := phone.ValidPhone(req.Phone); err != nil {
		c.JSON(http.StatusBadRequest, e.ID(e.IUserSignup, e.MParamsErr, e.ParamsInvalidPhone, detail))
		return
	}

	ok, err := CheckcodeVerify(c, req.Phone, req.Source, common.SignupPurpose, req.VerifyCode)
	if err != nil {
		c.JSON(http.StatusInternalServerError, e.IP(e.IUserSignup, e.MCheckcodeErr, e.CheckcodeCheckErr, err))
		return
	} else if !ok {
		c.JSON(http.StatusBadRequest, e.I(e.IUserSignup, e.MCheckcodeErr, e.CheckcodeMismatch))
		return
	}

	hashedPassword, err := password.Encrypt(req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, e.IP(e.IUserSignup, e.MPasswordErr, e.PasswordEncryptErr, err))
		return
	}

	passwordLevel, err := password.PasswordStrength(req.Password)
	if err != nil {
		c.JSON(http.StatusBadRequest, e.IP(e.IUserSignup, e.MPasswordErr, e.PasswordLevelErr, err))
		return
	}

	user := &model.User{
		Phone:          req.Phone,
		NickName:       req.NickName,
		HashedPassword: hashedPassword,
		PasswordLevel:  passwordLevel,
		RegisterSource: req.Source,
	}

	user, err = model.CreateUser(db.Mysql, user)
	if err != nil {
		if msgErr, ok := err.(e.Message); ok && msgErr.Code.IsSubError(e.MUserErr, e.UserAlreadyExist) {
			c.JSON(http.StatusBadRequest, e.IP(e.IUserSignup, e.MUserErr, e.UserAlreadyExist, err))
		} else {
			c.JSON(http.StatusInternalServerError, e.IP(e.IUserSignup, e.MUserErr, e.UserCreateErr, err))
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": e.OK,
	})
	return
}

func UserLoginHandler(c *gin.Context) {
	var req UserLoginRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, e.IP(e.IUserLogin, e.MParamsErr, e.ParamsBindErr, err))
		return
	}

	detail := fmt.Sprintf("%+v", req)
	if !linq.From(common.SourceRange).Contains(req.Source) {
		c.JSON(http.StatusBadRequest, e.ID(e.IUserLogin, e.MParamsErr, e.ParamsInvalidSource, detail))
		return
	}
	if _, err := phone.ValidPhone(req.Phone); err != nil {
		c.JSON(http.StatusBadRequest, e.ID(e.IUserLogin, e.MParamsErr, e.ParamsInvalidPhone, detail))
		return
	}

	user, err := model.GetUserByPhone(db.Mysql, req.Phone)
	if err != nil {
		c.JSON(http.StatusInternalServerError, e.IP(e.IUserLogin, e.MUserErr, e.UserGetErr, err))
		return
	} else if user == nil {
		c.JSON(http.StatusBadRequest, e.I(e.IUserLogin, e.MUserErr, e.UserNotExist))
		return
	}

	var loginKey = login.LoginKey{Phone: req.Phone, Purpose: common.SigninPurpose, Source: req.Source}
	login, err := loginKey.GetLogin()
	if err != nil {
		c.JSON(http.StatusInternalServerError, e.IP(e.IUserLogin, e.MLoginErr, e.LoginGetErr, err))
		return
	} else if login == nil {
		login, err = loginKey.CreateLogin(common.Config.Login.TTL.D())
		if err != nil || login == nil {
			c.JSON(http.StatusInternalServerError, e.IP(e.IUserLogin, e.MLoginErr, e.LoginCreateErr, err))
			return
		}
	}

	// 登录次数检查
	needCaptcha := false
	if login.RequestTimes >= common.Config.Login.MaxCaptchaTImes {
		needCaptcha = true
	}
	if needCaptcha {
		leftTimes := fmt.Sprintf("%d", login.GetLeftTimes())
		if req.CaptchaId == "" || req.CaptchaValue == "" {
			c.JSON(http.StatusBadRequest, e.ID(e.IUserLogin, e.MCaptchaErr, e.CaptchaRequired, leftTimes))
			return
		}

		if !captcha.VerifyCaptcha(req.CaptchaId, req.CaptchaValue) {
			c.JSON(http.StatusBadRequest, e.ID(e.IUserLogin, e.MCaptchaErr, e.CaptchaMismatch, leftTimes))
			return
		}
	}
	if login.RequestTimes >= common.Config.Login.MaxRequestTimes {
		c.JSON(http.StatusBadRequest, e.I(e.IUserLogin, e.MLoginErr, e.LogDumpRequestErr))
		return
	}
	login.RequestTimes++
	login.Save()

	if err := password.Verify(req.Password, user.HashedPassword); err != nil {
		c.JSON(http.StatusBadRequest, e.IP(e.IUserLogin, e.MPasswordErr, e.PasswordInvalid, err))
		return
	}

	accessToken, err := token.TokenCreateAndStore(user.ID, req.Source, common.Config.Token.AccessTokenTTL.D())
	if err != nil {
		c.JSON(http.StatusInternalServerError, e.IP(e.IUserLogin, e.MTokenErr, e.TokenCreateErr, err))
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":       e.OK,
		"user_token": accessToken,
	})
	return
}

func UserLogoutHandler(c *gin.Context) {
	userToken := c.GetHeader(common.AuthHeaderKey)
	payload, err := middleware.CheckAccessToken(userToken)

	if err == nil && payload != nil {
		token.TokenClean(int(payload.UserId), payload.LoginSource)
	}

	c.JSON(http.StatusOK, gin.H{
		"code": e.OK,
	})
	return
}
