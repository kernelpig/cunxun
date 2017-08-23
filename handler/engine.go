package handler

import (
	"github.com/meiqia/chi"
	chi_middleware "github.com/meiqia/chi/middleware"

	"wangqingang/cunxun/middleware"
)

func AccountEngine() chi.Router {
	router := chi.NewRouter()
	router.Use(middleware.Logger)
	router.Use(middleware.Recovery)
	router.Use(chi_middleware.StripSlashes)

	registerAccountRouter(router)

	return router
}

func registerAccountRouter(router chi.Router) {
	router.Route("/account/v1", func(r chi.Router) {
		// 图形验证码
		r.Get("/get_captcha_image", GetCaptchaImageHandler)
		r.Post("/get_captcha_token", CreateCaptchaHandler)

		// 短信校验码
		r.Post("/send_verify_code", SendVerifyCodeHandler)
		r.Post("/check_verify_code", CheckVerifyCodeHandler)

		// 用户操作
		r.Post("/reset_password", ResetPasswordHandler)
		r.Post("/sign_up", SignUpHandler)
		r.Post("/sign_in", SignInHandler)

		r.Route("/", func(r chi.Router) {
			r.Use(middleware.AuthRequired)
			r.Post("/sign_out", SignOutHandler)
		})

		// 内部接口
		r.Get("/internal/verify_token", VerifyTokenHandler)
		r.Get("/internal/get_phone", GetPhoneHandler)

		// 测试接口
		r.Get("/debug/ping", DebugPingHandler)
		r.Get("/debug/get_captcha_value", DebugGetCaptchaImageHandler)
	})
}
