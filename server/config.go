package server

import (
	"confdis/go/confdis"
	"fmt"
	"github.com/ActiveState/log"
	"io/ioutil"
	"net/url"
)

// Config refers to Stackato configuration under a specific
// group, such as "dea" or "cluster".
type Config struct {
	name    string
	changes chan error
	*confdis.ConfDis
}

func NewConfig(group string, s interface{}) (*Config, error) {
	addr, pass, db, err := getStackatoRedisAddr()
	if err != nil {
		return nil, err
	}

	redis, err := NewRedisClient(addr, pass, db)
	if err != nil {
		return nil, err
	}

	c, err := confdis.New(redis, group, s)
	if err != nil {
		return nil, err
	}
	gc := &Config{group, nil, c}
	go gc.monitor()
	return gc, nil
}

// GetChangesChannel returns a channel of (always) nil values that
// updates upon config changes.
func (g *Config) GetChangesChannel() chan error {
	// XXX: not bothering to lock this, yet.
	if g.changes == nil {
		g.changes = make(chan error)
	}
	return g.changes
}

// monitor monitors config changes, and exits abruptly upon on any
// error.
func (g *Config) monitor() {
	for err := range g.Changes {
		if err != nil {
			log.Fatalf("Error reading config for %s: %v",
				g.name, err)
		}
		if g.changes != nil {
			g.changes <- err
		}
	}
}

// getStackatoRedisAddr returns the redis connection address, password
// and database of the Stackato redis instance storing configuration.
func getStackatoRedisAddr() (string, string, int64, error) {
	uridata, err := ioutil.ReadFile("/s/etc/kato/redis_uri")
	if err != nil {
		return "", "", -1, err
	}
	u, err := url.Parse(string(uridata))
	if err != nil {
		return "", "", -1, err
	}

	// extract database number from Path
	var database int64
	fmt.Sscanf(u.Path, "/%d", &database)

	var pass string
	if u.User != nil {
		var haspass bool
		pass, haspass = u.User.Password()
		if !haspass {
			pass = ""
		}
	}

	return u.Host, pass, database, nil
}
