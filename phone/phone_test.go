package phone

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"wangqingang/cunxun/test"
)

func TestValidPhone(t *testing.T) {
	assert := assert.New(t)
	err := ValidPhone(test.GenFakePhone())
	assert.Nil(err)
}
