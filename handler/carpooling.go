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
