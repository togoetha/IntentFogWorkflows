package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"processor/config"
	"time"
)

func ProcessMessage(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var message Message
	if err := decoder.Decode(&message); err != nil {
		panic(err)
	}

	go func(message Message) {
		start := time.Now()
		bubbleSort(config.Cfg.DefaultWorkloadSize)
		logger(fmt.Sprintf("Bubble sort for %d took %dms", message.MessageId, time.Since(start).Milliseconds()))
		sendNextRESTMessage(message.MessageId, message.StartTime, message.History)
	}(message)
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

func sendNextRESTMessage(id int, times []int64, history []string) {
	data := generateMessage(config.Cfg.PayloadSize)
	data.MessageId = id
	data.History = append(history, InstanceName)
	data.StartTime = append(times, time.Now().UnixMicro())
	jsonData, err := json.Marshal(data)

	for _, targetIP := range TargetIPs {
		go func(targetIP string) {
			serviceUrl := fmt.Sprintf(config.Cfg.PushServiceURL, targetIP)
			_, err = http.Post(serviceUrl, "application/json",
				bytes.NewBuffer(jsonData))

			if err != nil {
				logger(fmt.Sprintf("Failed to write to service %s\n", serviceUrl))
			}
		}(targetIP)
	}
}
