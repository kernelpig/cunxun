package model

import (
	e "wangqingang/cunxun/error"
)

type Article struct {
	ID         int    `column:"id"`
	ColumnId   int    `column:"column_id"`
	Title      string `column:"title"`
	Content    string `column:"content"`
	CreaterUid int    `column:"creater_uid"`
	UpdaterUid int    `column:"updater_uid"`
}

func GetArticleByID(db sqlExec, articleID int) (*Article, error) {
	u := &Article{}
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
