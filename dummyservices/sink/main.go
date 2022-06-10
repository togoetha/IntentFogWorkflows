package main

import (
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"sink/config"
	"strings"
	"time"
)

var loglines []string
var InstanceName string

func main() {
	argsWithoutProg := os.Args[1:]
	cfgFile := "defaultconfig.json"
	InstanceName = "sink"

	if len(argsWithoutProg) > 0 {
		cfgFile = argsWithoutProg[0]
		InstanceName = argsWithoutProg[1]
	}

	loglines = []string{}
	config.LoadConfig(cfgFile)

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

func finishMessage(msg Message) {
	totalTime := time.Since(time.UnixMicro(msg.Hops[0].ExitTime))
	//now := time.Now().UnixMicro()
	strTimes := []string{}
	for idx, hop := range msg.Hops {
		if idx == 0 {
			strTimes = append(strTimes, fmt.Sprintf("%s;;", hop.NodeId)) //append(strTimes, fmt.Sprintf("%s, %d", hop.NodeId, hop.ExitTime))
		} else if idx == len(msg.Hops)-1 {
			strTimes = append(strTimes, fmt.Sprintf("%s;%f;%f", hop.NodeId, float32(hop.EntryTime-msg.Hops[idx-1].ExitTime)/1000.0, 0.0)) //fmt.Sprintf("-> %fms -> %s, %d", float32(hop.EntryTime-msg.Hops[idx-1].ExitTime)/1000.0, hop.NodeId, hop.EntryTime))
		} else {
			strTimes = append(strTimes, fmt.Sprintf("%s;%f;%f", hop.NodeId, float32(hop.EntryTime-msg.Hops[idx-1].ExitTime)/1000.0, float32(hop.ExitTime-hop.EntryTime)/1000.0)) //fmt.Sprintf("-> %fms -> %s, %d to %d, %fms processing", float32(hop.EntryTime-msg.Hops[idx-1].ExitTime)/1000.0, hop.NodeId, hop.EntryTime, hop.ExitTime, float32(hop.ExitTime-hop.EntryTime)/1000.0))
		}
	}
	logline := fmt.Sprintf("%s;%f;%s\n", msg.MessageId, float32(totalTime.Microseconds())/1000.0, strings.Join(strTimes, ";")) //fmt.Sprintf("Message id %d took %fms chain %s\n", msg.MessageId, float32(totalTime.Microseconds())/1000.0, strings.Join(strTimes, " "))
	//fmt.Println(logline)
	loglines = append(loglines, logline)

	if len(loglines) == 20 {
		f, err := os.OpenFile("/usr/bin/output.txt", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0666)
		if err != nil {
			panic(err)
		}

		f.WriteString("MessageId;TotalTime;Hop1;HopLatency1;HopProcessing1\n")
		for _, logline := range loglines {
			if _, err = f.WriteString(logline); err != nil {
				panic(err)
			}
		}
		f.Close()
		loglines = []string{}
	}
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
