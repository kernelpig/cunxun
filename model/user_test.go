package model

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"git.meiqia.com/business_platform/account/common"
	"git.meiqia.com/business_platform/account/db"
	"git.meiqia.com/business_platform/account/utils/password"
)

const (
	testName           = "wangqingang"
	testPhone          = "186531293004"
	testEmail          = "403726259@qq.com"
	testNewEmail       = "18653193004@163.com"
	testPassword       = "123abcABC!@#"
	testInvalidId      = "32767"
	testRegisterSourde = common.WebSource
	testInvitatedCode  = "sdfsfdsfsdsdfdfsfsdf"
)

func init() {
	common.InitConfig("../conf/config.dev.toml")
	db.InitMysql(common.Config.Mysql)
	db.InitRedis(common.Config.Redis)
}

func flushTable() {
	db.Mysql.Exec("truncate account")
}

func TestCreateAccount(t *testing.T) {
	flushTable()

	assert := assert.New(t)

	account, err := GetAccountById(db.Mysql, testInvalidId)
	assert.Nil(err)
	assert.Nil(account)

	account, err = GetAccountByPhone(db.Mysql, testPhone)
	assert.Nil(err)
	assert.Nil(account)

	hashedPassword, err := password.Encrypt(testPassword)
	assert.Nil(err)
	id, err := CreateAccount(db.Mysql, testPhone, Account{
		ID:             testPhone,
		Phone:          testPhone,
		RegisterSource: testRegisterSourde,
		HashedPassword: hashedPassword,
	})
	assert.Nil(err)

	account, err = GetAccountById(db.Mysql, id)
	assert.Nil(err)
	assert.NotNil(account)
	assert.Equal(testPhone, account.Phone)
	assert.Nil(password.Verify(testPassword, account.HashedPassword))

	account, err = GetAccountByPhone(db.Mysql, testPhone)
	assert.Nil(err)
	assert.NotNil(account)

	_, err = UpdateAccountById(db.Mysql, id, map[string]interface{}{"hashed_password": "1234", "password_level": 2})
	assert.Nil(err)

	err = DeleteAccountById(db.Mysql, account.ID)
	assert.Nil(err)

	account, err = GetAccountById(db.Mysql, account.ID)
	assert.Nil(err)
	assert.Nil(account)

}
