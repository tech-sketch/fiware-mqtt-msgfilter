/*
Package conf : configuration variables

	license: Apache license 2.0
	copyright: Nobuyuki Matsui <nobuyuki.matsui@gmail.com>
*/
package conf

import (
	"os"
	"regexp"
	"strconv"
)

const (
	listenPort          = "LISTEN_PORT"
	defaultListenPort   = "5001"
	etcdEndpoint        = "ETCD_ENDPOINT"
	defaultEtcdEndpoint = "http://127.0.0.1:2379"
	etcdEndpointRe      = `http://.+:(\d+)`
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
	r := regexp.MustCompile(etcdEndpointRe)
	if !r.MatchString(etcdEndpoint) {
		etcdEndpoint = defaultEtcdEndpoint
	}
	g := r.FindStringSubmatch(etcdEndpoint)
	ip, err := strconv.Atoi(g[1])
	if err != nil || ip < 1 || 65535 < ip {
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
