package server

import (
	"github.com/ActiveState/log"
	"github.com/apcera/nats"
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
	return client
}
