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
	list, err := model.GetColumnList(db.Mysql)
	if err != nil {
		c.JSON(http.StatusInternalServerError, e.IP(e.IColumnGetAll, e.MColumnErr, e.ColumnGetAllErr, err))
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": e.OK,
		"list": list,
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

	column := &model.Column{
		Name:       req.Name,
		CreaterUid: int(currentCtx.Payload.UserId),
	}
	if _, err := model.CreateColumn(db.Mysql, column); err != nil {
		if msgErr, ok := err.(e.Message); ok && msgErr.Code.IsSubError(e.MColumnErr, e.ColumnAlreadyExist) {
			c.JSON(http.StatusBadRequest, e.IP(e.IColumnCreate, e.MColumnErr, e.ColumnAlreadyExist, err))
		} else {
			c.JSON(http.StatusInternalServerError, e.IP(e.IColumnCreate, e.MColumnErr, e.ColumnCreateErr, err))
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"code":      e.OK,
		"column_id": column.ID,
	})
}

func ColumnUpdateByIdHandler(c *gin.Context) {
	var req ColumnUpdateRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, e.IP(e.IColumnUpdateById, e.MParamsErr, e.ParamsBindErr, err))
		return
	}
	columnID, err := strconv.ParseInt(c.Param("column_id"), 10, 64)
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
	if _, err := model.UpdateColumnById(db.Mysql, int(columnID), column); err != nil {
		c.JSON(http.StatusInternalServerError, e.IP(e.IColumnUpdateById, e.MColumnErr, e.ColumnUpdateById, err))
	}

	c.JSON(http.StatusOK, gin.H{
		"code": e.OK,
	})
}

func ColumnDeleteByIdHandler(c *gin.Context) {
	columnID, err := strconv.ParseInt(c.Param("column_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, e.IP(e.IColumnDeleteById, e.MParamsErr, e.ParamsInvalidColumnID, err))
		return
	}
	currentCtx := middleware.GetCurrentAuth(c)
	if currentCtx == nil {
		c.JSON(http.StatusBadRequest, e.I(e.IColumnDeleteById, e.MAuthErr, e.AuthGetCurrentErr))
		return
	}
	if _, err := model.DeleteColumnById(db.Mysql, int(columnID)); err != nil {
		c.JSON(http.StatusInternalServerError, e.IP(e.IColumnDeleteById, e.MColumnErr, e.ColumnUpdateById, err))
	}

	c.JSON(http.StatusOK, gin.H{
		"code": e.OK,
	})
}
