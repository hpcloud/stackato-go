package server

import (
	"github.com/ActiveState/log"
	"sync"
)

type NodeInfo struct {
	Roles map[string]string `json:"roles"` // role name -> role status
}

// NodeConfig corresponds to Kato's /node config, which is a hash of ipaddr => NodeInfo
type NodeConfig map[string]NodeInfo

var nodeConfig *Config
var onceNodeConfig sync.Once

func GetNodeConfig() *NodeConfig {
	onceNodeConfig.Do(createNodeConfig)
	return nodeConfig.GetConfig().(*NodeConfig)
}

func createNodeConfig() {
	var err error
	nodeConfig, err = NewConfig("node", NodeConfig{})
	if err != nil {
		log.Fatal(err)
	}
}
