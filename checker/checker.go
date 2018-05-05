/*
Package checker : authorize and authenticate HTTP Request using HTTP Header.

	license: Apache license 2.0
	copyright: Nobuyuki Matsui <nobuyuki.matsui@gmail.com>
*/
package checker

import (
	"context"
	"fmt"
	"time"

	"github.com/coreos/etcd/client"

	"github.com/tech-sketch/fiware-mqtt-msgfilter/conf"
	"github.com/tech-sketch/fiware-mqtt-msgfilter/utils"
)

/*
Checker : a struct to check message duplication using etcd
*/
type Checker struct {
	client client.Client
	config *conf.Config
}

/*
NewChecker : a factory method to create Checker.
*/
func NewChecker(config *conf.Config) (*Checker, error) {
	cfg := client.Config{
		Endpoints:               []string{config.EtcdEndpoint},
		Transport:               client.DefaultTransport,
		HeaderTimeoutPerRequest: time.Second,
	}
	c, err := client.New(cfg)
	if err != nil {
		return nil, err
	}

	checker := &Checker{
		client: c,
		config: config,
	}
	return checker, nil
}

/*
IsDuplicate : check whether the artument message is duplicated.
*/
func (c *Checker) IsDuplicate(message string) (bool, error) {
	logger := utils.NewLogger("isDuplicate")

	lockKey := fmt.Sprintf("/lock/%s", message)
	logger.Debugf("lockKey = %s", lockKey)

	m, err := newMutex(lockKey, c.config.LockTTL, c.client)
	if err != nil {
		logger.Errorf("newMutex failed: %s", err.Error())
		return true, err
	}
	err = m.Lock()
	if err != nil {
		logger.Errorf("mutex.Lock failed: %s", err.Error())
		return true, err
	}
	defer m.Unlock()

	dataKey := fmt.Sprintf("/data/%s", message)
	logger.Debugf("dataKey = %s", dataKey)

	kapi := client.NewKeysAPI(c.client)
	_, err = kapi.Get(context.Background(), dataKey, nil)
	if err != nil {
		setOptions := &client.SetOptions{
			PrevExist: client.PrevNoExist,
			TTL:       time.Second * time.Duration(c.config.DataTTL),
		}
		_, err = kapi.Set(context.Background(), dataKey, "duplicate", setOptions)
		if err != nil {
			logger.Errorf("etcd set failed: %s", err.Error())
			return true, err
		}
		logger.Debugf("%s is not duplicate", message)
		return false, nil
	}
	logger.Debugf("%s is duplicate", message)
	return true, nil
}
