/*
Package checker : authorize and authenticate HTTP Request using HTTP Header.

	license: Apache license 2.0
	copyright: Nobuyuki Matsui <nobuyuki.matsui@gmail.com>
*/
package checker

import (
	"context"
	"errors"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/coreos/etcd/client"

	"github.com/tech-sketch/fiware-mqtt-msgfilter/utils"
)

const (
	defaultTTL   = 60
	defaultTry   = 3
	deleteAction = "delete"
	expireAction = "expire"
)

// GetNewKeysAPI : injection point to mock etcd api
var GetNewKeysAPI = client.NewKeysAPI

// GetMutexID : injection point to mock mutexId
var GetMutexID = func(hostname string) string {
	return fmt.Sprintf("%v-%v-%v", hostname, os.Getpid(), time.Now().Format("20060102-15:04:05.999999999"))
}

// A Mutex is a mutual exclusion lock which is distributed across a cluster.
// original: https://github.com/zieckey/etcdsync
type mutex struct {
	key    string
	id     string // The identity of the caller
	client client.Client
	kapi   client.KeysAPI
	ctx    context.Context
	ttl    time.Duration
	mutex  *sync.Mutex
	logger *utils.Logger
}

// newMutex creates a Mutex with the given key which must be the same
// across the cluster nodes.
// machines are the ectd cluster addresses
func newMutex(key string, ttl int, c client.Client) (*mutex, error) {
	hostname, err := os.Hostname()
	if err != nil {
		return nil, err
	}

	if len(key) == 0 {
		return nil, errors.New("wrong lock key")
	}

	if key[0] != '/' {
		key = "/" + key
	}

	if ttl < 1 {
		ttl = defaultTTL
	}

	return &mutex{
		key:    key,
		id:     GetMutexID(hostname),
		client: c,
		kapi:   GetNewKeysAPI(c),
		ctx:    context.TODO(),
		ttl:    time.Second * time.Duration(ttl),
		mutex:  new(sync.Mutex),
		logger: utils.NewLogger("mutex"),
	}, nil
}

// Lock locks m.
// If the lock is already in use, the calling goroutine
// blocks until the mutex is available.
func (m *mutex) Lock() (err error) {
	m.mutex.Lock()
	for try := 1; try <= defaultTry; try++ {
		err = m.lock()
		if err == nil {
			return nil
		}

		m.logger.Debugf("Lock node %v ERROR %v", m.key, err)
		if try < defaultTry {
			m.logger.Debugf("Try to lock node %v again", m.key, err)
		}
	}
	return err
}

func (m *mutex) lock() (err error) {
	m.logger.Debugf("Trying to create a node : key=%v", m.key)
	setOptions := &client.SetOptions{
		PrevExist: client.PrevNoExist,
		TTL:       m.ttl,
	}
	for {
		resp, err := m.kapi.Set(m.ctx, m.key, m.id, setOptions)
		if err == nil {
			m.logger.Debugf("Create node %v (%v) OK [%q]", m.key, m.id, resp)
			return nil
		}
		m.logger.Debugf("Create node %v failed [%v]", m.key, err)
		e, ok := err.(client.Error)
		if !ok {
			return err
		}

		if e.Code != client.ErrorCodeNodeExist {
			return err
		}

		// Get the already node's value.
		resp, err = m.kapi.Get(m.ctx, m.key, nil)
		if err != nil {
			return err
		}
		m.logger.Debugf("Get node %v OK", m.key)
		watcherOptions := &client.WatcherOptions{
			AfterIndex: resp.Index,
			Recursive:  false,
		}
		watcher := m.kapi.Watcher(m.key, watcherOptions)
		for {
			m.logger.Debugf("Watching %v ...", m.key)
			resp, err = watcher.Next(m.ctx)
			if err != nil {
				return err
			}

			m.logger.Debugf("Received an event : %q", resp)
			if resp.Action == deleteAction || resp.Action == expireAction {
				// break this for-loop, and try to create the node again.
				break
			}
		}
	}
}

// Unlock unlocks m.
// It is a run-time error if m is not locked on entry to Unlock.
//
// A locked Mutex is not associated with a particular goroutine.
// It is allowed for one goroutine to lock a Mutex and then
// arrange for another goroutine to unlock it.
func (m *mutex) Unlock() (err error) {
	defer m.mutex.Unlock()
	for i := 1; i <= defaultTry; i++ {
		var resp *client.Response
		resp, err = m.kapi.Delete(m.ctx, m.key, nil)
		if err == nil {
			m.logger.Debugf("Delete %v OK", m.key)
			return nil
		}
		m.logger.Debugf("Delete %v falied: %q", m.key, resp)
		e, ok := err.(client.Error)
		if ok && e.Code == client.ErrorCodeKeyNotFound {
			return nil
		}
	}
	return err
}
