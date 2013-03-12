package server

import (
	"github.com/ActiveState/log"
	"github.com/apcera/nats"
	"time"
)

// NewNatsClient connects to the NATS server of the Stackato cluster
func NewNatsClient() *nats.EncodedConn {
	log.Infof("Connecting to NATS %s\n", Config.NatsUri)
	nc, err := nats.Connect(Config.NatsUri)
	if err != nil {
		log.Fatal(err)
	}
	log.Infof("Connected to NATS %s\n", Config.NatsUri)
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
