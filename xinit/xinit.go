package main

import (
	"flag"
	"math/rand"
	"time"

	"wangqingang/cunxun/avatar"
	"wangqingang/cunxun/common"
	"wangqingang/cunxun/db"
	"wangqingang/cunxun/script"
)

func main() {
	configPath := flag.String("config", "", "config file's path")
	flag.Parse()

	common.InitConfig(*configPath)

	avatar.InitAvatar(common.Config.User.DefaultAvatarDir, common.Config.User.DefaultAvatarFile)
	db.InitRedis(common.Config.Redis)
	db.InitMysql(common.Config.Mysql)

	// 1. 创建超级用户
	if err := script.UserCreateSuperAdmin(common.Config.User); err != nil {
		panic(err)
	}

	rand.Seed(time.Now().UTC().UnixNano())
}
