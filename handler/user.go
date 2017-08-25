package handler

import (
	"net/http"

	"github.com/ahmetb/go-linq"
	"github.com/gin-gonic/gin"

	"wangqingang/cunxun/captcha"
	"wangqingang/cunxun/common"
	"wangqingang/cunxun/db"
	"wangqingang/cunxun/login"
	"wangqingang/cunxun/model"
	"wangqingang/cunxun/password"
	"wangqingang/cunxun/phone"
	"wangqingang/cunxun/token"
)

func UserSignupHandler(c *gin.Context) {
	var req UserSignupRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": common.AccountBindFailed,
		})
		return
	}

	if !linq.From(common.SourceRange).Contains(req.Source) {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": common.AccountInvalidSource,
		})
		return
	}

	if err := phone.ValidPhone(req.Phone); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": common.AccountInvalidPhone,
		})
		return
	}

	ok, err := CheckcodeVerify(c, req.Phone, req.Source, common.SignupPurpose, req.VerifyCode)
	if err != nil {
		return
	} else if !ok {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": common.AccountVerifyCodeNotMatch,
		})
		return
	}

	hashedPassword, err := password.Encrypt(req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": common.AccountInternalError,
		})
		return
	}

	passwordLevel := password.PasswordStrength(req.Password)
	if passwordLevel == password.LevelIllegal {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": common.AccountPasswordLevelIllegal,
		})
		return
	}

	user := &model.User{
		Phone:          req.Phone,
		HashedPassword: hashedPassword,
		PasswordLevel:  passwordLevel,
		RegisterSource: req.Source,
	}

	user, err = model.CreateUser(db.Mysql, user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    common.AccountDBError,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": common.OK,
	})
	return
}

func UserLoginHandler(c *gin.Context) {
	var req UserLoginRequest

	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": common.AccountBindFailed,
		})
		return
	}
	if !linq.From(common.SourceRange).Contains(req.Source) {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": common.AccountInvalidSource,
		})
		return
	}
	if err := phone.ValidPhone(req.Phone); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": common.AccountInvalidPhone,
		})
		return
	}
	user, err := model.GetUserByPhone(db.Mysql, req.Phone)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": common.AccountDBError,
		})
		return
	} else if user == nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": common.AccountAccountNotExist,
		})
		return
	}

	var loginKey = login.LoginKey{Phone: req.Phone, Purpose: common.SigninPurpose, Source: req.Source}
	login, err := loginKey.GetLogin()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": common.AccountDBError,
		})
		return
	} else if login == nil {
		login, err = loginKey.CreateLogin(common.Config.Login.TTL.D())
		if err != nil || login == nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"code": common.AccountDBError,
			})
			return
		}
	}

	// 登录次数检查
	needCaptcha := false
	if login.RequestTimes >= common.Config.Login.MaxCaptchaTImes {
		needCaptcha = true
	}
	if needCaptcha {
		if req.CaptchaId == "" || req.CaptchaValue == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":       common.AccountNeedCaptcha,
				"left_times": login.GetLeftTimes(),
			})
			return
		}

		if !captcha.VerifyCaptcha(req.CaptchaId, req.CaptchaValue) {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":       common.AccountCaptchaNotMatch,
				"left_times": login.GetLeftTimes(),
			})
			return
		}
	}
	if login.RequestTimes >= common.Config.Login.MaxRequestTimes {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": common.AccountRequestLimit,
		})
		return
	}
	login.RequestTimes++
	login.Save()

	if password.Verify(req.Password, user.HashedPassword) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": common.AccountInvalidPassword,
		})
		return
	}

	accessToken, err := token.TokenCreateAndStore(user.ID, req.Source)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": common.AccountInternalError,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":          common.OK,
		"account_token": accessToken,
	})
	return
}

func UserLogoutHandler(c *gin.Context) {

}
