package handler

import (
	"net/http"

	linq "github.com/ahmetb/go-linq"

	"wangqingang/cunxun/common"
	"wangqingang/cunxun/db"
	"wangqingang/cunxun/model"
	"wangqingang/cunxun/sms"
	"wangqingang/cunxun/utils"
	"wangqingang/cunxun/utils/captcha"
	"wangqingang/cunxun/utils/render"
)

func SendVerifyCodeHandler(w http.ResponseWriter, r *http.Request) {
	var req SendVerifyCodeRequest
	if err := utils.BindJSON(r, &req); err != nil {
		render.BindError(w, r, err)
		return
	}

	if !linq.From(common.PurposeRange).Contains(req.Purpose) {
		render.InvalidPurpose(w, r)
		return
	}

	if !linq.From(common.SourceRange).Contains(req.Source) {
		render.InvalidSource(w, r)
		return
	}

	if err := utils.ValidPhone(req.Phone); err != nil {
		render.JSON(w, r, http.StatusBadRequest, map[string]interface{}{
			"code":    common.AccountInvalidPhone,
			"message": err.Error(),
		})
		return
	}

	// 图形验证码校验
	if !captcha.VerifyCaptcha(req.CaptchaToken, req.CaptchaValue) {
		render.CaptchaNotMatch(w, r)
		return
	}

	var verifyKey = model.VerifyKey{Purpose: req.Purpose, Source: req.Source, Phone: req.Phone}
	verify, err := verifyKey.GetVerify()
	if err != nil {
		render.InternelError(w, r, err)
		return
	} else if verify == nil {
		verifyCode := utils.RandomDigits(6)
		verify, err = verifyKey.CreateVerify(verifyCode)
		if err != nil || verify == nil {
			render.InternelError(w, r, err)
			return
		}
	} else {
		if verify.SendTimes >= common.Config.Verify.MaxSendTimes {
			render.RequestLimit(w, r)
			return
		} else if verify.CheckTimes >= common.Config.Verify.MaxCheckTimes {
			render.RequestLimit(w, r)
			return
		}
	}

	// 注册时判断用户
	if req.Purpose == common.SignupPurpose {
		account, err := model.GetAccountByPhone(db.Mysql, req.Phone)
		if err != nil {
			render.InternelError(w, r, err)
			return
		} else if account != nil {
			render.AccountAlreadyExist(w, r)
			return
		}
	} else if req.Purpose == common.ResetPasswordPurpose {
		account, err := model.GetAccountByPhone(db.Mysql, req.Phone)
		if err != nil {
			render.InternelError(w, r, err)
			return
		} else if account == nil {
			render.AccountNotExist(w, r)
			return
		}
	}

	// 发送短信校验码: 短信验证码在生存周期内，不管请求发送几次，都使用同一个验证码，产品需求!
	if common.Config.ReleaseMode {
		_, err = sms.SendVerifyCode(common.Config.Sms, GetMeiqiaTenantId(), req.Purpose, req.Phone, verify.VerifyCode)
		if err != nil {
			render.InternelError(w, r, err)
			return
		}
	}
	verify.SendTimes++
	verify.Save()

	if !common.Config.ReleaseMode {
		renderDebugWithVerifyCode(w, r, verify.VerifyCode)
		return

	}

	render.Success(w, r)
	return
}

// 注册验证处理
func CheckVerifyCodeHandler(w http.ResponseWriter, r *http.Request) {
	var req CheckVerifyCodeRequest
	if err := utils.BindJSON(r, &req); err != nil {
		render.BindError(w, r, err)
		return
	}
	// purpose校验
	if !linq.From(common.PurposeRange).Contains(req.Purpose) {
		render.InvalidPurpose(w, r)
		return
	}
	// source校验
	if !linq.From(common.SourceRange).Contains(req.Source) {
		render.InvalidSource(w, r)
		return
	}

	if err := utils.ValidPhone(req.Phone); err != nil {
		render.JSON(w, r, http.StatusBadRequest, map[string]interface{}{
			"code":    common.AccountInvalidPhone,
			"message": err.Error(),
		})
		return
	}

	// confirm获取
	var verifyKey = model.VerifyKey{Phone: req.Phone, Purpose: req.Purpose, Source: req.Source}
	verify, err := verifyKey.GetVerify()
	if err != nil {
		render.InternelError(w, r, err)
		return
	} else if verify == nil {
		render.InvalidVerifyCode(w, r)
		return
	}

	if verify.CheckTimes >= common.Config.Verify.MaxCheckTimes-1 { // -1 是为其他需要验证码的业务接口(如 signup/reset_password)预留一次
		render.RequestLimit(w, r)
		return
	}

	ok, err := verify.Verify(req.VerifyCode)
	if err != nil {
		render.InternelError(w, r, err)
		return
	} else if !ok {
		render.InvalidVerifyCode(w, r)
		return
	}

	render.Success(w, r)
	return
}

