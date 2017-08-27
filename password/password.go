package password

import (
	"strings"
	"unicode"

	"golang.org/x/crypto/bcrypt"

	e "wangqingang/cunxun/error"
)

const (
	// default cost in golang.org/x/crypto/bcrypt is 10, we use a higher cost
	cost = 12
)

// 密码级别
const (
	LevelIllegal = iota
	LevelWeak
	LevelNormal
	LevelStrong
)

// Encrypt 使用bcrypt生成密码的哈希
func Encrypt(passwd string) (string, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(passwd), cost)
	if err != nil {
		return "", e.SE(e.MPasswordErr, e.PasswordEncryptErr, err)
	}
	return string(hashed), nil
}

// Verify 验证密码与哈希是否匹配
func Verify(passwd string, hash string) error {
	if err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(passwd)); err != nil {
		return e.SE(e.MPasswordErr, e.PasswordVerifyErr, err)
	}
	return nil
}

func isSpecialChar(r rune) bool {
	return unicode.IsGraphic(r) && !unicode.IsNumber(r) && !unicode.IsLetter(r)
}

// 计算密码强度
func PasswordStrength(password string) (int, error) {
	length := len(password)
	hasNumber := strings.IndexFunc(password, unicode.IsNumber) >= 0
	hasLetter := strings.IndexFunc(password, unicode.IsLetter) >= 0
	hasSpecialChar := strings.IndexFunc(password, isSpecialChar) >= 0

	weakLevel := length >= 8 && length <= 12 && (hasNumber && hasLetter || hasNumber && hasSpecialChar || hasLetter && hasSpecialChar)
	normalLevel := length > 12 && (hasNumber && hasLetter || hasNumber && hasSpecialChar || hasLetter && hasSpecialChar) ||
		length >= 8 && length <= 12 && hasNumber && hasLetter && hasSpecialChar
	strongLevel := length > 12 && hasNumber && hasLetter && hasSpecialChar

	if strongLevel {
		return LevelStrong, nil
	} else if normalLevel {
		return LevelNormal, nil
	} else if weakLevel {
		return LevelWeak, nil
	}
	return LevelIllegal, e.SE(e.MPasswordErr, e.PasswordLevelErr, nil)
}
