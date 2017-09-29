package script

import (
	"os"
	"path"
	"time"

	"io/ioutil"
	"path/filepath"

	"wangqingang/cunxun/common"
	"wangqingang/cunxun/db"
	e "wangqingang/cunxun/error"
	"wangqingang/cunxun/model"
	"wangqingang/cunxun/password"
)

const (
	createDbFormat = "create database ?"
	sqlFileType    = ".sql"
)

const (
	columnNews       = "资讯"
	columnBar        = "贴吧"
	columnRental     = "租房"
	columnCarpooling = "拼车"
)

func CreateDatabase() error {
	_, err := db.Mysql.Exec(createDbFormat, common.Config.Mysql.DatabaseName)
	if err != nil {
		return e.SP(e.MMysqlErr, e.MysqlCreateDatabase, err)
	}
	return nil
}

func CreateTables(sqlPathDir string) error {
	return filepath.Walk(sqlPathDir, func(p string, f os.FileInfo, err error) error {
		if err != nil {
			return e.SP(e.MMysqlErr, e.MysqlWalkSqlUnkownErr, err)
		} else if f.IsDir() {
			return e.S(e.MMysqlErr, e.MysqlWalkSqlNotSupportSubDir)
		} else if path.Ext(p) != sqlFileType {
			return e.S(e.MMysqlErr, e.MysqlWalkSqlUnsupportType)
		}
		bytes, err := ioutil.ReadFile(p)
		if err != nil {
			return e.S(e.MMysqlErr, e.MysqlWalkSqlReadFileErr)
		}
		if _, err := db.Mysql.Exec(string(bytes)); err != nil {
			return e.S(e.MMysqlErr, e.MysqlWalkSqlExecute)
		}
		return nil
	})
}

func CreateColumns(user *model.User) ([]*model.Column, error) {
	columns := []string{columnNews, columnBar, columnRental, columnCarpooling}
	models := make([]*model.Column, 0)
	for _, c := range columns {
		x := &model.Column{
			Name:       c,
			CreaterUid: user.ID,
			CreatedAt:  time.Now(),
		}
		x, err := model.CreateColumn(db.Mysql, x)
		if err != nil {
			return nil, err
		}
		models = append(models, x)
	}
	return models, nil
}

func CreateSuperAdmin() (*model.User, error) {
	conf := common.Config.User
	hashedPassword, err := password.Encrypt(conf.SuperAdminPassword)
	if err != nil {
		return nil, err
	}
	passwordLevel, err := password.PasswordStrength(conf.SuperAdminPassword)
	if err != nil {
		return nil, err
	}
	user := &model.User{
		Phone:          "86 " + conf.SuperAdminPhone,
		NickName:       "admin",
		HashedPassword: hashedPassword,
		PasswordLevel:  passwordLevel,
		RegisterSource: "web",
		Role:           model.UserRoleSuperAdmin,
	}
	user, err = model.CreateUser(db.Mysql, user)
	if err != nil {
		if msgErr, ok := err.(e.Message); ok && msgErr.Code.IsSubError(e.MUserErr, e.UserAlreadyExist) {
			user, err := model.GetUserByPhone(db.Mysql, user.Phone)
			if err != nil {
				return nil, err
			}
			return user, nil
		}
		return nil, err
	}
	return user, nil
}
