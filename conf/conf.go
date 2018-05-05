/*
Package conf : configuration variables

	license: Apache license 2.0
	copyright: Nobuyuki Matsui <nobuyuki.matsui@gmail.com>
*/
package conf

import (
	"os"
	"strconv"
)

const (
	listenPort          = "LISTEN_PORT"
	defaultListenPort   = "5001"
	etcdEndpoint        = "ETCD_ENDPOINT"
	defaultEtcdEndpoint = "http://127.0.0.1:2379"
	lockTTL             = "LOCK_TTL"
	defaultLockTTL      = "10"
	dataTTL             = "DATA_TTL"
	defaultDataTTL      = "600"
)

/*
Config : a struct to hold configuration variables
*/
type Config struct {
	ListenPort   string
	EtcdEndpoint string
	LockTTL      int
	DataTTL      int
}

/*
NewConfig : a factory method to create Config.
*/
func NewConfig() *Config {
	port := os.Getenv(listenPort)
	if len(port) == 0 {
		port = defaultListenPort
	}
	intPort, err := strconv.Atoi(port)
	if err != nil || intPort < 1 || 65535 < intPort {
		port = defaultListenPort
	}

	etcdEndpoint := os.Getenv(etcdEndpoint)
	if len(etcdEndpoint) == 0 {
		etcdEndpoint = defaultEtcdEndpoint
	}

	return &Config{
		ListenPort:   ":" + port,
		EtcdEndpoint: etcdEndpoint,
		LockTTL:      envToPositiveInt(lockTTL, defaultLockTTL),
		DataTTL:      envToPositiveInt(dataTTL, defaultDataTTL),
	}
}

func envToPositiveInt(envKey string, defVar string) int {
	strEnvVar := os.Getenv(envKey)
	if len(strEnvVar) == 0 {
		strEnvVar = defVar
	}
	envVar, err := strconv.Atoi(strEnvVar)
	if err != nil || envVar < 0 {
		envVar, _ = strconv.Atoi(defVar)
	}
	return envVar
}
