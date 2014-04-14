package server

import (
	"fmt"
	"github.com/ActiveState/log"
	"net"
)

func LocalIPMust() string {
	ip, err := LocalIP()
	if err != nil {
		log.Fatalf("Unable to determine local IP: %v", err)
	}
	return ip
}

func NodeIPMust() string {
	ip, err := NodeIP()
	if err != nil {
		log.Fatalf("Unable to determine node IP: %v", err)
	}
	return ip
}

// LocalIP returns the externally visible address of localhost
func LocalIP() (string, error) {
	ip, err := localIP()
	if err != nil {
		return "", err
	}
	return ip.String(), nil
}

func localIP() (net.IP, error) {
	tt, err := net.Interfaces()
	if err != nil {
		return nil, err
	}
	for _, t := range tt {
		aa, err := t.Addrs()
		if err != nil {
			return nil, err
		}
		for _, a := range aa {
			ipnet, ok := a.(*net.IPNet)
			if !ok {
				continue
			}
			v4 := ipnet.IP.To4()
			if v4 == nil || v4[0] == 127 { // loopback address
				continue
			}
			return v4, nil
		}
	}
	return nil, fmt.Errorf("no interfaces")
}

// NodeIP returns the externally visible address of the docker host. If the
// current process is not running in a docker container, returns LocalIP()
func NodeIP() (string, error) {
	if InsideDocker() {
		return GetDockerHostIp()
	} else {
		return LocalIP()
	}
}
