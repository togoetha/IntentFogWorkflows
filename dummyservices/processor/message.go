package main

import (
	"strconv"
)

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

func generateMessage(payloadSize int) Message {
	data := ""
	for i := 0; i < payloadSize; i++ {
		data += strconv.Itoa(i % 10)
	}

	message := Message{
		Payload: data,
	}

	return message
}
