package handler

import (
	"net/http"
	"strconv"

	"wangqingang/cunxun/model"
	"wangqingang/cunxun/common"
)

func CreateCaptchaHandler(w http.ResponseWriter, r *http.Request) {
	token := captcha.GenCaptcha(common.Config.Captcha.DefaultLength)
	render.SuccessWithCaptchaToken(w, r, token)
	return
}

func GetCaptchaImageHandler(w http.ResponseWriter, r *http.Request) {
	token := r.URL.Query().Get("captcha_token")

	width, err := strconv.ParseInt(r.URL.Query().Get("width"), 10, 64)
	if err != nil {
		width = int64(common.Config.Captcha.DefaultWidth)
	}

	height, err := strconv.ParseInt(r.URL.Query().Get("height"), 10, 64)
	if err != nil {
		height = int64(common.Config.Captcha.DefaultHeight)
	}

	data, err := captcha.GetCaptchaImage(token, int(width), int(height))
	if err != nil {
		render.JSON(w, r, http.StatusInternalServerError, map[string]interface{}{
			"code":    common.AccountInternalError,
			"message": "GetCaptchaImage " + err.Error(),
		})
		return
	}

	if !common.Config.ReleaseMode {
		render.JSON(w, r, http.StatusOK, map[string]interface{}{
			"code":          common.OK,
			"message":       common.Errors[common.OK],
			"captcha_value": GetCaptchaValue(token),
		})
		return
	}

	render.PNG(w, r, http.StatusOK, data)
	return
}
