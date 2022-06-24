package main

import (
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"processor/config"
	"processor/message"
	"strconv"
	"strings"
)

var InstanceName string
var WorkloadSize int
var Targets []Target

//var LoadBalanceMode bool
var TotalMessages int

//var MessageCounts map[string]int

type Target struct {
	IPQuota       map[string]float64
	MessageCounts map[string]int
}

var messageData string

func main() {
	argsWithoutProg := os.Args[1:]
	cfgFile := "defaultconfig.json"
	InstanceName = "processor"

	Targets = []Target{}
	TotalMessages = 0
	if len(argsWithoutProg) > 0 {
		cfgFile = argsWithoutProg[0]
		InstanceName = argsWithoutProg[1]
		WorkloadSize, _ = strconv.Atoi(argsWithoutProg[2])
		targets := argsWithoutProg[3:]
		//TargetIPs = []Target{}
		for _, target := range targets {
			ipquotas := make(map[string]float64)
			ips := strings.Split(target, ",")
			for _, ip := range ips {
				parts := strings.Split(ip, ":")
				num, _ := strconv.ParseFloat(parts[1], 64)
				ipquotas[parts[0]] = num
			}
			Targets = append(Targets, Target{IPQuota: ipquotas, MessageCounts: make(map[string]int)})
		}
	}

	/*start := time.Now()
	bubbles := 100
	for i := 0; i < bubbles; i++ {
		bubbleSort(1000)
	}*/

	//fmt.Printf("%d bubbles took %f s\n", bubbles, float32(time.Since(start).Milliseconds())/1000.0)

	config.LoadConfig(cfgFile)

	data := []byte{}
	for i := 0; i < config.Cfg.PayloadSize; i++ {
		data = append(data, byte(70+i%10))
	}
	messageData = string(data)

	/*testMsg := generateMessage()
	testMsg.MessageId = "10"
	testMsg.Workload = 200
	testMsg.Hops = []message.NodeData{{NodeId: "test", EntryTime: time.Now().UnixMicro(), ExitTime: time.Now().UnixMicro()}}

	json, _ := easyjson.Marshal(testMsg)
	fmt.Println(string(json))*/

	/*for i := 0; i < 1000; i++ {
		sendNextRESTMessage(Message{Workload: 20, MessageId: "test", Hops: []NodeData{{NodeId: "test"}}})
	}*/

	//fmt.Println(time.Now().UnixNano())
	if config.Cfg.ServiceMode {
		router := NewRouter()
		log.Fatal(http.ListenAndServe(":8080", router))
	} else {
		processMessages()
	}
}

func getTlsConfig() *tls.Config {
	return &tls.Config{
		ClientAuth: tls.RequestClientCert,
	}
}

func execCmdBash(dfCmd string) (string, error) {
	logger(fmt.Sprintf("Executing %s\n", dfCmd))
	cmd := exec.Command("sh", "-c", dfCmd)
	stdout, err := cmd.Output()

	if err != nil {
		println(err.Error())
		return "", err
	}
	//fmt.Println(string(stdout))
	return string(stdout), nil
}

func logger(line string) {
	f, err := os.OpenFile("/usr/bin/output.txt", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		//panic(err)
	}

	defer f.Close()
	if _, err = f.WriteString(line); err != nil {
		//panic(err)
	}
}

func generateMessage() message.Message {
	/*data := ""
	for i := 0; i < payloadSize; i++ {
		data += strconv.Itoa(i % 10)
	}*/

	message := message.Message{
		Payload: messageData,
	}

	return message
}
