package main

import (
	"wangqingang/cunxun/common"
	"wangqingang/cunxun/db"
	e "wangqingang/cunxun/error"
	"wangqingang/cunxun/model"
	"wangqingang/cunxun/password"
)

func UserCreateSuperAdmin(conf *common.UserConfig) error {
	hashedPassword, err := password.Encrypt(conf.SuperAdminPassword)
	if err != nil {
		return e.SP(e.MPasswordErr, e.PasswordEncryptErr, err)
	}
	passwordLevel, err := password.PasswordStrength(conf.SuperAdminPassword)
	if err != nil {
		return e.SP(e.MPasswordErr, e.PasswordLevelErr, err)
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
			return nil
		}
		return e.SP(e.MUserErr, e.UserCreateErr, err)
	}
	return nil
}
