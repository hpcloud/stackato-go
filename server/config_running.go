package server

import (
	"github.com/ActiveState/log"
	"os"
	"sync"
)

// runningInfo is a map of process name to process pid
type runningInfo map[string]int

// RunningConfig corresponds to Kato's /running config, which is a hash of
// ipaddr => runningInfo
type RunningConfig map[string]runningInfo

// MarkRunning marks the current process, identified by name `name`, as
// "running" for kato-status.
func MarkRunning(name string) {
	// XXX: it is not safe to write to /running yet.
	// https://bugs.activestate.com/show_bug.cgi?id=102768#c10
	return

	onceRunningConfig.Do(createRunningConfig)
	nodeid, err := GetClusterConfig().CurrentNodeId()
	if err != nil {
		log.Fatal("Unable to get current node id: %v", err)
	}
	pid := os.Getpid()
	err = runningConfig.AtomicSave(func(i interface{}) error {
		config := i.(*RunningConfig)
		if *config == nil {
			*config = make(RunningConfig)
		}
		if _, ok := (*config)[nodeid]; !ok {
			(*config)[nodeid] = make(map[string]int)
		}
		(*config)[nodeid][name] = pid
		return nil
	})
	log.Infof("PID is %v", pid)
	if err != nil {
		log.Fatal("Error setting running status: %v", err)
	}
}

var runningConfig *Config
var onceRunningConfig sync.Once

func GetRunningConfig() *RunningConfig {
	onceRunningConfig.Do(createRunningConfig)
	return runningConfig.GetConfig().(*RunningConfig)
}

func createRunningConfig() {
	var err error
	runningConfig, err = NewConfig("running", RunningConfig{})
	if err != nil {
		log.Fatal(err)
	}
}
