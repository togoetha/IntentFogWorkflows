package main

import (
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"processor/config"
	"strconv"
)

var InstanceName string
var WorkloadSize int
var TargetIPs []Target
var LoadBalanceMode bool
var TotalMessages int
var MessageCounts map[string]int

type Target struct {
	ip    string
	quota float64
}

func main() {
	argsWithoutProg := os.Args[1:]
	cfgFile := "defaultconfig.json"
	InstanceName = "processor"
	TargetIPs = []Target{{ip: "127.0.0.1"}}
	MessageCounts = make(map[string]int)
	TotalMessages = 0
	LoadBalanceMode = false
	if len(argsWithoutProg) > 0 {
		cfgFile = argsWithoutProg[0]
		InstanceName = argsWithoutProg[1]
		WorkloadSize, _ = strconv.Atoi(argsWithoutProg[2])
		targets := argsWithoutProg[3:]
		TargetIPs = []Target{}
		for idx := 0; idx < len(targets); idx += 2 {
			quota, _ := strconv.ParseFloat(targets[idx+1], 32)
			if quota != 0 {
				LoadBalanceMode = true
			}
			MessageCounts[targets[idx]] = 0
			TargetIPs = append(TargetIPs, Target{ip: targets[idx], quota: quota})
		}
	}

	/*start := time.Now()
	bubbles := 100
	for i := 0; i < bubbles; i++ {
		bubbleSort(1000)
	}*/

	//fmt.Printf("%d bubbles took %f s\n", bubbles, float32(time.Since(start).Milliseconds())/1000.0)

	config.LoadConfig(cfgFile)

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
		panic(err)
	}

	defer f.Close()
	if _, err = f.WriteString(line); err != nil {
		panic(err)
	}
}
