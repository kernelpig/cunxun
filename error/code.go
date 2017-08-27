/**
* 1. 错误码共4字节, 包括分配的模块标识, 接口标识, 子模块标识, 子模块错误, 各占1字节空间
* 2. 逻辑与组合在一起, 从高位到低位依次为: 模块标识, 接口标识, 子模块标识, 子模块错误
* 3. 暂未涉及到多种CPU, 如果后续涉及会采用大段字节序传输
* 4. 0x00错误码为正常状态, 包括各个子模块, 故成功错误码为0x00000000
**/
package error

import (
	"encoding/json"
	"fmt"
)

// 特殊错误码
const (
	IG        = 0x00
	OK        = 0x00000000
	Exception = 0xffffffff
)

// 错误码处理
const (
	ServiceErrMask   = 0xff000000
	InterfaceErrMask = 0x00ff0000
	SubModuleErrMask = 0x0000ff00
	SubErrorMask     = 0x000000ff
)

// 分配的模块标识
const (
	_ServiceErrMin = iota
	SCunxun

	// 注意: 请在此处增加错误码, 已废弃的请保留不要删除

	_ServiceErrMax
)

// 接口标识
const (
	_InterfaceErrMin = iota
	ICaptchaCreate
	ICaptchaGetImage
	ICheckcodeSend
	ICheckcodeCheck
	IDebugPing
	IDebugCaptchaGetValue
	IDebugCheckcodeGetValue
	IUserSignup
	IUserLogin
	IUserLogout

	// 注意: 请在此处增加错误码, 已废弃的请保留不要删除!

	_InterfaceErrMax
)

// 子模块标识
const (
	_SubModuleErrMin = iota
	MConfigErr
	MLogErr
	MUserErr
	MLoginErr
	MCheckcodeErr
	MTokenErr
	MCaptchaErr
	MPasswordErr
	MMysqlErr
	MRedisErr
	MSmsErr
	MParamsErr
	MUtilsErr
	MAuthErr
	MOthersErr
	MPhoneErr

	// 注意: 请在此处增加错误码, 已废弃的请保留不要删除!

	_SubModuleErrMax
)

// 配置错误
const (
	_ConfigErrMin = iota
	ConfigLoadErr
	ConfigParseErr
	ConfigParseTimeErr

	// 注意: 请在此处增加错误码, 已废弃的请保留不要删除!

	_ConfigErrMax
)

// Log日志错误
const (
	_LogErrMin = iota
	LogDumpRequestErr

	// 注意: 请在此处增加错误码, 已废弃的请保留不要删除!

	_LogErrMax
)

//  model错误
const (
	_UserErrMin = iota
	UserGetErr
	UserAlreadyExist
	UserNotExist
	UserCreateErr

	// 注意: 请在此处增加错误码, 已废弃的请保留不要删除!

	_UserErrMax
)

// 登录限制错误
const (
	_LoginErrMin = iota
	LoginSaveErr
	LoginGetErr
	LoginCreateErr

	// 注意: 请在此处增加错误码, 已废弃的请保留不要删除!

	_LoginErrMax
)

// 短信验证码检查错误
const (
	_CheckcodeErrMin = iota
	CheckcodeSaveErr
	CheckcodeGetErr
	CheckcodeCreateErr
	CheckcodeRequestLimit
	CheckcodeCheckLimit
	CheckcodeMismatch
	CheckcodeNotFound
	CheckcodeCheckErr

	// 注意: 请在此处增加错误码, 已废弃的请保留不要删除!

	_CheckcodeErrMax
)

// Access Token错误
const (
	_TokenErrMin = iota
	TokenInitPubKeyErr
	TokenInitPriKeyErr
	TokenSaveErr
	TokenIsEmpty
	TokenDecryptErr
	TokenExpired
	TokenGetErr
	TokenInvalid
	TokenEcryptErr
	TokenInvalidVersion
	TokenBase64DecodeErr
	TokenSignErr
	TokenSignVerifyErr
	TokenCreateErr

	// 注意: 请在此处增加错误码, 已废弃的请保留不要删除!

	_TokenErrMax
)

// 图形验证码错误
const (
	_CaptchaErrMin = iota
	CaptchaWriteImageErr
	CaptchaWriteAudioErr
	CaptchaGetErr
	CaptchaMismatch
	CaptchaRequired

	// 注意: 请在此处增加错误码, 已废弃的请保留不要删除!

	_CaptchaErrMax
)

// 密码错误
const (
	_PasswordErrMin = iota
	PasswordEncryptErr
	PasswordVerifyErr
	PasswordLevelErr
	PasswordInvalid

	// 注意: 请在此处增加错误码, 已废弃的请保留不要删除!

	_PasswordErrMax
)

// Mysql服务错误
const (
	_MysqlErrMin = iota
	MysqlConnectErr
	MysqlSelectErr
	MysqlUpdateErr
	MysqlInsertErr
	MysqlDeleteErr
	MysqlRowAffectErr
	MysqlRowScanErr

	// 注意: 请在此处增加错误码, 已废弃的请保留不要删除!

	_MysqlErrMax
)

// redis服务错误
const (
	_RedisErrMin = iota
	RedisValueMarshalErr
	RedisValueUnmarshalErr
	RedisSetErr
	RedisGetErr
	RedisDelErr

	// 注意: 请在此处增加错误码, 已废弃的请保留不要删除!

	_RedisErrMax
)

// 短信服务错误
const (
	_SmsErrMin = iota
	SmsConnectErr
	SmsInvalidPurpose
	SmsSendErr

	// 注意: 请在此处增加错误码, 已废弃的请保留不要删除!

	_SmsErrMax
)

// 参数错误
const (
	_ParamsErrMin = iota
	ParamsInvalidPhone
	ParamsInvalidSource
	ParamsInvalidPurpose
	ParamsBindErr

	// 注意: 请在此处增加错误码, 已废弃的请保留不要删除!

	_ParamErrMax
)

// 认证错误
const (
	_AuthErrMin = iota
	AuthTokenEmpty

	// 注意: 请在此处增加错误码, 已废弃的请保留不要删除!

	_AuthErrMax
)

// utils错误
const (
	_UtilsErrMin = iota

	// 注意: 请在此处增加错误码, 已废弃的请保留不要删除!

	_UtilsErrMax
)

// 其他杂项错误
const (
	_OthersErrMin = iota

	// 注意: 请在此处增加错误码, 已废弃的请保留不要删除!

	_OthersErrMax
)

// 手机号码错误
const (
	_PhoneErrMin = iota
	PhoneEmpty
	PhoneFormatErr
	PhoneInvalidCountryCode
	PhoneUnknownRegion
	PhoneRegionMismatch
	PhoneParseNumberErr

	// 注意: 请在此处增加错误码, 已废弃的请保留不要删除!

	_PhoneErrMax
)

type Message struct {
	Code        int       `json:"code"`
	Service     string    `json:"service"`
	Interface   string    `json:"interface"`
	SubModule   string    `json:"sub_module"`
	SubError    string    `json:"sub_error"`
	ErrorDetail string    `json:"error_detail"`
	ErrorStack  []Message `json:"error_stack"`
}

func (m Message) Error() string {
	jsonMessage, err := json.Marshal(&m)
	if err != nil {
		return fmt.Sprintf("{\"code\": 0x%08x, \"detail\": \"%s\"}", Exception, err.Error())
	}
	return fmt.Sprintf("%s", jsonMessage)
}

func SE(idxSubModuleErr, idxSubError int, detail error) error {
	return MD2E(SCunxun, IG, idxSubModuleErr, idxSubError, detail)
}

func IE(idxInterfaceErr, idxSubModuleErr, idxSubError int, detail error) error {
	return MD2E(SCunxun, idxInterfaceErr, idxSubModuleErr, idxSubError, detail)
}

// 生成带有详细信息的错误信息
func MD2E(idxServiceErr, idxInterfaceErr, idxSubModuleErr, idxSubError int, detail error) error {
	code := C(idxServiceErr, idxInterfaceErr, idxSubModuleErr, idxSubError)
	message := Message{Code: code}
	message.ErrorStack = make([]Message, 0)

	if idxServiceErr >= _ServiceErrMax || idxServiceErr < _ServiceErrMin {
		message.Service = "invalid service error code"
	} else {
		message.Service = ServiceErrs[idxServiceErr]
	}
	if idxInterfaceErr >= _InterfaceErrMax || idxInterfaceErr < _InterfaceErrMin {
		message.Interface = "invalid interface error code"
	} else {
		message.Interface = InterfaceErrs[idxInterfaceErr]
	}
	if idxSubModuleErr >= _SubModuleErrMax || idxSubModuleErr < _SubModuleErrMin {
		message.SubModule = "invalid sub module error code"
	} else {
		message.SubModule = SubModuleErrs[idxSubModuleErr]
	}
	if idxSubError >= len(SubErrors[idxSubModuleErr]) {
		message.SubError = "invalid sub error code"
	} else {
		message.SubError = SubErrors[idxSubModuleErr][idxSubError]
	}

	// 如果是Message类型的Error, 则拷贝其error stack, 并追加本次错误
	// 如果是Normal类型的Error, 则写入到本次的detail信息
	if detail != nil {
		if detailMessage, ok := detail.(Message); ok {
			message.ErrorStack = append(message.ErrorStack, detailMessage)
			for _, element := range detailMessage.ErrorStack {
				message.ErrorStack = append(message.ErrorStack, element)
			}
		} else {
			message.ErrorDetail = detail.Error()
		}
	}
	return message
}

// 生成错误信息, M - Marshal
func M2E(idxServiceErr, idxInterfaceErr, idxSubModuleErr, idxSubError int) error {
	return MD2E(idxServiceErr, idxInterfaceErr, idxSubModuleErr, idxSubError, nil)
}

// 生成错误码, C - Code
func C(idxServiceErr, idxInterfaceErr, idxSubModuleErr, idxSubError int) int {
	return (idxServiceErr << 24 & ServiceErrMask) | (idxInterfaceErr << 16 & InterfaceErrMask) |
		(idxSubModuleErr << 8 & SubModuleErrMask) | (idxSubError << 0 & SubErrorMask)
}
