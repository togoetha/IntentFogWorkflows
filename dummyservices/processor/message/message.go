package message

type Message struct {
	Hops      []NodeData `json:"hops"`
	Workload  int        `json:"workload"`
	Payload   string     `json:"payload"`
	MessageId string     `json:"messageId"`
}

type NodeData struct {
	NodeId    string `json:"history"`
	EntryTime int64  `json:"startTime"`
	ExitTime  int64  `json:"endTime"`
}
