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
	regColumnRouter(router)
	regArticleRouter(router)

	return router
}

func regUserRouter(router *gin.Engine) {
	group := router.Group("/u")
	group.POST("/signup", UserSignupHandler)
	group.POST("/login", UserLoginHandler)
	group.POST("/logout", UserLogoutHandler)
}

func regCaptchaRouter(router *gin.Engine) {
	group := router.Group("/captcha")
	group.GET("/:captcha_id", CaptchaGetImageHandler)
	group.POST("/", CaptchaCreateHandler)
}

func regCheckCodeRouter(router *gin.Engine) {
	group := router.Group("/checkcode")
	group.POST("/send", CheckcodeSendHandler)
	group.POST("/check", CheckcodeVerifyHandler)
}

func regDebugRouter(router *gin.Engine) {
	group := router.Group("/debug")
	group.GET("/ping", DebugPingHandler)
	group.GET("/captcha/:captcha_id", DebugCaptchaGetValueHandler)
	group.GET("/checkcode", DebugCheckcodeGetValueHandler)
}

func regColumnRouter(router *gin.Engine) {
	group := router.Group("/column")
	group.GET("/", ColumnGetAllHandler)
	group.Use(middleware.AuthMiddleware())
	{
		group.POST("/", ColumnCreateHandler)
	}
}

func regArticleRouter(router *gin.Engine) {
	group := router.Group("/article")
	group.GET("/", ArticleGetListHandler)
	group.GET("/:article_id", ArticleGetHandler)
	group.Use(middleware.AuthMiddleware())
	{
		group.POST("/", ArticleCreateHandler)
	}
}
