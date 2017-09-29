package db

import (
	"gopkg.in/redis.v4"

	"wangqingang/cunxun/common"
)

// Redis 全局的Redis Client
var Redis *redis.Client

// InitRedis 根据配置初始化Redis Client
func InitRedis(config *common.RedisConfig) error {
	options := &redis.Options{
		Addr:         config.Addr,
		Password:     config.Password,
		PoolSize:     config.PoolSize,
		DB:           config.DB,
		DialTimeout:  config.DialTimeout.D(),
		ReadTimeout:  config.ReadTimeout.D(),
		WriteTimeout: config.WriteTimeout.D(),
	}
	Redis = redis.NewClient(options)
	return nil
}
