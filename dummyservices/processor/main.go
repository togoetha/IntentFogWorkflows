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

func main() {
	argsWithoutProg := os.Args[1:]
	cfgFile := "defaultconfig.json"
	InstanceName = "processor"
	Targets = []Target{}
	//MessageCounts = make(map[string]int)
	TotalMessages = 0
	//LoadBalanceMode = false
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
