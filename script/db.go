package script

import (
	"fmt"
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
	sqlFileType = ".sql"
)

func CreateTables(sqlPathDir string) error {
	return filepath.Walk(sqlPathDir, func(p string, f os.FileInfo, err error) error {
		if err != nil {
			return e.SP(e.MMysqlErr, e.MysqlWalkSqlUnkownErr, err)
		} else if f.IsDir() && p == sqlPathDir {
			// 为当前目录, 啥都不做
			return nil
		} else if f.IsDir() && p != sqlPathDir {
			// 为子目录, 提示错误
			return e.S(e.MMysqlErr, e.MysqlWalkSqlNotSupportSubDir)
		} else if path.Ext(p) != sqlFileType {
			// 当前目录有非sql文件提示
			return e.S(e.MMysqlErr, e.MysqlWalkSqlUnsupportType)
		}
		// 正常的sql文件读取并执行
		bytes, err := ioutil.ReadFile(p)
		if err != nil {
			return e.SP(e.MMysqlErr, e.MysqlWalkSqlReadFileErr, err)
		}
		if !common.Config.ReleaseMode {
			fmt.Println("_SQL: ", string(bytes))
		}
		if _, err := db.Mysql.Exec(string(bytes)); err != nil {
			return e.SP(e.MMysqlErr, e.MysqlWalkSqlExecute, err)
		}
		return nil
	})
}

func CreateColumns(user *model.User) ([]*model.Column, error) {
	columnIds := []int{
		model.ColumnIdNews,
		model.ColumnIdBar,
		model.ColumnIdRental,
		model.ColumnIdCarpooling,
	}
	models := make([]*model.Column, 0)
	for _, id := range columnIds {
		x := &model.Column{
			ID:         uint64(id),
			Name:       model.ColumnsOfOrigin[id],
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
	avatarURL := "http://" + path.Join(common.Config.Oss.Domain, common.Config.Avatar.DirPrefix, common.Config.Avatar.DefaultAvatarFile)
	user := &model.User{
		Phone:          "86 " + conf.SuperAdminPhone,
		NickName:       "admin",
		HashedPassword: hashedPassword,
		PasswordLevel:  passwordLevel,
		RegisterSource: "web",
		Role:           model.UserRoleSuperAdmin,
		Avatar:         avatarURL,
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
