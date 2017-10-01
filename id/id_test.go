package id

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGenerate(t *testing.T) {
	assert := assert.New(t)

	err := InitIdGenerator()
	assert.Nil(err)

	for i := 0; i < 10; i++ {
		id, err := Generate()
		assert.Nil(err)
		assert.NotZero(id)
		t.Log(id)
	}
}
