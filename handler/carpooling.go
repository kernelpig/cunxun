package handler

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	"wangqingang/cunxun/common"
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
	if req.PeopleCount > common.Config.Carpooling.DefaultMaxSeat {
		c.JSON(http.StatusBadRequest, e.I(e.ICarpoolingCreate, e.MParamsErr, e.ParamsInvalidCarpoolingSeat))
		return
	}
	departTime := time.Unix(req.DepartTime, 0)
	if departTime.Before(time.Now()) {
		c.JSON(http.StatusBadRequest, e.I(e.ICarpoolingCreate, e.MParamsErr, e.ParamsInvalidCarpoolingDepart))
		return
	}
	currentCtx := middleware.GetCurrentAuth(c)
	if currentCtx == nil {
		c.JSON(http.StatusBadRequest, e.I(e.ICarpoolingCreate, e.MAuthErr, e.AuthGetCurrentErr))
		return
	}
	carpooling := &model.Carpooling{
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		FromCity:    req.FromCity,
		ToCity:      req.ToCity,
		DepartTime:  departTime,
		PeopleCount: req.PeopleCount,
		Contact:     req.Contact,
		Status:      model.CarpoolingEnable,
		Remark:      req.Remark,
		CreaterUid:  currentCtx.Payload.UserId,
		UpdaterUid:  currentCtx.Payload.UserId,
	}
	if _, err := model.CreateCarpooling(db.Mysql, carpooling); err != nil {
		if msgErr, ok := err.(e.Message); ok && msgErr.Code.IsSubError(e.MCarpoolingErr, e.CarpoolingAlreadyExist) {
			c.JSON(http.StatusBadRequest, e.IP(e.ICarpoolingCreate, e.MCarpoolingErr, e.CarpoolingAlreadyExist, err))
		} else {
			c.JSON(http.StatusInternalServerError, e.IP(e.ICarpoolingCreate, e.MCarpoolingErr, e.CarpoolingCreateErr, err))
		}
		return
	}

	c.JSON(http.StatusOK, CreateResponse{Code: e.OK, Id: FormatId(carpooling.ID)})
}

func CarpoolingGetHandler(c *gin.Context) {
	carpoolingID, err := strconv.ParseUint(c.Param("carpooling_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, e.IP(e.ICarpoolingGetById, e.MParamsErr, e.ParamsInvalidMultiForm, err))
		return
	}
	carpooling, err := model.GetCarpoolingByID(db.Mysql, carpoolingID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, e.IP(e.ICarpoolingGetById, e.MCarpoolingErr, e.CarpoolingGetErr, err))
		return
	}
	c.JSON(http.StatusOK, CarpoolingGetListResponse{
		Code: e.OK,
		End:  true,
		List: m2rCarpoolingList([]*model.CarpoolingDetailView{carpooling}),
	})
}

// TODO: 热贴需要支持时间范围过滤
func CarpoolingGetListHandler(c *gin.Context) {
	createrUid, err := strconv.ParseUint(c.Query("creater_uid"), 10, 64)
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
	c.JSON(http.StatusOK, CarpoolingGetListResponse{
		Code: e.OK,
		End:  isOver,
		List: m2rCarpoolingList(list),
	})
}

func CarpoolingUpdateByIdHandler(c *gin.Context) {
	var req CarpoolingUpdateRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, e.IP(e.ICarpoolingUpdateById, e.MParamsErr, e.ParamsBindErr, err))
		return
	}
	carpoolingID, err := strconv.ParseUint(c.Param("carpooling_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, e.IP(e.ICarpoolingUpdateById, e.MParamsErr, e.ParamsInvalidCarpoolingID, err))
		return
	}
	if req.PeopleCount > common.Config.Carpooling.DefaultMaxSeat {
		c.JSON(http.StatusBadRequest, e.I(e.ICarpoolingUpdateById, e.MParamsErr, e.ParamsInvalidCarpoolingSeat))
		return
	}
	departTime := time.Unix(req.DepartTime, 0)
	if departTime.Before(time.Now()) {
		c.JSON(http.StatusBadRequest, e.I(e.ICarpoolingUpdateById, e.MParamsErr, e.ParamsInvalidCarpoolingDepart))
		return
	}
	currentCtx := middleware.GetCurrentAuth(c)
	if currentCtx == nil {
		c.JSON(http.StatusBadRequest, e.I(e.ICarpoolingUpdateById, e.MAuthErr, e.AuthGetCurrentErr))
		return
	}
	Carpooling := &model.Carpooling{
		UpdatedAt:   time.Now(),
		FromCity:    req.FromCity,
		ToCity:      req.ToCity,
		DepartTime:  departTime,
		PeopleCount: req.PeopleCount,
		Contact:     req.Contact,
		Status:      model.CarpoolingEnable,
		Remark:      req.Remark,
		UpdaterUid:  currentCtx.Payload.UserId,
	}
	if currentCtx.Payload.Role == model.UserRoleAdmin || currentCtx.Payload.Role == model.UserRoleSuperAdmin {
		// 管理员操作
		if _, err := model.UpdateCarpoolingById(db.Mysql, carpoolingID, Carpooling); err != nil {
			c.JSON(http.StatusInternalServerError, e.IP(e.ICarpoolingUpdateById, e.MCarpoolingErr, e.CarpoolingUpdateByIdErr, err))
			return
		}
	} else {
		// 创建者操作
		if _, err := model.UpdateCarpoolingByIdOfSelf(db.Mysql, carpoolingID, currentCtx.Payload.UserId, Carpooling); err != nil {
			c.JSON(http.StatusInternalServerError, e.IP(e.ICarpoolingUpdateById, e.MCarpoolingErr, e.CarpoolingUpdateByIdSelfErr, err))
			return
		}
	}
	c.JSON(http.StatusOK, OKResponse{
		Code: e.OK,
	})
}

func CarpoolingDeleteByIdHandler(c *gin.Context) {
	carpoolingID, err := strconv.ParseUint(c.Param("carpooling_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, e.IP(e.ICarpoolingDeleteById, e.MParamsErr, e.ParamsInvalidCarpoolingID, err))
		return
	}
	currentCtx := middleware.GetCurrentAuth(c)
	if currentCtx == nil {
		c.JSON(http.StatusBadRequest, e.I(e.ICarpoolingDeleteById, e.MAuthErr, e.AuthGetCurrentErr))
		return
	}
	if currentCtx.Payload.Role == model.UserRoleAdmin || currentCtx.Payload.Role == model.UserRoleSuperAdmin {
		if _, err := model.DeleteCarpoolingById(db.Mysql, carpoolingID); err != nil {
			c.JSON(http.StatusInternalServerError, e.IP(e.ICarpoolingDeleteById, e.MCarpoolingErr, e.CarpoolingDeleteByIdErr, err))
			return
		}
	} else {
		if _, err := model.DeleteCarpoolingByIdOfSelf(db.Mysql, carpoolingID, currentCtx.Payload.UserId); err != nil {
			c.JSON(http.StatusInternalServerError, e.IP(e.ICarpoolingDeleteById, e.MCarpoolingErr, e.CarpoolingDeleteByIdSelfErr, err))
			return
		}
	}

	c.JSON(http.StatusOK, OKResponse{
		Code: e.OK,
	})
}
