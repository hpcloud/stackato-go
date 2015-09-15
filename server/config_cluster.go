package server

import (
	"sync"

	"github.com/hpcloud/log"
)

type ClusterConfig struct {
	MbusIp   string `json:"mbusip"`
	Endpoint string `json:"endpoint"`
}

var clusterConfig *Config
var onceClusterConfig sync.Once

func GetClusterConfig() *ClusterConfig {
	onceClusterConfig.Do(createClusterConfig)
	return clusterConfig.GetConfig().(*ClusterConfig)
}

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

func (c ClusterConfig) CurrentNodeId() (string, error) {
	if c.IsMicro() {
		return "127.0.0.1", nil
	} else {
		return LocalIP()
	}
}

// GetMbusIP returns a docker-friendly value for MbusIP
func (c *ClusterConfig) GetMbusIP() string {
	if InsideDocker() {
		if c.MbusIp == "127.0.0.1" {
			dockerHost, err := GetDockerHostIp()
			if err != nil {
				log.Fatal(err)
			}
			return dockerHost
		}
	}
	return c.MbusIp
}
