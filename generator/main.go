package main

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"generator/config"
	"log"
	"net/http"
	"os"
	"os/exec"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

func main() {
	argsWithoutProg := os.Args[1:]
	cfgFile := "defaultconfig.json"
	if len(argsWithoutProg) > 0 {
		cfgFile = argsWithoutProg[0]
	}

	config.LoadConfig(cfgFile)

	generate()
}

var running bool

func generate() {
	frequency := 1000 / config.Cfg.MessageFrequency
	id := 1
	running = true

	for running {
		go func() {
			message := generateMessage(config.Cfg.PayloadSize)
			message.MessageId = id
			message.SortSize = config.Cfg.DefaultWorkloadSize

			if config.Cfg.ServiceMode {

			} else {
				sendMqttMessage(message)
			}
		}()

		id++
		time.Sleep(time.Duration(frequency) * time.Millisecond)
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

func sendRESTMessage(message Message) {
	jsonData, err := json.Marshal(message)

	_, err = http.Post(config.Cfg.PushServiceURL, "application/json",
		bytes.NewBuffer(jsonData))

	if err != nil {
		log.Fatal(err)
	}
}

func sendMqttMessage(message Message) {
	fmt.Println("Sending content to MQTT")
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