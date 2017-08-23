package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"wangqingang/cunxun/captcha"
	"wangqingang/cunxun/common"
)

func CreateCaptchaHandler(c *gin.Context) {
	id := captcha.GenCaptcha(common.Config.Captcha.DefaultLength)
	c.JSON(http.StatusOK, gin.H{
		"code":       common.OK,
		"captcha_id": id,
	})
	return
}

func GetCaptchaImageHandler(c *gin.Context) {
	id := c.Param("captcha_id")

	width, err := strconv.ParseInt(c.Query("width"), 10, 64)
	if err != nil {
		width = int64(common.Config.Captcha.DefaultWidth)
	}

	height, err := strconv.ParseInt(c.Query("height"), 10, 64)
	if err != nil {
		height = int64(common.Config.Captcha.DefaultHeight)
	}

	data, err := captcha.GetCaptchaImage(id, int(width), int(height))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    common.AccountInternalError,
			"message": "GetCaptchaImage " + err.Error(),
		})
		return
	}

	c.Data(http.StatusOK, "image/png", data)
	return
}
