package main

import (
	"bytes"
	"encoding/json"
	"fmt"
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
		bubbleSort(config.Cfg.DefaultWorkloadSize)
		sendNextRESTMessage(message.MessageId, message.StartTime)
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

func sendNextRESTMessage(id int, time int64) {
	data := generateMessage(config.Cfg.PayloadSize)
	data.MessageId = id
	data.StartTime = time
	jsonData, err := json.Marshal(data)

	for _, targetIP := range TargetIPs {
		serviceUrl := fmt.Sprintf(config.Cfg.PushServiceURL, targetIP)
		_, err = http.Post(serviceUrl, "application/json",
			bytes.NewBuffer(jsonData))

		if err != nil {
			fmt.Printf("Failed to write to service %s\n", serviceUrl)
		}
	}
}
