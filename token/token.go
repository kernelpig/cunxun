package token

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"time"

	"gopkg.in/redis.v4"

	"wangqingang/cunxun/common"
	"wangqingang/cunxun/db"
	"wangqingang/cunxun/token/token_lib"
)

func InitKeyPem(publicKeyFile, privateKeyFile string) {

	publickKeyPem, err := ioutil.ReadFile(publicKeyFile)
	if err != nil {
		panic(err)
	}

	privateKeyPem, err := ioutil.ReadFile(privateKeyFile)
	if err != nil {
		panic(err)
	}

	token_lib.InitPrivateKey(privateKeyPem)
	token_lib.InitPublicKeys(publickKeyPem)
}

type TokenKey struct {
	AccountId string `json:"account_id"`
	Source    string `json:"source"`
}

type Token struct {
	TokenKey
	Token            string `json:"token"`
	CreatedTimestamp int64  `json:"create_ts"`
	TTL              int64  `json:"ttl"` // 单位为秒
}

func (k *TokenKey) GetTokenKey() string {
	return fmt.Sprintf("%s:%s:%s", common.ModuleName, k.AccountId, k.Source)
}

func (t *Token) Save() error {

	value, err := json.Marshal(t)
	if err != nil {
		return err
	}

	key := t.GetTokenKey()
	expire := t.CreatedTimestamp + t.TTL - time.Now().Unix()
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

func (t *TokenKey) Clean() error {
	key := t.GetTokenKey()
	if err := db.Redis.Del(key).Err(); err != nil {
		return err
	}
	return nil
}

func (k *TokenKey) CreateToken(token string, ttl time.Duration) (*Token, error) {
	t := &Token{
		TokenKey:         *k,
		Token:            token,
		CreatedTimestamp: time.Now().Unix(),
		TTL:              int64(ttl.Seconds()),
	}
	err := t.Save()
	if err != nil {
		return nil, err
	}
	return t, nil
}

func (k *TokenKey) GetToken() (*Token, error) {
	key := k.GetTokenKey()
	bs, err := db.Redis.Get(key).Bytes()
	if err == redis.Nil {
		// key不存在则返回redis.Nil
		return nil, nil
	} else if err != nil {
		return nil, err
	}
	t := &Token{}
	if err = json.Unmarshal(bs, t); err != nil {
		return nil, err
	}
	return t, nil
}

func RemoveAllTokenOfAccount(accountId string) {
	for _, s := range common.SourceRange {
		tokenKey := TokenKey{AccountId: accountId, Source: s}
		tokenKey.Clean()
	}
}
