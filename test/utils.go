package test

import (
	"fmt"
	"math"
	"math/rand"

	"wangqingang/cunxun/common"
)

const (
	TestWebSource     = common.WebSource
	TestSignupPurpose = common.SignupPurpose
)

func GenFakePhone() string {
	max := int(math.Pow10(8))
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
