package model

/*
* 登录处理实体：
* 1. 登录次数限制，只有登录密码校验失败时才增加次数；
* 2. 登录超过5次后需要图形验证码；
 */

import (
	"encoding/json"
	"fmt"
	"time"

	"gopkg.in/redis.v4"

	"wangqingang/cunxun/common"
	"wangqingang/cunxun/db"
)

type LoginKey struct {
	Phone   string `json:"phone"`
	Purpose string `json:"purpose"`
	Source  string `json:"source"`
}

type Login struct {
	LoginKey
	RequestTimes     int   `json:"request_times"` // 生存周期内请求次数记录
	CreatedTimestamp int64 `json:"create_ts"`
	TTL              int64 `json:"ttl"`
}

func (k *LoginKey) GetLoginKey() string {
	return fmt.Sprintf("%s:%s:%s:%s", common.ModuleName, k.Phone, k.Purpose, k.Source)
}

func (login *Login) GetLeftTimes() int {
	if login.RequestTimes < common.Config.Login.MaxRequestTimes {
		return common.Config.Login.MaxRequestTimes - login.RequestTimes
	}
	return 0
}

// 保存限制信息
func (login *Login) Save() error {

	value, err := json.Marshal(login)
	if err != nil {
		return err
	}

	key := login.GetLoginKey()
	expire := login.CreatedTimestamp + login.TTL - time.Now().Unix()
	if expire <= 0 {
		db.Redis.Del(key) // 超时删除
		return nil
	}

	err = db.Redis.Set(key, value, time.Duration(expire)*time.Second).Err()
	if err != nil {
		return err
	}

	return nil
}

// 清除限制信息
func (login *Login) Clean() error {
	key := login.GetLoginKey()
	if err := db.Redis.Del(key).Err(); err != nil {
		return err
	}
	return nil
}

// 创建限制信息
func (k *LoginKey) CreateLogin(ttl time.Duration) (*Login, error) {
	login := &Login{
		LoginKey:         *k,
		RequestTimes:     0,
		CreatedTimestamp: time.Now().Unix(),
		TTL:              int64(ttl.Seconds()),
	}
	err := login.Save()
	if err != nil {
		return nil, err
	}
	return login, nil
}

// 获取限制信息
func (k *LoginKey) GetLogin() (*Login, error) {
	key := k.GetLoginKey()
	bs, err := db.Redis.Get(key).Bytes()
	if err == redis.Nil {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	login := &Login{}
	if err = json.Unmarshal(bs, login); err != nil {
		return nil, err
	}
	return login, nil
}
