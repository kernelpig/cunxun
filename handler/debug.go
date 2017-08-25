package handler

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"wangqingang/cunxun/captcha"
	"wangqingang/cunxun/checkcode"
	"wangqingang/cunxun/common"
)

func DebugPingHandler(c *gin.Context) {
	c.String(http.StatusOK, fmt.Sprintf("pong"))
	return
}

func DebugCaptchaGetValueHandler(c *gin.Context) {
	id := c.Param("captcha_id")

	c.JSON(http.StatusOK, gin.H{
		"code":          common.OK,
		"captcha_value": captcha.GetCaptchaValue(id, false),
	})

	return
}

func DebugCheckcodeGetValueHandler(c *gin.Context) {
	key := checkcode.CheckCodeKey{
		Phone:   c.Query("phone"),
		Purpose: c.Query("purpose"),
		Source:  c.Query("source"),
	}
	checkcode, err := key.GetCheckcode()
	if err != nil || checkcode == nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": common.AccountInternalError,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code":      common.OK,
		"checkcode": checkcode.Code,
	})
}
