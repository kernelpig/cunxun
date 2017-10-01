package handler

type UserSignupResponse struct {
	Code   int    `json:"code"`
	UserId string `json:"user_id"`
}

type ColumnCreateResponse struct {
	Code     int    `json:"code"`
	ColumnId string `json:"column_id"`
}

type ArticleCreateResponse struct {
	Code      int    `json:"code"`
	ArticleId string `json:"article_id"`
}

type CommentCreateResponse struct {
	Code      int    `json:"code"`
	CommentId string `json:"article_id"`
}

type CarpoolingCreateResponse struct {
	Code         int    `json:"code"`
	CarpoolingId string `json:"carpooling_id"`
}
