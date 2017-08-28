package error

import (
	"fmt"
	"testing"
)

func TestC(t *testing.T) {
	errPassword := fmt.Errorf("invalid password")

	dbDetail := fmt.Sprintf("invalid user name")
	errMysqlSelect := SE(MMysqlErr, MysqlSelectErr, dbDetail, errPassword)
	t.Log(errMysqlSelect)

	userNameDetail := fmt.Sprintf("name: wangdalian")
	errUserAlreadyExsits := SE(MUserErr, UserAlreadyExist, userNameDetail, errMysqlSelect)
	t.Log(errUserAlreadyExsits)

	errInterface := IE(IUserSignup, MUserErr, UserAlreadyExist, userNameDetail, errUserAlreadyExsits)
	t.Log(errInterface)
}
