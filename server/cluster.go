package server

import (
	"github.com/ActiveState/log"
	"fmt"
	"sync"
)

type ClusterConfig struct {
	MbusIp   string `json:"mbusip"`
	Endpoint string `json:"endpoint"`
}

func (c *ClusterConfig) GetNatsUri() string {
	// HACK: Ideally we should be reading NatsUri from
	// cloud_controller config (mbus). we take a shortcut here in
	// order to not have to create a separate ConfDis instance for
	// cloud_controller config (and having to watch it). This will
	// have to change if we switch to clustered version of NATS.
	return fmt.Sprintf("nats://%s:4222/", c.MbusIp)
}

var clusterConfig *Config

func GetClusterConfig() *ClusterConfig {
	once.Do(createClusterConfig)
	return clusterConfig.Config.(*ClusterConfig)
}

var once sync.Once
func createClusterConfig() {
	var err error
	clusterConfig, err = NewConfig("cluster", ClusterConfig{})
	if err != nil {
		log.Fatal(err)
	}
}

// IsMicro returns true if the cluster is configured as a micro cloud.
func (c *ClusterConfig) IsMicro() bool {
	return c.MbusIp == "127.0.0.1"
}
