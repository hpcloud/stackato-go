package server

import (
	"github.com/ActiveState/log"
	"sync"
)

type ClusterConfig struct {
	MbusIp   string `json:"mbusip"`
	Endpoint string `json:"endpoint"`
}

var clusterConfig *Config

func GetClusterConfig() *ClusterConfig {
	once.Do(createClusterConfig)
	return clusterConfig.GetConfig().(*ClusterConfig)
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
