package server

import (
	"net"
	"time"

	"github.com/hpcloud/log"
	"github.com/vmihailenco/redis"
)

// NewRedisClient connects to redis after ensuring that the server is
// indeed running.
func NewRedisClient(addr, password string, database int64) (*redis.Client, error) {
	// Bug #97459 -- is the redis client library faking connection for
	// the down server?
	conn, err := net.Dial("tcp", convertLoopbackIP(addr))
	if err != nil {
		return nil, err
	}
	conn.Close()

	return redis.NewTCPClient(addr, password, database), nil
}

func NewRedisClientMust(addr, password string, database int64) *redis.Client {
	client, err := NewRedisClient(addr, password, database)
	if err != nil {
		log.Fatalf("Unable to connect to redis; %v", err)
	}
	return client
}

func NewRedisClientRetry(addr, password string, database, retries int64) (*redis.Client, error) {
	var client *redis.Client
	var err error

	if retries < 0 {
		// Default retry: up to10 minutes
		retries = 600
	}

	for attempt := int64(0); attempt < retries; attempt++ {
		if attempt > 0 {
			log.Warnf("Retrying (#%d) failed connection to redis (%v) ...",
				attempt+1, err)
			time.Sleep(time.Second)
		}
		client, err = NewRedisClient(addr, password, database)
		if err == nil {
			return client, nil
		}
	}
	if err == nil {
		log.Fatal("impossible")
	}
	return client, err
}
