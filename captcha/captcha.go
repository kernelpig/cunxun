package captcha

import (
	"bytes"
	"fmt"
	"time"

	"github.com/dchest/captcha"

	"wangqingang/cunxun/db"
	e "wangqingang/cunxun/error"
)

type CaptchaStore struct {
	ttl time.Duration
}

var Store CaptchaStore

const CaptcharKeyPrefix = "captcha"

func (c CaptchaStore) Set(id string, digits []byte) {
	key := fmt.Sprintf("%s:%s", CaptcharKeyPrefix, id)
	db.Redis.Set(key, digits, c.ttl)
}

func (c CaptchaStore) Get(id string, clear bool) (digits []byte) {
	key := fmt.Sprintf("%s:%s", CaptcharKeyPrefix, id)
	digits, err := db.Redis.Get(key).Bytes()
	if err != nil {
		return nil
	}

	if clear {
		db.Redis.Del(key)
	}
	return digits
}

func (c CaptchaStore) GetValue(id string, clear bool) string {
	bs := c.Get(id, clear)
	for i, _ := range bs {
		bs[i] += 48
	}
	return string(bs)
}

func InitCaptcha(ttl time.Duration) {
	if ttl > 0 {
		Store.ttl = ttl
	} else {
		Store.ttl = 0
	}
	captcha.SetCustomStore(captcha.Store(Store))
}

func VerifyCaptcha(id, value string) bool {
	defer func() {
		captcha.Reload(id)
	}()

	return captcha.VerifyString(id, value)
}

func GenCaptcha(length int) (id string) {
	return captcha.NewLen(length)
}

func GetCaptchaImage(id string, width, height int) ([]byte, error) {
	var buf bytes.Buffer
	if err := captcha.WriteImage(&buf, id, width, height); err != nil {
		return nil, e.SE(e.MCaptchaErr, e.CaptchaWriteImageErr, err)
	}
	return buf.Bytes(), nil
}

func GetCaptchaValue(id string, clear bool) string {
	return Store.GetValue(id, clear)
}

func ReloadCaptcha(id string) {
	captcha.Reload(id)
}
