package test

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"wangqingang/cunxun/captcha"
	"wangqingang/cunxun/common"
	"wangqingang/cunxun/db"
	"wangqingang/cunxun/token/token_lib"
)

const (
	testUserSQLPath          = "../sql/user.sql"
	testArticleSQLPath       = "../sql/article.sql"
	testColumnSQLPath        = "../sql/column.sql"
	testCommentSQLPath       = "../sql/comment.sql"
	testCommentListSQLPath   = "../sql/commentlistview.sql"
	testArticleDetailSQLPath = "../sql/articledetailview.sql"
	testArticleListSQLPath   = "../sql/articlelistview.sql"
	testPrivateKeyPath       = "../conf/ecdsa_prv.pem"
	testPublicKeyPath        = "../conf/ecdsa_pub.pem"
)

func init() {
	common.InitConfig("../conf/config.toml")

	db.InitRedis(common.Config.Redis)
	db.InitMysql(common.Config.Mysql)

	captcha.InitCaptcha(common.Config.Captcha.TTL.D())

	token_lib.InitKeyPem(testPublicKeyPath, testPrivateKeyPath)
	rand.Seed(time.Now().UTC().UnixNano())
}

func initMySQLView(t *testing.T, sqlPath, viewName string) {
	assert := assert.New(t)

	f, err := ioutil.ReadFile(sqlPath)
	assert.Nil(err)
	assert.NotNil(f)

	_, err = db.Mysql.Exec(fmt.Sprintf("DROP VIEW IF EXISTS `%s`;", viewName))
	assert.Nil(err)

	_, err = db.Mysql.Exec(string(f))
	assert.Nil(err)
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
func InitTestCaseEnv(t *testing.T) {
	initRedis(t)
	initMySQLTable(t, testUserSQLPath, "user")
	initMySQLTable(t, testColumnSQLPath, "column")
	initMySQLTable(t, testArticleSQLPath, "article")
	initMySQLTable(t, testCommentSQLPath, "comment")
	initMySQLView(t, testCommentListSQLPath, "commentlistview")
	initMySQLView(t, testArticleDetailSQLPath, "articledetailview")
	initMySQLView(t, testArticleListSQLPath, "articlelistview")
}
