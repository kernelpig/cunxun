package common

// 错误码分类：1xx(参数错误)；2xx(操作错误); 3xx(外部依赖错误); 4xx(内部错误)
const (
	OK                          = 0
	AccountBindFailed           = 102100
	AccountInvalidPurpose       = 102101
	AccountInvalidSource        = 102102
	AccountAccountNotExist      = 102103
	AccountVerifyCodeNotMatch   = 102104
	AccountCaptchaNotMatch      = 102105
	AccountInvalidPassword      = 102106
	AccountNeedCaptcha          = 102107
	AccountAccountAlreadyExist  = 102108
	AccountInvalidPhone         = 102109
	AccountPasswordSameWithOld  = 102110
	AccountPasswordLevelIllegal = 102111
	AccountRequestLimit         = 102201
	AccountInvalidToken         = 102202
	AccountGenerateIdFailed     = 102301
	AccountDBError              = 102302
	AccountInternalError        = 102401
)


var Errors = make(map[int]string)

func init() {
	Errors[OK] = "操作成功"

	Errors[AccountInternalError] = "服务错误: "
	Errors[AccountBindFailed] = "参数解析错误: "
	Errors[AccountInvalidPurpose] = "无效目的操作"
	Errors[AccountInvalidSource] = "无效来源"
	Errors[AccountAccountNotExist] = "账户不存在"
	Errors[AccountVerifyCodeNotMatch] = "校验码不匹配"
	Errors[AccountCaptchaNotMatch] = "验证码不匹配"
	Errors[AccountInvalidPassword] = "密码错误: 剩余%d次机会"
	Errors[AccountNeedCaptcha] = "需要图形验证码"
	Errors[AccountAccountAlreadyExist] = "账户已存在"
	Errors[AccountPasswordSameWithOld] = "密码与原来相同"
	Errors[AccountPasswordLevelIllegal] = "密码等级非法"

	Errors[AccountRequestLimit] = "达到最大发送限制，请稍后重试"
	Errors[AccountInvalidToken] = "无效token："
}