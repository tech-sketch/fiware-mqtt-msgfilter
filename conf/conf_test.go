/*
Package conf : configuration variables

	license: Apache license 2.0
	copyright: Nobuyuki Matsui <nobuyuki.matsui@gmail.com>
*/
package conf

import (
	"fmt"
	"os"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewConfigNoEnv(t *testing.T) {
	assert := assert.New(t)

	l, _ := strconv.Atoi(defaultLockTTL)
	d, _ := strconv.Atoi(defaultDataTTL)

	expected := &Config{
		ListenPort:   ":" + defaultListenPort,
		EtcdEndpoint: defaultEtcdEndpoint,
		LockTTL:      l,
		DataTTL:      d,
	}

	config := NewConfig()
	assert.Equal(expected, config)
}

func TestNewConfigWithEnv(t *testing.T) {
	assert := assert.New(t)

	dl := ":" + defaultListenPort
	listenPortCases := []struct {
		port     string
		expected string
	}{
		{port: "8888", expected: ":8888"},
		{port: "", expected: dl},
		{port: " ", expected: dl},
		{port: "invalid", expected: dl},
		{port: "-1", expected: dl},
		{port: "65536", expected: dl},
		{port: "nil", expected: dl},
	}

	etcdEndpointCases := []struct {
		endpoint string
		expected string
	}{
		{endpoint: "http://test.example.com:1234", expected: "http://test.example.com:1234"},
		{endpoint: "", expected: defaultEtcdEndpoint},
		{endpoint: " ", expected: defaultEtcdEndpoint},
		{endpoint: "invalid", expected: defaultEtcdEndpoint},
		{endpoint: "http://x:-1", expected: defaultEtcdEndpoint},
		{endpoint: "http://x:65536", expected: defaultEtcdEndpoint},
		{endpoint: "nil", expected: defaultEtcdEndpoint},
	}

	l, _ := strconv.Atoi(defaultLockTTL)
	lockTTLCases := []struct {
		lockTTL  string
		expected int
	}{
		{lockTTL: "0", expected: 0},
		{lockTTL: "123", expected: 123},
		{lockTTL: "", expected: l},
		{lockTTL: " ", expected: l},
		{lockTTL: "invalid", expected: l},
		{lockTTL: "-1", expected: l},
	}

	d, _ := strconv.Atoi(defaultDataTTL)
	dataTTLCases := []struct {
		dataTTL  string
		expected int
	}{
		{dataTTL: "0", expected: 0},
		{dataTTL: "123", expected: 123},
		{dataTTL: "", expected: d},
		{dataTTL: " ", expected: d},
		{dataTTL: "invalid", expected: d},
		{dataTTL: "-1", expected: d},
	}

	for _, p := range listenPortCases {
		for _, e := range etcdEndpointCases {
			for _, l := range lockTTLCases {
				for _, d := range dataTTLCases {
					caseName := fmt.Sprintf(`port:%s|endpoint:%s|lockTTL:%s|dataTTL:%s`, p.port, e.endpoint, l.lockTTL, d.dataTTL)
					t.Run(caseName, func(t *testing.T) {
						if p.port != "nil" {
							os.Setenv(listenPort, p.port)
						}
						if e.endpoint != "nil" {
							os.Setenv(etcdEndpoint, e.endpoint)
						}
						if l.lockTTL != "nil" {
							os.Setenv(lockTTL, l.lockTTL)
						}
						if d.dataTTL != "nil" {
							os.Setenv(dataTTL, d.dataTTL)
						}
						expected := &Config{
							ListenPort:   p.expected,
							EtcdEndpoint: e.expected,
							LockTTL:      l.expected,
							DataTTL:      d.expected,
						}
						config := NewConfig()
						assert.Equal(expected, config)

						os.Unsetenv(listenPort)
						os.Unsetenv(etcdEndpoint)
						os.Unsetenv(lockTTL)
						os.Unsetenv(dataTTL)
					})
				}
			}
		}
	}
}
