package common

import (
	"time"

	"github.com/BurntSushi/toml"

	e "wangqingang/cunxun/error"
)

// Configs 全局配置信息
// Duration 类型支持的单位h-小时，m-分钟，s-秒
type Configs struct {
	ReleaseMode bool
	Listen      string
	Token       *TokenConfig
	Gomaxprocs  int
	Mysql       *MysqlConfig
	Redis       *RedisConfig
	Log         *LogConfig
	Captcha     *CaptchaConfig
	Checkcode   *CheckcodeConfig
	Login       *LoginConfig
	Sms         *SmsConfig
	User        *UserConfig
	Comment     *CommentConfig
	Oss         *OssConfig
}

type TokenConfig struct {
	TokenLibVersion int
	AccessTokenTTL  Duration
	PrivateKeyPath  string
	PublicKeyPath   string
}

type MysqlConfig struct {
	Dsn     string
	MaxIdle int
	MaxOpen int
}

type CaptchaConfig struct {
	DefaultWidth  int // image width
	DefaultHeight int // image heigth
	DefaultLength int // value length
	TTL           Duration
}

type CheckcodeConfig struct {
	DefaultLength int      // 短信验证码长度，默认为6位
	MaxSendTimes  int      // 周期内最大发送次数，默认为5次
	MaxCheckTimes int      // 周期内最大检验次数，默认为5次
	TTL           Duration // 检测周期，默认为10分钟
}

type LoginConfig struct {
	TTL             Duration
	MaxRequestTimes int // 周期内最大错误登录次数
	MaxCaptchaTImes int // 周期内N次后需要验证码
}

// RedisConfig Redis相关配置信息
type RedisConfig struct {
	Addr         string
	Password     string
	PoolSize     int
	DB           int
	DialTimeout  Duration
	ReadTimeout  Duration
	WriteTimeout Duration
}

// LogConfig 日志相关配置信息
type LogConfig struct {
	Level string
}

// Sms服务相关配置
type SmsConfig struct {
	AliAccessId     string
	AliAccessSecret string
}

// Oss服务相关配置
type OssConfig struct {
	Endpoint         string
	AliAccessId      string
	AliAccessSecret  string
	Bucket           string
	Domain           string
	PermitImageType  []string
	DialTimeout      Duration
	ReadWriteTimeout Duration
}

// 用户相关配置
type UserConfig struct {
	DefaultAvatarFile  string
	DefaultAvatarDir   string
	SuperAdminPhone    string
	SuperAdminPassword string
}

// 评论相关配置
type CommentConfig struct {
	DefaultMaxLength int
}

// Config 全局配置信息
var Config *Configs

// InitConfig 加载配置
func InitConfig(path string) {
	config, err := loadConfig(path)
	if err != nil {
		panic(e.SP(e.MConfigErr, e.ConfigLoadErr, err))
	}
	Config = config
}

func loadConfig(path string) (*Configs, error) {
	config := new(Configs)
	if _, err := toml.DecodeFile(path, config); err != nil {
		return nil, e.SP(e.MConfigErr, e.ConfigParseErr, err)
	}

	return config, nil
}

// Duration 配置中使用的时长
type Duration struct {
	time.Duration
}

// UnmarshalText 将字符串形式的时长信息转换为Duration类型
func (d *Duration) UnmarshalText(text []byte) (err error) {
	d.Duration, err = time.ParseDuration(string(text))
	if err != nil {
		return e.SP(e.MConfigErr, e.ConfigParseTimeErr, err)
	}
	return nil
}

// D 从Duration struct中取出time.Duration类型的值
func (d *Duration) D() time.Duration {
	return d.Duration
}
