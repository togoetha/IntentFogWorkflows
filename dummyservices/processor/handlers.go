package main

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"processor/config"
)

func ProcessMessage(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var message Message
	if err := decoder.Decode(&message); err != nil {
		panic(err)
	}

	go func() {
		bubbleSort(message.SortSize)
		sendNextRESTMessage(message.MessageId)
	}()
}

func bubbleSort(n int) []int {
	numbers := []int{}
	for i := 0; i < n; i++ {
		numbers = append(numbers, n-i)
	}

	for i := 0; i < n-1; i++ {
		for j := 0; j < n-i-1; j++ {
			if numbers[j] > numbers[j+1] {
				numbers[j], numbers[j+1] = numbers[j+1], numbers[j]
			}
		}
	}
	return numbers
}

func sendNextRESTMessage(id int) {
	data := generateMessage(config.Cfg.PayloadSize)
	data.SortSize = config.Cfg.DefaultWorkloadSize
	data.MessageId = id
	jsonData, err := json.Marshal(data)

	_, err = http.Post(config.Cfg.PushServiceURL, "application/json",
		bytes.NewBuffer(jsonData))

	if err != nil {
		log.Fatal(err)
	}
}
