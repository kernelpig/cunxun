package handler

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"wangqingang/cunxun/common"
	"wangqingang/cunxun/model/captcha"
)

func DebugPingHandler(c *gin.Context) {
	c.String(http.StatusOK, fmt.Sprintf("pong"))
	return
}

func DebugGetCaptchaImageHandler(c *gin.Context) {
	id := c.Param("id")

	c.JSON(http.StatusOK, gin.H{
		"code":          common.OK,
		"captcha_value": captcha.GetCaptchaValue(id, false),
	})

	return
}
