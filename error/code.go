/**
* 1. 错误码共4字节, 包括分配的模块标识, 接口标识, 子模块标识, 子模块错误, 各占1字节空间
* 2. 逻辑与组合在一起, 从高位到低位依次为: 模块标识, 接口标识, 子模块标识, 子模块错误
* 3. 暂未涉及到多种CPU, 如果后续涉及会采用大段字节序传输
* 4. 0x00错误码为正常状态, 包括各个子模块, 故成功错误码为0x00000000

* PS 今天感觉有点蛋疼, 修改这个错误码花费了我整整一天的时间, 刚才突然想到, 没有必要这么复
*    杂, 因为前端/移动端只关心业务相关的错误码, 你内部怎么管理你错误码是自己的事情, 只要
*    对外的错误码即可, 错误码栈可以参考一下开源的处理.
**/
package error

import (
	"encoding/json"
	"fmt"
	"strconv"
)

// 特殊错误码
const (
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
	MysqlLastInsertErr
	MysqlDuplicateErr

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

type Code struct {
	ServiceIndex   int
	InterfaceIndex int
	SubModuleIndex int
	SubErrorIndex  int
	Code           int
}

func (c Code) MarshalJSON() ([]byte, error) {
	return []byte(strconv.Itoa(c.Code)), nil
}

type Message struct {
	Code        Code      `json:"code"`
	Service     string    `json:"service"`
	Interface   string    `json:"interface"`
	SubModule   string    `json:"sub_module"`
	SubError    string    `json:"sub_error"`
	Detail string    `json:"detail"`
	Previous  	*Message `json:"previous"`
}

func (m Message) Error() string {
	jsonMessage, err := json.Marshal(&m)
	if err != nil {
		return fmt.Sprintf("{\"code\": 0x%08x, \"detail\": \"%s\"}", Exception, err.Error())
	}
	return fmt.Sprintf("%s", jsonMessage)
}

// 生成错误码, C - Code
func (c *Code) C() int {
	c.Code = (c.ServiceIndex << 24 & ServiceErrMask) | (c.InterfaceIndex << 16 & InterfaceErrMask) |
		(c.SubModuleIndex << 8 & SubModuleErrMask) | (c.SubErrorIndex << 0 & SubErrorMask)
	return c.Code
}

// 生成带有详细信息的错误信息
func (c *Code) MD2E(detail string, lastErr error) error {
	c.C()
	message := Message{Code: *c}
	message.Detail = detail

	if c.ServiceIndex >= _ServiceErrMax || c.ServiceIndex < _ServiceErrMin {
		message.Service = "invalid service error code"
	} else {
		message.Service = ServiceErrs[c.ServiceIndex]
	}
	if c.InterfaceIndex >= _InterfaceErrMax || c.InterfaceIndex < _InterfaceErrMin {
		message.Interface = "invalid interface error code"
	} else {
		message.Interface = InterfaceErrs[c.InterfaceIndex]
	}
	if c.SubModuleIndex >= _SubModuleErrMax || c.SubModuleIndex < _SubModuleErrMin {
		message.SubModule = "invalid sub module error code"
	} else {
		message.SubModule = SubModuleErrs[c.SubModuleIndex]
	}
	if c.SubErrorIndex >= len(SubErrors[c.SubModuleIndex]) {
		message.SubError = "invalid sub error code"
	} else {
		message.SubError = SubErrors[c.SubModuleIndex][c.SubErrorIndex]
	}
	if lastErr != nil {
		if lastErrMsg, ok := lastErr.(Message); ok {
			message.Previous = &lastErrMsg
		} else {
			message.Previous = &Message{
				Code: Code{Code: Exception},
				Detail: lastErr.Error(),
			}
		}
	}
	return message
}

func (c *Code) IsSubError(subModuleIndex, subErrorIndex int) bool {
	return c.SubModuleIndex == subModuleIndex && c.SubErrorIndex == subErrorIndex
}

// 子模块错误
func S(idxSubModuleErr, idxSubError int) error {
	return SE(idxSubModuleErr, idxSubError, "", nil)
}

// 子模块错误, subModule error with previous
func SP(idxSubModuleErr, idxSubError int, previous error) error {
	return SE(idxSubModuleErr, idxSubError, "", previous)
}

// 子模块错误, subModule error with detail
func SD(idxSubModuleErr, idxSubError int, detail string) error {
	return SE(idxSubModuleErr, idxSubError, detail, nil)
}

// 子模块扩展错误, subModule extent error with datail and previous
func SE(idxSubModuleErr, idxSubError int, detail string, previous error) error {
	code := Code{
		ServiceIndex:   SCunxun,
		InterfaceIndex: 0x00,
		SubModuleIndex: idxSubModuleErr,
		SubErrorIndex:  idxSubError,
	}
	return code.MD2E(detail, previous)
}

// 接口错误
func I(idxInterfaceErr, idxSubModuleErr, idxSubError int) error {
	return IE(idxInterfaceErr, idxSubModuleErr, idxSubError, "", nil)
}

// 接口错误, interface error with previous
func IP(idxInterfaceErr, idxSubModuleErr, idxSubError int, previous error) error {
	return IE(idxInterfaceErr, idxSubModuleErr, idxSubError, "", previous)
}

// 接口错误, interface error with detail
func ID(idxInterfaceErr, idxSubModuleErr, idxSubError int, detail string) error {
	return IE(idxInterfaceErr, idxSubModuleErr, idxSubError, detail, nil)
}

// 接口扩展错误, interface extent error with datail and previous
func IE(idxInterfaceErr, idxSubModuleErr, idxSubError int, detail string, previous error) error {
	code := Code{
		ServiceIndex:   SCunxun,
		InterfaceIndex: idxInterfaceErr,
		SubModuleIndex: idxSubModuleErr,
		SubErrorIndex:  idxSubError,
	}
	return code.MD2E(detail, previous)
}
