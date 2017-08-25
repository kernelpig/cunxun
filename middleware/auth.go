package middleware

import (
	"errors"
	"net/http"
	"time"

	"wangqingang/cunxun/common"
	"wangqingang/cunxun/token"
	"wangqingang/cunxun/token/token_lib"

	"github.com/gin-gonic/gin"
)

type AuthContext struct {
	Payload *token_lib.Payload
}

func CheckAccessToken(authToken string) (*token_lib.Payload, error) {
	if authToken == "" {
		return nil, errors.New("user token is empty")
	}

	payload, err := token_lib.Decrypt(authToken)
	if err != nil {
		return payload, err
	}

	// payload.ttl 单位为分钟
	ttlDuration := time.Duration(payload.TTL) * time.Minute

	// 转为秒检测超时
	if uint64(payload.IssueTime)+uint64(ttlDuration.Seconds()) <= uint64(time.Now().Unix()) {
		return payload, errors.New("user token expired")
	}

	tokenKey := token.TokenKey{UserId: int(payload.UserId), Source: payload.LoginSource}
	token, err := tokenKey.GetToken()
	if err != nil || token == nil {
		return payload, errors.New("user token db get failed")
	}

	return payload, nil
}

func authMiddleware() gin.HandlerFunc {
	fn := func(c *gin.Context) {
		authToken := c.GetHeader(common.AuthHeaderKey)
		payload, err := CheckAccessToken(authToken)
		if err != nil || payload == nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"code": common.AccountInvalidToken,
			})
			return
		}

		authContext := AuthContext{Payload: payload}
		c.Set(common.CurrentUser, authContext)
		c.Next()
	}

	return fn
}
