package handler

import (
	"errors"
	"net/http"
	"time"

	linq "github.com/ahmetb/go-linq"

	"wangqingang/cunxun/common"
	"wangqingang/cunxun/db"
	"wangqingang/cunxun/middleware"
	"wangqingang/cunxun/model"
	"wangqingang/cunxun/tenant"
	"wangqingang/cunxun/token"
	"wangqingang/cunxun/token/token_lib"
	"wangqingang/cunxun/utils"
	"wangqingang/cunxun/utils/captcha"
	"wangqingang/cunxun/utils/password"
	"wangqingang/cunxun/utils/render"
)

func checkVerifyCode(w http.ResponseWriter, r *http.Request, phone, source, purpose, code string) (bool, error) {
	var verifyKey = model.VerifyKey{
		Phone:   phone,
		Purpose: purpose,
		Source:  source,
	}
	verify, err := verifyKey.GetVerify()
	if err != nil {
		render.InternelError(w, r, err)
		return false, err
	} else if verify == nil {
		return false, nil
	} else {
		if verify.CheckTimes >= common.Config.Verify.MaxCheckTimes {
			render.RequestLimit(w, r)
			return false, errors.New("verify code expired")
		}
	}

	ok, err := verify.Verify(code)
	if err != nil {
		render.InternelError(w, r, err)
		return false, err
	}

	return ok, nil
}

func SignUpHandler(w http.ResponseWriter, r *http.Request) {
	log := middleware.GetAccessLog(r)
	log.Request = ""

	var req SignUpRequest
	if err := utils.BindJSON(r, &req); err != nil {
		render.BindError(w, r, err)
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

	ok, err := checkVerifyCode(w, r, req.Phone, req.Source, common.SignupPurpose, req.VerifyCode)
	if err != nil {
		return
	} else if !ok {
		render.InvalidVerifyCode(w, r)
		return
	}

	hashedPassword, err := password.Encrypt(req.Password)
	if err != nil {
		render.InternelError(w, r, err)
		return
	}

	passwordLevel := password.PasswordStrength(req.Password)
	if passwordLevel == password.LevelIllegal {
		render.PasswordLevelIllegal(w, r)
		return
	}

	account := model.Account{
		Phone:          req.Phone,
		RegisterSource: req.Source,
		HashedPassword: hashedPassword,
		PasswordLevel:  passwordLevel,
	}

	accountID, err := utils.GenID()
	if err != nil {
		render.JSON(w, r, http.StatusInternalServerError, map[string]interface{}{
			"code":    common.AccountGenerateIdFailed,
			"message": "GenID: " + err.Error(),
		})
		return
	}
	_, err = model.CreateAccount(db.Mysql, accountID, account)
	if err != nil {
		render.JSON(w, r, http.StatusInternalServerError, map[string]interface{}{
			"code":    common.AccountDBError,
			"message": "CreateAccount: " + err.Error(),
		})
		return
	}

	render.Success(w, r)
	return
}

func ResetPasswordHandler(w http.ResponseWriter, r *http.Request) {
	log := middleware.GetAccessLog(r)
	log.Request = ""

	var req ResetPasswordRequest
	if err := utils.BindJSON(r, &req); err != nil {
		render.BindError(w, r, err)
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

	ok, err := checkVerifyCode(w, r, req.Phone, req.Source, common.ResetPasswordPurpose, req.VerifyCode)
	if err != nil {
		return
	} else if !ok {
		render.InvalidVerifyCode(w, r)
		return
	}

	account, err := model.GetAccountByPhone(db.Mysql, req.Phone)
	if err != nil {
		render.JSON(w, r, http.StatusInternalServerError, map[string]interface{}{
			"code":    common.AccountDBError,
			"message": "GetAccountByPhone: " + err.Error(),
		})
		return
	} else if account == nil {
		render.AccountNotExist(w, r)
		return
	}

	hashedPassword, err := password.Encrypt(req.NewPassword)
	if err != nil {
		render.InternelError(w, r, err)
		return
	}
	if password.Verify(req.NewPassword, account.HashedPassword) == nil {
		render.PassSameWithOld(w, r)
		return
	}
	passwordLevel := password.PasswordStrength(req.NewPassword)
	if passwordLevel == password.LevelIllegal {
		render.PasswordLevelIllegal(w, r)
		return
	}

	_, err = model.UpdateAccountById(db.Mysql, account.ID, map[string]interface{}{
		"password_level":  passwordLevel,
		"hashed_password": hashedPassword,
	})
	if err != nil {
		render.JSON(w, r, http.StatusInternalServerError, map[string]interface{}{
			"code":    common.AccountDBError,
			"message": "UpdateAccountById: " + err.Error(),
		})
		return
	}

	// 清除本帐户所有的account token, 及access token/refresh token
	token.RemoveAllTokenOfAccount(account.ID)
	tenant.RemoveAllTokenOfAccount(account.ID)

	render.Success(w, r)
	return
}

func SignInHandler(w http.ResponseWriter, r *http.Request) {
	log := middleware.GetAccessLog(r)
	log.Request = ""

	var req SignInRequest
	if err := utils.BindJSON(r, &req); err != nil {
		render.BindError(w, r, err)
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

	account, err := model.GetAccountByPhone(db.Mysql, req.Phone)
	if err != nil {
		render.JSON(w, r, http.StatusInternalServerError, map[string]interface{}{
			"code":    common.AccountDBError,
			"message": "GetAccountByPhone: " + err.Error(),
		})
		return
	} else if account == nil {
		render.AccountNotExist(w, r)
		return
	}

	var loginKey = model.LoginKey{Phone: req.Phone, Purpose: common.SigninPurpose, Source: req.Source}
	login, err := loginKey.GetLogin()
	if err != nil {
		render.JSON(w, r, http.StatusInternalServerError, map[string]interface{}{
			"code":    common.AccountDBError,
			"message": "GetLogin: " + err.Error(),
		})
		return
	} else if login == nil {
		login, err = loginKey.CreateLogin(common.Config.Login.TTL.D())
		if err != nil || login == nil {
			render.JSON(w, r, http.StatusInternalServerError, map[string]interface{}{
				"code":    common.AccountDBError,
				"message": "CreateLogin: " + err.Error(),
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
		if req.CaptchaToken == "" || req.CaptchaValue == "" {
			render.NeedCaptchaWithInfo(w, r, login.GetLeftTimes())
			return
		}

		if !captcha.VerifyCaptcha(req.CaptchaToken, req.CaptchaValue) {
			render.CaptchaNotMatchWithInfo(w, r, login.GetLeftTimes())
			return
		}
	}
	if login.RequestTimes >= common.Config.Login.MaxRequestTimes {
		render.RequestLimit(w, r)
		return
	}
	login.RequestTimes++
	login.Save()

	if password.Verify(req.Password, account.HashedPassword) != nil {
		render.InvalidPassword(w, r, login.GetLeftTimes())
		return
	}

	accessToken, err := genTokenAndStore(account.ID, req.Source)
	if err != nil {
		render.InternelError(w, r, err)
		return
	}

	render.SuccessWithToken(w, r, accessToken)
	return
}

func SignOutHandler(w http.ResponseWriter, r *http.Request) {
	authContext := getAuthContext(r)

	if err := cleanAccountToken(&authContext); err != nil {
		render.InternelError(w, r, err)
		return
	}

	render.Success(w, r)
	return
}

func genTokenAndStore(accountId, source string) (string, error) {

	issueTime := time.Now()

	// payload中ttl单位为分钟
	accessToken, err := token_lib.Encrypt(common.Config.Token.TokenLibVersion, &token_lib.Payload{
		IssueTime:   uint32(uint64(issueTime.Unix())),
		TTL:         uint16(common.Config.Token.AccessTokenTTL.Minutes()),
		AccountId:   accountId,
		LoginSource: source,
	})
	if err != nil {
		return "", err
	}

	tokenKey := token.TokenKey{
		AccountId: accountId,
		Source:    source,
	}
	_, err = tokenKey.CreateToken(accessToken, common.Config.Token.AccessTokenTTL.D())
	if err != nil {
		return "", err
	}

	return accessToken, nil
}

func cleanAccountToken(authContext *middleware.AuthContext) error {
	tokenKey := token.TokenKey{
		AccountId: authContext.CurrentAccount.ID,
		Source:    authContext.Payload.LoginSource,
	}

	token, err := tokenKey.GetToken()
	if err != nil {
		return err
	} else if token != nil {
		token.Clean()
	}

	return nil
}
