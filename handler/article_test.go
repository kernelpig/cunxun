package handler

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/gavv/httpexpect"
	"github.com/stretchr/testify/assert"

	"wangqingang/cunxun/common"
	"wangqingang/cunxun/error"
	"wangqingang/cunxun/model"
	"wangqingang/cunxun/test"
)

func testArticleCreate(t *testing.T, e *httpexpect.Expect, xToken string, request *ArticleCreateRequest) string {
	assert := assert.New(t)
	resp := e.POST("/api/article/").
		WithHeader(common.AuthHeaderKey, xToken).
		WithJSON(request).
		Expect().Status(http.StatusOK)

	object := &CreateResponse{}
	err := json.Unmarshal([]byte(resp.Body().Raw()), object)
	assert.Nil(err)
	assert.Equal(error.OK, object.Code)

	return object.Id
}

func testArticleCreateHandler(t *testing.T, e *httpexpect.Expect) {
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
	testArticleCreate(t, e, xToken, createArticleRequest)
}

func testArticleGetList(t *testing.T, e *httpexpect.Expect, columnID string, pageNum, pageSize int) []*Article {
	assert := assert.New(t)

	resp := e.GET("/api/article/").
		WithQuery("creater_uid", 0).
		WithQuery("column_id", columnID).
		WithQuery("order_by", model.OrderByCommentCount).
		WithQuery("page_num", pageNum).
		WithQuery("page_size", pageSize).
		Expect().Status(http.StatusOK)

	object := &ArticleGetListResponse{}
	err := json.Unmarshal([]byte(resp.Body().Raw()), &object)
	assert.Nil(err)
	assert.Equal(error.OK, object.Code)

	return object.List
}

func testArticleGetListHandler(t *testing.T, e *httpexpect.Expect) {
	test.InitTestCaseEnv(t)
	assert := assert.New(t)

	xToken := testNormalUserLogin(t, e)
	xSuperToken := testSuperAdminLogin(t, e)

	createColumnRequest := &ColumnCreateRequest{
		Name: test.GenRandString(),
	}
	columnID := testColumnCreate(t, e, xSuperToken, createColumnRequest)

	for i := 0; i < 5; i++ {
		createArticleRequest := &ArticleCreateRequest{
			ColumnId: columnID,
			Title:    test.GenRandString(),
			Content:  test.GenRandString() + test.GenRandString(),
		}
		testArticleCreate(t, e, xToken, createArticleRequest)
	}

	list := testArticleGetList(t, e, columnID, 1, 10)
	assert.Equal(5, len(list))
}

func testArticleGet(t *testing.T, e *httpexpect.Expect, articleID string) *Article {
	assert := assert.New(t)

	resp := e.GET("/api/article/{article_id}").
		WithPath("article_id", articleID).
		Expect().Status(http.StatusOK)

	object := &ArticleGetListResponse{}
	err := json.Unmarshal([]byte(resp.Body().Raw()), object)
	assert.Nil(err)
	assert.Equal(error.OK, object.Code)
	assert.Equal(1, len(object.List))

	return object.List[0]
}

func testArticleGetHandler(t *testing.T, e *httpexpect.Expect) {
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
	article := testArticleGet(t, e, articleID)

	assert.NotNil(article)
	assert.Equal(articleID, article.ID)
}

func testArticleUpdateById(t *testing.T, e *httpexpect.Expect, xToken string, articleId string, request *ArticleUpdateRequest) {
	resp := e.PUT("/api/article/{article_id}").
		WithPath("article_id", articleId).
		WithHeader(common.AuthHeaderKey, xToken).
		WithJSON(request).
		Expect().Status(http.StatusOK)

	respObj := resp.JSON().Object()
	respObj.Value("code").Number().Equal(error.OK)
}

func testArticleUpdateByIdHandler(t *testing.T, e *httpexpect.Expect) {
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
	articleId := testArticleCreate(t, e, xToken, createArticleRequest)

	updateArticleRequest := &ArticleUpdateRequest{
		ColumnId: columnID,
		Title:    test.GenRandString(),
		Content:  test.GenRandString() + test.GenRandString(),
	}
	testArticleUpdateById(t, e, xToken, articleId, updateArticleRequest)
}

func testArticleDeleteById(t *testing.T, e *httpexpect.Expect, xToken string, articleId string) {
	resp := e.DELETE("/api/article/{article_id}").
		WithPath("article_id", articleId).
		WithHeader(common.AuthHeaderKey, xToken).
		Expect().Status(http.StatusOK)

	respObj := resp.JSON().Object()
	respObj.Value("code").Number().Equal(error.OK)
}

func testArticleDeleteByIdHandler(t *testing.T, e *httpexpect.Expect) {
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
	articleId := testArticleCreate(t, e, xToken, createArticleRequest)

	testArticleDeleteById(t, e, xToken, articleId)
}
