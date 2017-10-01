package handler

type UserSignupResponse struct {
	Code   int    `json:"code"`
	UserId uint64 `json:"user_id"`
}

type ColumnCreateResponse struct {
	Code     int    `json:"code"`
	ColumnId uint64 `json:"column_id"`
}

type ArticleCreateResponse struct {
	Code      int    `json:"code"`
	ArticleId uint64 `json:"article_id"`
}

type CommentCreateResponse struct {
	Code      int    `json:"code"`
	CommentId uint64 `json:"article_id"`
}

type CarpoolingCreateResponse struct {
	Code         int    `json:"code"`
	CarpoolingId uint64 `json:"carpooling_id"`
}
