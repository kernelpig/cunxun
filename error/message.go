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
			_ConfigErrMin:       "Basic configuration error information.",
			ConfigLoadErr:       "Failed to load profile.",
			ConfigParseErr:      "Parsing configuration file failed.",
			ConfigParseTimeErr:  "Parse time configuration field failed.",
			ConfigLoadAvatarErr: "Load avatar image file failed.",
		},
		MLogErr: {
			_LogErrMin:        "Basic log error information.",
			LogDumpRequestErr: "Dump request parameter failed.",
		},
		MUserErr: {
			_UserErrMin:         "Basic user error information.",
			UserGetErr:          "Failed to obtain user information.",
			UserAlreadyExist:    "User already exists.",
			UserNotExist:        "User does not exist.",
			UserCreateErr:       "Failed to create user.",
			UserAvatarDecodeErr: "Failed to decode user's avatar.",
			UserNotPermit:       "The user does not have permission.",
			UserGetListErr:      "Failed to get user list.",
			UserUpdateErr:       "Failed to update user information.",
			UserDeleteErr:       "Failed to delete user information.",
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
			_TokenErrMin:            "Basic token error information.",
			TokenInitPubKeyErr:      "Public key initialization failed.",
			TokenInitPriKeyErr:      "Private key initialization failed.",
			TokenSaveErr:            "Token save failed.",
			TokenIsEmpty:            "The token is empty.",
			TokenDecryptErr:         "Token decrypting failure.",
			TokenExpired:            "The token has expired.",
			TokenGetErr:             "Failed to obtain token.",
			TokenInvalid:            "Invalid token.",
			TokenEcryptErr:          "Token encrypting failure.",
			TokenInvalidVersion:     "Invalid token version.",
			TokenBase64DecodeErr:    "Base64 decoding failed.",
			TokenSignErr:            "Token signature failure.",
			TokenSignVerifyErr:      "Token signature check failed.",
			TokenCreateErr:          "Failed to create token.",
			TokenReadPubKeyFileErr:  "Failed to read public key file.",
			TokenReadPriKeyFileErr:  "Failed to read private key file.",
			TokenInvalidPubKey:      "Invalid public key.",
			TokenParsePubKey:        "Failed to parse public key.",
			TokenParsePriKey:        "Failed to parse private key.",
			TokenWritePriKeyFileErr: "Failed to write private key to file.",
			TokenWritePubKeyFileErr: "Failed to write publick key to file.",
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
			_MysqlErrMin:                 "Basic mysql error information.",
			MysqlConnectErr:              "Database connection failed.",
			MysqlSelectErr:               "Query data failed.",
			MysqlUpdateErr:               "Failed to update data.",
			MysqlInsertErr:               "Failed to insert data.",
			MysqlDeleteErr:               "Failed to delete data.",
			MysqlRowAffectErr:            "Failed to get the row number affected.",
			MysqlRowScanErr:              "Failed to Scan data into model.",
			MysqlLastInsertErr:           "Failed to get the row number last inserted.",
			MysqlDuplicateErr:            "Unique key dumplicate.",
			MysqlInvalidPageNum:          "Invalid page number.",
			MysqlNoEnoughModelBuf:        "No enough model buffer size.",
			MysqlInvalidOrderType:        "Invalid order type.",
			MysqlCreateDatabase:          "Failed to create database.",
			MysqlWalkSqlUnkownErr:        "Failed to walk sql.",
			MysqlWalkSqlNotSupportSubDir: "NOT support sub directory in walking sql directory.",
			MysqlWalkSqlUnsupportType:    "NOT support sql file type in walking sql directory.",
			MysqlWalkSqlReadFileErr:      "Failed to read sql file in walking sql directory.",
			MysqlWalkSqlExecute:          "Failed to execute sql file in walking sql directory.",
		},
		MRedisErr: {
			_RedisErrMin:           "Basic redis error information.",
			RedisValueMarshalErr:   "JSON generation failed.",
			RedisValueUnmarshalErr: "JSON parsing failed.",
			RedisSetErr:            "Saving information to redis failed.",
			RedisGetErr:            "Failed to obtain information from redis.",
			RedisDelErr:            "Deleting redis data failed.",
			RedisKeyNotExist:       "Key not exist in redis data.",
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
			_ParamsErrMin:             "Basic params error information.",
			ParamsInvalidPhone:        "Invalid cell number.",
			ParamsInvalidSource:       "Invalid source.",
			ParamsInvalidPurpose:      "Invalid purpose.",
			ParamsBindErr:             "Parameter parsing error.",
			ParamsInvalidColumnID:     "Invalid column id.",
			ParamsInvalidPageNum:      "Invalid page number.",
			ParamsInvalidPageSize:     "Invalid page size.",
			ParamsInvalidArticleID:    "Invalid article id.",
			ParamsInvalidCommentID:    "Invalid comment id.",
			ParamsInvalidOrderBy:      "Invalid order by.",
			ParamsInvalidUserId:       "Invalid user id.",
			ParamsCommentLengthLimit:  "Comment length limit.",
			ParamsInvalidMultiForm:    "Invalid multipart form data.",
			ParamsInvalidCarpoolingID: "Invalid carpooling id.",
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
			_ColumnErrMin:        "Basic column error information.",
			ColumnGetErr:         "Failed to obtain column information.",
			ColumnAlreadyExist:   "Column already exists.",
			ColumnNotExist:       "Column does not exist.",
			ColumnCreateErr:      "Failed to create column.",
			ColumnGetOnePageErr:  "Failed to obtain one page column.",
			ColumnGetAllErr:      "Failed to obtain all column information.",
			ColumnUpdateById:     "Failed to update column information by id.",
			ColumnDeleteErr:      "Failed to delete column information by id.",
			ColumnUpdateByIdSelf: "Failed to update column information by creater user.",
			ColumnDeleteById:     "Failed to delete column information by id.",
			ColumnDeleteByIdSelf: "Failed to delete column information by creater user.",
			ColumnUpdateErr:      "Failed to update column information.",
		},
		MArticleErr: {
			_ArticleErrMin:           "Basic article error information.",
			ArticleGetErr:            "Failed to obtain article information.",
			ArticleAlreadyExist:      "Article already exists.",
			ArticleNotExist:          "Article does not exist.",
			ArticleCreateErr:         "Failed to create article.",
			ArticleGetListErr:        "Failed to get article list.",
			ArticleDeleteErr:         "Failed to delete article list.",
			ArticleUpdateByIdErr:     "Failed to update article information by id.",
			ArticleUpdateByIdSelfErr: "Failed to update article information by created user.",
			ArticleDeleteByIdErr:     "Failed to delete article information by id.",
			ArticleDeleteByIdSelfErr: "Failed to delete article information by created user.",
		},
		MCommentErr: {
			_CommentErrMin:        "Basic comment error information.",
			CommentGetErr:         "Failed to obtain comment information.",
			CommentAlreadyExist:   "Comment already exists.",
			CommentNotExist:       "Comment does not exist.",
			CommentCreateErr:      "Failed to create comment.",
			CommentGetListErr:     "Failed to get comment list.",
			CommentDeleteErr:      "Failed to delete comment list.",
			CommentUpdateErr:      "Failed to update comment information.",
			CommentUpdateByIdSelf: "Failed to update comment information by creater user.",
			CommentDeleteByIdSelf: "Failed to delete comment information by creater user.",
		},
		MImageErr: {
			ImageGetErr:   "Failed to get image.",
			ImageNotFound: "Image file not found.",
			ImageSaveErr:  "Failed to save image file.",
			ImageReadErr:  "Failed to read image file.",
		},
		MOssErr: {
			OssClientInitErr:       "Failed to init oss client.",
			OssBucketGetErr:        "Failed to get oss bucket.",
			OssPutObjectByBytesErr: "Failed to put object by bytes.",
		},
		MCarpoolingErr: {
			CarpoolingGetErr:       "Failed to get carpooling information.",
			CarpoolingAlreadyExist: "The carpooling information is already exist.",
			CarpoolingCreateErr:    "Failed to create carpooling information.",
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
		MCommentErr:      "comment",            // 评论错误
		MImageErr:        "image",              // 图片错误
		MOssErr:          "oss",                // oss存储错误
		MCarpoolingErr:   "carpooling",         // 拼车错误
	}
}

func initInterfaceErr() {
	InterfaceErrs = [_InterfaceErrMax]string{
		_InterfaceErrMin:        "Same with the top stack.",
		ICaptchaCreate:          "post /api/captcha/",
		ICaptchaGetImage:        "get /api/captcha/:captcha_id",
		ICheckcodeSend:          "post /api/checkcode/send",
		ICheckcodeCheck:         "post /api/checkcode/verify",
		IDebugPing:              "get /api/debug/ping",
		IDebugCaptchaGetValue:   "get /api/debug/captcha/:captcha_id",
		IDebugCheckcodeGetValue: "get /api/debug/checkcode?*",
		IUserSignup:             "get /api/u/signup",
		IUserLogin:              "post /api/u/login",
		IUserLogout:             "post /api/u/logout",
		IColumnCreate:           "post /api/column/",
		IArticleCreate:          "post /api/article/",
		IColumnGetAll:           "get /api/column/",
		IArticleGetList:         "get /api/article/",
		IArticleGet:             "get /api/article/:article_id",
		ICommentCreate:          "post /api/comment/",
		ICommentGetList:         "get /api/comment/",
		ICommentGet:             "get /api/comment/:comment_id",
		IUserGetAvatar:          "get /api/u/:user_id/avatar",
		IColumnUpdateById:       "put /api/column/:column_id",
		IColumnDeleteById:       "delete /api/column/:column_id",
		IArticleUpdateById:      "put /api/article/:article_id",
		IArticleDeleteById:      "delete /api/article/:article_id",
		ICommentUpdateById:      "put /api/comment/:comment_id",
		ICommentDeleteById:      "delete /api/comment/:comment_id",
		IUserGetInfo:            "get /api/u/:user_id",
		IUserGetList:            "get /api/u",
		IUserCreate:             "post /api/u",
		IUserUpdateById:         "put /api/u/:user_id",
		IUserDeleteById:         "delete /api/u/:user_id",
		IImageCreate:            "post /api/image/",
		ICarpoolingCreate:       "post /api/carpooling/",
		ICarpoolingGetById:      "get /api/carpooling/:id",
	}
}
