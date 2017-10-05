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
	Avatar   string `json:"avatar"`
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

type Article struct {
	ID           string    `json:"id"`
	ColumnId     string    `json:"column_id"`
	Title        string    `json:"title"`
	Content      string    `json:"content"`
	CreaterUid   string    `json:"creater_uid"`
	UpdaterUid   string    `json:"updater_uid"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	Nickname     string    `json:"nickname"`
	CommentCount int       `json:"comment_count"`
}

type ArticleGetListResponse struct {
	Code int        `json:"code"`
	End  bool       `json:"end"`
	List []*Article `json:"list"`
}

type Carpooling struct {
	ID          string    `json:"id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	CreaterUid  string    `json:"creater_uid"`
	UpdaterUid  string    `json:"updater_uid"`
	FromCity    string    `json:"from_city"`
	ToCity      string    `json:"to_city"`
	DepartTime  time.Time `json:"depart_time"`
	PeopleCount int       `json:"people_count"`
	Contact     string    `json:"contact"`
	Status      int       `json:"status"`
	Remark      string    `json:"remark"`
	Nickname    string    `json:"nickname"`
}

type CarpoolingGetListResponse struct {
	Code int           `json:"code"`
	End  bool          `json:"end"`
	List []*Carpooling `json:"list"`
}

type Comment struct {
	ID         string    `json:"id"`
	RelateId   string    `json:"relate_id"`
	Content    string    `json:"content"`
	CreaterUid string    `json:"creater_uid"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
	Nickname   string    `json:"nickname"`
}

type CommentGetListResponse struct {
	Code int        `json:"code"`
	End  bool       `json:"end"`
	List []*Comment `json:"list"`
}

func FormatId(id uint64) string {
	return strconv.FormatUint(id, 10)
}

func m2rUser(m *model.User) *User {
	if m == nil {
		return nil
	}
	r := &User{}
	r.ID = FormatId(m.ID)
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
	r.ID = FormatId(m.ID)
	r.Name = m.Name
	r.CreaterUid = FormatId(m.CreaterUid)
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

func m2rArticle(m *model.ArticleDetailView) *Article {
	if m == nil {
		return nil
	}
	r := &Article{}
	r.ID = FormatId(m.ID)
	r.ColumnId = FormatId(m.ColumnId)
	r.Title = m.Title
	r.Content = m.Content
	r.CreaterUid = FormatId(m.CreaterUid)
	r.UpdaterUid = FormatId(m.UpdaterUid)
	r.CreatedAt = m.CreatedAt
	r.UpdatedAt = m.UpdatedAt
	r.Nickname = m.Nickname
	r.CommentCount = m.CommentCount
	return r
}

func m2rArticleList(ms []*model.ArticleDetailView) []*Article {
	if ms == nil {
		return nil
	}
	rs := make([]*Article, 0)
	for _, m := range ms {
		rs = append(rs, m2rArticle(m))
	}
	return rs
}

func m2rCarpooling(m *model.CarpoolingDetailView) *Carpooling {
	if m == nil {
		return nil
	}
	r := &Carpooling{}
	r.ID = FormatId(m.ID)
	r.CreatedAt = m.CreatedAt
	r.UpdatedAt = m.UpdatedAt
	r.CreaterUid = FormatId(m.CreaterUid)
	r.UpdaterUid = FormatId(m.UpdaterUid)
	r.FromCity = m.FromCity
	r.ToCity = m.ToCity
	r.DepartTime = m.DepartTime
	r.PeopleCount = m.PeopleCount
	r.Contact = m.Contact
	r.Status = m.Status
	r.Remark = m.Remark
	r.Nickname = m.Nickname
	return r
}

func m2rCarpoolingList(ms []*model.CarpoolingDetailView) []*Carpooling {
	if ms == nil {
		return nil
	}
	rs := make([]*Carpooling, 0)
	for _, m := range ms {
		rs = append(rs, m2rCarpooling(m))
	}
	return rs
}

func m2rComment(m *model.CommentDetailView) *Comment {
	if m == nil {
		return nil
	}
	r := &Comment{}
	r.ID = FormatId(m.ID)
	r.RelateId = FormatId(m.RelateId)
	r.Content = m.Content
	r.CreaterUid = FormatId(m.CreaterUid)
	r.CreatedAt = m.CreatedAt
	r.UpdatedAt = m.UpdatedAt
	r.Nickname = m.Nickname
	return r
}

func m2rCommentList(ms []*model.CommentDetailView) []*Comment {
	if ms == nil {
		return nil
	}
	rs := make([]*Comment, 0)
	for _, m := range ms {
		rs = append(rs, m2rComment(m))
	}
	return rs
}
