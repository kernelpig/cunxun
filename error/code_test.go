package error

import (
	"fmt"
	"testing"
)

func TestC(t *testing.T) {
	errDB := fmt.Errorf("invalid user name")
	errMysqlSelect := SE(MMysqlErr, MysqlSelectErr, errDB)
	t.Log(errMysqlSelect)

	errUserAlreadyExsits := SE(MUserErr, UserAlreadyExist, errMysqlSelect)
	t.Log(errUserAlreadyExsits)

	errInterface := IE(IUserSignup, MUserErr, UserAlreadyExist, errUserAlreadyExsits)
	t.Log(errInterface)
}
