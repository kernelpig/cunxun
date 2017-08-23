package handler

import (
	"github.com/gin-gonic/gin"
)

func ServerEngine() *gin.Engine {
	router := gin.New()

	if router == nil {
		panic("create server failed")
	}

	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	regCaptchaRouter(router)
	regCheckCodeRouter(router)
	regDebugRouter(router)

	return router
}

func regCaptchaRouter(router *gin.Engine) {
	captchaGroup := router.Group("/captcha")
	captchaGroup.GET("/:captcha_id", GetCaptchaImageHandler)
	captchaGroup.POST("/", CreateCaptchaHandler)
}

func regCheckCodeRouter(router *gin.Engine) {
	checkCodeGroup := router.Group("/checkcode")
	checkCodeGroup.POST("/")
}

func regDebugRouter(router *gin.Engine) {
	debugGroup := router.Group("/debug")
	debugGroup.GET("/ping", DebugPingHandler)
	debugGroup.GET("/captcha/:captcha_id", DebugGetCaptchaValueHandler)
}
