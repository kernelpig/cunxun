package phone

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"wangqingang/cunxun/test"
)

func TestValidPhone(t *testing.T) {
	assert := assert.New(t)
	err := ValidPhone(test.GenFakePhone())
	assert.Nil(err)
}
