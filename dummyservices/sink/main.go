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

func main() {
	argsWithoutProg := os.Args[1:]
	cfgFile := "defaultconfig.json"
	if len(argsWithoutProg) > 0 {
		cfgFile = argsWithoutProg[0]
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
	timetaken := time.Since(time.UnixMicro(msg.StartTime))
	now := time.Now().UnixMicro()
	logline := fmt.Sprintf("Message id %d start %d now %d took %d ms chain %s\n", msg.MessageId, msg.StartTime, now, timetaken.Microseconds()/1000.0, strings.Join(msg.History, ">"))
	//fmt.Println(logline)
	loglines = append(loglines, logline)

	if len(loglines) == 20 {
		f, err := os.OpenFile("/usr/bin/output.txt", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0666)
		if err != nil {
			panic(err)
		}

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
