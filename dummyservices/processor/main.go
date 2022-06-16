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
	"time"
)

var InstanceName string
var WorkloadSize int
var TargetIPs []string

func main() {
	argsWithoutProg := os.Args[1:]
	cfgFile := "defaultconfig.json"
	InstanceName = "processor"
	TargetIPs = []string{"127.0.0.1"}
	if len(argsWithoutProg) > 0 {
		cfgFile = argsWithoutProg[0]
		InstanceName = argsWithoutProg[1]
		WorkloadSize, _ = strconv.Atoi(argsWithoutProg[2])
		TargetIPs = argsWithoutProg[3:]
	}

	start := time.Now()
	bubbles := 1000
	for i := 0; i < bubbles; i++ {
		bubbleSort(1000)
	}

	fmt.Printf("%d bubbles took %f s\n", bubbles, float32(time.Since(start).Milliseconds())/1000.0)

	config.LoadConfig(cfgFile)

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
