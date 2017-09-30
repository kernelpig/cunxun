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
	regCommentRouter(router)
	regImageRouter(router)
	regCarpoolingRouter(router)

	return router
}

func regUserRouter(router *gin.Engine) {
	group := router.Group("/api/u")
	group.POST("/signup", UserSignupHandler)
	group.POST("/login", UserLoginHandler)
	group.POST("/logout", UserLogoutHandler)
	group.GET("/:user_id/avatar", UserGetAvatarHandler)
	group.Use(middleware.AuthMiddleware())
	{
		group.GET("/", UserGetListHandler)
		group.GET("/:user_id", UserGetInfoHandler)
		group.POST("/", UserCreateHandler)
		group.PUT("/:user_id", UserUpdateHandler)
		group.DELETE("/:user_id", UserDeleteByIdHandler)
	}
}

func regCaptchaRouter(router *gin.Engine) {
	group := router.Group("/api/captcha")
	group.GET("/:captcha_id", CaptchaGetImageHandler)
	group.POST("/", CaptchaCreateHandler)
}

func regCheckCodeRouter(router *gin.Engine) {
	group := router.Group("/api/checkcode")
	group.POST("/send", CheckcodeSendHandler)
	group.POST("/check", CheckcodeVerifyHandler)
}

func regDebugRouter(router *gin.Engine) {
	group := router.Group("/api/debug")
	group.GET("/ping", DebugPingHandler)
	group.GET("/captcha/:captcha_id", DebugCaptchaGetValueHandler)
	group.GET("/checkcode", DebugCheckcodeGetValueHandler)
}

func regColumnRouter(router *gin.Engine) {
	group := router.Group("/api/column")
	group.GET("/", ColumnGetListHandler)
	group.Use(middleware.AuthMiddleware())
	{
		group.POST("/", ColumnCreateHandler)
		group.PUT("/:column_id", ColumnUpdateByIdHandler)
		group.DELETE("/:column_id", ColumnDeleteByIdHandler)
	}
}

func regArticleRouter(router *gin.Engine) {
	group := router.Group("/api/article")
	group.GET("/", ArticleGetListHandler)
	group.GET("/:article_id", ArticleGetHandler)
	group.Use(middleware.AuthMiddleware())
	{
		group.POST("/", ArticleCreateHandler)
		group.PUT("/:article_id", ArticleUpdateByIdHandler)
		group.DELETE("/:article_id", ArticleDeleteByIdHandler)
	}
}

func regCarpoolingRouter(router *gin.Engine) {
	group := router.Group("/api/carpooling")
	//group.GET("/", ArticleGetListHandler)
	group.GET("/:article_id", CarpoolingGetHandler)
	group.Use(middleware.AuthMiddleware())
	{
		group.POST("/", CarpoolingCreateHandler)
		//group.PUT("/:article_id", ArticleUpdateByIdHandler)
		//group.DELETE("/:article_id", ArticleDeleteByIdHandler)
	}
}

func regCommentRouter(router *gin.Engine) {
	group := router.Group("/api/comment")
	group.GET("/", CommentGetListHandler)
	group.GET("/:comment_id", CommentGetHandler)
	group.Use(middleware.AuthMiddleware())
	{
		group.POST("/", CommentCreateHandler)
		group.PUT("/:comment_id", CommentUpdateByIdHandler)
		group.DELETE("/:comment_id", CommentDeleteByIdHandler)
	}
}

func regImageRouter(router *gin.Engine) {
	group := router.Group("/api/image")
	group.POST("/", ImageCreateHandler)
}
