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

func CommentCreateHandler(c *gin.Context) {
	var req CommentCreateRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, e.IP(e.ICommentCreate, e.MParamsErr, e.ParamsBindErr, err))
		return
	}
	if len([]rune(req.Content)) > common.Config.Comment.DefaultMaxLength {
		c.JSON(http.StatusBadRequest, e.I(e.ICommentCreate, e.MParamsErr, e.ParamsCommentLengthLimit))
		return
	}

	currentCtx := middleware.GetCurrentAuth(c)
	if currentCtx == nil {
		c.JSON(http.StatusBadRequest, e.I(e.ICommentCreate, e.MAuthErr, e.AuthGetCurrentErr))
		return
	}
	relateId, err := strconv.ParseUint(req.RelateId, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, e.I(e.ICommentCreate, e.MParamsErr, e.ParamsInvalidRelateID))
		return
	}
	comment := &model.Comment{
		RelateId:   relateId,
		Content:    req.Content,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
		CreaterUid: currentCtx.Payload.UserId,
	}
	if _, err := model.CreateComment(db.Mysql, comment); err != nil {
		if msgErr, ok := err.(e.Message); ok && msgErr.Code.IsSubError(e.MCommentErr, e.CommentAlreadyExist) {
			c.JSON(http.StatusBadRequest, e.IP(e.ICommentCreate, e.MCommentErr, e.CommentAlreadyExist, err))
		} else {
			c.JSON(http.StatusInternalServerError, e.IP(e.ICommentCreate, e.MCommentErr, e.CommentCreateErr, err))
		}
		return
	}

	c.JSON(http.StatusOK, CreateResponse{Code: e.OK, Id: FormatId(comment.ID)})
}

func CommentGetHandler(c *gin.Context) {
	commentID, err := strconv.ParseUint(c.Param("comment_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, e.IP(e.ICommentGet, e.MParamsErr, e.ParamsInvalidCommentID, err))
		return
	}
	comment, err := model.GetCommentByID(db.Mysql, commentID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, e.IP(e.ICommentGet, e.MCommentErr, e.CommentGetErr, err))
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": e.OK,
		"item": comment,
	})
}

// columnID, pageNum, pageSize
func CommentGetListHandler(c *gin.Context) {
	createrUid, err := strconv.ParseUint(c.Query("creater_uid"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, e.IP(e.IArticleGetList, e.MParamsErr, e.ParamsInvalidUserId, err))
		return
	}
	relateID, err := strconv.ParseUint(c.Query("relate_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, e.IP(e.ICommentGetList, e.MParamsErr, e.ParamsInvalidCommentID, err))
		return
	}
	pageNum, err := strconv.ParseInt(c.Query("page_num"), 10, 64)
	if err != nil || pageNum == 0 {
		c.JSON(http.StatusBadRequest, e.IP(e.ICommentGetList, e.MParamsErr, e.ParamsInvalidPageNum, err))
		return
	}
	pageSize, err := strconv.ParseInt(c.Query("page_size"), 10, 64)
	if err != nil || pageSize == 0 {
		c.JSON(http.StatusBadRequest, e.IP(e.ICommentGetList, e.MParamsErr, e.ParamsInvalidPageSize, err))
		return
	}
	where := map[string]interface{}{"relate_id": relateID, "creater_uid": createrUid}
	list, isOver, err := model.GetCommentList(db.Mysql, where, int(pageSize), int(pageNum))
	if err != nil {
		c.JSON(http.StatusInternalServerError, e.IP(e.ICommentGetList, e.MCommentErr, e.CommentGetListErr, err))
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": e.OK,
		"end":  isOver,
		"list": list,
	})
}

func CommentUpdateByIdHandler(c *gin.Context) {
	var req CommentUpdateRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, e.IP(e.ICommentUpdateById, e.MParamsErr, e.ParamsBindErr, err))
		return
	}
	commentID, err := strconv.ParseUint(c.Param("comment_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, e.IP(e.ICommentUpdateById, e.MParamsErr, e.ParamsInvalidCommentID, err))
		return
	}
	if len([]rune(req.Content)) > common.Config.Comment.DefaultMaxLength {
		c.JSON(http.StatusBadRequest, e.I(e.ICommentUpdateById, e.MParamsErr, e.ParamsCommentLengthLimit))
		return
	}
	currentCtx := middleware.GetCurrentAuth(c)
	if currentCtx == nil {
		c.JSON(http.StatusBadRequest, e.I(e.ICommentUpdateById, e.MAuthErr, e.AuthGetCurrentErr))
		return
	}
	relateId, err := strconv.ParseUint(req.RelateId, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, e.I(e.ICommentUpdateById, e.MParamsErr, e.ParamsInvalidRelateID))
		return
	}
	comment := &model.Comment{
		RelateId:  relateId,
		Content:   req.Content,
		UpdatedAt: time.Now(),
	}
	if currentCtx.Payload.Role == model.UserRoleAdmin || currentCtx.Payload.Role == model.UserRoleSuperAdmin {
		if _, err := model.UpdateCommentById(db.Mysql, commentID, comment); err != nil {
			c.JSON(http.StatusInternalServerError, e.IP(e.ICommentUpdateById, e.MCommentErr, e.CommentUpdateErr, err))
			return
		}
	} else {
		if _, err := model.UpdateCommentByIdOfSelf(db.Mysql, commentID, currentCtx.Payload.UserId, comment); err != nil {
			c.JSON(http.StatusInternalServerError, e.IP(e.ICommentUpdateById, e.MCommentErr, e.CommentUpdateByIdSelf, err))
			return
		}
	}
	c.JSON(http.StatusOK, gin.H{
		"code": e.OK,
	})
}

func CommentDeleteByIdHandler(c *gin.Context) {
	commentID, err := strconv.ParseUint(c.Param("comment_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, e.IP(e.ICommentDeleteById, e.MParamsErr, e.ParamsInvalidCommentID, err))
		return
	}
	currentCtx := middleware.GetCurrentAuth(c)
	if currentCtx == nil {
		c.JSON(http.StatusBadRequest, e.I(e.ICommentDeleteById, e.MAuthErr, e.AuthGetCurrentErr))
		return
	}
	if currentCtx.Payload.Role == model.UserRoleAdmin || currentCtx.Payload.Role == model.UserRoleSuperAdmin {
		if _, err := model.DeleteCommentById(db.Mysql, commentID); err != nil {
			c.JSON(http.StatusInternalServerError, e.IP(e.ICommentDeleteById, e.MCommentErr, e.CommentDeleteErr, err))
			return
		}
	} else {
		if _, err := model.DeleteCommentByIdOfSelf(db.Mysql, commentID, currentCtx.Payload.UserId); err != nil {
			c.JSON(http.StatusInternalServerError, e.IP(e.ICommentDeleteById, e.MCommentErr, e.CommentDeleteByIdSelf, err))
			return
		}
	}
	c.JSON(http.StatusOK, gin.H{
		"code": e.OK,
	})
}
