package handler

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"wangqingang/cunxun/captcha"
	"wangqingang/cunxun/checkcode"
	e "wangqingang/cunxun/error"
)

func DebugPingHandler(c *gin.Context) {
	c.String(http.StatusOK, fmt.Sprintf("pong"))
	return
}

func DebugCaptchaGetValueHandler(c *gin.Context) {
	id := c.Param("captcha_id")

	c.JSON(http.StatusOK, gin.H{
		"code":          e.OK,
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
		c.JSON(http.StatusInternalServerError, e.IP(e.IDebugCheckcodeGetValue, e.MCheckcodeErr, e.CheckcodeGetErr, err))
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code":      e.OK,
		"checkcode": checkcode.Code,
	})
}
