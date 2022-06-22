package main

import (
	"encoding/json"
	"fmt"
	"processor/config"
	"processor/message"
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

func handleMessage(client mqtt.Client, mess mqtt.Message) {
	msg := message.Message{}
	json.Unmarshal(mess.Payload(), &msg)

	go func() {
		//bubbleSort(config.Cfg.DefaultWorkloadSize)
		//sendNextMqttMessage(msg.MessageId)
	}()
}

func sendNextMqttMessage(id int, time int64) {
	fmt.Println("Sending content to MQTT")
	client := *getClient()
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		fmt.Println("Connect fucked")
		fmt.Println(token.Error())
	}

	data := generateMessage()
	//data.StartTime = []int64{time}
	//data.MessageId = id

	payload, err := json.Marshal(data)
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
	client.Disconnect(50)
}
