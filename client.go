package client

import ()

type Client interface {
	PushApp(name string, dir string) error
}
