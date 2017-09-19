package common

const (
	ModuleName    = "cunxun"
	CurrentUser   = "current_user"
	AuthHeaderKey = "Authorization"
)

const (
	SignupPurpose        = "signup"
	SigninPurpose        = "signin"
	UpdatePhonePurpose   = "update_phone"
	UpdateEmailPurpose   = "update_email"
	ResetPasswordPurpose = "reset_password"
)

const (
	WebSource = "web"
	AppSource = "app"
)

const (
	// 默认头像标识, TEXT字段NULL处理问题
	UserAvatarDefulatFlag = "default avatar"
)

var SourceRange []string
var PurposeRange []string

func init() {
	SourceRange = []string{WebSource, AppSource}
	PurposeRange = []string{SignupPurpose, SigninPurpose, UpdatePhonePurpose, UpdateEmailPurpose, ResetPasswordPurpose}
}
