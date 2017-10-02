package model

import (
	"time"

	e "wangqingang/cunxun/error"
	"wangqingang/cunxun/id"
)

type Comment struct {
	ID         uint64    `json:"id" column:"id"`
	RelateId   uint64    `json:"relate_id" column:"relate_id"`
	Content    string    `json:"content" column:"content"`
	CreaterUid uint64    `json:"creater_uid" column:"creater_uid"`
	CreatedAt  time.Time `json:"created_at" column:"created_at"`
	UpdatedAt  time.Time `json:"updated_at" column:"updated_at"`
}

type CommentDetailView struct {
	ID         uint64    `json:"id" column:"id"`
	RelateId   uint64    `json:"relate_id" column:"relate_id"`
	Content    string    `json:"content" column:"content"`
	CreaterUid uint64    `json:"creater_uid" column:"creater_uid"`
	CreatedAt  time.Time `json:"created_at" column:"created_at"`
	UpdatedAt  time.Time `json:"updated_at" column:"updated_at"`
	Nickname   string    `json:"nickname" column:"nickname"`
}

func CreateComment(db sqlExec, comment *Comment) (*Comment, error) {
	// 未设置ID使用自动生成的id, 主要是考虑到人工设置特殊的ID场景
	if comment.ID == 0 {
		id, err := id.Generate()
		if err != nil {
			return nil, err
		}
		comment.ID = id
	}
	_, err := SQLInsert(db, comment)
	if err != nil {
		if isDBDuplicateErr(err) {
			return nil, e.SP(e.MCommentErr, e.CommentAlreadyExist, err)
		}
		return nil, e.SP(e.MCommentErr, e.CommentCreateErr, err)
	}
	return comment, nil
}

func GetCommentByID(db sqlExec, commentID uint64) (*CommentDetailView, error) {
	u := &CommentDetailView{}
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

func GetCommentList(db sqlExec, where map[string]interface{}, pageSize, pageNum int) ([]*CommentDetailView, bool, error) {
	var list []*CommentDetailView

	// 初始化缓冲区
	var modelBuf = make([]interface{}, 0)
	for i := 0; i < pageSize; i++ {
		modelBuf = append(modelBuf, &CommentDetailView{})
	}

	// 每次只取pageSize个
	isOver, err := SQLQueryRows(db, &modelBuf, where, OrderByIgnore, pageSize, pageNum)
	if err != nil {
		return nil, true, e.SP(e.MCommentErr, e.CommentGetListErr, err)
	}
	for _, item := range modelBuf {
		if model, ok := item.(*CommentDetailView); ok {
			list = append(list, model)
		}
	}
	return list, isOver, nil
}

func UpdateCommentList(db sqlExec, wheres map[string]interface{}, valueWillSet *Comment) (int64, error) {
	count, err := SQLUpdate(db, valueWillSet, wheres)
	if err != nil {
		return 0, e.SP(e.MCommentErr, e.CommentGetErr, err)
	} else {
		return count, nil
	}
}

func DeleteCommentList(db sqlExec, wheres map[string]interface{}) (int64, error) {
	count, err := SQLDelete(db, &Comment{}, wheres)
	if err != nil {
		return 0, e.SP(e.MCommentErr, e.CommentDeleteErr, err)
	} else {
		return count, nil
	}
}

func UpdateCommentById(db sqlExec, commentId uint64, valueWillSet *Comment) (int64, error) {
	return UpdateCommentList(db, map[string]interface{}{"id": commentId}, valueWillSet)
}

func UpdateCommentByIdOfSelf(db sqlExec, commentId, userId uint64, valueWillSet *Comment) (int64, error) {
	return UpdateCommentList(db, map[string]interface{}{"id": commentId, "creater_uid": userId}, valueWillSet)
}

func DeleteCommentById(db sqlExec, commentId uint64) (int64, error) {
	return DeleteCommentList(db, map[string]interface{}{"id": commentId})
}

func DeleteCommentByIdOfSelf(db sqlExec, commentId, userId uint64) (int64, error) {
	return DeleteCommentList(db, map[string]interface{}{"id": commentId, "creater_uid": userId})
}
