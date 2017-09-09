package error

var ServiceErrs [_ServiceErrMax]string
var InterfaceErrs [_InterfaceErrMax]string
var SubModuleErrs [_SubModuleErrMax]string
var SubErrors [_SubModuleErrMax][]string

func init() {
	initServiceErrs()
	initInterfaceErr()
	initSubModuleErrs()
	initSubErrors()
}

func initServiceErrs() {
	ServiceErrs = [_ServiceErrMax]string{
		SCunxun: "cunxun",
	}
}

func initSubErrors() {
	SubErrors = [_SubModuleErrMax][]string{
		_SubModuleErrMin: {
			_SubModuleErrMin: "Basic sub modules error information.",
		},
		MConfigErr: {
			_ConfigErrMin:      "Basic configuration error information.",
			ConfigLoadErr:      "Failed to load profile.",
			ConfigParseErr:     "Parsing configuration file failed.",
			ConfigParseTimeErr: "Parse time configuration field failed.",
		},
		MLogErr: {
			_LogErrMin:        "Basic log error information.",
			LogDumpRequestErr: "Dump request parameter failed.",
		},
		MUserErr: {
			_UserErrMin:      "Basic user error information.",
			UserGetErr:       "Failed to obtain user information.",
			UserAlreadyExist: "User already exists.",
			UserNotExist:     "User does not exist.",
			UserCreateErr:    "Failed to create user.",
		},
		MLoginErr: {
			_LoginErrMin:   "Basic logon error information.",
			LoginSaveErr:   "Failed to save logon information.",
			LoginGetErr:    "Failed to obtain logon information.",
			LoginCreateErr: "Failed to create logon information.",
		},
		MCheckcodeErr: {
			_CheckcodeErrMin:      "Basic checkcode error information.",
			CheckcodeSaveErr:      "Failed to save check code information.",
			CheckcodeGetErr:       "Failed to retrieve check code information.",
			CheckcodeCreateErr:    "Failed to create verification code information.",
			CheckcodeRequestLimit: "Check code creation request limit.",
			CheckcodeCheckLimit:   "Check code check request limit.",
			CheckcodeMismatch:     "Check code mismatch.",
			CheckcodeNotFound:     "Check code not found.",
			CheckcodeCheckErr:     "Failed to verify check code.",
		},
		MTokenErr: {
			_TokenErrMin:         "Basic token error information.",
			TokenInitPubKeyErr:   "Public key initialization failed.",
			TokenInitPriKeyErr:   "Private key initialization failed.",
			TokenSaveErr:         "Token save failed.",
			TokenIsEmpty:         "The token is empty.",
			TokenDecryptErr:      "Token decrypting failure.",
			TokenExpired:         "The token has expired.",
			TokenGetErr:          "Failed to obtain token.",
			TokenInvalid:         "Invalid token.",
			TokenEcryptErr:       "Token encrypting failure.",
			TokenInvalidVersion:  "Invalid token version.",
			TokenBase64DecodeErr: "Base64 decoding failed.",
			TokenSignErr:         "Token signature failure.",
			TokenSignVerifyErr:   "Token signature check failed.",
			TokenCreateErr:       "Failed to create token.",
		},
		MCaptchaErr: {
			_CaptchaErrMin:       "Basic captcha error information.",
			CaptchaGetErr:        "Failed to generate graphic verification code picture.",
			CaptchaWriteImageErr: "Failed to obtain graphical verification code.",
			CaptchaWriteAudioErr: "Generate graphics validation code audio failure.",
			CaptchaMismatch:      "The graphics verification code does not match.",
			CaptchaRequired:      "Graphical verification code is required.",
		},
		MPasswordErr: {
			_PasswordErrMin:    "Basic password error information.",
			PasswordEncryptErr: "Password encryption failed.",
			PasswordVerifyErr:  "Failed to verify password.",
			PasswordLevelErr:   "Invalid password level.",
			PasswordInvalid:    "Invalid password.",
		},
		MMysqlErr: {
			_MysqlErrMin:          "Basic mysql error information.",
			MysqlConnectErr:       "Database connection failed.",
			MysqlSelectErr:        "Query data failed.",
			MysqlUpdateErr:        "Failed to update data.",
			MysqlInsertErr:        "Failed to insert data.",
			MysqlDeleteErr:        "Failed to delete data.",
			MysqlRowAffectErr:     "Failed to get the row number affected.",
			MysqlRowScanErr:       "Failed to Scan data into model.",
			MysqlLastInsertErr:    "Failed to get the row number last inserted.",
			MysqlDuplicateErr:     "Unique key dumplicate.",
			MysqlInvalidPageNum:   "Invalid page number.",
			MysqlNoEnoughModelBuf: "No enough model buffer size.",
		},
		MRedisErr: {
			_RedisErrMin:           "Basic redis error information.",
			RedisValueMarshalErr:   "JSON generation failed.",
			RedisValueUnmarshalErr: "JSON parsing failed.",
			RedisSetErr:            "Saving information to redis failed.",
			RedisGetErr:            "Failed to obtain information from redis.",
			RedisDelErr:            "Deleting redis data failed.",
		},
		MSmsErr: {
			_SmsErrMin:        "Basic sms error information.",
			SmsConnectErr:     "Failed to connect SMS service.",
			SmsInvalidPurpose: "Invalid purpose.",
			SmsSendErr:        "SMS failed.",
			SmsReadResponse:   "Failed to read response content.",
			SmsDecodeResponse: "Response content decoding failure.",
		},
		MParamsErr: {
			_ParamsErrMin:        "Basic params error information.",
			ParamsInvalidPhone:   "Invalid cell number.",
			ParamsInvalidSource:  "Invalid source.",
			ParamsInvalidPurpose: "Invalid purpose.",
			ParamsBindErr:        "Parameter parsing error.",
		},
		MAuthErr: {
			_AuthErrMin:       "Basic auth error information.",
			AuthTokenEmpty:    "Auth token is empty.",
			AuthGetCurrentErr: "Failed to obtain current authentication information.",
		},
		MUtilsErr: {
			_UtilsErrMin: "Basic utils error information.",
		},
		MOthersErr: {
			_OthersErrMin: "Basic others error information.",
		},
		MPhoneErr: {
			_PhoneErrMin:            "Basic phone error information.",
			PhoneEmpty:              "The cell phone number is empty.",
			PhoneFormatErr:          "Phone number parsing failed.",
			PhoneInvalidCountryCode: "Invalid country code.",
			PhoneUnknownRegion:      "Unknown region.",
			PhoneRegionMismatch:     "Region mismatch.",
			PhoneParseNumberErr:     "Failed to parse phone number.",
		},
		MColumnErr: {
			_ColumnErrMin:       "Basic column error information.",
			ColumnGetErr:        "Failed to obtain column information.",
			ColumnAlreadyExist:  "Column already exists.",
			ColumnNotExist:      "Column does not exist.",
			ColumnCreateErr:     "Failed to create column.",
			ColumnGetOnePageErr: "Failed to obtain one page column.",
			ColumnGetAllErr:     "Failed to obtain all column information.",
		},
		MArticleErr: {
			_ArticleErrMin:      "Basic article error information.",
			ArticleGetErr:       "Failed to obtain article information.",
			ArticleAlreadyExist: "Article already exists.",
			ArticleNotExist:     "Article does not exist.",
			ArticleCreateErr:    "Failed to create article.",
		},
	}
}

func initSubModuleErrs() {
	SubModuleErrs = [_SubModuleErrMax]string{
		_SubModuleErrMin: "Invalid sub module", // 占位
		MConfigErr:       "config",             // 配置错误
		MLogErr:          "log",                // 日志错误
		MUserErr:         "user model",         // 用户model错误
		MLoginErr:        "login",              // 登录限制错误
		MCheckcodeErr:    "checkcode",          // 短信验证码检查错误
		MTokenErr:        "token",              // Access Token错误
		MCaptchaErr:      "captcha",            // 图形验证码错误
		MPasswordErr:     "password",           // 密码错误
		MMysqlErr:        "mysql",              // Mysql服务错误
		MRedisErr:        "redis",              // redis服务错误
		MSmsErr:          "sms",                // 短信服务错误
		MParamsErr:       "params",             // 参数错误
		MAuthErr:         "auth",               // 认证错误
		MUtilsErr:        "utils",              // utils错误
		MOthersErr:       "others",             // 其他杂项错误
		MPhoneErr:        "phone",              // 手机号码错误
		MColumnErr:       "column",             // 栏目错误
		MArticleErr:      "article",            // 文章错误
	}
}

func initInterfaceErr() {
	InterfaceErrs = [_InterfaceErrMax]string{
		_InterfaceErrMin:        "Same with the top stack.",
		ICaptchaCreate:          "post /captcha/",
		ICaptchaGetImage:        "get /captcha/:captcha_id",
		ICheckcodeSend:          "post /checkcode/send",
		ICheckcodeCheck:         "post /checkcode/verify",
		IDebugPing:              "get /debug/ping",
		IDebugCaptchaGetValue:   "get /debug/captcha/:captcha_id",
		IDebugCheckcodeGetValue: "get /debug/checkcode?*",
		IUserSignup:             "get /u/signup",
		IUserLogin:              "post /u/login",
		IUserLogout:             "post /u/logout",
		IColumnCreate:           "post /column/",
		IArticleCreate:          "post /article/",
		IColumnGetAll:           "get /column/",
	}
}
