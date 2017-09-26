package main

import (
	"flag"
	"math/rand"
	"runtime"
	"time"

	"wangqingang/cunxun/avatar"
	"wangqingang/cunxun/captcha"
	"wangqingang/cunxun/common"
	"wangqingang/cunxun/db"
	"wangqingang/cunxun/handler"
	"wangqingang/cunxun/token/token_lib"
)

func main() {
	configPath := flag.String("config", "", "config file's path")
	flag.Parse()

	common.InitConfig(*configPath)
	if common.Config.Gomaxprocs >= 1 {
		runtime.GOMAXPROCS(common.Config.Gomaxprocs)
	}

	avatar.InitAvatar(common.Config.User.DefaultAvatarDir, common.Config.User.DefaultAvatarFile)
	db.InitRedis(common.Config.Redis)
	db.InitMysql(common.Config.Mysql)
	captcha.InitCaptcha(common.Config.Captcha.TTL.D())
	token_lib.InitKeyPem(common.Config.Token.PublicKeyPath, common.Config.Token.PrivateKeyPath)

	// TODO: initial log
	rand.Seed(time.Now().UTC().UnixNano())

	router := handler.ServerEngine()
	router.Run(common.Config.Listen)
}
