package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"wangqingang/cunxun/common"
	"wangqingang/cunxun/model/captcha"
)

func CreateCaptchaHandler(c *gin.Context) {
	id := captcha.GenCaptcha(common.Config.Captcha.DefaultLength)
	c.JSON(http.StatusOK, gin.H{
		"code": common.OK,
		"id":   id,
	})
	return
}

func GetCaptchaImageHandler(c *gin.Context) {
	id := c.Param("id")

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
		c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"code":    common.AccountInternalError,
			"message": "GetCaptchaImage " + err.Error(),
		})
		return
	}

	c.Data(http.StatusOK, "image/png", data)
	return
}
