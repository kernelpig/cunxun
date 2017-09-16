package model

import (
	"time"

	e "wangqingang/cunxun/error"
)

type Comment struct {
	ID         int       `json:"id" column:"id"`
	ArticleId  int       `json:"article_id" column:"article_id"`
	Content    string    `json:"content" column:"content"`
	CreaterUid int       `json:"creater_uid" column:"creater_uid"`
	CreatedAt  time.Time `json:"created_at" column:"created_at"`
	UpdatedAt  time.Time `json:"updated_at" column:"updated_at"`
}

type CommentListView struct {
	ID         int       `json:"id" column:"id"`
	ArticleId  int       `json:"article_id" column:"article_id"`
	Content    string    `json:"content" column:"content"`
	CreaterUid int       `json:"creater_uid" column:"creater_uid"`
	CreatedAt  time.Time `json:"created_at" column:"created_at"`
	UpdatedAt  time.Time `json:"updated_at" column:"updated_at"`
	Nickname   string    `json:"nickname" column:"nickname"`
}

func CreateComment(db sqlExec, comment *Comment) (*Comment, error) {
	id, err := SQLInsert(db, comment)
	if err != nil {
		if isDBDuplicateErr(err) {
			return nil, e.SP(e.MCommentErr, e.CommentAlreadyExist, err)
		}
		return nil, e.SP(e.MCommentErr, e.CommentCreateErr, err)
	}
	comment.ID = int(id)
	return comment, nil
}

func GetCommentByID(db sqlExec, commentID int) (*Comment, error) {
	u := &Comment{}
	isFound, err := SQLQueryRow(db, u, map[string]interface{}{
		"id": commentID,
	})
	if err != nil {
		return nil, e.SP(e.MCommentErr, e.CommentGetErr, err)
	} else if !isFound {
		return nil, nil
	} else {
		return u, nil
	}
}

func GetCommentList(db sqlExec, where map[string]interface{}, pageSize, pageNum int) ([]*CommentListView, bool, error) {
	var list []*CommentListView

	// 初始化缓冲区
	var modelBuf = make([]interface{}, 0)
	for i := 0; i < pageSize; i++ {
		modelBuf = append(modelBuf, &CommentListView{})
	}

	// 每次只取pageSize个
	isOver, err := SQLQueryRows(db, &modelBuf, where, OrderByIgnore, pageSize, pageNum)
	if err != nil {
		return nil, true, e.SP(e.MCommentErr, e.CommentGetListErr, err)
	}
	for _, item := range modelBuf {
		if model, ok := item.(*CommentListView); ok {
			list = append(list, model)
		}
	}
	return list, isOver, nil
}
