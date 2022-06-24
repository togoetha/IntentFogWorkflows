package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"processor/config"
	"processor/message"
	"time"

	"github.com/mailru/easyjson"
)

func ProcessMessage(w http.ResponseWriter, r *http.Request) {
	var msg message.Message

	/*defer func() {
		if r := recover(); r != nil {
			logger(fmt.Sprintf("Recovered %i", r))
		}
	}()*/

	/*if err := easyjson.UnmarshalFromReader(r.Body, &msg); err != nil {
		//panic(err)
		logger(err.Error())
		return
	} */ /*else {
		fmt.Printf("Message %s received payload %d", msg.MessageId, len(msg.Payload))
	}*/

	decoder := json.NewDecoder(r.Body)

	if err := decoder.Decode(&msg); err != nil {
		//panic(err)
		logger(err.Error())
		return
	}

	logger(fmt.Sprintf("%d Message %s received\n", time.Now().UnixMilli(), msg.MessageId))
	msg.Hops = append(msg.Hops, message.NodeData{NodeId: InstanceName, EntryTime: time.Now().UnixMicro()})

	go func(msg message.Message) {
		start := time.Now()
		bubbleSort(WorkloadSize)
		logger(fmt.Sprintf("%d Bubble sort for %s took %dms\n", time.Now().UnixMilli(), msg.MessageId, time.Since(start).Milliseconds()))
		sendNextRESTMessage(msg)
	}(msg)
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

func sendNextRESTMessage(orig message.Message) {
	data := generateMessage()
	data.MessageId = orig.MessageId
	data.Payload = orig.Payload
	data.Workload = orig.Workload
	data.Hops = orig.Hops
	logger(fmt.Sprintf("%d Send next message for id %s\n", time.Now().UnixMilli(), orig.MessageId))

	//if LoadBalanceMode {
	//fmt.Println("Load balance mode")
	for _, target := range Targets {
		tIP := ""
		for ip, quota := range target.IPQuota {
			if quota == 0 || TotalMessages == 0 {
				tIP = ip
				break
			} else {
				current := float64(target.MessageCounts[ip]) / float64(TotalMessages)
				logger(fmt.Sprintf("%d Messages processed %d, target %s current %f quota %f\n", time.Now().UnixMilli(), TotalMessages, ip, current, quota))
				//fmt.Printf("Messages processed %d, target %s current %f quota %f\n", TotalMessages, ip, current, quota)
				if current < quota {
					tIP = ip
					break
				}
			}
		}
		logger(fmt.Sprintf("%d Sending to %s\n", time.Now().UnixMilli(), tIP))
		//fmt.Printf("Sending to %s\n", tIP)
		go func(target string) {
			serviceUrl := fmt.Sprintf(config.Cfg.PushServiceURL, target)
			data.Hops[len(data.Hops)-1].ExitTime = time.Now().UnixMicro()
			jsonData, err := easyjson.Marshal(data)
			_, err = http.Post(serviceUrl, "application/json",
				bytes.NewBuffer(jsonData))

			if err != nil {
				logger(fmt.Sprintf("%d Failed to write to service %s\n", time.Now().UnixMilli(), serviceUrl))
				//fmt.Printf("Failed to write to service %s\n", serviceUrl)
			} else {
				logger(fmt.Sprintf("%d Message %s sent\n", time.Now().UnixMilli(), data.MessageId))
			}

		}(tIP)
		target.MessageCounts[tIP] += 1
	}
	TotalMessages++
	/*} else {
		for _, targetIP := range Targets {
			go func(target Target) {
				serviceUrl := fmt.Sprintf(config.Cfg.PushServiceURL, target.ip)
				data.Hops[len(data.Hops)-1].ExitTime = time.Now().UnixMicro()
				jsonData, err := json.Marshal(data)
				_, err = http.Post(serviceUrl, "application/json",
					bytes.NewBuffer(jsonData))

				if err != nil {
					logger(fmt.Sprintf("Failed to write to service %s\n", serviceUrl))
				}
			}(targetIP)
		}
	}*/
}
