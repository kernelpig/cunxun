package db

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"

	"wangqingang/cunxun/common"
	e "wangqingang/cunxun/error"
)

// Mysql 全局的MySQL连接池
var Mysql *sql.DB

// InitMysql 根据配置信息连接MySQL并设置参数
func InitMysql(config *common.MysqlConfig) {
	pool, err := sql.Open("mysql", config.Dsn)
	if err != nil {
		panic(e.SP(e.MMysqlErr, e.MysqlConnectErr, err))
	}

	if config.MaxIdle > 0 {
		pool.SetMaxIdleConns(config.MaxIdle)
	}
	if config.MaxOpen > 0 {
		pool.SetMaxOpenConns(config.MaxOpen)
	}

	Mysql = pool
}
