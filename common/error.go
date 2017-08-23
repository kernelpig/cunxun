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
