package handler

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"net/http"
	"path"
	"strconv"
	"strings"
	"time"

	"github.com/ahmetb/go-linq"
	"github.com/gin-gonic/gin"
	"github.com/satori/go.uuid"

	"wangqingang/cunxun/captcha"
	"wangqingang/cunxun/common"
	"wangqingang/cunxun/db"
	e "wangqingang/cunxun/error"
	"wangqingang/cunxun/id"
	"wangqingang/cunxun/login"
	"wangqingang/cunxun/middleware"
	"wangqingang/cunxun/model"
	"wangqingang/cunxun/oss"
	"wangqingang/cunxun/password"
	"wangqingang/cunxun/phone"
	"wangqingang/cunxun/token"
)

// 携带的数据格式为: data:image/png;base64,BASE64编码内容
func parseAvatar(data string) (string, string, error) {
	if len(data) <= 2 {
		return "", "", e.SD(e.MUserErr, e.UserAvatarDecodeErr, "invalid length")
	}
	s := strings.Split(data, ";")
	if len(s) != 2 || len(s[0]) == 0 || len(s[1]) == 0 {
		return "", "", e.SD(e.MUserErr, e.UserAvatarDecodeErr, "split by ; error")
	}
	s1 := strings.Split(s[0], ":")
	if len(s1) != 2 || len(s1[0]) == 0 || len(s1[1]) == 0 {
		return "", "", e.SD(e.MUserErr, e.UserAvatarDecodeErr, "split by : error")
	}
	s2 := strings.Split(s[1], ",")
	if len(s2) != 2 || len(s2[0]) == 0 || len(s2[1]) == 0 {
		return "", "", e.SD(e.MUserErr, e.UserAvatarDecodeErr, "split by , error")
	}
	s3 := strings.Split(s1[1], "/")
	if len(s3) != 2 || s3[0] != "image" || len(s3[1]) == 0 {
		return "", "", e.SD(e.MUserErr, e.UserAvatarDecodeErr, "split by / error")
	}
	return s3[1], s2[1], nil
}

// 获取头像文件路径
func getAvatarFilePath(reqAvatar string, id uint64) string {
	var fileName string = common.Config.Avatar.DefaultAvatarFile
	if reqAvatar != "" {
		fileName = FormatId(id)
	}
	return path.Join(common.Config.Avatar.DirPrefix, fileName)
}

// 写入头像文件, 并生成头像文件名, DB不存储路径, 读取时根据配置读取
func writeAvatarFile(filePath, reqAvatar string) (string, error) {
	if reqAvatar == "" {
		return "", e.SD(e.MUserErr, e.UserAvatarDecodeErr, "avatar param invalid")
	}
	_, contentBase64, err := parseAvatar(reqAvatar)
	if err != nil {
		return "", err
	}
	contentBytes, err := base64.StdEncoding.DecodeString(contentBase64)
	if err != nil {
		return "", e.SP(e.MUserErr, e.UserAvatarDecodeErr, err)
	}
	link := oss.PutImageByFileAsync(filePath, bytes.NewBuffer(contentBytes))
	return link, nil
}

func UserSignupHandler(c *gin.Context) {
	var req UserSignupRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, e.IP(e.IUserSignup, e.MParamsErr, e.ParamsBindErr, err))
		return
	}

	detail := fmt.Sprintf("%+v", req)
	if !linq.From(common.SourceRange).Contains(req.Source) {
		c.JSON(http.StatusBadRequest, e.ID(e.IUserSignup, e.MParamsErr, e.ParamsInvalidSource, detail))
		return
	}

	if _, err := phone.ValidPhone(req.Phone); err != nil {
		c.JSON(http.StatusBadRequest, e.ID(e.IUserSignup, e.MParamsErr, e.ParamsInvalidPhone, detail))
		return
	}

	ok, err := CheckcodeVerify(c, req.Phone, req.Source, common.SignupPurpose, req.VerifyCode)
	if err != nil {
		c.JSON(http.StatusInternalServerError, e.IP(e.IUserSignup, e.MCheckcodeErr, e.CheckcodeCheckErr, err))
		return
	} else if !ok {
		c.JSON(http.StatusBadRequest, e.I(e.IUserSignup, e.MCheckcodeErr, e.CheckcodeMismatch))
		return
	}

	hashedPassword, err := password.Encrypt(req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, e.IP(e.IUserSignup, e.MPasswordErr, e.PasswordEncryptErr, err))
		return
	}

	passwordLevel, err := password.PasswordStrength(req.Password)
	if err != nil {
		c.JSON(http.StatusBadRequest, e.IP(e.IUserSignup, e.MPasswordErr, e.PasswordLevelErr, err))
		return
	}
	id, err := id.Generate()
	if err != nil {
		c.JSON(http.StatusBadRequest, e.IP(e.IUserSignup, e.MIdGeneratorErr, e.IdGeneratorErr, err))
		return
	}

	filePath := getAvatarFilePath(req.Avatar, id)
	avatarLink, err := writeAvatarFile(filePath, req.Avatar)
	if err != nil {
		filePath = getAvatarFilePath("", id)
		avatarLink = path.Join(common.Config.Oss.Domain, filePath)
	}

	user := &model.User{
		Phone:          req.Phone,
		NickName:       req.NickName,
		HashedPassword: hashedPassword,
		PasswordLevel:  passwordLevel,
		RegisterSource: req.Source,
		Avatar:         avatarLink,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	user, err = model.CreateUser(db.Mysql, user)
	if err != nil {
		if msgErr, ok := err.(e.Message); ok && msgErr.Code.IsSubError(e.MUserErr, e.UserAlreadyExist) {
			c.JSON(http.StatusBadRequest, e.IP(e.IUserSignup, e.MUserErr, e.UserAlreadyExist, err))
		} else {
			c.JSON(http.StatusInternalServerError, e.IP(e.IUserSignup, e.MUserErr, e.UserCreateErr, err))
		}
		return
	}

	c.JSON(http.StatusOK, CreateResponse{Code: e.OK, Id: strconv.FormatUint(user.ID, 10)})
	return
}

func UserUpdateHandler(c *gin.Context) {
	var req UserUpdateRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, e.IP(e.IUserUpdateById, e.MParamsErr, e.ParamsBindErr, err))
		return
	}
	userID, err := strconv.ParseUint(c.Param("user_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, e.IP(e.IUserUpdateById, e.MParamsErr, e.ParamsInvalidColumnID, err))
		return
	}
	currentCtx := middleware.GetCurrentAuth(c)
	if currentCtx == nil {
		c.JSON(http.StatusBadRequest, e.I(e.IUserUpdateById, e.MAuthErr, e.AuthGetCurrentErr))
		return
	}
	if currentCtx.Payload.Role != model.UserRoleSuperAdmin {
		c.JSON(http.StatusBadRequest, e.I(e.IUserUpdateById, e.MUserErr, e.UserNotPermit))
		return
	}
	update := &model.User{
		Role:      req.Role,
		UpdatedAt: time.Now(),
	}
	if _, err := model.UpdateUserById(db.Mysql, userID, update); err != nil {
		c.JSON(http.StatusInternalServerError, e.IP(e.IUserUpdateById, e.MUserErr, e.UserCreateErr, err))
		return
	}
	c.JSON(http.StatusOK, OKResponse{Code: e.OK})
	return
}

func UserCreateHandler(c *gin.Context) {
	var req UserCreateRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, e.IP(e.IUserCreate, e.MParamsErr, e.ParamsBindErr, err))
		return
	}
	if req.NickName == "" {
		req.NickName = uuid.NewV4().String()
	}
	if _, err := phone.ValidPhone(req.Phone); err != nil {
		c.JSON(http.StatusBadRequest, e.ID(e.IUserCreate, e.MParamsErr, e.ParamsInvalidPhone, req.Phone))
		return
	}
	hashedPassword, err := password.Encrypt(req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, e.IP(e.IUserCreate, e.MPasswordErr, e.PasswordEncryptErr, err))
		return
	}
	passwordLevel, err := password.PasswordStrength(req.Password)
	if err != nil {
		c.JSON(http.StatusBadRequest, e.IP(e.IUserCreate, e.MPasswordErr, e.PasswordLevelErr, err))
		return
	}
	currentCtx := middleware.GetCurrentAuth(c)
	if currentCtx == nil {
		c.JSON(http.StatusBadRequest, e.I(e.IUserCreate, e.MAuthErr, e.AuthGetCurrentErr))
		return
	}
	if currentCtx.Payload.Role != model.UserRoleSuperAdmin {
		c.JSON(http.StatusBadRequest, e.I(e.IUserCreate, e.MUserErr, e.UserNotPermit))
		return
	}
	user := &model.User{
		Phone:          req.Phone,
		NickName:       req.NickName,
		HashedPassword: hashedPassword,
		PasswordLevel:  passwordLevel,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}
	user, err = model.CreateUser(db.Mysql, user)
	if err != nil {
		if msgErr, ok := err.(e.Message); ok && msgErr.Code.IsSubError(e.MUserErr, e.UserAlreadyExist) {
			c.JSON(http.StatusBadRequest, e.IP(e.IUserCreate, e.MUserErr, e.UserAlreadyExist, err))
		} else {
			c.JSON(http.StatusInternalServerError, e.IP(e.IUserCreate, e.MUserErr, e.UserCreateErr, err))
		}
		return
	}
	c.JSON(http.StatusOK, CreateResponse{Code: e.OK, Id: strconv.FormatUint(user.ID, 10)})
	return
}

func UserLoginHandler(c *gin.Context) {
	var req UserLoginRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, e.IP(e.IUserLogin, e.MParamsErr, e.ParamsBindErr, err))
		return
	}

	detail := fmt.Sprintf("%+v", req)
	if !linq.From(common.SourceRange).Contains(req.Source) {
		c.JSON(http.StatusBadRequest, e.ID(e.IUserLogin, e.MParamsErr, e.ParamsInvalidSource, detail))
		return
	}
	if _, err := phone.ValidPhone(req.Phone); err != nil {
		c.JSON(http.StatusBadRequest, e.ID(e.IUserLogin, e.MParamsErr, e.ParamsInvalidPhone, detail))
		return
	}

	user, err := model.GetUserByPhone(db.Mysql, req.Phone)
	if err != nil {
		c.JSON(http.StatusInternalServerError, e.IP(e.IUserLogin, e.MUserErr, e.UserGetErr, err))
		return
	} else if user == nil {
		c.JSON(http.StatusBadRequest, e.I(e.IUserLogin, e.MUserErr, e.UserNotExist))
		return
	}

	var loginKey = login.LoginKey{Phone: req.Phone, Purpose: common.SigninPurpose, Source: req.Source}
	login, err := loginKey.GetLogin()
	if err != nil {
		c.JSON(http.StatusInternalServerError, e.IP(e.IUserLogin, e.MLoginErr, e.LoginGetErr, err))
		return
	} else if login == nil {
		login, err = loginKey.CreateLogin(common.Config.Login.TTL.D())
		if err != nil || login == nil {
			c.JSON(http.StatusInternalServerError, e.IP(e.IUserLogin, e.MLoginErr, e.LoginCreateErr, err))
			return
		}
	}

	// 登录次数检查
	needCaptcha := false
	if login.RequestTimes >= common.Config.Login.MaxCaptchaTImes {
		needCaptcha = true
	}
	if needCaptcha {
		leftTimes := fmt.Sprintf("%d", login.GetLeftTimes())
		if req.CaptchaId == "" || req.CaptchaValue == "" {
			c.JSON(http.StatusBadRequest, e.ID(e.IUserLogin, e.MCaptchaErr, e.CaptchaRequired, leftTimes))
			return
		}

		if !captcha.VerifyCaptcha(req.CaptchaId, req.CaptchaValue) {
			c.JSON(http.StatusBadRequest, e.ID(e.IUserLogin, e.MCaptchaErr, e.CaptchaMismatch, leftTimes))
			return
		}
	}
	if login.RequestTimes >= common.Config.Login.MaxRequestTimes {
		c.JSON(http.StatusBadRequest, e.I(e.IUserLogin, e.MLoginErr, e.LogDumpRequestErr))
		return
	}
	login.RequestTimes++
	login.Save()

	if err := password.Verify(req.Password, user.HashedPassword); err != nil {
		c.JSON(http.StatusBadRequest, e.IP(e.IUserLogin, e.MPasswordErr, e.PasswordInvalid, err))
		return
	}

	accessToken, err := token.TokenCreateAndStore(user.ID, user.Role, req.Source, common.Config.Token.AccessTokenTTL.D())
	if err != nil {
		c.JSON(http.StatusInternalServerError, e.IP(e.IUserLogin, e.MTokenErr, e.TokenCreateErr, err))
		return
	}

	c.JSON(http.StatusOK, UserLoginResponse{
		Code:      e.OK,
		UserRole:  user.Role,
		UserId:    strconv.FormatUint(user.ID, 10),
		UserToken: accessToken,
	})
	return
}

func UserLogoutHandler(c *gin.Context) {
	userToken := c.GetHeader(common.AuthHeaderKey)
	payload, err := middleware.CheckAccessToken(userToken)

	if err == nil && payload != nil {
		token.TokenClean(payload.UserId, payload.LoginSource)
	}

	c.JSON(http.StatusOK, OKResponse{
		Code: e.OK,
	})
	return
}

func UserGetInfoHandler(c *gin.Context) {
	userId, err := strconv.ParseUint(c.Param("user_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, e.IP(e.IUserGetInfo, e.MParamsErr, e.ParamsInvalidUserId, err))
		return
	}
	user, err := model.GetUserByID(db.Mysql, userId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, e.IP(e.IUserGetInfo, e.MUserErr, e.UserGetErr, err))
		return
	} else if user == nil {
		detail := fmt.Sprintf("userid: %d", userId)
		c.JSON(http.StatusBadRequest, e.ID(e.IUserGetInfo, e.MUserErr, e.UserNotExist, detail))
		return
	}

	c.JSON(http.StatusOK, UserGetInfoResponse{
		Code:     e.OK,
		UserId:   strconv.FormatUint(user.ID, 10),
		Nickname: user.NickName,
		Phone:    user.Phone,
		Avatar:   user.Avatar,
	})
}

func UserGetListHandler(c *gin.Context) {
	pageNum, err := strconv.ParseInt(c.Query("page_num"), 10, 64)
	if err != nil || pageNum == 0 {
		c.JSON(http.StatusBadRequest, e.IP(e.IArticleGetList, e.MParamsErr, e.ParamsInvalidPageNum, err))
		return
	}
	pageSize, err := strconv.ParseInt(c.Query("page_size"), 10, 64)
	if err != nil || pageSize == 0 {
		c.JSON(http.StatusBadRequest, e.IP(e.IArticleGetList, e.MParamsErr, e.ParamsInvalidPageSize, err))
		return
	}
	currentCtx := middleware.GetCurrentAuth(c)
	if currentCtx == nil {
		c.JSON(http.StatusBadRequest, e.I(e.IUserGetList, e.MAuthErr, e.AuthGetCurrentErr))
		return
	}
	if currentCtx.Payload.Role != model.UserRoleSuperAdmin {
		c.JSON(http.StatusBadRequest, e.I(e.IUserGetList, e.MUserErr, e.UserNotPermit))
		return
	}
	list, isOver, err := model.GetUserList(db.Mysql, map[string]interface{}{}, c.Query("order_by"), int(pageSize), int(pageNum))
	if err != nil {
		c.JSON(http.StatusInternalServerError, e.IP(e.IUserGetList, e.MColumnErr, e.ColumnGetAllErr, err))
		return
	}
	c.JSON(http.StatusOK, UserGetListResponse{
		Code: e.OK,
		End:  isOver,
		List: m2rUserList(list),
	})
}

func UserDeleteByIdHandler(c *gin.Context) {
	userID, err := strconv.ParseUint(c.Param("user_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, e.IP(e.IUserDeleteById, e.MParamsErr, e.ParamsInvalidUserId, err))
		return
	}
	currentCtx := middleware.GetCurrentAuth(c)
	if currentCtx == nil {
		c.JSON(http.StatusBadRequest, e.I(e.IUserDeleteById, e.MAuthErr, e.AuthGetCurrentErr))
		return
	}
	if currentCtx.Payload.Role != model.UserRoleSuperAdmin {
		c.JSON(http.StatusBadRequest, e.I(e.IUserDeleteById, e.MUserErr, e.UserNotPermit))
		return
	}
	if _, err := model.DeleteUserById(db.Mysql, userID); err != nil {
		c.JSON(http.StatusInternalServerError, e.IP(e.IUserDeleteById, e.MUserErr, e.UserDeleteErr, err))
		return
	}
	c.JSON(http.StatusOK, OKResponse{
		Code: e.OK,
	})
}
