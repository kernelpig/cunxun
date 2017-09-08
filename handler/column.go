package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"wangqingang/cunxun/db"
	e "wangqingang/cunxun/error"
	"wangqingang/cunxun/middleware"
	"wangqingang/cunxun/model"
)

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
		"code": e.OK,
	})
}
