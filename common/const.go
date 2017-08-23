package common

const (
	CurrentUser = "current_user"
	AuthHeader  = "Authorization"
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

var SourceRange []string
var PurposeRange []string

func init() {
	SourceRange = []string{WebSource, AppSource}
	PurposeRange = []string{SignupPurpose, SigninPurpose, UpdatePhonePurpose, UpdateEmailPurpose, ResetPasswordPurpose}
}
