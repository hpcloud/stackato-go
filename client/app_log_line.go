package client

type AppLogLine struct {
	Text          string `json:"text"`
	Source        string `json:"source"`
	Filename      string `json:"filename"`
	InstanceIndex int    `json:"instance"`
	Timestamp     int64  `json:"timestamp"`
	NodeID        string `json:"nodeid"`
}
