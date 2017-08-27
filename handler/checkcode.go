package handler

import (
	"errors"
	"net/http"

	"github.com/ahmetb/go-linq"
	"github.com/gin-gonic/gin"

	"wangqingang/cunxun/captcha"
	"wangqingang/cunxun/checkcode"
	"wangqingang/cunxun/common"
	"wangqingang/cunxun/db"
	"wangqingang/cunxun/model"
	"wangqingang/cunxun/phone"
	"wangqingang/cunxun/sms"
)

func CheckcodeSendHandler(c *gin.Context) {
	var req CheckcodeSendRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": common.AccountBindFailed,
		})
		return
	}

	if !linq.From(common.PurposeRange).Contains(req.Purpose) {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": common.AccountInvalidPurpose,
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

	// 图形验证码校验
	if !captcha.VerifyCaptcha(req.CaptchaId, req.CaptchaValue) {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": common.AccountCaptchaNotMatch,
		})
		return
	}

	var checkcodeKey = checkcode.CheckCodeKey{Purpose: req.Purpose, Source: req.Source, Phone: req.Phone}
	verify, err := checkcodeKey.GetCheckcode()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": common.AccountInternalError,
		})
		return
	} else if verify == nil {
		verify, err = checkcodeKey.CreateCheckCode(common.Config.Checkcode.TTL.D())
		if err != nil || verify == nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"code": common.AccountInternalError,
			})
			return
		}
	} else {
		if verify.SendTimes >= common.Config.Checkcode.MaxSendTimes {
			c.JSON(http.StatusBadRequest, gin.H{
				"code": common.AccountRequestLimit,
			})
			return
		} else if verify.CheckTimes >= common.Config.Checkcode.MaxCheckTimes {
			c.JSON(http.StatusBadRequest, gin.H{
				"code": common.AccountRequestLimit,
			})
			return
		}
	}

	// 注册时判断用户
	if req.Purpose == common.SignupPurpose {
		user, err := model.GetUserByPhone(db.Mysql, req.Phone)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"code": common.AccountInternalError,
			})
			return
		} else if user != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"code": common.UserAlreadyExist,
			})
			return
		}
	} else if req.Purpose == common.ResetPasswordPurpose {
		user, err := model.GetUserByPhone(db.Mysql, req.Phone)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"code": common.AccountInternalError,
			})
			return
		} else if user == nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"code": common.UserNotExist,
			})
			return
		}
	}

	// 发送短信校验码: 短信验证码在生存周期内，不管请求发送几次，都使用同一个验证码，产品需求!
	if common.Config.ReleaseMode {
		_, err = sms.SendCheckcode(common.Config.Sms, req.Purpose, req.Phone, verify.Code)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"code": common.AccountInternalError,
			})
			return
		}
	}
	verify.SendTimes++
	verify.Save()

	c.JSON(http.StatusOK, gin.H{
		"code": common.OK,
	})
	return
}

// 注册验证处理
func CheckcodeVerifyHandler(c *gin.Context) {
	var req CheckVerifyCodeRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": common.AccountBindFailed,
		})
		return
	}
	// purpose校验
	if !linq.From(common.PurposeRange).Contains(req.Purpose) {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": common.AccountInvalidPurpose,
		})
		return
	}
	// source校验
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

	// confirm获取
	var checkcodeKey = checkcode.CheckCodeKey{Phone: req.Phone, Purpose: req.Purpose, Source: req.Source}
	verify, err := checkcodeKey.GetCheckcode()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": common.AccountInternalError,
		})
		return
	} else if verify == nil {
		c.JSON(http.StatusBadRequest, gin.H{
			// todo: 确认此错误码值
			"code": common.AccountVerifyCodeNotMatch,
		})
		return
	}

	if verify.CheckTimes >= common.Config.Checkcode.MaxCheckTimes-1 { // -1 是为其他需要验证码的业务接口(如 signup/reset_password)预留一次
		c.JSON(http.StatusBadRequest, gin.H{
			"code": common.AccountRequestLimit,
		})
		return
	}

	ok, err := verify.Check(req.VerifyCode)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": common.AccountInternalError,
		})
		return
	} else if !ok {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": common.AccountVerifyCodeNotMatch,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": common.OK,
	})
	return
}

func CheckcodeVerify(c *gin.Context, phone, source, purpose, code string) (bool, error) {
	var key = checkcode.CheckCodeKey{
		Phone:   phone,
		Purpose: purpose,
		Source:  source,
	}
	verify, err := key.GetCheckcode()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": common.AccountInternalError,
		})
		return false, err
	} else if verify == nil {
		return false, nil
	} else {
		if verify.CheckTimes >= common.Config.Checkcode.MaxCheckTimes {
			c.JSON(http.StatusBadRequest, gin.H{
				"code": common.AccountRequestLimit,
			})
			return false, errors.New("verify code expired")
		}
	}

	ok, err := verify.Check(code)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": common.AccountInternalError,
		})
		return false, err
	}

	return ok, nil
}
