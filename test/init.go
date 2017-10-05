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
	"wangqingang/cunxun/id"
	"wangqingang/cunxun/token/token_lib"
)

const (
	testTableUser       = "../sql/table/user.sql"
	testTableArticle    = "../sql/table/article.sql"
	testTableColumn     = "../sql/table/column.sql"
	testTableComment    = "../sql/table/comment.sql"
	testTableCarpooling = "../sql/table/carpooling.sql"
)

const (
	testViewCommentDetail     = "../sql/view/commentdetailview.sql"
	testViewArticleDetail     = "../sql/view/articledetailview.sql"
	testViewColumnDetail      = "../sql/view/columndetailview.sql"
	testViewCarpoolingDeatail = "../sql/view/carpoolingdetailview.sql"
)

const (
	testPrivateKeyPath = "../conf/ecdsa_prv.pem"
	testPublicKeyPath  = "../conf/ecdsa_pub.pem"
)

func init() {
	common.InitConfig("../conf/config.dev.toml")

	if err := id.InitIdGenerator(); err != nil {
		panic(err)
	}
	if err := db.InitRedis(common.Config.Redis); err != nil {
		panic(err)
	}
	if err := db.InitMysql(common.Config.Mysql); err != nil {
		panic(err)
	}
	if err := captcha.InitCaptcha(common.Config.Captcha.TTL.D()); err != nil {
		panic(err)
	}
	if err := token_lib.InitKeyPem(testPublicKeyPath, testPrivateKeyPath); err != nil {
		panic(err)
	}
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

	initMySQLTable(t, testTableUser, "user")
	initMySQLTable(t, testTableColumn, "column")
	initMySQLTable(t, testTableArticle, "article")
	initMySQLTable(t, testTableComment, "comment")
	initMySQLTable(t, testTableCarpooling, "carpooling")

	initMySQLView(t, testViewCommentDetail, "commentdetailview")
	initMySQLView(t, testViewArticleDetail, "articledetailview")
	initMySQLView(t, testViewColumnDetail, "columndetailview")
	initMySQLView(t, testViewCarpoolingDeatail, "carpoolingdetailview")
}
