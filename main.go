package main

import (
	"flag"
	"math/rand"
	"runtime"
	"time"

	"wangqingang/cunxun/common"
	"wangqingang/cunxun/db"
	"wangqingang/cunxun/handler"
	"wangqingang/cunxun/model/captcha"
)

func main() {
	configPath := flag.String("config", "", "config file's path")
	flag.Parse()

	common.InitConfig(*configPath)
	if common.Config.Gomaxprocs >= 1 {
		runtime.GOMAXPROCS(common.Config.Gomaxprocs)
	}

	db.InitRedis(common.Config.Redis)
	db.InitMysql(common.Config.Mysql)
	captcha.InitCaptcha(common.Config.Captcha.TTL.D())

	// TODO: init log

	rand.Seed(time.Now().UTC().UnixNano())

	router := handler.ServerEngine()
	router.Run(common.Config.Listen)
}
