package main

import (
	"flag"
	"math/rand"
	"runtime"
	"time"

	"wangqingang/cunxun/captcha"
	"wangqingang/cunxun/common"
	"wangqingang/cunxun/db"
	"wangqingang/cunxun/handler"
	"wangqingang/cunxun/id"
	"wangqingang/cunxun/oss"
	"wangqingang/cunxun/script"
	"wangqingang/cunxun/token/token_lib"
)

const (
	cmdKeyConfig = "config"
	cmdKeyInit   = "init"
	cmdKeyStart  = "start"
)

const (
	cmdHelpConfig = "config file's path"
	cmdHelpInit   = "init database etc, default off"
	cmdHelpStart  = "start server, default on"
)

func cmdConfigHandler(config string) {
	if err := common.InitConfig(config); err != nil {
		panic(err)
	}
	if common.Config.Gomaxprocs >= 1 {
		runtime.GOMAXPROCS(common.Config.Gomaxprocs)
	}
	if err := id.InitIdGenerator(); err != nil {
		panic(err)
	}
	if err := db.InitRedis(common.Config.Redis); err != nil {
		panic(err)
	}
	if err := db.InitMysql(common.Config.Mysql); err != nil {
		panic(err)
	}
	if err := captcha.InitCaptcha(common.Config.Captcha.TTL.D()); err != nil {
		panic(err)
	}
	if err := token_lib.InitKeyPem(common.Config.Token.PublicKeyPath, common.Config.Token.PrivateKeyPath); err != nil {
		panic(err)
	}
	if err := oss.InitOss(); err != nil {
		panic(err)
	}
	rand.Seed(time.Now().UTC().UnixNano())
}

func cmdInitHandler(init bool) {
	if init {
		if err := script.InitScript(); err != nil {
			panic(err)
		}
	}
}

func cmdStartHandler(start bool) {
	if start {
		router := handler.ServerEngine()
		router.Run(common.Config.Listen)
	}
}

func main() {
	config := flag.String(cmdKeyConfig, "", cmdHelpConfig)
	init := flag.Bool(cmdKeyInit, false, cmdHelpInit)
	start := flag.Bool(cmdKeyStart, true, cmdHelpStart)

	flag.Parse()

	cmdConfigHandler(*config)
	cmdInitHandler(*init)
	cmdStartHandler(*start)
}
