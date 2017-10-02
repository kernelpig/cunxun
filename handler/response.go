package handler

import (
	"strconv"
	"time"
	"wangqingang/cunxun/model"
)

type OKResponse struct {
	Code int `json:"code"`
}

type CreateResponse struct {
	Code int    `json:"code"`
	Id   string `json:"id"`
}

type UserLoginResponse struct {
	Code      int    `json:"code"`
	UserRole  int    `json:"user_role"`
	UserId    string `json:"user_id"`
	UserToken string `json:"user_token"`
}

type UserGetInfoResponse struct {
	Code     int    `json:"code"`
	UserId   string `json:"user_id"`
	Nickname string `json:"nickname"`
	Phone    string `json:"phone"`
}

type User struct {
	ID             string    `json:"id"`
	Phone          string    `json:"phone"`
	NickName       string    `json:"nickname"`
	HashedPassword string    `json:"hashed_password"`
	PasswordLevel  int       `json:"password_level"`
	RegisterSource string    `json:"register_source"`
	Avatar         string    `json:"avatar"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
	Role           int       `json:"role"`
}

type UserGetListResponse struct {
	Code int     `json:"code"`
	End  bool    `json:"end"`
	List []*User `json:"list"`
}

type ColumnGetListResponse struct {
	Code int       `json:"code"`
	End  bool      `json:"end"`
	List []*Column `json:"list"`
}

type ImageCreateResponse struct {
	Code int    `json:"code"`
	Link string `json:"link"`
}

type Column struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	CreaterUid  string    `json:"creater_uid"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Nickname    string    `json:"nickname"`
	ColumnCount int       `json:"column_count"`
}

type ArticleListItem struct {
	ID           uint64    `json:"id"`
	ColumnId     uint64    `json:"column_id"`
	Title        string    `json:"title"`
	Content      string    `json:"content"`
	CreaterUid   uint64    `json:"creater_uid"`
	UpdaterUid   uint64    `json:"updater_uid"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	Nickname     string    `json:"nickname"`
	CommentCount int       `json:"comment_count"`
}

func m2rUser(m *model.User) *User {
	if m == nil {
		return nil
	}
	r := &User{}
	r.ID = strconv.FormatUint(m.ID, 10)
	r.Phone = m.Phone
	r.NickName = m.NickName
	r.HashedPassword = m.HashedPassword
	r.PasswordLevel = m.PasswordLevel
	r.RegisterSource = m.RegisterSource
	r.Avatar = m.Avatar
	r.CreatedAt = m.CreatedAt
	r.UpdatedAt = m.UpdatedAt
	r.Role = m.Role
	return r
}

func m2rUserList(ms []*model.User) []*User {
	if ms == nil {
		return nil
	}
	rs := make([]*User, 0)
	for _, m := range ms {
		rs = append(rs, m2rUser(m))
	}
	return rs
}

func m2rColumn(m *model.ColumnDetailView) *Column {
	if m == nil {
		return nil
	}
	r := &Column{}
	r.ID = strconv.FormatUint(m.ID, 10)
	r.Name = m.Name
	r.CreaterUid = strconv.FormatUint(m.CreaterUid, 10)
	r.CreatedAt = m.CreatedAt
	r.UpdatedAt = m.UpdatedAt
	r.Nickname = m.Nickname
	r.ColumnCount = m.ColumnCount
	return r
}

func m2rColumnList(ms []*model.ColumnDetailView) []*Column {
	if ms == nil {
		return nil
	}
	rs := make([]*Column, 0)
	for _, m := range ms {
		rs = append(rs, m2rColumn(m))
	}
	return rs
}
