package test

import (
	"fmt"
	"math"
	"math/rand"

	"time"
	"wangqingang/cunxun/common"
)

func GenFakePhone() string {
	max := int(math.Pow10(9))
	return fmt.Sprintf("+86 186%08d", rand.Intn(max))
}

func GenFakeSource() string {
	index := rand.Intn(len(common.SourceRange))
	return common.SourceRange[index]
}

func GenFakePurpose() string {
	index := rand.Intn(len(common.PurposeRange))
	return common.PurposeRange[index]
}

func GenFakeCheckcode() string {
	max := int(math.Pow10(7))
	return fmt.Sprintf("%06d", rand.Intn(max))
}

func GenFakeTime(duration time.Duration) func() time.Time {
	n := time.Now().Add(duration)
	return func() time.Time {
		return n
	}
}
