package main

import (
	"flag"
	"runtime"
	"math/rand"
	"time"
	"net/http"
	"fmt"

	"wangqingang/cunxun/common"
	"wangqingang/cunxun/db"
	"wangqingang/cunxun/model"
	"wangqingang/cunxun/handler"
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

	if err := http.ListenAndServe(common.Config.Listen, handler.AccountEngine()); err != nil {
		fmt.Println(err.Error())
	}
}
