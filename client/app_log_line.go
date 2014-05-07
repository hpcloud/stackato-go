package client

type AppLogLine struct {
	Text      string `json:"text"`
	Source    string `json:"text"`
	Filename  string `json:"filename"`
	Timestamp int64  `json:"timestamp"`
	NodeID    string `json:"nodeid"`
}
