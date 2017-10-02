package handler

import (
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"github.com/gavv/httpexpect"
	"github.com/stretchr/testify/assert"

	"wangqingang/cunxun/common"
	"wangqingang/cunxun/error"
	"wangqingang/cunxun/model"
	"wangqingang/cunxun/test"
)

func testCarpoolingCreate(t *testing.T, e *httpexpect.Expect, xToken string, request *CarpoolingCreateRequest) string {
	assert := assert.New(t)
	resp := e.POST("/api/carpooling/").
		WithHeader(common.AuthHeaderKey, xToken).
		WithJSON(request).
		Expect().Status(http.StatusOK)

	object := &CreateResponse{}
	err := json.Unmarshal([]byte(resp.Body().Raw()), object)
	assert.Nil(err)
	assert.Equal(error.OK, object.Code)

	return object.Id
}

func testCarpoolingCreateHandler(t *testing.T, e *httpexpect.Expect) {
	test.InitTestCaseEnv(t)

	xToken := testNormalUserLogin(t, e)

	createCarpoolingRequest := &CarpoolingCreateRequest{
		FromCity:    test.GenRandString(),
		ToCity:      test.GenRandString(),
		DepartTime:  time.Now().Add(time.Duration(5) * time.Second).Unix(),
		PeopleCount: test.GenRandInt(5),
		Contact:     test.GenRandString(),
		Remark:      test.GenRandString() + test.GenRandString(),
	}
	testCarpoolingCreate(t, e, xToken, createCarpoolingRequest)
}

func testCarpoolingGetList(t *testing.T, e *httpexpect.Expect, pageNum, pageSize int) []*Carpooling {
	assert := assert.New(t)

	resp := e.GET("/api/carpooling/").
		WithQuery("creater_uid", 0).
		WithQuery("order_by", model.OrderByCreateDate).
		WithQuery("page_num", pageNum).
		WithQuery("page_size", pageSize).
		Expect().Status(http.StatusOK)

	object := &CarpoolingGetListResponse{}
	err := json.Unmarshal([]byte(resp.Body().Raw()), &object)
	assert.Nil(err)
	assert.Equal(error.OK, object.Code)

	return object.List
}

func testCarpoolingGetListHandler(t *testing.T, e *httpexpect.Expect) {
	test.InitTestCaseEnv(t)
	assert := assert.New(t)

	xToken := testNormalUserLogin(t, e)

	for i := 0; i < 5; i++ {
		createCarpoolingRequest := &CarpoolingCreateRequest{
			FromCity:    test.GenRandString(),
			ToCity:      test.GenRandString(),
			DepartTime:  time.Now().Add(time.Duration(5) * time.Second).Unix(),
			PeopleCount: test.GenRandInt(5),
			Contact:     test.GenRandString(),
			Remark:      test.GenRandString() + test.GenRandString(),
		}
		testCarpoolingCreate(t, e, xToken, createCarpoolingRequest)
	}

	list := testCarpoolingGetList(t, e, 1, 10)
	assert.Equal(5, len(list))
}

func testCarpoolingUpdateById(t *testing.T, e *httpexpect.Expect, xToken string, carpoolingId string, request *CarpoolingUpdateRequest) {
	resp := e.PUT("/api/carpooling/{carpooling_id}").
		WithPath("carpooling_id", carpoolingId).
		WithHeader(common.AuthHeaderKey, xToken).
		WithJSON(request).
		Expect().Status(http.StatusOK)

	respObj := resp.JSON().Object()
	respObj.Value("code").Number().Equal(error.OK)
}

func testCarpoolingUpdateByIdHandler(t *testing.T, e *httpexpect.Expect) {
	test.InitTestCaseEnv(t)

	xToken := testNormalUserLogin(t, e)

	createCarpoolingRequest := &CarpoolingCreateRequest{
		FromCity:    test.GenRandString(),
		ToCity:      test.GenRandString(),
		DepartTime:  time.Now().Add(time.Duration(5) * time.Second).Unix(),
		PeopleCount: test.GenRandInt(5),
		Contact:     test.GenRandString(),
		Remark:      test.GenRandString() + test.GenRandString(),
	}
	carpoolingId := testCarpoolingCreate(t, e, xToken, createCarpoolingRequest)

	updateCarpoolingRequest := &CarpoolingUpdateRequest{
		FromCity:    test.GenRandString(),
		ToCity:      test.GenRandString(),
		DepartTime:  time.Now().Add(time.Duration(5) * time.Second).Unix(),
		PeopleCount: test.GenRandInt(5),
		Contact:     test.GenRandString(),
		Remark:      test.GenRandString() + test.GenRandString(),
	}
	testCarpoolingUpdateById(t, e, xToken, carpoolingId, updateCarpoolingRequest)
}

func testCarpoolingDeleteById(t *testing.T, e *httpexpect.Expect, xToken string, carpoolingId string) {
	resp := e.DELETE("/api/carpooling/{carpooling_id}").
		WithPath("carpooling_id", carpoolingId).
		WithHeader(common.AuthHeaderKey, xToken).
		Expect().Status(http.StatusOK)

	respObj := resp.JSON().Object()
	respObj.Value("code").Number().Equal(error.OK)
}

func testCarpoolingDeleteByIdHandler(t *testing.T, e *httpexpect.Expect) {
	test.InitTestCaseEnv(t)

	xToken := testNormalUserLogin(t, e)

	createCarpoolingRequest := &CarpoolingCreateRequest{
		FromCity:    test.GenRandString(),
		ToCity:      test.GenRandString(),
		DepartTime:  time.Now().Add(time.Duration(5) * time.Second).Unix(),
		PeopleCount: test.GenRandInt(5),
		Contact:     test.GenRandString(),
		Remark:      test.GenRandString() + test.GenRandString(),
	}
	carpoolingId := testCarpoolingCreate(t, e, xToken, createCarpoolingRequest)

	testCarpoolingDeleteById(t, e, xToken, carpoolingId)
}

func testCarpoolingGet(t *testing.T, e *httpexpect.Expect, carpoolingID string) *Carpooling {
	assert := assert.New(t)

	resp := e.GET("/api/carpooling/{carpooling_id}").
		WithPath("carpooling_id", carpoolingID).
		Expect().Status(http.StatusOK)

	object := &CarpoolingGetListResponse{}
	err := json.Unmarshal([]byte(resp.Body().Raw()), &object)
	assert.Nil(err)
	assert.Equal(error.OK, object.Code)
	assert.Equal(1, len(object.List))

	return object.List[0]
}

func testCarpoolingGetHandler(t *testing.T, e *httpexpect.Expect) {
	test.InitTestCaseEnv(t)
	assert := assert.New(t)

	xToken := testNormalUserLogin(t, e)

	createCarpoolingRequest := &CarpoolingCreateRequest{
		FromCity:    test.GenRandString(),
		ToCity:      test.GenRandString(),
		DepartTime:  time.Now().Add(time.Duration(5) * time.Second).Unix(),
		PeopleCount: test.GenRandInt(5),
		Contact:     test.GenRandString(),
		Remark:      test.GenRandString() + test.GenRandString(),
	}
	carpoolingID := testCarpoolingCreate(t, e, xToken, createCarpoolingRequest)
	carpooling := testCarpoolingGet(t, e, carpoolingID)

	assert.NotNil(carpooling)
	assert.Equal(carpoolingID, carpooling.ID)
}
