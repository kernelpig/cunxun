package handler

import (
	"net/http"

	linq "github.com/ahmetb/go-linq"

	"wangqingang/cunxun/common"
	"wangqingang/cunxun/db"
	"wangqingang/cunxun/middleware"
	"wangqingang/cunxun/model"
	"wangqingang/cunxun/sms"
	"wangqingang/cunxun/utils"
	"wangqingang/cunxun/utils/captcha"
	"wangqingang/cunxun/utils/render"
)

func renderDebugWithVerifyCode(w http.ResponseWriter, r *http.Request, verifyCode string) {
	render.JSON(w, r, http.StatusOK, map[string]interface{}{
		"code":        common.OK,
		"message":     common.Errors[common.OK],
		"verify_code": verifyCode,
	})
}

// 获取当前用户
func getAuthContext(r *http.Request) middleware.AuthContext {
	return r.Context().Value(common.CurrentAccount).(middleware.AuthContext)
}
