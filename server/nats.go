package server

import (
	"fmt"
	"github.com/ActiveState/log"
	"github.com/apcera/nats"
	"time"
)

// NewNatsClient connects to the NATS server of the Stackato cluster
func NewNatsClient(retries int) *nats.EncodedConn {
	natsUri, err := getNatsUri()
	if err != nil {
		log.Fatalf("Unable to get Nats URI: %v", err)
	}
	log.Infof("Connecting to NATS server %s\n", natsUri)

	var nc *nats.Conn

	for attempt := 0; attempt < retries; attempt++ {
		nc, err = nats.Connect(natsUri)
		if err != nil {
			if (attempt + 1) == retries {
				log.Fatal(err)
			}
			log.Warnf("NATS connection error (%v); retrying after 1 second..",
				err)
			time.Sleep(time.Second)
		}
	}

	log.Infof("Connected to NATS server %s\n", natsUri)
	client, err := nats.NewEncodedConn(nc, "json")
	if err != nil {
		log.Fatal(err)
	}

	// Diagnosing Bug #97856 by periodically checking if we are still
	// connected to NATS.
	go func() {
		log.Info("Periodically checking NATS connectivity")
		for _ = range time.Tick(1 * time.Minute) {
			if nc.IsClosed() {
				log.Fatal("Connection to NATS has been closed (in the last minute)")
			}
		}
	}()

	return client
}

func getNatsUri() (string, error) {
	var ipaddr string
	// Use non-lookback address on a micro cloud to connect from docker
	// container to host.
	if InsideDocker() && GetClusterConfig().IsMicro() {
		var err error
		ipaddr, err = GetDockerHostIp()
		if ipaddr == "" {
			return "", err
		}
	} else {
		ipaddr = GetClusterConfig().MbusIp
	}

	// HACK: Ideally we should be reading NatsUri from
	// cloud_controller config (mbus). we take a shortcut here in
	// order to not have to create a separate ConfDis instance for
	// cloud_controller config (and having to watch it). This will
	// have to change if we switch to clustered version of NATS.
	return fmt.Sprintf("nats://%s:4222/", ipaddr), nil
}
