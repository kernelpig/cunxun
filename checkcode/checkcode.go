package checkcode

import (
	"encoding/json"
	"fmt"
	"math"
	"math/rand"
	"time"

	"gopkg.in/redis.v4"

	"wangqingang/cunxun/common"
	"wangqingang/cunxun/db"
)

type CheckCodeKey struct {
	Phone   string `json:"phone"`
	Purpose string `json:"purpose"`
	Source  string `json:"source"`
}

type CheckCode struct {
	CheckCodeKey
	SendTimes        int    `json:"send_times"`
	CheckTimes       int    `json:"check_times"`
	Code             string `json:"verify_code"`
	CreatedTimestamp int64  `json:"created_ts"`
}

func (c *CheckCode) Check(code string) (bool, error) {
	if c.Code == code {
		return true, nil
	}

	c.CheckTimes++
	return false, c.Save()
}

func (c *CheckCode) Save() error {
	value, err := json.Marshal(c)
	if err != nil {
		return err
	}

	key := c.GetRedisKey()
	expire := c.CreatedTimestamp + int64(common.Config.Checkcode.TTL.Seconds()) - time.Now().Unix()
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

func (c *CheckCode) Clean() error {
	key := c.GetRedisKey()
	if err := db.Redis.Del(key).Err(); err != nil {
		return err
	}

	return nil
}

func (k *CheckCodeKey) GetRedisKey() string {
	return fmt.Sprintf("%s:%s:%s:%s", common.ModuleName, k.Phone, k.Purpose, k.Source)
}

func genCode() string {
	max := int64(math.Pow10(common.Config.Checkcode.DefaultLength + 1))
	codeFormat := fmt.Sprintf("%%0%dd", common.Config.Checkcode.DefaultLength)
	return fmt.Sprintf(codeFormat, rand.Int63n(max))
}

func (k *CheckCodeKey) CreateCheckCode() (*CheckCode, error) {
	checkcode := &CheckCode{
		CheckCodeKey:     *k,
		SendTimes:        0,
		CheckTimes:       0,
		Code:             genCode(),
		CreatedTimestamp: time.Now().Unix(),
	}

	err := checkcode.Save()
	if err != nil {
		return nil, err
	}
	return checkcode, nil
}

func (k *CheckCodeKey) GetCheckcode() (*CheckCode, error) {
	key := k.GetRedisKey()
	bs, err := db.Redis.Get(key).Bytes()

	if err == redis.Nil {
		// key不存在,不返回错误
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	c := &CheckCode{}
	if err = json.Unmarshal(bs, c); err != nil {
		return nil, err
	}
	return c, nil
}
