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
	message.Hops = append(message.Hops, NodeData{NodeId: InstanceName, EntryTime: time.Now().UnixMicro()})

	go func(message Message) {
		start := time.Now()
		bubbleSort(message.Workload)
		logger(fmt.Sprintf("Bubble sort for %d took %dms", message.MessageId, time.Since(start).Milliseconds()))
		sendNextRESTMessage(message)
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

func sendNextRESTMessage(orig Message) {
	data := generateMessage(config.Cfg.PayloadSize)
	data.MessageId = orig.MessageId
	data.Payload = orig.Payload
	data.Workload = orig.Workload
	data.Hops = orig.Hops

	for _, targetIP := range TargetIPs {
		go func(targetIP string) {
			serviceUrl := fmt.Sprintf(config.Cfg.PushServiceURL, targetIP)
			data.Hops[len(data.Hops)-1].ExitTime = time.Now().UnixMicro()
			jsonData, err := json.Marshal(data)
			_, err = http.Post(serviceUrl, "application/json",
				bytes.NewBuffer(jsonData))

			if err != nil {
				logger(fmt.Sprintf("Failed to write to service %s\n", serviceUrl))
			}
		}(targetIP)
	}
}
