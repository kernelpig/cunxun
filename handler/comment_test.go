package handler

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/gavv/httpexpect"
	"github.com/stretchr/testify/assert"

	"wangqingang/cunxun/common"
	"wangqingang/cunxun/error"
	"wangqingang/cunxun/test"
)

func testCommentCreate(t *testing.T, e *httpexpect.Expect, xToken string, request *CommentCreateRequest) string {
	assert := assert.New(t)
	resp := e.POST("/api/comment/").
		WithHeader(common.AuthHeaderKey, xToken).
		WithJSON(request).
		Expect().Status(http.StatusOK)

	object := &CreateResponse{}
	err := json.Unmarshal([]byte(resp.Body().Raw()), object)
	assert.Nil(err)
	assert.Equal(error.OK, object.Code)

	return object.Id
}

func testCommentCreateHandler(t *testing.T, e *httpexpect.Expect) {
	test.InitTestCaseEnv(t)

	xToken := testNormalUserLogin(t, e)
	xSuperToken := testSuperAdminLogin(t, e)

	createColumnRequest := &ColumnCreateRequest{
		Name: test.GenRandString(),
	}
	columnID := testColumnCreate(t, e, xSuperToken, createColumnRequest)

	createArticleRequest := &ArticleCreateRequest{
		ColumnId: columnID,
		Title:    test.GenRandString(),
		Content:  test.GenRandString() + test.GenRandString(),
	}
	articleID := testArticleCreate(t, e, xToken, createArticleRequest)

	createCommentRequest := &CommentCreateRequest{
		RelateId: articleID,
		Content:  test.GenRandString() + test.GenRandString(),
	}
	testCommentCreate(t, e, xToken, createCommentRequest)
}

func testCommentGetList(t *testing.T, e *httpexpect.Expect, relateID string, pageNum, pageSize int) []*Comment {
	assert := assert.New(t)

	resp := e.GET("/api/comment/").
		WithQuery("creater_uid", 0).
		WithQuery("relate_id", relateID).
		WithQuery("page_num", pageNum).
		WithQuery("page_size", pageSize).
		Expect().Status(http.StatusOK)

	object := &CommentGetListResponse{}
	err := json.Unmarshal([]byte(resp.Body().Raw()), &object)
	assert.Nil(err)
	assert.Equal(error.OK, object.Code)

	return object.List
}

func testCommentGetListHandler(t *testing.T, e *httpexpect.Expect) {
	test.InitTestCaseEnv(t)
	assert := assert.New(t)

	xToken := testNormalUserLogin(t, e)
	xSuperToken := testSuperAdminLogin(t, e)

	createColumnRequest := &ColumnCreateRequest{
		Name: test.GenRandString(),
	}
	columnID := testColumnCreate(t, e, xSuperToken, createColumnRequest)

	createArticleRequest := &ArticleCreateRequest{
		ColumnId: columnID,
		Title:    test.GenRandString(),
		Content:  test.GenRandString() + test.GenRandString(),
	}
	articleID := testArticleCreate(t, e, xToken, createArticleRequest)

	for i := 0; i < 5; i++ {
		createCommentRequest := &CommentCreateRequest{
			RelateId: articleID,
			Content:  test.GenRandString() + test.GenRandString(),
		}
		testCommentCreate(t, e, xToken, createCommentRequest)
	}

	list := testCommentGetList(t, e, articleID, 1, 10)
	assert.Equal(5, len(list))
}

func testCommentGet(t *testing.T, e *httpexpect.Expect, commentID string) *Comment {
	assert := assert.New(t)

	resp := e.GET("/api/comment/{comment_id}").
		WithPath("comment_id", commentID).
		Expect().Status(http.StatusOK)

	object := &CommentGetListResponse{}
	err := json.Unmarshal([]byte(resp.Body().Raw()), &object)
	assert.Nil(err)
	assert.Equal(error.OK, object.Code)
	assert.Equal(1, len(object.List))

	return object.List[0]
}

func testCommentGetHandler(t *testing.T, e *httpexpect.Expect) {
	test.InitTestCaseEnv(t)
	assert := assert.New(t)

	xToken := testNormalUserLogin(t, e)
	xSuperToken := testSuperAdminLogin(t, e)

	createColumnRequest := &ColumnCreateRequest{
		Name: test.GenRandString(),
	}
	columnID := testColumnCreate(t, e, xSuperToken, createColumnRequest)

	createArticleRequest := &ArticleCreateRequest{
		ColumnId: columnID,
		Title:    test.GenRandString(),
		Content:  test.GenRandString() + test.GenRandString(),
	}
	articleID := testArticleCreate(t, e, xToken, createArticleRequest)

	createCommentRequest := &CommentCreateRequest{
		RelateId: articleID,
		Content:  test.GenRandString() + test.GenRandString(),
	}
	commentID := testCommentCreate(t, e, xToken, createCommentRequest)
	comment := testCommentGet(t, e, commentID)

	assert.NotNil(comment)
	assert.Equal(commentID, comment.ID)
}

func testCommentUpdateById(t *testing.T, e *httpexpect.Expect, xToken string, commentId string, request *CommentCreateRequest) {
	resp := e.PUT("/api/comment/{comment_id}").
		WithPath("comment_id", commentId).
		WithHeader(common.AuthHeaderKey, xToken).
		WithJSON(request).
		Expect().Status(http.StatusOK)

	respObj := resp.JSON().Object()
	respObj.Value("code").Number().Equal(error.OK)
}

func testCommentUpdateByIdHandler(t *testing.T, e *httpexpect.Expect) {
	test.InitTestCaseEnv(t)

	xToken := testNormalUserLogin(t, e)
	xSuperToken := testSuperAdminLogin(t, e)

	createColumnRequest := &ColumnCreateRequest{
		Name: test.GenRandString(),
	}
	columnID := testColumnCreate(t, e, xSuperToken, createColumnRequest)

	createArticleRequest := &ArticleCreateRequest{
		ColumnId: columnID,
		Title:    test.GenRandString(),
		Content:  test.GenRandString() + test.GenRandString(),
	}
	articleID := testArticleCreate(t, e, xToken, createArticleRequest)

	createCommentRequest := &CommentCreateRequest{
		RelateId: articleID,
		Content:  test.GenRandString() + test.GenRandString(),
	}
	commentId := testCommentCreate(t, e, xToken, createCommentRequest)

	updateCommentRequest := &CommentCreateRequest{
		RelateId: articleID,
		Content:  test.GenRandString() + test.GenRandString(),
	}
	testCommentUpdateById(t, e, xToken, commentId, updateCommentRequest)
}

func testCommentDeleteById(t *testing.T, e *httpexpect.Expect, xToken string, commentId string) {
	resp := e.DELETE("/api/comment/{comment_id}").
		WithPath("comment_id", commentId).
		WithHeader(common.AuthHeaderKey, xToken).
		Expect().Status(http.StatusOK)

	respObj := resp.JSON().Object()
	respObj.Value("code").Number().Equal(error.OK)
}

func testCommentDeleteByIdHandler(t *testing.T, e *httpexpect.Expect) {
	test.InitTestCaseEnv(t)

	xToken := testNormalUserLogin(t, e)
	xSuperToken := testSuperAdminLogin(t, e)

	createColumnRequest := &ColumnCreateRequest{
		Name: test.GenRandString(),
	}
	columnID := testColumnCreate(t, e, xSuperToken, createColumnRequest)

	createArticleRequest := &ArticleCreateRequest{
		ColumnId: columnID,
		Title:    test.GenRandString(),
		Content:  test.GenRandString() + test.GenRandString(),
	}
	articleID := testArticleCreate(t, e, xToken, createArticleRequest)

	createCommentRequest := &CommentCreateRequest{
		RelateId: articleID,
		Content:  test.GenRandString() + test.GenRandString(),
	}
	commentId := testCommentCreate(t, e, xToken, createCommentRequest)

	testCommentDeleteById(t, e, xToken, commentId)
}
