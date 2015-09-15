package client

type Client interface {
	PushApp(name string, dir string) error
}
