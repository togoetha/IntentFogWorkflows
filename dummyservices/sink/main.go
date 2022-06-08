package main

import (
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"sink/config"
	"time"
)

func main() {
	argsWithoutProg := os.Args[1:]
	cfgFile := "defaultconfig.json"
	if len(argsWithoutProg) > 0 {
		cfgFile = argsWithoutProg[0]
	}

	config.LoadConfig(cfgFile)

	fmt.Println(time.Now().UnixNano())
	if config.Cfg.ServiceMode {
		router := NewRouter()
		log.Fatal(http.ListenAndServe(":8081", router))
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
	fmt.Printf("Executing %s\n", dfCmd)
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
	timetaken := time.Since(time.UnixMicro(msg.StartTime))
	fmt.Printf("Message id %d took %d ms\n", msg.MessageId, timetaken.Microseconds()/1000.0)
}
