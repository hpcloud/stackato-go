package server

import (
	"confdis/go/confdis"
	"github.com/ActiveState/log"
)

// GroupConfig refers to Stackato configuration under a specific
// group, such as "dea" or "cluster".
type GroupConfig struct {
	name    string
	changes chan error
	*confdis.ConfDis
}

func NewGroupConfig(group string, s interface{}) (*GroupConfig, error) {
	// TODO: use cluster config to determine redis location.
	c, err := confdis.New("localhost:5454", group, s)
	if err != nil {
		return nil, err
	}
	gc := &GroupConfig{group, nil, c}
	go gc.monitor()
	return gc, nil
}

// GetChangesChannel returns a channel of (always) nil values that
// updates upon config changes.
func (g *GroupConfig) GetChangesChannel() chan error {
	// XXX: not bothering to lock this, yet.
	if g.changes == nil {
		g.changes = make(chan error)
	}
	return g.changes
}

// monitor monitors config changes, and exits abruptly upon on any
// error.
func (g *GroupConfig) monitor() {
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
