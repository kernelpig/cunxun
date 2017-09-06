package middleware

import (
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"wangqingang/cunxun/common"
	e "wangqingang/cunxun/error"
	"wangqingang/cunxun/token"
	"wangqingang/cunxun/token/token_lib"
)

type AuthContext struct {
	Payload *token_lib.Payload
}

func CheckAccessToken(authToken string) (*token_lib.Payload, error) {
	if authToken == "" {
		return nil, e.SP(e.MTokenErr, e.TokenIsEmpty, errors.New("auth middleware or logout"))
	}

	payload, err := token_lib.Decrypt(authToken)
	if err != nil {
		return payload, e.SP(e.MTokenErr, e.TokenDecryptErr, err)
	}

	// payload.ttl 单位为分钟
	ttlDuration := time.Duration(payload.TTL) * time.Minute

	// 转为秒检测超时
	if uint64(payload.IssueTime)+uint64(ttlDuration.Seconds()) <= uint64(time.Now().Unix()) {
		return payload, e.SP(e.MTokenErr, e.TokenExpired, errors.New("auth middlewareor logout"))
	}

	tokenKey := token.TokenKey{UserId: int(payload.UserId), Source: payload.LoginSource}
	token, err := tokenKey.GetToken()
	if err != nil || token == nil {
		return payload, e.SP(e.MTokenErr, e.TokenGetErr, err)
	}

	return payload, nil
}

func AuthMiddleware() gin.HandlerFunc {
	fn := func(c *gin.Context) {
		authToken := c.GetHeader(common.AuthHeaderKey)
		payload, err := CheckAccessToken(authToken)
		if err != nil || payload == nil {
			c.JSON(http.StatusBadRequest, e.SP(e.MTokenErr, e.TokenInvalid, err))
			return
		}

		authContext := AuthContext{Payload: payload}
		c.Set(common.CurrentUser, authContext)
		c.Next()
	}

	return fn
}
