package utils

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/meiqia/chi/render"

	"wangqingang/cunxun/common"
)

func JSON(w http.ResponseWriter, r *http.Request, statusCode int, resp interface{}) {
	render.Status(r, statusCode)
	render.JSON(w, r, resp)
}

func PNG(w http.ResponseWriter, r *http.Request, statusCode int, v []byte) {
	render.Status(r, statusCode)
	w.Header().Set("Content-Type", "image/png")
	w.Write(v)
}

func WAV(w http.ResponseWriter, r *http.Request, statusCode int, v []byte) {
	render.Status(r, statusCode)
	w.Header().Set("Content-Type", "audio/wav")
	w.Write(v)
}

func InvalidToken(w http.ResponseWriter, r *http.Request) {
	JSON(w, r, http.StatusUnauthorized, map[string]interface{}{
		"code":    common.AccountInvalidToken,
		"message": common.Errors[common.AccountInvalidToken],
	})
}

func InvalidPassword(w http.ResponseWriter, r *http.Request, leftTimes int) {
	JSON(w, r, http.StatusBadRequest, map[string]interface{}{
		"code":       common.AccountInvalidPassword,
		"message":    fmt.Sprintf(common.Errors[common.AccountInvalidPassword], leftTimes),
		"left_times": leftTimes,
	})
}

func SuccessWithToken(w http.ResponseWriter, r *http.Request, token string) {
	JSON(w, r, http.StatusOK, map[string]interface{}{
		"code":    common.OK,
		"message": common.Errors[common.OK],
		"token":   token,
	})
}

func SuccessWithCaptchaToken(w http.ResponseWriter, r *http.Request, captchaToken string) {
	JSON(w, r, http.StatusOK, map[string]interface{}{
		"code":          common.OK,
		"message":       common.Errors[common.OK],
		"captcha_token": captchaToken,
	})
}

func CaptchaNotMatchWithInfo(w http.ResponseWriter, r *http.Request, leftTimes int) {
	JSON(w, r, http.StatusBadRequest, map[string]interface{}{
		"code":       common.AccountCaptchaNotMatch,
		"message":    common.Errors[common.AccountCaptchaNotMatch],
		"left_times": leftTimes,
	})
}

func NeedCaptchaWithInfo(w http.ResponseWriter, r *http.Request, leftTimes int) {
	JSON(w, r, http.StatusBadRequest, map[string]interface{}{
		"code":       common.AccountNeedCaptcha,
		"message":    common.Errors[common.AccountNeedCaptcha],
		"left_times": leftTimes,
	})
}

func AccountNotExist(w http.ResponseWriter, r *http.Request) {
	JSON(w, r, http.StatusBadRequest, map[string]interface{}{
		"code":    common.AccountAccountNotExist,
		"message": common.Errors[common.AccountAccountNotExist],
	})
}

func AccountAlreadyExist(w http.ResponseWriter, r *http.Request) {
	JSON(w, r, http.StatusBadRequest, map[string]interface{}{
		"code":    common.AccountAccountAlreadyExist,
		"message": common.Errors[common.AccountAccountAlreadyExist],
	})
}

func BindError(w http.ResponseWriter, r *http.Request, err error) {
	JSON(w, r, http.StatusBadRequest, map[string]interface{}{
		"code":    common.AccountBindFailed,
		"message": common.Errors[common.AccountBindFailed] + err.Error(),
	})
}

func InvalidPurpose(w http.ResponseWriter, r *http.Request) {
	JSON(w, r, http.StatusBadRequest, map[string]interface{}{
		"code":    common.AccountInvalidPurpose,
		"message": common.Errors[common.AccountInvalidPurpose],
	})
}

func InvalidSource(w http.ResponseWriter, r *http.Request) {
	JSON(w, r, http.StatusBadRequest, map[string]interface{}{
		"code":    common.AccountInvalidSource,
		"message": common.Errors[common.AccountInvalidSource],
	})
}

func InvalidVerifyCode(w http.ResponseWriter, r *http.Request) {
	JSON(w, r, http.StatusBadRequest, map[string]interface{}{
		"code":    common.AccountVerifyCodeNotMatch,
		"message": common.Errors[common.AccountVerifyCodeNotMatch],
	})
}

func CaptchaNotMatch(w http.ResponseWriter, r *http.Request) {
	JSON(w, r, http.StatusBadRequest, map[string]interface{}{
		"code":    common.AccountCaptchaNotMatch,
		"message": common.Errors[common.AccountCaptchaNotMatch],
	})
}

func InternelError(w http.ResponseWriter, r *http.Request, err error) {
	if err == nil {
		err = errors.New("未知错误")
	}
	JSON(w, r, http.StatusInternalServerError, map[string]interface{}{
		"code":    common.AccountInternalError,
		"message": common.Errors[common.AccountInternalError] + err.Error(),
	})
}

func PassSameWithOld(w http.ResponseWriter, r *http.Request) {
	JSON(w, r, http.StatusBadRequest, map[string]interface{}{
		"code":    common.AccountPasswordSameWithOld,
		"message": common.Errors[common.AccountPasswordSameWithOld],
	})
}

func Success(w http.ResponseWriter, r *http.Request) {
	JSON(w, r, http.StatusOK, map[string]interface{}{
		"code":    common.OK,
		"message": common.Errors[common.OK],
	})
}

func RequestLimit(w http.ResponseWriter, r *http.Request) {
	JSON(w, r, http.StatusBadRequest, map[string]interface{}{
		"code":    common.AccountRequestLimit,
		"message": common.Errors[common.AccountRequestLimit],
	})
}

func PasswordLevelIllegal(w http.ResponseWriter, r *http.Request) {
	JSON(w, r, http.StatusBadRequest, map[string]interface{}{
		"code":    common.AccountPasswordLevelIllegal,
		"message": common.Errors[common.AccountPasswordLevelIllegal],
	})
}
