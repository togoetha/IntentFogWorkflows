package main

import (
	"encoding/json"
	"net/http"
	"time"
)

func ProcessMessage(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var message Message
	if err := decoder.Decode(&message); err != nil {
		panic(err)
	}

	go func() {
		message.Hops = append(message.Hops, NodeData{NodeId: InstanceName, EntryTime: time.Now().UnixMicro()})
		//bubbleSort(config.Cfg.DefaultWorkload)
		finishMessage(message)
	}()
}

/*func bubbleSort(n int) []int {
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
}*/
