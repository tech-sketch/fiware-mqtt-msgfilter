/*
Package checker : authorize and authenticate HTTP Request using HTTP Header.

	license: Apache license 2.0
	copyright: Nobuyuki Matsui <nobuyuki.matsui@gmail.com>
*/
package checker

import (
	"testing"
	"time"

	"github.com/coreos/etcd/client"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/tech-sketch/fiware-mqtt-msgfilter/mock"
)

type setUpObj struct {
	client client.Client
	kapi   *mock.MockKeysAPI
	ctrl   *gomock.Controller
}

func setUpMutex(t *testing.T) (*setUpObj, func()) {
	t.Helper()
	ctrl := gomock.NewController(t)
	kapi := mock.NewMockKeysAPI(ctrl)

	GetNewKeysAPI = func(c client.Client) client.KeysAPI {
		return kapi
	}

	cfg := client.Config{
		Endpoints:               []string{"http://127.0.0.1:2379"},
		Transport:               client.DefaultTransport,
		HeaderTimeoutPerRequest: time.Second,
	}
	c, err := client.New(cfg)
	assert.NotNil(t, c)
	assert.NoError(t, err)

	obj := &setUpObj{
		client: c,
		kapi:   kapi,
		ctrl:   ctrl,
	}
	tearDown := func() {
		ctrl.Finish()
	}
	return obj, tearDown
}

func TestMutexLockSuccess(t *testing.T) {
	assert := assert.New(t)
	obj, tearDown := setUpMutex(t)
	defer tearDown()

	mutex, err := newMutex("key", 60, obj.client)
	assert.NotNil(mutex)
	assert.NoError(err)

	options := &client.SetOptions{
		PrevExist: client.PrevNoExist,
		TTL:       time.Second * time.Duration(60),
	}
	obj.kapi.EXPECT().Set(mutex.ctx, "/key", mutex.id, options).Return(nil, nil)
	err = mutex.Lock()
	assert.NoError(err)

	obj.kapi.EXPECT().Delete(mutex.ctx, "/key", nil).Return(nil, nil)
	err = mutex.Unlock()
	assert.NoError(err)
}

func TestMutexLockSuccessAfterExpire(t *testing.T) {
	assert := assert.New(t)
	obj, tearDown := setUpMutex(t)
	defer tearDown()

	mutex, err := newMutex("key", 60, obj.client)
	assert.NotNil(mutex)
	assert.NoError(err)

	options := &client.SetOptions{
		PrevExist: client.PrevNoExist,
		TTL:       time.Second * time.Duration(60),
	}
	nodeExist := client.Error{
		Code:    client.ErrorCodeNodeExist,
		Message: "test message",
		Cause:   "test cause",
		Index:   0,
	}
	resp := &client.Response{
		Action:    "expire",
		Node:      nil,
		PrevNode:  nil,
		Index:     0,
		ClusterID: "",
	}
	watcher := mock.NewMockWatcher(obj.ctrl)

	gomock.InOrder(
		obj.kapi.EXPECT().Set(mutex.ctx, "/key", mutex.id, options).Return(nil, nodeExist),
		obj.kapi.EXPECT().Get(mutex.ctx, "/key", nil).Return(resp, nil),
		obj.kapi.EXPECT().Watcher("/key", &client.WatcherOptions{AfterIndex: 0, Recursive: false}).Return(watcher),
		watcher.EXPECT().Next(mutex.ctx).Return(resp, nil),
		obj.kapi.EXPECT().Set(mutex.ctx, "/key", mutex.id, options).Return(nil, nil),
	)
	err = mutex.Lock()
	assert.NoError(err)

	obj.kapi.EXPECT().Delete(mutex.ctx, "/key", nil).Return(nil, nil)
	err = mutex.Unlock()
	assert.NoError(err)
}
