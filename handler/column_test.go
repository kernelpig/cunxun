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

func testColumnCreate(t *testing.T, e *httpexpect.Expect, xToken string, request *ColumnCreateRequest) uint64 {
	assert := assert.New(t)
	resp := e.POST("/api/column/").
		WithHeader(common.AuthHeaderKey, xToken).
		WithJSON(request).
		Expect().Status(http.StatusOK)

	object := &ColumnCreateResponse{}
	err := json.Unmarshal([]byte(resp.Body().Raw()), object)
	assert.Nil(err)
	assert.Equal(error.OK, object.Code)
	assert.NotZero(object.ColumnId)

	return object.ColumnId
}

func testColumnCreateHandler(t *testing.T, e *httpexpect.Expect) {
	test.InitTestCaseEnv(t)

	xSuperToken := testSuperAdminLogin(t, e)

	createRequest := &ColumnCreateRequest{
		Name: test.GenRandString(),
	}
	testColumnCreate(t, e, xSuperToken, createRequest)
}

func testColumnGetList(t *testing.T, e *httpexpect.Expect) []*model.Column {
	assert := assert.New(t)

	resp := e.GET("/api/column").
		WithQuery("creater_uid", "0").
		Expect().Status(http.StatusOK)

	respObj := resp.JSON().Object()
	respObj.Value("code").Number().Equal(error.OK)

	var result struct {
		Code int             `json:"code"`
		List []*model.Column `json:"list"`
	}
	err := json.Unmarshal([]byte(resp.Body().Raw()), &result)
	assert.Nil(err)

	return result.List
}

func testColumnGetListHandler(t *testing.T, e *httpexpect.Expect) {
	test.InitTestCaseEnv(t)
	assert := assert.New(t)

	xSuperToken := testSuperAdminLogin(t, e)
	for i := 0; i < 2; i++ {
		createRequest := &ColumnCreateRequest{
			Name: test.GenRandString(),
		}
		testColumnCreate(t, e, xSuperToken, createRequest)
	}

	list := testColumnGetList(t, e)
	assert.Equal(2, len(list))
}

func testColumnUpdateById(t *testing.T, e *httpexpect.Expect, xToken string, columnId uint64, request *ColumnUpdateRequest) {
	resp := e.PUT("/api/column/{column_id}").
		WithPath("column_id", columnId).
		WithHeader(common.AuthHeaderKey, xToken).
		WithJSON(request).
		Expect().Status(http.StatusOK)

	respObj := resp.JSON().Object()
	respObj.Value("code").Number().Equal(error.OK)
}

func testColumnUpdateByIdHandler(t *testing.T, e *httpexpect.Expect) {
	test.InitTestCaseEnv(t)
	assert := assert.New(t)

	xSuperToken := testSuperAdminLogin(t, e)
	createRequest := &ColumnCreateRequest{
		Name: test.GenRandString(),
	}
	columnId := testColumnCreate(t, e, xSuperToken, createRequest)
	assert.NotZero(columnId)

	updateRequest := &ColumnUpdateRequest{
		Name: test.GenRandString(),
	}
	testColumnUpdateById(t, e, xSuperToken, columnId, updateRequest)
}

func testColumnDeleteById(t *testing.T, e *httpexpect.Expect, xToken string, columnId uint64) {
	resp := e.DELETE("/api/column/{column_id}").
		WithPath("column_id", columnId).
		WithHeader(common.AuthHeaderKey, xToken).
		Expect().Status(http.StatusOK)

	respObj := resp.JSON().Object()
	respObj.Value("code").Number().Equal(error.OK)
}

func testColumnDeleteByIdHandler(t *testing.T, e *httpexpect.Expect) {
	test.InitTestCaseEnv(t)
	assert := assert.New(t)

	xSuperToken := testSuperAdminLogin(t, e)
	createRequest := &ColumnCreateRequest{
		Name: test.GenRandString(),
	}
	columnId := testColumnCreate(t, e, xSuperToken, createRequest)
	assert.NotZero(columnId)

	testColumnDeleteById(t, e, xSuperToken, columnId)
}
