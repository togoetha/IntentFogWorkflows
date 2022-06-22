package main

import (
	"net/http"
	"sink/message"
	"time"

	easyjson "github.com/mailru/easyjson"
)

func ProcessMessage(w http.ResponseWriter, r *http.Request) {
	//decoder := json.NewDecoder(r.Body)
	var msg message.Message
	if err := easyjson.UnmarshalFromReader(r.Body, &msg); err != nil {
		panic(err)
	}

	go func() {
		msg.Hops = append(msg.Hops, message.NodeData{NodeId: InstanceName, EntryTime: time.Now().UnixMicro()})
		//bubbleSort(config.Cfg.DefaultWorkload)
		finishMessage(msg)
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
