package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

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

	currentCtx := middleware.GetCurrentAuth(c)
	if currentCtx == nil {
		c.JSON(http.StatusBadRequest, e.I(e.ICommentCreate, e.MAuthErr, e.AuthGetCurrentErr))
		return
	}

	comment := &model.Comment{
		ArticleId:  req.ArticleId,
		Content:    req.Content,
		CreaterUid: int(currentCtx.Payload.UserId),
	}
	if _, err := model.CreateComment(db.Mysql, comment); err != nil {
		if msgErr, ok := err.(e.Message); ok && msgErr.Code.IsSubError(e.MCommentErr, e.CommentAlreadyExist) {
			c.JSON(http.StatusBadRequest, e.IP(e.ICommentCreate, e.MCommentErr, e.CommentAlreadyExist, err))
		} else {
			c.JSON(http.StatusInternalServerError, e.IP(e.ICommentCreate, e.MCommentErr, e.CommentCreateErr, err))
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":       e.OK,
		"comment_id": comment.ID,
	})
}

func CommentGetHandler(c *gin.Context) {
	commentID, err := strconv.ParseInt(c.Param("comment_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, e.IP(e.ICommentGet, e.MParamsErr, e.ParamsInvalidCommentID, err))
		return
	}
	comment, err := model.GetCommentByID(db.Mysql, int(commentID))
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
	articleID, err := strconv.ParseInt(c.Query("article_id"), 10, 64)
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
	list, isOver, err := model.GetCommentList(db.Mysql, map[string]interface{}{"article_id": articleID}, int(pageSize), int(pageNum))
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
