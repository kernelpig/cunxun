package phone

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"wangqingang/cunxun/test"
)

func TestValidPhone(t *testing.T) {
	assert := assert.New(t)

	phoneStr := test.GenFakePhone()
	phone, err := ValidPhone(phoneStr)

	assert.Nil(err)
	t.Log(phone)
}
