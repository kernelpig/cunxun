package model

import (
	"time"

	e "wangqingang/cunxun/error"
	"wangqingang/cunxun/id"
)

type Article struct {
	ID         uint64    `json:"id" column:"id"`
	ColumnId   uint64    `json:"column_id" column:"column_id"`
	Title      string    `json:"title" column:"title"`
	Content    string    `json:"content" column:"content"`
	CreaterUid uint64    `json:"creater_uid" column:"creater_uid"`
	UpdaterUid uint64    `json:"updater_uid" column:"updater_uid"`
	CreatedAt  time.Time `json:"created_at" column:"created_at"`
	UpdatedAt  time.Time `json:"updated_at" column:"updated_at"`
}

type ArticleDetailView struct {
	ID           uint64    `json:"id" column:"id"`
	ColumnId     uint64    `json:"column_id" column:"column_id"`
	Title        string    `json:"title" column:"title"`
	Content      string    `json:"content" column:"content"`
	CreaterUid   uint64    `json:"creater_uid" column:"creater_uid"`
	UpdaterUid   uint64    `json:"updater_uid" column:"updater_uid"`
	CreatedAt    time.Time `json:"created_at" column:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" column:"updated_at"`
	Nickname     string    `json:"nickname" column:"nickname"`
	CommentCount int       `json:"comment_count" column:"comment_count"`
}

func GetArticleByID(db sqlExec, articleID uint64) (*ArticleDetailView, error) {
	u := &ArticleDetailView{}
	isFound, err := SQLQueryRow(db, u, map[string]interface{}{
		"id": articleID,
	})
	if err != nil {
		return nil, e.SP(e.MArticleErr, e.ArticleGetErr, err)
	} else if !isFound {
		return nil, nil
	} else {
		return u, nil
	}
}

func CreateArticle(db sqlExec, article *Article) (*Article, error) {
	// 未设置ID使用自动生成的id, 主要是考虑到人工设置特殊的ID场景
	if article.ID == 0 {
		id, err := id.Generate()
		if err != nil {
			return nil, err
		}
		article.ID = id
	}
	_, err := SQLInsert(db, article)
	if err != nil {
		if isDBDuplicateErr(err) {
			return nil, e.SP(e.MArticleErr, e.ArticleAlreadyExist, err)
		}
		return nil, e.SP(e.MArticleErr, e.ArticleCreateErr, err)
	}
	return article, nil
}

func GetArticleList(db sqlExec, where map[string]interface{}, orderBy string, pageSize, pageNum int) ([]*ArticleDetailView, bool, error) {
	var list []*ArticleDetailView

	// 初始化缓冲区
	var modelBuf = make([]interface{}, 0)
	for i := 0; i < pageSize; i++ {
		modelBuf = append(modelBuf, &ArticleDetailView{})
	}

	// 每次只取pageSize个
	isOver, err := SQLQueryRows(db, &modelBuf, where, []string{orderBy}, pageSize, pageNum)
	if err != nil {
		return nil, true, e.SP(e.MArticleErr, e.ArticleGetListErr, err)
	}
	for _, item := range modelBuf {
		if model, ok := item.(*ArticleDetailView); ok {
			list = append(list, model)
		}
	}
	return list, isOver, nil
}

func UpdateArticleList(db sqlExec, wheres map[string]interface{}, valueWillSet *Article) (int64, error) {
	count, err := SQLUpdate(db, valueWillSet, wheres)
	if err != nil {
		return 0, e.SP(e.MArticleErr, e.ArticleGetErr, err)
	} else {
		return count, nil
	}
}

func DeleteArticleList(db sqlExec, wheres map[string]interface{}) (int64, error) {
	count, err := SQLDelete(db, &Article{}, wheres)
	if err != nil {
		return 0, e.SP(e.MArticleErr, e.ArticleDeleteErr, err)
	} else {
		return count, nil
	}
}

func UpdateArticleById(db sqlExec, articleId uint64, valueWillSet *Article) (int64, error) {
	return UpdateArticleList(db, map[string]interface{}{"id": articleId}, valueWillSet)
}

func UpdateArticleByIdOfSelf(db sqlExec, articleId, userID uint64, valueWillSet *Article) (int64, error) {
	return UpdateArticleList(db, map[string]interface{}{"id": articleId, "creater_uid": userID}, valueWillSet)
}

func DeleteArticleById(db sqlExec, articleId uint64) (int64, error) {
	return DeleteArticleList(db, map[string]interface{}{"id": articleId})
}

func DeleteArticleByIdOfSelf(db sqlExec, articleId, userId uint64) (int64, error) {
	return DeleteArticleList(db, map[string]interface{}{"id": articleId, "creater_uid": userId})
}