package handler

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	"wangqingang/cunxun/db"
	e "wangqingang/cunxun/error"
	"wangqingang/cunxun/middleware"
	"wangqingang/cunxun/model"
)

func CarpoolingCreateHandler(c *gin.Context) {
	var req CarpoolingCreateRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, e.IP(e.ICarpoolingCreate, e.MParamsErr, e.ParamsBindErr, err))
		return
	}
	currentCtx := middleware.GetCurrentAuth(c)
	if currentCtx == nil {
		c.JSON(http.StatusBadRequest, e.I(e.ICarpoolingCreate, e.MAuthErr, e.AuthGetCurrentErr))
		return
	}
	Carpooling := &model.Carpooling{
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		FromCity:    req.FromCity,
		ToCity:      req.ToCity,
		DepartTIme:  req.DepartTIme,
		PeopleCount: req.PeopleCount,
		Status:      model.CarpoolingEnable,
		Remark:      req.Remark,
		CreaterUid:  int(currentCtx.Payload.UserId),
		UpdaterUid:  int(currentCtx.Payload.UserId),
	}
	if _, err := model.CreateCarpooling(db.Mysql, Carpooling); err != nil {
		if msgErr, ok := err.(e.Message); ok && msgErr.Code.IsSubError(e.MCarpoolingErr, e.CarpoolingAlreadyExist) {
			c.JSON(http.StatusBadRequest, e.IP(e.ICarpoolingCreate, e.MCarpoolingErr, e.CarpoolingAlreadyExist, err))
		} else {
			c.JSON(http.StatusInternalServerError, e.IP(e.ICarpoolingCreate, e.MCarpoolingErr, e.CarpoolingCreateErr, err))
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":          e.OK,
		"carpooling_id": Carpooling.ID,
	})
}

func CarpoolingGetHandler(c *gin.Context) {
	CarpoolingID, err := strconv.ParseInt(c.Param("carpooling_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, e.IP(e.ICarpoolingGetById, e.MParamsErr, e.ParamsInvalidMultiForm, err))
		return
	}
	Carpooling, err := model.GetCarpoolingByID(db.Mysql, int(CarpoolingID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, e.IP(e.ICarpoolingGetById, e.MCarpoolingErr, e.CarpoolingGetErr, err))
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": e.OK,
		"item": Carpooling,
	})
}

// TODO: 热贴需要支持时间范围过滤
func CarpoolingGetListHandler(c *gin.Context) {
	createrUid, err := strconv.ParseInt(c.Query("creater_uid"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, e.IP(e.ICarpoolingGetList, e.MParamsErr, e.ParamsInvalidUserId, err))
		return
	}
	pageNum, err := strconv.ParseInt(c.Query("page_num"), 10, 64)
	if err != nil || pageNum == 0 {
		c.JSON(http.StatusBadRequest, e.IP(e.ICarpoolingGetList, e.MParamsErr, e.ParamsInvalidPageNum, err))
		return
	}
	pageSize, err := strconv.ParseInt(c.Query("page_size"), 10, 64)
	if err != nil || pageSize == 0 {
		c.JSON(http.StatusBadRequest, e.IP(e.ICarpoolingGetList, e.MParamsErr, e.ParamsInvalidPageSize, err))
		return
	}
	where := map[string]interface{}{"creater_uid": createrUid}
	list, isOver, err := model.GetCarpoolingList(db.Mysql, where, c.Query("order_by"), int(pageSize), int(pageNum))
	if err != nil {
		c.JSON(http.StatusInternalServerError, e.IP(e.ICarpoolingGetList, e.MCarpoolingErr, e.CarpoolingGetListErr, err))
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": e.OK,
		"end":  isOver,
		"list": list,
	})
}

func CarpoolingUpdateByIdHandler(c *gin.Context) {
	var req CarpoolingUpdateRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, e.IP(e.ICarpoolingUpdateById, e.MParamsErr, e.ParamsBindErr, err))
		return
	}
	CarpoolingID, err := strconv.ParseInt(c.Param("carpooling_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, e.IP(e.ICarpoolingUpdateById, e.MParamsErr, e.ParamsInvalidCarpoolingID, err))
		return
	}
	currentCtx := middleware.GetCurrentAuth(c)
	if currentCtx == nil {
		c.JSON(http.StatusBadRequest, e.I(e.ICarpoolingUpdateById, e.MAuthErr, e.AuthGetCurrentErr))
		return
	}
	Carpooling := &model.Carpooling{
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		FromCity:    req.FromCity,
		ToCity:      req.ToCity,
		DepartTIme:  req.DepartTIme,
		PeopleCount: req.PeopleCount,
		Status:      model.CarpoolingEnable,
		Remark:      req.Remark,
		UpdaterUid:  int(currentCtx.Payload.UserId),
	}
	if currentCtx.Payload.Role == model.UserRoleAdmin || currentCtx.Payload.Role == model.UserRoleSuperAdmin {
		// 管理员操作
		if _, err := model.UpdateCarpoolingById(db.Mysql, int(CarpoolingID), Carpooling); err != nil {
			c.JSON(http.StatusInternalServerError, e.IP(e.ICarpoolingUpdateById, e.MCarpoolingErr, e.CarpoolingUpdateByIdErr, err))
			return
		}
	} else {
		// 创建者操作
		if _, err := model.UpdateCarpoolingByIdOfSelf(db.Mysql, int(CarpoolingID), int(currentCtx.Payload.UserId), Carpooling); err != nil {
			c.JSON(http.StatusInternalServerError, e.IP(e.ICarpoolingUpdateById, e.MCarpoolingErr, e.CarpoolingUpdateByIdSelfErr, err))
			return
		}
	}
	c.JSON(http.StatusOK, gin.H{
		"code": e.OK,
	})
}
