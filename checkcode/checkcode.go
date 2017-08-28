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
	e "wangqingang/cunxun/error"
)

type CheckCodeKey struct {
	Phone   string `json:"phone"`
	Purpose string `json:"purpose"`
	Source  string `json:"source"`
}

type CheckCode struct {
	CheckCodeKey
	SendTimes        int           `json:"send_times"`
	CheckTimes       int           `json:"check_times"`
	Code             string        `json:"verify_code"`
	CreatedTimestamp time.Time     `json:"created_ts"`
	TTL              time.Duration `json:"ttl"`
}

func (c *CheckCode) Check(code string) (bool, error) {
	if c.Code == code {
		return true, nil
	}

	c.CheckTimes++
	return false, e.SP(e.MCheckcodeErr, e.CheckcodeSaveErr, c.Save())
}

func (c *CheckCode) Save() error {
	value, err := json.Marshal(c)
	if err != nil {
		return e.SP(e.MRedisErr, e.RedisValueMarshalErr, err)
	}

	key := c.GetRedisKey()
	expire := c.CreatedTimestamp.Add(c.TTL).Sub(time.Now())
	if expire <= 0 {
		db.Redis.Del(key) // 超时删除
		return nil
	}

	err = db.Redis.Set(key, value, expire).Err()
	if err != nil {
		return e.SP(e.MRedisErr, e.RedisSetErr, err)
	}

	return nil
}

func (c *CheckCode) Clean() error {
	key := c.GetRedisKey()
	if err := db.Redis.Del(key).Err(); err != nil {
		return e.SP(e.MRedisErr, e.RedisDelErr, err)
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

func (k *CheckCodeKey) CreateCheckCode(ttl time.Duration) (*CheckCode, error) {
	checkcode := &CheckCode{
		CheckCodeKey:     *k,
		SendTimes:        0,
		CheckTimes:       0,
		Code:             genCode(),
		CreatedTimestamp: time.Now(),
		TTL:              ttl,
	}

	err := checkcode.Save()
	if err != nil {
		return nil, e.SP(e.MCheckcodeErr, e.CheckcodeSaveErr, err)
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
		return nil, e.SP(e.MRedisErr, e.RedisGetErr, err)
	}

	c := &CheckCode{}
	if err = json.Unmarshal(bs, c); err != nil {
		return nil, e.SP(e.MRedisErr, e.RedisValueUnmarshalErr, err)
	}
	return c, nil
}
