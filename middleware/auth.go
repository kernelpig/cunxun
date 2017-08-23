package middleware

import (
	"context"
	"errors"
	"net/http"
	"time"

	"wangqingang/cunxun/common"
	"wangqingang/cunxun/db"
	"wangqingang/cunxun/model"
	"wangqingang/cunxun/token"
	"wangqingang/cunxun/token/token_lib"
	"wangqingang/cunxun/utils/render"
)

func AuthRequired(next http.Handler) http.Handler {
	return authMiddleware(next)
}

type AuthContext struct {
	CurrentAccount *model.Account
	Payload        *token_lib.Payload
}

func CheckAccessToken(authToken string) (*model.Account, *token_lib.Payload, error) {
	if authToken == "" {
		return nil, nil, errors.New("account token is empty")
	}

	payload, err := token_lib.Decrypt(authToken)
	if err != nil {
		return nil, payload, err
	}

	// payload.ttl 单位为分钟
	ttlDuration := time.Duration(payload.TTL) * time.Minute

	// 转为秒检测超时
	if uint64(payload.IssueTime)+uint64(ttlDuration.Seconds()) <= uint64(time.Now().Unix()) {
		return nil, payload, errors.New("account token expired")
	}

	tokenKey := token.TokenKey{AccountId: payload.AccountId, Source: payload.LoginSource}
	token, err := tokenKey.GetToken()
	if err != nil || token == nil {
		return nil, payload, errors.New("account token db get failed")
	}

	account, err := model.GetAccountById(db.Mysql, payload.AccountId)
	if err != nil || account == nil {
		return nil, payload, errors.New("account token payload invalid")
	}

	return account, payload, nil
}

func authMiddleware(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		authToken := r.Header.Get(common.AuthHeader)
		account, payload, err := CheckAccessToken(authToken)
		if err != nil {
			render.InvalidToken(w, r)
			return
		}

		authContext := AuthContext{CurrentAccount: account, Payload: payload}
		ctx := context.WithValue(r.Context(), common.CurrentAccount, authContext)
		next.ServeHTTP(w, r.WithContext(ctx))
	}

	return http.HandlerFunc(fn)
}
