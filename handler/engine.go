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
	regDebugRouter(router)

	return router
}

func regCaptchaRouter(router *gin.Engine) {
	captchaGroup := router.Group("/captcha")
	captchaGroup.GET("/:id", GetCaptchaImageHandler)
	captchaGroup.POST("/", CreateCaptchaHandler)
}

func regDebugRouter(router *gin.Engine) {
	debugGroup := router.Group("/debug")
	debugGroup.GET("/ping", DebugPingHandler)
	debugGroup.GET("/captcha/:id", DebugGetCaptchaImageHandler)
}
