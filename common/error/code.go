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

const (
	ServiceErrMask   = 0xff000000
	InterfaceErrMask = 0x00ff0000
	SubModuleErrMask = 0x0000ff00
	SubErrorMask     = 0x000000ff
	OK               = 0x00000000
	Exception        = 0xffffffff
)

// 分配的模块标识
const (
	_ServiceErrMin = iota
	SCunxun

	// 注意: 请在此处增加错误码, 已废弃的请保留不要删除!

	_ServiceErrMax
)

// 接口标识
const (
	_InterfaceErrMin = iota
	ICaptchaCreate
	ICaptchaGetImage
	ICheckcodeSend
	ICheckcode
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

	// 注意: 请在此处增加错误码, 已废弃的请保留不要删除!

	_UserErrMax
)

// 登录限制错误
const (
	_LoginErrMin = iota
	LoginSaveErr

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

	// 注意: 请在此处增加错误码, 已废弃的请保留不要删除!

	_CaptchaErrMax
)

// 密码错误
const (
	_PasswordErrMin = iota
	PasswordEncryptErr

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

	// 注意: 请在此处增加错误码, 已废弃的请保留不要删除!

	_ParamErrMax
)

// 参数错误
const (
	_UtilsErrMin = iota

	// 注意: 请在此处增加错误码, 已废弃的请保留不要删除!

	_UtilsErrMax
)

// 其他杂项错误
const (
	_OthersErrMin = iota
	PhoneEmpty
	PhoneFormatErr
	PhoneInvalidCountryCode
	PhoneUnknownRegion
	PhoneRegionMismatch

	// 注意: 请在此处增加错误码, 已废弃的请保留不要删除!

	_OthersErrMax
)

type Message struct {
	Code        int    `json:"code"`
	Service     string `json:"service"`
	Interface   string `json:"interface"`
	SubModule   string `json:"sub_module"`
	SubError    string `json:"sub_error"`
	ErrorDetail string `json:"error_detail"`
}

// 生成带有详细信息的错误信息
func MDS(code int, detail string) string {
	message := Message{Code: code}
	message.ErrorDetail = detail

	serviceIndex := code & ServiceErrMask >> 24
	if serviceIndex >= _ServiceErrMax || serviceIndex <= _ServiceErrMin {
		message.Service = "invalid service error code"
	} else {
		message.Service = ServiceErrs[serviceIndex]
	}
	interfaceIndex := code & InterfaceErrMask >> 16
	if interfaceIndex >= _InterfaceErrMax || interfaceIndex <= _InterfaceErrMin {
		message.Interface = "invalid interface error code"
	} else {
		message.Interface = InterfaceErrs[interfaceIndex]
	}
	subModuleIndex := code & SubModuleErrMask >> 8
	if subModuleIndex >= _SubModuleErrMax || subModuleIndex <= _SubModuleErrMin {
		message.SubModule = "invalid sub module error code"
	} else {
		message.SubModule = SubModuleErrs[subModuleIndex]
	}
	idxSubError := code & SubErrorMask >> 0
	if idxSubError >= len(SubErrors[subModuleIndex]) {
		message.SubError = "invalid sub error code"
	} else {
		message.SubError = SubErrors[subModuleIndex][idxSubError]
	}
	fmt.Println("11111", serviceIndex, interfaceIndex, subModuleIndex, idxSubError)
	jsonMessage, err := json.Marshal(&message)
	if err != nil {
		return fmt.Sprintf("{\"code\": 0x%08x, \"detail\": \"%s\"}", Exception, err.Error())
	}
	return string(jsonMessage)
}

// 生成错误信息, M - Marshal
func M(code int) string {
	return MDS(code, "")
}

// 生成带有详细信息的错误信息
func MDE(code int, err error) string {
	if err != nil {
		return MDS(code, err.Error())
	}
	return M(code)
}

// 生成错误码, C - Code
func C(idxServiceErr, idxInterfaceErr, idxSubModuleErr, idxSubError int) int {
	return (idxServiceErr << 24 & ServiceErrMask) | (idxInterfaceErr << 16 & InterfaceErrMask) |
		(idxSubModuleErr << 8 & SubModuleErrMask) | (idxSubError << 0 & SubErrorMask)
}
