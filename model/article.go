package model

import (
	"time"

	e "wangqingang/cunxun/error"
)

type Article struct {
	ID         int       `json:"id" column:"id"`
	ColumnId   int       `json:"column_id" column:"column_id"`
	Title      string    `json:"title" column:"title"`
	Content    string    `json:"content" column:"content"`
	CreaterUid int       `json:"creater_uid" column:"creater_uid"`
	UpdaterUid int       `json:"updater_uid" column:"updater_uid"`
	CreatedAt  time.Time `json:"created_at" column:"created_at"`
	UpdatedAt  time.Time `json:"updated_at" column:"updated_at"`
}

type ArticleDetailView struct {
	ID           int       `json:"id" column:"id"`
	ColumnId     int       `json:"column_id" column:"column_id"`
	Title        string    `json:"title" column:"title"`
	Content      string    `json:"content" column:"content"`
	CreaterUid   int       `json:"creater_uid" column:"creater_uid"`
	UpdaterUid   int       `json:"updater_uid" column:"updater_uid"`
	CreatedAt    time.Time `json:"created_at" column:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" column:"updated_at"`
	Nickname     string    `json:"nickname" column:"nickname"`
	CommentCount int       `json:"comment_count" column:"comment_count"`
}

type ArticleListView struct {
	ID           int       `json:"id" column:"id"`
	ColumnId     int       `json:"column_id" column:"column_id"`
	Title        string    `json:"title" column:"title"`
	Content      string    `json:"content" column:"content"`
	CreaterUid   int       `json:"creater_uid" column:"creater_uid"`
	UpdaterUid   int       `json:"updater_uid" column:"updater_uid"`
	CreatedAt    time.Time `json:"created_at" column:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" column:"updated_at"`
	Nickname     string    `json:"nickname" column:"nickname"`
	CommentCount int       `json:"comment_count" column:"comment_count"`
}

func GetArticleByID(db sqlExec, articleID int) (*ArticleDetailView, error) {
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
	id, err := SQLInsert(db, article)
	if err != nil {
		if isDBDuplicateErr(err) {
			return nil, e.SP(e.MArticleErr, e.ArticleAlreadyExist, err)
		}
		return nil, e.SP(e.MArticleErr, e.ArticleCreateErr, err)
	}
	article.ID = int(id)
	return article, nil
}

func GetArticleList(db sqlExec, where map[string]interface{}, orderBy string, pageSize, pageNum int) ([]*ArticleListView, bool, error) {
	var list []*ArticleListView

	// 初始化缓冲区
	var modelBuf = make([]interface{}, 0)
	for i := 0; i < pageSize; i++ {
		modelBuf = append(modelBuf, &ArticleListView{})
	}

	// 每次只取pageSize个
	isOver, err := SQLQueryRows(db, &modelBuf, where, orderBy, pageSize, pageNum)
	if err != nil {
		return nil, true, e.SP(e.MArticleErr, e.ArticleGetListErr, err)
	}
	for _, item := range modelBuf {
		if model, ok := item.(*ArticleListView); ok {
			list = append(list, model)
		}
	}
	return list, isOver, nil
}