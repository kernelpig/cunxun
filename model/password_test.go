package model

import (
	"testing"

	"fmt"
	"wangqingang/cunxun/common"
	"wangqingang/cunxun/utils"
	"github.com/stretchr/testify/assert"
	"math/rand"
)

func Test(t *testing.T) {
	plaintext := "hello world"
	hash, err := Encrypt(plaintext)
	if err != nil {
		t.Error(err)
	}
	err = Verify(plaintext, hash)
	if err != nil {
		t.Error(err)
	}

	// 密码错误
	err = Verify("wrong password", hash)
	if err == nil {
		t.Error("should return error when password is wrong")
	}

}

func TestPasswordStrength(t *testing.T) {
	as := assert.New(t)

	as.Equal(LevelIllegal, PasswordStrength("12345678"))
	as.Equal(LevelWeak, PasswordStrength("12345678_"))
	as.Equal(LevelWeak, PasswordStrength("1234567890aa"))
	as.Equal(LevelNormal, PasswordStrength("123456789a_"))
	as.Equal(LevelNormal, PasswordStrength("1234567890a_"))
	as.Equal(LevelStrong, PasswordStrength("1234567890a_1"))
}

func BenchmarkEncrypt(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Encrypt("hello world")
	}
}

func BenchmarkDecrypt(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Verify("hello world", "$2a$12$YBgqbKhnmNC73wBnL.NDveRkwPp2qN.I4lCYuBYxQhswqja9vr93a")
	}
}

var (
	letterRunes = []rune("0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ!@#$%^&*()")
)

func randString(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func TestGenPassword(t *testing.T) {
	phones := []string{"+86 15677777700"}

	common.InitConfig("../../conf/config.dev.toml")
	for i := 0; i < len(phones); i++ {
		id, err := utils.GenID()
		if err != nil {
			fmt.Println(err)
		}
		randStr := randString(13)
		level := PasswordStrength(randStr)
		pass, err := Encrypt(randStr)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println("%s: %s: %s: %d: %s", id, phones[i], pass, level, randStr)
	}
}
