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

func ArticleGetHandler(c *gin.Context) {
	articleID, err := strconv.ParseInt(c.Param("article_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, e.IP(e.IArticleGet, e.MParamsErr, e.ParamsInvalidArticleID, err))
		return
	}
	article, err := model.GetArticleByID(db.Mysql, int(articleID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, e.IP(e.IArticleGet, e.MArticleErr, e.ArticleGetErr, err))
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": e.OK,
		"item": article,
	})
}

// columnID, orderBy, pageNum, pageSize
// TODO: 热贴需要支持时间范围过滤
func ArticleGetListHandler(c *gin.Context) {
	columnID, err := strconv.ParseInt(c.Query("column_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, e.IP(e.IArticleGetList, e.MParamsErr, e.ParamsInvalidColumnID, err))
		return
	}
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
	list, isOver, err := model.GetArticleList(db.Mysql, map[string]interface{}{"column_id": columnID}, c.Query("order_by"), int(pageSize), int(pageNum))
	if err != nil {
		c.JSON(http.StatusInternalServerError, e.IP(e.IArticleGetList, e.MArticleErr, e.ArticleGetListErr, err))
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": e.OK,
		"end":  isOver,
		"list": list,
	})
}