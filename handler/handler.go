package handler

import (
	"github.com/gin-gonic/gin"

	"wangqingang/cunxun/middleware"
)

func ServerEngine() *gin.Engine {
	router := gin.New()

	if router == nil {
		panic("create server failed")
	}

	router.Use(middleware.CrossMiddleware())
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	regCaptchaRouter(router)
	regCheckCodeRouter(router)
	regUserRouter(router)
	regDebugRouter(router)

	return router
}

func regUserRouter(router *gin.Engine) {
	userGroup := router.Group("/u")
	userGroup.POST("/signup", UserSignupHandler)
	userGroup.POST("/login", UserLoginHandler)
	userGroup.POST("/logout", UserLogoutHandler)
}

func regCaptchaRouter(router *gin.Engine) {
	captchaGroup := router.Group("/captcha")
	captchaGroup.GET("/:captcha_id", CaptchaGetImageHandler)
	captchaGroup.POST("/", CaptchaCreateHandler)
}

func regCheckCodeRouter(router *gin.Engine) {
	checkCodeGroup := router.Group("/checkcode")
	checkCodeGroup.POST("/send", CheckcodeSendHandler)
	checkCodeGroup.POST("/check", CheckcodeVerifyHandler)
}

func regDebugRouter(router *gin.Engine) {
	debugGroup := router.Group("/debug")
	debugGroup.GET("/ping", DebugPingHandler)
	debugGroup.GET("/captcha/:captcha_id", DebugCaptchaGetValueHandler)
	debugGroup.GET("/checkcode", DebugCheckcodeGetValueHandler)
}
