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

func ColumnGetListHandler(c *gin.Context) {
	createrUid, err := strconv.ParseUint(c.Query("creater_uid"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, e.IP(e.IArticleGetList, e.MParamsErr, e.ParamsInvalidUserId, err))
		return
	}
	list, err := model.GetColumnList(db.Mysql, map[string]interface{}{"creater_uid": createrUid})
	if err != nil {
		c.JSON(http.StatusInternalServerError, e.IP(e.IColumnGetAll, e.MColumnErr, e.ColumnGetAllErr, err))
		return
	}
	c.JSON(http.StatusOK, ColumnGetListResponse{
		Code: e.OK,
		End:  true,
		List: m2rColumnList(list),
	})
}

func ColumnCreateHandler(c *gin.Context) {
	var req ColumnCreateRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, e.IP(e.IColumnCreate, e.MParamsErr, e.ParamsBindErr, err))
		return
	}
	currentCtx := middleware.GetCurrentAuth(c)
	if currentCtx == nil {
		c.JSON(http.StatusBadRequest, e.I(e.IColumnCreate, e.MAuthErr, e.AuthGetCurrentErr))
		return
	}
	if currentCtx.Payload.Role != model.UserRoleSuperAdmin {
		c.JSON(http.StatusBadRequest, e.I(e.ICommentCreate, e.MUserErr, e.UserNotPermit))
		return
	}
	column := &model.Column{
		Name:       req.Name,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
		CreaterUid: currentCtx.Payload.UserId,
	}
	if _, err := model.CreateColumn(db.Mysql, column); err != nil {
		if msgErr, ok := err.(e.Message); ok && msgErr.Code.IsSubError(e.MColumnErr, e.ColumnAlreadyExist) {
			c.JSON(http.StatusBadRequest, e.IP(e.IColumnCreate, e.MColumnErr, e.ColumnAlreadyExist, err))
		} else {
			c.JSON(http.StatusInternalServerError, e.IP(e.IColumnCreate, e.MColumnErr, e.ColumnCreateErr, err))
		}
	}
	c.JSON(http.StatusOK, CreateResponse{Code: e.OK, Id: FormatId(column.ID)})
}

func ColumnUpdateByIdHandler(c *gin.Context) {
	var req ColumnUpdateRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, e.IP(e.IColumnUpdateById, e.MParamsErr, e.ParamsBindErr, err))
		return
	}
	columnID, err := strconv.ParseUint(c.Param("column_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, e.IP(e.IColumnUpdateById, e.MParamsErr, e.ParamsInvalidColumnID, err))
		return
	}
	currentCtx := middleware.GetCurrentAuth(c)
	if currentCtx == nil {
		c.JSON(http.StatusBadRequest, e.I(e.IColumnUpdateById, e.MAuthErr, e.AuthGetCurrentErr))
		return
	}
	column := &model.Column{
		Name:      req.Name,
		UpdatedAt: time.Now(),
	}
	if currentCtx.Payload.Role == model.UserRoleAdmin || currentCtx.Payload.Role == model.UserRoleSuperAdmin {
		if _, err := model.UpdateColumnById(db.Mysql, columnID, column); err != nil {
			c.JSON(http.StatusInternalServerError, e.IP(e.IColumnUpdateById, e.MColumnErr, e.ColumnUpdateById, err))
			return
		}
	} else {
		if _, err := model.UpdateColumnByIdOfSelf(db.Mysql, columnID, currentCtx.Payload.UserId, column); err != nil {
			c.JSON(http.StatusInternalServerError, e.IP(e.IColumnUpdateById, e.MColumnErr, e.ColumnUpdateByIdSelf, err))
			return
		}
	}
	c.JSON(http.StatusOK, OKResponse{
		Code: e.OK,
	})
}

func ColumnDeleteByIdHandler(c *gin.Context) {
	columnID, err := strconv.ParseUint(c.Param("column_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, e.IP(e.IColumnDeleteById, e.MParamsErr, e.ParamsInvalidColumnID, err))
		return
	}
	currentCtx := middleware.GetCurrentAuth(c)
	if currentCtx == nil {
		c.JSON(http.StatusBadRequest, e.I(e.IColumnDeleteById, e.MAuthErr, e.AuthGetCurrentErr))
		return
	}
	if currentCtx.Payload.Role == model.UserRoleAdmin || currentCtx.Payload.Role == model.UserRoleSuperAdmin {
		if _, err := model.DeleteColumnById(db.Mysql, columnID); err != nil {
			c.JSON(http.StatusInternalServerError, e.IP(e.IColumnDeleteById, e.MColumnErr, e.ColumnDeleteById, err))
			return
		}
	} else {
		if _, err := model.DeleteColumnByIdOfSelf(db.Mysql, columnID, currentCtx.Payload.UserId); err != nil {
			c.JSON(http.StatusInternalServerError, e.IP(e.IColumnDeleteById, e.MColumnErr, e.ColumnDeleteByIdSelf, err))
			return
		}
	}
	c.JSON(http.StatusOK, OKResponse{
		Code: e.OK,
	})
}
