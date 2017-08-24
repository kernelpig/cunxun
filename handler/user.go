package handler

import (
	"net/http"

	"github.com/ahmetb/go-linq"
	"github.com/gin-gonic/gin"

	"wangqingang/cunxun/common"
	"wangqingang/cunxun/db"
	"wangqingang/cunxun/model"
	"wangqingang/cunxun/password"
	"wangqingang/cunxun/phone"
)

func UserSignupHandler(c *gin.Context) {
	var req UserSignupRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": common.AccountBindFailed,
		})
		return
	}

	if !linq.From(common.SourceRange).Contains(req.Source) {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": common.AccountInvalidSource,
		})
		return
	}

	if err := phone.ValidPhone(req.Phone); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": common.AccountInvalidPhone,
		})
		return
	}

	ok, err := CheckcodeVerify(c, req.Phone, req.Source, common.SignupPurpose, req.VerifyCode)
	if err != nil {
		return
	} else if !ok {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": common.AccountVerifyCodeNotMatch,
		})
		return
	}

	hashedPassword, err := password.Encrypt(req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": common.AccountInternalError,
		})
		return
	}

	passwordLevel := password.PasswordStrength(req.Password)
	if passwordLevel == password.LevelIllegal {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": common.AccountPasswordLevelIllegal,
		})
		return
	}

	user := &model.User{
		Phone:          req.Phone,
		HashedPassword: hashedPassword,
		PasswordLevel:  passwordLevel,
		RegisterSource: req.Source,
	}

	user, err = model.CreateUser(db.Mysql, user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    common.AccountDBError,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": common.OK,
	})
	return
}

func UserLoginHandler(c *gin.Context) {

}

func UserLogoutHandler(c *gin.Context) {

}
