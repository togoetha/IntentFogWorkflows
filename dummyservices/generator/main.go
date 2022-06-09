package main

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"generator/config"
	"net/http"
	"os"
	"os/exec"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

var InstanceName string
var TargetIPs []string

func main() {
	argsWithoutProg := os.Args[1:]
	cfgFile := "defaultconfig.json"
	if len(argsWithoutProg) > 0 {
		cfgFile = argsWithoutProg[0]
		InstanceName = argsWithoutProg[1]
		TargetIPs = argsWithoutProg[2:]
	}

	config.LoadConfig(cfgFile)

	generate()
}

var running bool

func generate() {
	frequency := 1000 / config.Cfg.MessageFrequency
	id := 1
	//running = true

	loglines := []string{}
	for i := 0; i < config.Cfg.Messages; i++ {
		go func() {
			message := generateMessage(config.Cfg.PayloadSize)
			message.MessageId = id
			message.StartTime = time.Now().UnixMicro()
			//message.SortSize = config.Cfg.DefaultWorkloadSize

			if config.Cfg.ServiceMode {
				loglines = append(loglines, sendRESTMessage(message))
			} else {
				sendMqttMessage(message)
			}
		}()

		id++
		time.Sleep(time.Duration(frequency) * time.Millisecond)
	}
	f, err := os.OpenFile("/usr/bin/output.txt", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		panic(err)
	}

	defer f.Close()
	for _, logline := range loglines {
		if _, err = f.WriteString(logline); err != nil {
			panic(err)
		}
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

func sendRESTMessage(message Message) string {
	jsonData, err := json.Marshal(message)

	for _, targetIP := range TargetIPs {
		go func(targetIP string) {
			message.History = []string{targetIP}
			serviceUrl := fmt.Sprintf(config.Cfg.PushServiceURL, targetIP)
			_, err = http.Post(serviceUrl, "application/json",
				bytes.NewBuffer(jsonData))

			if err != nil {
				fmt.Printf("Failed to write to service %s\n", serviceUrl)
			}
		}(targetIP)
	}

	logline := fmt.Sprintf("Message id %d sent at %d\n", message.MessageId, message.StartTime)
	/*fmt.Println(logline)
	f, err := os.OpenFile("/usr/bin/output.txt", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		panic(err)
	}

	defer f.Close()
	if _, err = f.WriteString(logline); err != nil {
		panic(err)
	}*/
	return logline
}

func sendMqttMessage(message Message) {
	fmt.Println("Sending content to MQTT")
	client := *getClient()
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		fmt.Println("Connect fucked")
		fmt.Println(token.Error())
	}

	payload, err := json.Marshal(message)
	if err != nil {
		fmt.Println(err.Error())
	}

	fmt.Println("Publishing")
	token := client.Publish(config.Cfg.MqttTopicWrite, 0, false, payload)
	token.Wait()
	if token.Error() != nil {
		fmt.Println("Publish fucked")
		fmt.Println(token.Error())
	}

	fmt.Println("Published")
	client.Disconnect(250)
}

func getTlsConfig() *tls.Config {
	return &tls.Config{
		ClientAuth: tls.RequestClientCert,
	}
}

func getClient() *mqtt.Client {
	opts := mqtt.NewClientOptions()
	opts.AddBroker(config.Cfg.MqttBroker)
	opts.SetClientID(config.Cfg.MqttClientId).SetTLSConfig(getTlsConfig())
	opts.SetUsername(config.Cfg.MqttUser)
	opts.SetPassword(config.Cfg.MqttPass) //flexnet

	fmt.Printf("Connecting to %s\n", config.Cfg.MqttBroker)

	client := mqtt.NewClient(opts)
	return &client
}
