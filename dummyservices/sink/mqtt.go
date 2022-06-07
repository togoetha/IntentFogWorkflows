package main

import (
	"encoding/json"
	"fmt"
	"sink/config"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

var running bool

func processMessages() {
	running = true

	c := *getClient()
	if token := c.Subscribe(config.Cfg.MqttTopicRead, 0, handleMessage); token.Wait() && token.Error() != nil {
		fmt.Println(token.Error())
	}
	fmt.Print("Subscribe topic " + config.Cfg.MqttTopicRead + " success\n")

	for running {
		time.Sleep(50 * time.Millisecond)
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

func handleMessage(client mqtt.Client, message mqtt.Message) {
	msg := Message{}
	json.Unmarshal(message.Payload(), &msg)

	go func() {
		bubbleSort(msg.SortSize)
		//sendNextMqttMessage(msg.MessageId)
	}()
}
