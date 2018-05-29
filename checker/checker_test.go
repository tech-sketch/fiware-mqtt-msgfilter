/*
Package checker : authorize and authenticate HTTP Request using HTTP Header.

	license: Apache license 2.0
	copyright: Nobuyuki Matsui <nobuyuki.matsui@gmail.com>
*/
package checker

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/coreos/etcd/client"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/tech-sketch/fiware-mqtt-msgfilter/conf"
	"github.com/tech-sketch/fiware-mqtt-msgfilter/mock"
)

func setUpChecker(t *testing.T) (*mock.MockKeysAPI, func()) {
	t.Helper()
	ctrl := gomock.NewController(t)
	kapi := mock.NewMockKeysAPI(ctrl)

	GetNewKeysAPI = func(c client.Client) client.KeysAPI {
		return kapi
	}
	GetMutexID = func(_ string) string {
		return "mutexID"
	}

	tearDown := func() {
		ctrl.Finish()
	}
	return kapi, tearDown
}

func TestDuplicate(t *testing.T) {
	assert := assert.New(t)
	kapi, tearDown := setUpChecker(t)
	defer tearDown()

	config := conf.NewConfig()
	checker, err := NewChecker(config)

	assert.NotNil(checker)
	assert.NoError(err)

	options := &client.SetOptions{
		PrevExist: client.PrevNoExist,
		TTL:       time.Second * time.Duration(config.LockTTL),
	}

	gomock.InOrder(
		kapi.EXPECT().Set(context.TODO(), "/lock/test", "mutexID", options).Return(nil, nil),
		kapi.EXPECT().Get(context.Background(), "/data/test", nil).Return(nil, nil),
		kapi.EXPECT().Delete(context.TODO(), "/lock/test", nil).Return(nil, nil),
	)
	result, err := checker.IsDuplicate("test")
	assert.True(result)
	assert.NoError(err)
}

func TestNotDuplicate(t *testing.T) {
	assert := assert.New(t)
	kapi, tearDown := setUpChecker(t)
	defer tearDown()

	config := conf.NewConfig()
	checker, err := NewChecker(config)

	assert.NotNil(checker)
	assert.NoError(err)

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

	gomock.InOrder(
		kapi.EXPECT().Set(context.TODO(), "/lock/test", "mutexID", lockOptions).Return(nil, nil),
		kapi.EXPECT().Get(context.Background(), "/data/test", nil).Return(nil, keyNotFound),
		kapi.EXPECT().Set(context.Background(), "/data/test", "duplicate", dataOptions).Return(nil, nil),
		kapi.EXPECT().Delete(context.TODO(), "/lock/test", nil).Return(nil, nil),
	)
	result, err := checker.IsDuplicate("test")
	assert.False(result)
	assert.NoError(err)
}

func TestRaiseError1(t *testing.T) {
	assert := assert.New(t)
	kapi, tearDown := setUpChecker(t)
	defer tearDown()

	config := conf.NewConfig()
	checker, err := NewChecker(config)

	assert.NotNil(checker)
	assert.NoError(err)

	lockOptions := &client.SetOptions{
		PrevExist: client.PrevNoExist,
		TTL:       time.Second * time.Duration(config.LockTTL),
	}
	raisedError := errors.New("error")

	gomock.InOrder(
		kapi.EXPECT().Set(context.TODO(), "/lock/test", "mutexID", lockOptions).Return(nil, raisedError).AnyTimes(),
	)
	result, err := checker.IsDuplicate("test")
	assert.True(result)
	assert.Equal(raisedError, err)
}

func TestRaiseError2(t *testing.T) {
	assert := assert.New(t)
	kapi, tearDown := setUpChecker(t)
	defer tearDown()

	config := conf.NewConfig()
	checker, err := NewChecker(config)

	assert.NotNil(checker)
	assert.NoError(err)

	lockOptions := &client.SetOptions{
		PrevExist: client.PrevNoExist,
		TTL:       time.Second * time.Duration(config.LockTTL),
	}
	raisedError := errors.New("error")

	gomock.InOrder(
		kapi.EXPECT().Set(context.TODO(), "/lock/test", "mutexID", lockOptions).Return(nil, nil),
		kapi.EXPECT().Get(context.Background(), "/data/test", nil).Return(nil, raisedError),
		kapi.EXPECT().Delete(context.TODO(), "/lock/test", nil).Return(nil, nil),
	)
	result, err := checker.IsDuplicate("test")
	assert.True(result)
	assert.Equal(raisedError, err)
}

func TestRaiseError3(t *testing.T) {
	assert := assert.New(t)
	kapi, tearDown := setUpChecker(t)
	defer tearDown()

	config := conf.NewConfig()
	checker, err := NewChecker(config)

	assert.NotNil(checker)
	assert.NoError(err)

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
	raisedError := errors.New("error")

	gomock.InOrder(
		kapi.EXPECT().Set(context.TODO(), "/lock/test", "mutexID", lockOptions).Return(nil, nil),
		kapi.EXPECT().Get(context.Background(), "/data/test", nil).Return(nil, keyNotFound),
		kapi.EXPECT().Set(context.Background(), "/data/test", "duplicate", dataOptions).Return(nil, raisedError),
		kapi.EXPECT().Delete(context.TODO(), "/lock/test", nil).Return(nil, nil),
	)
	result, err := checker.IsDuplicate("test")
	assert.True(result)
	assert.Equal(raisedError, err)
}

func TestRaiseError4(t *testing.T) {
	assert := assert.New(t)
	kapi, tearDown := setUpChecker(t)
	defer tearDown()

	config := conf.NewConfig()
	checker, err := NewChecker(config)

	assert.NotNil(checker)
	assert.NoError(err)

	lockOptions := &client.SetOptions{
		PrevExist: client.PrevNoExist,
		TTL:       time.Second * time.Duration(config.LockTTL),
	}
	raisedError := errors.New("error")

	gomock.InOrder(
		kapi.EXPECT().Set(context.TODO(), "/lock/test", "mutexID", lockOptions).Return(nil, nil),
		kapi.EXPECT().Get(context.Background(), "/data/test", nil).Return(nil, nil),
		kapi.EXPECT().Delete(context.TODO(), "/lock/test", nil).Return(nil, raisedError).AnyTimes(),
	)
	result, err := checker.IsDuplicate("test")
	assert.True(result)
	assert.NoError(err)
}
