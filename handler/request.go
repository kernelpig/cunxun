package handler

type SendVerifyCodeRequest struct {
	Phone        string `json:"phone" binding:"required"`
	Purpose      string `json:"purpose" binding:"required"`
	Source       string `json:"source" binding:"required"`
	CaptchaToken string `json:"captcha_token" binding:"required"`
	CaptchaValue string `json:"captcha_value" binding:"required"`
}

type CheckVerifyCodeRequest struct {
	Phone      string `json:"phone" binding:"required"`
	Purpose    string `json:"purpose" binding:"required"`
	Source     string `json:"source" binding:"required"`
	VerifyCode string `json:"verify_code" binding:"required"`
}

type SignUpRequest struct {
	Phone      string `json:"phone" binding:"required"`
	Source     string `json:"source" binding:"required"`
	Password   string `json:"password" binding:"required,min=8,max=16"`
	VerifyCode string `json:"verify_code" binding:"required"`
}

type SignInRequest struct {
	Phone        string `json:"phone" binding:"required"`
	Source       string `json:"source" binding:"required"`
	Password     string `json:"password" binding:"required,min=8,max=20"`
	CaptchaToken string `json:"captcha_token" binding:"omitempty"`
	CaptchaValue string `json:"captcha_value" binding:"omitempty"`
}

type ResetPasswordRequest struct {
	Phone       string `json:"phone" binding:"required"`
	Source      string `json:"source" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,min=8,max=20"`
	VerifyCode  string `json:"verify_code" binding:"required"`
}
