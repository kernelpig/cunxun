package script

import (
	"path"
	"wangqingang/cunxun/common"
)

func InitScript() error {
	tablePath := path.Join(common.Config.Mysql.SqlPathDir, "/table")
	if err := CreateTables(tablePath); err != nil {
		return err
	}
	viewPath := path.Join(common.Config.Mysql.SqlPathDir, "/view")
	if err := CreateTables(viewPath); err != nil {
		return err
	}
	user, err := CreateSuperAdmin()
	if err != nil {
		return err
	}
	if _, err := CreateColumns(user); err != nil {
		return err
	}
	if err := CreateTokenKey(common.Config.Token.PublicKeyPath, common.Config.Token.PrivateKeyPath); err != nil {
		return err
	}
	return nil
}
