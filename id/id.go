package id

import (
	"github.com/sony/sonyflake"

	e "wangqingang/cunxun/error"
)

var IdGenerator *sonyflake.Sonyflake

func InitIdGenerator() error {
	idGenerator := sonyflake.NewSonyflake(sonyflake.Settings{})
	if idGenerator == nil {
		return e.S(e.MIdGeneratorErr, e.IdGeneratorInitErr)
	}
	IdGenerator = idGenerator
	return nil
}

func Generate() (uint64, error) {
	return IdGenerator.NextID()
}
