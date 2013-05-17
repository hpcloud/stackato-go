package server

import (
	"github.com/ActiveState/log"
	"strings"
)

type ClusterConfig struct {
	MbusIp   string `json:"mbusip"`
	Endpoint string `json:"endpoint"`
	// TODO: somehow use NatsUri from cc config
	NatsUri string
}

var clusterConfig *GroupConfig

func GetClusterConfig() *ClusterConfig {
	if clusterConfig == nil {
		log.Fatal("server.Init() not called")
	}
	return clusterConfig.Config.(*ClusterConfig)
}

func Init() {
	var err error
	clusterConfig, err = NewGroupConfig("cluster", ClusterConfig{})
	if err != nil {
		log.Fatal(err)
	}
}

// IsMicro returns true if the cluster is configured as a micro cloud.
func (c *ClusterConfig) IsMicro() bool {
	return strings.Contains(c.NatsUri, "/127.0.0.1:")
}
