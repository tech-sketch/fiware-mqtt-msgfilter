/*
Package router : routing http request and check message duplication using Checker.

	license: Apache license 2.0
	copyright: Nobuyuki Matsui <nobuyuki.matsui@gmail.com>
*/
package router

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/coreos/etcd/client"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/tech-sketch/fiware-mqtt-msgfilter/checker"
	"github.com/tech-sketch/fiware-mqtt-msgfilter/conf"
	"github.com/tech-sketch/fiware-mqtt-msgfilter/mock"
)

func setUp(t *testing.T) (func(string, string, string, string, string, bool) (*http.Response, error), func()) {
	t.Helper()
	gin.SetMode(gin.ReleaseMode)
	ctrl := gomock.NewController(t)
	kapi := mock.NewMockKeysAPI(ctrl)

	checker.GetNewKeysAPI = func(c client.Client) client.KeysAPI {
		return kapi
	}
	checker.GetMutexID = func(_ string) string {
		return "mutexID"
	}

	var ts *httptest.Server
	c := http.DefaultClient
	doRequest := func(method string, path string, ctype string, key string, jsonBody string, isDuplicate bool) (*http.Response, error) {
		config := conf.NewConfig()

		lockOptions := &client.SetOptions{
			PrevExist: client.PrevNoExist,
			TTL:       time.Second * time.Duration(config.LockTTL),
		}
		dataOptions := &client.SetOptions{
			PrevExist: client.PrevNoExist,
			TTL:       time.Second * time.Duration(config.DataTTL),
		}
		keyNotFound := client.Error{
			Code:    client.ErrorCodeKeyNotFound,
			Message: "test message",
			Cause:   "test cause",
			Index:   0,
		}

		if isDuplicate {
			gomock.InOrder(
				kapi.EXPECT().Set(context.TODO(), "/lock/"+key, "mutexID", lockOptions).Return(nil, nil),
				kapi.EXPECT().Get(context.Background(), "/data/"+key, nil).Return(nil, nil),
				kapi.EXPECT().Delete(context.TODO(), "/lock/"+key, nil).Return(nil, nil),
			)
		} else {
			gomock.InOrder(
				kapi.EXPECT().Set(context.TODO(), "/lock/"+key, "mutexID", lockOptions).Return(nil, nil),
				kapi.EXPECT().Get(context.Background(), "/data/"+key, nil).Return(nil, keyNotFound),
				kapi.EXPECT().Set(context.Background(), "/data/"+key, "duplicate", dataOptions).Return(nil, nil),
				kapi.EXPECT().Delete(context.TODO(), "/lock/"+key, nil).Return(nil, nil),
			)
		}

		handler, err := NewHandler(config)
		assert.NoError(t, err)
		ts = httptest.NewServer(handler.Engine)
		r, err := http.NewRequest(method, ts.URL+path, bytes.NewBuffer([]byte(jsonBody)))
		if err != nil {
			t.Errorf("NewRequest Error. %v", err)
		}
		if ctype != "" {
			r.Header.Add("content-type", ctype)
		}
		return c.Do(r)
	}
	tearDown := func() {
		ts.Close()
	}
	return doRequest, tearDown
}

func TestDistinctOK(t *testing.T) {
	assert := assert.New(t)
	doRequest, tearDown := setUp(t)
	defer tearDown()

	r, err := doRequest("POST", "/distinct/", "application/json", "a", `{"payload": "a"}`, false)
	assert.Nil(err)
	assert.Equal(http.StatusOK, r.StatusCode)
}

func TestDistinctConflict(t *testing.T) {
	assert := assert.New(t)
	doRequest, tearDown := setUp(t)
	defer tearDown()

	r, err := doRequest("POST", "/distinct/", "application/json", "a", `{"payload": "a"}`, true)
	assert.Nil(err)
	assert.Equal(http.StatusConflict, r.StatusCode)
}

func TestBadRequest(t *testing.T) {
	assert := assert.New(t)
	doRequest, tearDown := setUp(t)
	defer tearDown()

	testCases := []struct {
		cType string
		body  string
	}{
		{cType: "application/json", body: ""},
		{cType: "application/json", body: "payload=a"},
		{cType: "application/json", body: `{"x":"Y"}`},
		{cType: "application/x-www-form-urlencoded", body: `{"payload": "a"}`},
		{cType: "application/x-www-form-urlencoded", body: "payload=a"},
		{cType: "", body: `{"payload": "a"}`},
		{cType: "", body: "payload=a"},
	}
	for _, testCase := range testCases {
		r, err := doRequest("POST", "/distinct/", testCase.cType, "a", testCase.body, false)
		assert.Nil(err)
		assert.Equal(http.StatusBadRequest, r.StatusCode)
	}
}

func TestNotAllowdMethod(t *testing.T) {
	assert := assert.New(t)
	doRequest, tearDown := setUp(t)
	defer tearDown()

	for _, method := range []string{"GET", "PUT", "PATCH", "DELETE"} {
		r, err := doRequest(method, "/distinct/", "application/json", "a", `{"payload": "a"}`, false)
		assert.Nil(err)
		assert.Equal(http.StatusNotFound, r.StatusCode)
	}
}

func TestNotFoundPath(t *testing.T) {
	assert := assert.New(t)
	doRequest, tearDown := setUp(t)
	defer tearDown()

	for _, method := range []string{"GET", "POST", "PUT", "PATCH", "DELETE"} {
		r, err := doRequest(method, "/invalid/", "application/json", "a", `{"payload": "a"}`, false)
		assert.Nil(err)
		assert.Equal(http.StatusNotFound, r.StatusCode)
	}
}
