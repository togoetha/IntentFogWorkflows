package main

import (
	"strconv"
)

type Message struct {
	SortSize  int    `json:"sortSize"`
	Payload   string `json:"payload"`
	MessageId int    `json:"messageId"`
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