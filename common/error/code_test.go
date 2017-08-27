package error

import (
	"github.com/juju/errors"
	"testing"
)

func TestC(t *testing.T) {
	code := C(SCunxun, IUserSignup, MPasswordErr, CaptchaMismatch)
	t.Logf("0x%08x, 0x%08x, 0x%08x, 0x%08x", SCunxun, IUserSignup, MPasswordErr, CaptchaMismatch)
	t.Logf("0x%08x", code)
	t.Log(MDE(code, errors.New("error detail")))
}
