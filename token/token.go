package token

import (
	"encoding/json"
	"fmt"
	"time"

	"gopkg.in/redis.v4"

	"wangqingang/cunxun/common"
	"wangqingang/cunxun/db"
	e "wangqingang/cunxun/error"
	"wangqingang/cunxun/token/token_lib"
)

type TokenKey struct {
	UserId uint64 `json:"user_id"`
	Source string `json:"source"`
}

type Token struct {
	TokenKey
	Token            string        `json:"token"`
	CreatedTimestamp time.Time     `json:"create_ts"`
	TTL              time.Duration `json:"ttl"`
}

func (k *TokenKey) GetTokenKey() string {
	return fmt.Sprintf("%s:%d:%s", common.ModuleName, k.UserId, k.Source)
}

func (t *Token) Save() error {

	value, err := json.Marshal(t)
	if err != nil {
		return e.SP(e.MRedisErr, e.RedisValueMarshalErr, err)
	}

	key := t.GetTokenKey()
	expire := t.CreatedTimestamp.Add(t.TTL).Sub(time.Now())
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

func (t *TokenKey) Clean() error {
	key := t.GetTokenKey()
	if err := db.Redis.Del(key).Err(); err != nil {
		return e.SP(e.MRedisErr, e.RedisDelErr, err)
	}
	return nil
}

func (k *TokenKey) CreateToken(token string, ttl time.Duration) (*Token, error) {
	t := &Token{
		TokenKey:         *k,
		Token:            token,
		CreatedTimestamp: time.Now(),
		TTL:              ttl,
	}
	err := t.Save()
	if err != nil {
		return nil, e.SP(e.MTokenErr, e.TokenSaveErr, err)
	}
	return t, nil
}

func (k *TokenKey) GetToken() (*Token, error) {
	key := k.GetTokenKey()
	bs, err := db.Redis.Get(key).Bytes()
	if err == redis.Nil {
		// key不存在则返回redis.Nil
		return nil, e.SD(e.MRedisErr, e.RedisKeyNotExist, key)
	} else if err != nil {
		return nil, e.SP(e.MRedisErr, e.RedisGetErr, err)
	}
	t := &Token{}
	if err = json.Unmarshal(bs, t); err != nil {
		return nil, e.SP(e.MRedisErr, e.RedisValueUnmarshalErr, err)
	}
	return t, nil
}

func TokenCreateAndStore(userID uint64, userRole int, source string, ttl time.Duration) (string, error) {
	issueTime := time.Now()

	// payload中ttl单位为分钟
	accessToken, err := token_lib.Encrypt(common.Config.Token.TokenLibVersion, &token_lib.Payload{
		IssueTime:   uint32(uint64(issueTime.Unix())),
		TTL:         uint16(ttl.Minutes()),
		Role:        uint16(userRole),
		UserId:      userID,
		LoginSource: source,
	})
	if err != nil {
		return "", e.SP(e.MTokenErr, e.TokenCreateErr, err)
	}

	tokenKey := TokenKey{
		UserId: userID,
		Source: source,
	}
	_, err = tokenKey.CreateToken(accessToken, common.Config.Token.AccessTokenTTL.D())
	if err != nil {
		return "", e.SP(e.MTokenErr, e.TokenCreateErr, err)
	}

	return accessToken, nil
}

func TokenClean(userID uint64, source string) (*Token, error) {
	key := &TokenKey{
		UserId: userID,
		Source: source,
	}
	value, err := key.GetToken()
	if err != nil || value == nil {
		return nil, err
	}
	return value, value.Clean()
}
