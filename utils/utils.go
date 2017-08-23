package utils

import (
	"database/sql"
	"math/rand"
	"net/http"

	"github.com/meiqia/chi/render"
)

// NullStringToString NullString类型转换为String类型
func NullStringToString(s sql.NullString) string {
	if s.Valid {
		return s.String
	}
	return ""
}

const (
	Digits   = "0123456789"
	Alphabet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	Ascii    = Alphabet + Digits + "~!@#$%^&*()-_+={}[]\\|<,>.?/\"';:`"
)

func RandomDigits(length int) string {
	return randomString(length, []byte(Digits))
}

func randomString(length int, base []byte) string {
	bytes := make([]byte, length)
	maxIndex := len(base)
	for i := 0; i < length; i++ {
		index := rand.Intn(maxIndex)
		bytes[i] = byte(base[index])
	}

	return string(bytes)
}

func BindJSON(r *http.Request, obj interface{}) error {
	if err := render.Bind(r.Body, obj); err != nil {
		return err
	}

	return structValidator.ValidateStruct(obj)
}
