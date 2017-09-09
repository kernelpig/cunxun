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

func ArticleCreateHandler(c *gin.Context) {
	var req ArticleCreateRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, e.IP(e.IArticleCreate, e.MParamsErr, e.ParamsBindErr, err))
		return
	}

	currentCtx := middleware.GetCurrentAuth(c)
	if currentCtx == nil {
		c.JSON(http.StatusBadRequest, e.I(e.IArticleCreate, e.MAuthErr, e.AuthGetCurrentErr))
		return
	}

	article := &model.Article{
		ColumnId:   req.ColumnId,
		Title:      req.Title,
		Content:    req.Content,
		CreaterUid: int(currentCtx.Payload.UserId),
		UpdaterUid: int(currentCtx.Payload.UserId),
	}
	if _, err := model.CreateArticle(db.Mysql, article); err != nil {
		if msgErr, ok := err.(e.Message); ok && msgErr.Code.IsSubError(e.MArticleErr, e.ArticleAlreadyExist) {
			c.JSON(http.StatusBadRequest, e.IP(e.IArticleCreate, e.MArticleErr, e.ArticleAlreadyExist, err))
		} else {
			c.JSON(http.StatusInternalServerError, e.IP(e.IArticleCreate, e.MArticleErr, e.ArticleCreateErr, err))
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":       e.OK,
		"article_id": article.ID,
	})
}

// columnID, pageNum, pageSize
func ArticleGetListHandler(c *gin.Context) {
	columnID, err := strconv.ParseInt(c.Query("column_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, e.IP(e.IArticleGetList, e.MParamsErr, e.ParamsInvalidColumnID, err))
		return
	}
	pageNum, err := strconv.ParseInt(c.Query("page_num"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, e.IP(e.IArticleGetList, e.MParamsErr, e.ParamsInvalidPageNum, err))
		return
	}
	pageSize, err := strconv.ParseInt(c.Query("page_size"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, e.IP(e.IArticleGetList, e.MParamsErr, e.ParamsInvalidPageSize, err))
		return
	}
	list, err := model.GetArticleList(db.Mysql, map[string]interface{}{"column_id": columnID}, int(pageSize), int(pageNum))
	if err != nil {
		c.JSON(http.StatusInternalServerError, e.IP(e.IArticleGetList, e.MArticleErr, e.ArticleGetListErr, err))
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": e.OK,
		"list": list,
	})
}
