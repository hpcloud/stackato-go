package server

import (
	"fmt"
	"os"
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
