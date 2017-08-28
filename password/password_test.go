package password

import (
	"testing"

	"github.com/stretchr/testify/assert"
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

	level, err := PasswordStrength("12345678")
	as.NotNil(err)
	as.Equal(LevelIllegal, level)

	level, err = PasswordStrength("12345678_")
	as.Nil(err)
	as.Equal(LevelWeak, level)

	level, err = PasswordStrength("1234567890aa")
	as.Nil(err)
	as.Equal(LevelWeak, level)

	level, err = PasswordStrength("123456789a_")
	as.Nil(err)
	as.Equal(LevelNormal, level)

	level, err = PasswordStrength("1234567890a_1")
	as.Nil(err)
	as.Equal(LevelStrong, level)
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
