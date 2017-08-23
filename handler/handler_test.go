package handler

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gavv/httpexpect"
	"github.com/stretchr/testify/assert"

	"wangqingang/cunxun/common"
	"wangqingang/cunxun/db"
	"wangqingang/cunxun/model/captcha"
)

const (
	testUserSQLPath    = "../sql/user.sql"
	testArticleSQLPath = "../sql/article.sql"
	testColumnSQLPath  = "../sql/column.sql"
)

func init() {
	common.InitConfig("../conf/config.dev.toml")

	db.InitRedis(common.Config.Redis)
	db.InitMysql(common.Config.Mysql)

	captcha.InitCaptcha(common.Config.Captcha.TTL.D())
	rand.Seed(time.Now().UTC().UnixNano())

	ServerEngine()
}

// 主测试函数
func TestHandlers(t *testing.T) {
	server := httptest.NewServer(ServerEngine())
	defer server.Close()

	e := httpexpect.New(t, server.URL)

	testBaseHandler(t, e)
	testDebugHandler(t, e)
	testInternelHandler(t, e)
	testExceptions(t, e)
}

func initMySQLTable(t *testing.T, sqlPath, tableName string) {
	assert := assert.New(t)

	f, err := ioutil.ReadFile(sqlPath)
	assert.Nil(err)
	assert.NotNil(f)

	_, err = db.Mysql.Exec(fmt.Sprintf("DROP TABLE IF EXISTS `%s`;", tableName))
	assert.Nil(err)

	_, err = db.Mysql.Exec(string(f))
	assert.Nil(err)
}

// 清空redis缓存
func initRedis(t *testing.T) {
	assert := assert.New(t)
	err := db.Redis.FlushAll().Err()
	assert.Nil(err)
}

// 初始化测试例环境
func initTestCaseEnv(t *testing.T) {
	initRedis(t)
	initMySQLTable(t, testUserSQLPath, "user")
	initMySQLTable(t, testColumnSQLPath, "column")
	initMySQLTable(t, testArticleSQLPath, "article")
}

// 接口基础功能测试
func testBaseHandler(t *testing.T, e *httpexpect.Expect) {
	testCreateCaptchaHandler(t, e)
	testGetCaptchaImageHandler(t, e)
}

// debug接口测试
func testDebugHandler(t *testing.T, e *httpexpect.Expect) {
	testDebugPingHandler(t, e)
	testDebugGetCaptchaValue(t, e)
}

// 内部接口测试
func testInternelHandler(t *testing.T, e *httpexpect.Expect) {

}

// 异常测试
func testExceptions(t *testing.T, e *httpexpect.Expect) {

}
