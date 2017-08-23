package model

import (
	"encoding/json"
	"fmt"
	"time"

	redis "gopkg.in/redis.v4"

	"git.meiqia.com/business_platform/account/common"
	"git.meiqia.com/business_platform/account/db"
)

type VerifyKey struct {
	Phone   string `json:"phone"`
	Purpose string `json:"purpose"`
	Source  string `json:"source"`
}

type Verify struct {
	VerifyKey
	SendTimes        int    `json:"send_times"`
	CheckTimes       int    `json:"check_times"`
	VerifyCode       string `json:"verify_code"`
	CreatedTimestamp int64  `json:"created_ts"`
}

func (k *VerifyKey) GetRedisKey() string {
	return fmt.Sprintf("%s:%s:%s:%s", common.ModuleName, k.Phone, k.Purpose, k.Source)
}

func (c *Verify) Verify(verifyCode string) (bool, error) {
	if c.VerifyCode == verifyCode {
		return true, nil
	}

	c.CheckTimes++
	return false, c.Save()
}

func (c *Verify) Save() error {

	value, err := json.Marshal(c)
	if err != nil {
		return err
	}

	key := c.GetRedisKey()
	expire := c.CreatedTimestamp + int64(common.Config.Verify.TTL.Seconds()) - time.Now().Unix()
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

func (c *Verify) Clean() error {
	key := c.GetRedisKey()
	if err := db.Redis.Del(key).Err(); err != nil {
		return err
	}

	return nil
}

func (k *VerifyKey) CreateVerify(verifyCode string) (*Verify, error) {
	verify := &Verify{
		VerifyKey:        *k,
		SendTimes:        0,
		CheckTimes:       0,
		VerifyCode:       verifyCode,
		CreatedTimestamp: time.Now().Unix(),
	}

	err := verify.Save()
	if err != nil {
		return nil, err
	}
	return verify, nil
}

func (k *VerifyKey) GetVerify() (*Verify, error) {
	key := k.GetRedisKey()
	bs, err := db.Redis.Get(key).Bytes()
	if err == redis.Nil {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	c := &Verify{}
	if err = json.Unmarshal(bs, c); err != nil {
		return nil, err
	}

	return c, nil
}
