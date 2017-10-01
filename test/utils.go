package test

import (
	"fmt"
	"math"
	"math/rand"
	"strings"

	"github.com/satori/go.uuid"

	"wangqingang/cunxun/common"
)

const (
	TestWebSource     = common.WebSource
	TestSignupPurpose = common.SignupPurpose
)

func GenFakePhone() string {
	max := int(math.Pow10(8))
	return fmt.Sprintf("86 186%08d", rand.Intn(max))
}

func GenFakeSource() string {
	index := rand.Intn(len(common.SourceRange))
	return common.SourceRange[index]
}

func GenFakePurpose() string {
	index := rand.Intn(len(common.PurposeRange))
	return common.PurposeRange[index]
}

func GenRandString() string {
	return strings.Replace(uuid.NewV4().String(), "-", "", -1)
}

// 注意: 不包括0
func GenRandInt(max int) int {
	if max <= 1 {
		return 1
	}
	return rand.Intn(max-1) + 1
}

func GenRandId(max int) uint64 {
	if max <= 1 {
		return 1
	}
	return uint64(rand.Intn(max-1) + 1)
}

func GenFakePassword() string {
	return GenRandString()[16:]
}
