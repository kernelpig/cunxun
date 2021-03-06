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
	columnId, err := strconv.ParseUint(req.ColumnId, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, e.IP(e.IArticleCreate, e.MParamsErr, e.ParamsInvalidColumnID, err))
		return
	}
	if columnId == model.ColumnIdNews {
		if currentCtx.Payload.Role != model.UserRoleAdmin && currentCtx.Payload.Role != model.UserRoleSuperAdmin {
			c.JSON(http.StatusBadRequest, e.I(e.IArticleCreate, e.MUserErr, e.UserNotPermit))
			return
		}
	}
	article := &model.Article{
		ColumnId:   columnId,
		Title:      req.Title,
		Content:    req.Content,
		Priority:   req.Priority,
		CreaterUid: currentCtx.Payload.UserId,
		UpdaterUid: currentCtx.Payload.UserId,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}
	if _, err := model.CreateArticle(db.Mysql, article); err != nil {
		if msgErr, ok := err.(e.Message); ok && msgErr.Code.IsSubError(e.MArticleErr, e.ArticleAlreadyExist) {
			c.JSON(http.StatusBadRequest, e.IP(e.IArticleCreate, e.MArticleErr, e.ArticleAlreadyExist, err))
		} else {
			c.JSON(http.StatusInternalServerError, e.IP(e.IArticleCreate, e.MArticleErr, e.ArticleCreateErr, err))
		}
		return
	}

	c.JSON(http.StatusOK, CreateResponse{Code: e.OK, Id: FormatId(article.ID)})
}

func ArticleGetHandler(c *gin.Context) {
	articleID, err := strconv.ParseUint(c.Param("article_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, e.IP(e.IArticleGet, e.MParamsErr, e.ParamsInvalidArticleID, err))
		return
	}
	article, err := model.GetArticleByID(db.Mysql, articleID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, e.IP(e.IArticleGet, e.MArticleErr, e.ArticleGetErr, err))
		return
	}
	c.JSON(http.StatusOK, ArticleGetListResponse{
		Code: e.OK,
		End:  true,
		List: m2rArticleList([]*model.ArticleDetailView{article}),
	})
}

// columnID, orderBy, pageNum, pageSize, 为默认值则忽略此查询条件
// TODO: 热贴需要支持时间范围过滤
func ArticleGetListHandler(c *gin.Context) {
	createrUid, err := strconv.ParseUint(c.Query("creater_uid"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, e.IP(e.IArticleGetList, e.MParamsErr, e.ParamsInvalidUserId, err))
		return
	}
	columnID, err := strconv.ParseUint(c.Query("column_id"), 10, 64)
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
	where := map[string]interface{}{"column_id": columnID, "creater_uid": createrUid}
	orderBys := make([]string, 0)
	if c.Query("order_by") == model.OrderByCommentCount {
		orderBys = append(orderBys, model.OrderByCommentCount)
	} else {
		orderBys = append(orderBys, model.OrderByPriority)
	}
	list, isOver, err := model.GetArticleList(db.Mysql, where, orderBys, int(pageSize), int(pageNum))
	if err != nil {
		c.JSON(http.StatusInternalServerError, e.IP(e.IArticleGetList, e.MArticleErr, e.ArticleGetListErr, err))
		return
	}
	c.JSON(http.StatusOK, ArticleGetListResponse{
		Code: e.OK,
		End:  isOver,
		List: m2rArticleList(list),
	})
}

func ArticleUpdateByIdHandler(c *gin.Context) {
	var req ArticleUpdateRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, e.IP(e.IArticleUpdateById, e.MParamsErr, e.ParamsBindErr, err))
		return
	}
	articleID, err := strconv.ParseUint(c.Param("article_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, e.IP(e.IArticleUpdateById, e.MParamsErr, e.ParamsInvalidArticleID, err))
		return
	}
	columnID, err := strconv.ParseUint(req.ColumnId, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, e.IP(e.IArticleUpdateById, e.MParamsErr, e.ParamsInvalidArticleID, err))
		return
	}
	currentCtx := middleware.GetCurrentAuth(c)
	if currentCtx == nil {
		c.JSON(http.StatusBadRequest, e.I(e.IArticleUpdateById, e.MAuthErr, e.AuthGetCurrentErr))
		return
	}
	article := &model.Article{
		ColumnId:   columnID,
		Title:      req.Title,
		Content:    req.Content,
		Priority:   req.Priority,
		UpdatedAt:  time.Now(),
		UpdaterUid: currentCtx.Payload.UserId,
	}
	if currentCtx.Payload.Role == model.UserRoleAdmin || currentCtx.Payload.Role == model.UserRoleSuperAdmin {
		// 管理员操作
		if _, err := model.UpdateArticleById(db.Mysql, articleID, article); err != nil {
			c.JSON(http.StatusInternalServerError, e.IP(e.IArticleUpdateById, e.MArticleErr, e.ArticleUpdateByIdErr, err))
			return
		}
	} else {
		// 创建者操作
		if _, err := model.UpdateArticleByIdOfSelf(db.Mysql, articleID, currentCtx.Payload.UserId, article); err != nil {
			c.JSON(http.StatusInternalServerError, e.IP(e.IArticleUpdateById, e.MArticleErr, e.ArticleUpdateByIdSelfErr, err))
			return
		}
	}
	c.JSON(http.StatusOK, OKResponse{
		Code: e.OK,
	})
}

func ArticleDeleteByIdHandler(c *gin.Context) {
	articleID, err := strconv.ParseUint(c.Param("article_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, e.IP(e.IArticleDeleteById, e.MParamsErr, e.ParamsInvalidArticleID, err))
		return
	}
	currentCtx := middleware.GetCurrentAuth(c)
	if currentCtx == nil {
		c.JSON(http.StatusBadRequest, e.I(e.IArticleDeleteById, e.MAuthErr, e.AuthGetCurrentErr))
		return
	}
	if currentCtx.Payload.Role == model.UserRoleAdmin || currentCtx.Payload.Role == model.UserRoleSuperAdmin {
		if _, err := model.DeleteArticleById(db.Mysql, articleID); err != nil {
			c.JSON(http.StatusInternalServerError, e.IP(e.IArticleDeleteById, e.MArticleErr, e.ArticleDeleteByIdErr, err))
			return
		}
	} else {
		if _, err := model.DeleteArticleByIdOfSelf(db.Mysql, articleID, currentCtx.Payload.UserId); err != nil {
			c.JSON(http.StatusInternalServerError, e.IP(e.IArticleDeleteById, e.MArticleErr, e.ArticleDeleteByIdSelfErr, err))
			return
		}
	}

	c.JSON(http.StatusOK, OKResponse{
		Code: e.OK,
	})
}
