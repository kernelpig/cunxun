package handler

type CheckcodeSendRequest struct {
	Phone        string `json:"phone" binding:"required"`
	Purpose      string `json:"purpose" binding:"required"`
	Source       string `json:"source" binding:"required"`
	CaptchaId    string `json:"captcha_id" binding:"required"`
	CaptchaValue string `json:"captcha_value" binding:"required"`
}

type CheckVerifyCodeRequest struct {
	Phone      string `json:"phone" binding:"required"`
	Purpose    string `json:"purpose" binding:"required"`
	Source     string `json:"source" binding:"required"`
	VerifyCode string `json:"verify_code" binding:"required"`
}

type UserSignupRequest struct {
	Phone      string `json:"phone" binding:"required"`
	Source     string `json:"source" binding:"required"`
	NickName   string `json:"nickname" binding:"required,min=1,max=32"`
	Password   string `json:"password" binding:"required,min=8,max=16"`
	VerifyCode string `json:"verify_code" binding:"required"`
	Avatar     string `json:"avatar" binding:"omitempty"`
}

type UserCreateRequest struct {
	Phone    string `json:"phone" binding:"required"`
	Password string `json:"password" binding:"required,min=8,max=16"`
	NickName string `json:"nickname" binding:"omitempty"`
	Role     int    `json:"role"`
}

type UserUpdateRequest struct {
	Role int `json:"role"`
}

type UserLoginRequest struct {
	Phone        string `json:"phone" binding:"required"`
	Source       string `json:"source" binding:"required"`
	Password     string `json:"password" binding:"required,min=8,max=20"`
	CaptchaId    string `json:"captcha_id" binding:"omitempty"`
	CaptchaValue string `json:"captcha_value" binding:"omitempty"`
}

type ResetPasswordRequest struct {
	Phone       string `json:"phone" binding:"required"`
	Source      string `json:"source" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,min=8,max=20"`
	VerifyCode  string `json:"verify_code" binding:"required"`
}

type ColumnCreateRequest struct {
	Name string `json:"name"`
}

type ColumnUpdateRequest struct {
	Name string `json:"name" binding:"omitempty"`
}

type CarpoolingCreateRequest struct {
	FromCity    string `json:"from_city" binding:"required"`
	ToCity      string `json:"to_city" binding:"required"`
	DepartTime  int64  `json:"depart_time" binding:"required"`
	PeopleCount int    `json:"people_count" binding:"required"`
	Contact     string `json:"contact" binding:"required"`
	Remark      string `json:"remark" binding:"required"`
}

type CarpoolingUpdateRequest struct {
	FromCity    string `json:"from_city" binding:"required"`
	ToCity      string `json:"to_city" binding:"required"`
	DepartTime  int64  `json:"depart_time" binding:"required"`
	PeopleCount int    `json:"people_count" binding:"required"`
	Contact     string `json:"contact" binding:"required"`
	Remark      string `json:"remark" binding:"required"`
}

type ArticleCreateRequest struct {
	ColumnId int    `json:"column_id" binding:"required"`
	Title    string `json:"title" binding:"required"`
	Content  string `json:"content" binding:"required"`
}

type ArticleUpdateRequest struct {
	ColumnId int    `json:"column_id" binding:"omitempty"`
	Title    string `json:"title" binding:"omitempty"`
	Content  string `json:"content" binding:"omitempty"`
}

type CommentCreateRequest struct {
	ArticleId int    `json:"article_id" binding:"required"`
	Content   string `json:"content" binding:"required"`
}

type CommentUpdateRequest struct {
	ArticleId int    `json:"article_id" binding:"omitempty"`
	Content   string `json:"content" binding:"omitempty"`
}
