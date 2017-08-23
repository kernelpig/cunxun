package handler


import (
	"net/http"

	"wangqingang/cunxun/common"
	"wangqingang/cunxun/utils/render"
)

func DebugPingHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("ok\n"))
	return
}

func DebugGetCaptchaImageHandler(w http.ResponseWriter, r *http.Request) {
	token := r.URL.Query().Get("captcha_token")

	render.JSON(w, r, http.StatusOK, map[string]interface{}{
		"code":          common.OK,
		"message":       common.Errors[common.OK],
		"captcha_value": GetCaptchaValue(token),
	})

	return
}
