package server

import (
	"github.com/ActiveState/log"
	"github.com/vmihailenco/redis"
	"net"
	"time"
)

// NewRedisClient connects to redis after ensuring that the server is
// indeed running.
func NewRedisClient(addr, password string, database int64) (*redis.Client, error) {
	// Bug #97459 -- is the redis client library faking connection for
	// the down server?
	conn, err := net.Dial("tcp", InjectDockerHostIp(addr))
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

	for attempt := int64(0); attempt < retries; attempt++ {
		if attempt > 0 {
			log.Warnf("Redis error (%v); waiting a second before retrying...",
				err)
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
