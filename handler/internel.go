package handler


import (
	"net/http"

	"wangqingang/cunxun/common"
	"wangqingang/cunxun/db"
	"wangqingang/cunxun/middleware"
	"wangqingang/cunxun/model"
	"wangqingang/cunxun/sms"
	"wangqingang/cunxun/utils/render"
)

func VerifyTokenHandler(w http.ResponseWriter, r *http.Request) {
	accessToken := r.URL.Query().Get("token")
	account, payload, err := middleware.CheckAccessToken(accessToken)
	if err != nil {
		render.JSON(w, r, http.StatusOK, map[string]interface{}{
			"code":       common.OK,
			"result":     false,
			"account_id": "",
			"device":     "",
		})
		return
	}

	render.JSON(w, r, http.StatusOK, map[string]interface{}{
		"code":       common.OK,
		"result":     true,
		"account_id": account.ID,
		"device":     payload.LoginSource,
	})
	return
}

func GetPhoneHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")

	if id == "" {
		render.BindError(w, r, errors.New("id is empty"))
		return
	}

	account, err := model.GetAccountById(db.Mysql, id)
	if err != nil {
		render.InternelError(w, r, err)
		return
	} else if account == nil {
		render.AccountNotExist(w, r)
		return
	}

	render.JSON(w, r, http.StatusOK, map[string]interface{}{
		"code":    common.OK,
		"message": common.Errors[common.OK],
		"phone":   account.Phone,
	})
	return
}
