package server

import (
	"fmt"
	"os"
	"strings"
)

// InsideDocker returns true if the process is run in a docker container.
func InsideDocker() bool {
	return os.Getenv("STACKATO_DOCKER") != ""
}

func GetDockerHostIp() (string, error) {
	ipaddr := os.Getenv("DOCKER_HOST")
	if ipaddr == "" {
		return "", fmt.Errorf("DOCKER_HOST env not set")
	} else {
		return ipaddr, nil
	}
}

func convertLoopbackIP(addr string) string {
	// Avoid going through docker's userland proxy by not using 127.0.0.1
	if strings.Contains(addr, "127.0.0.1") {
		addr = strings.Replace(addr, "127.0.0.1", LocalIPMust(), 1)
	}
	return addr
}
