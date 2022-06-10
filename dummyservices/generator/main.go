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
	InstanceName = "generator"
	TargetIPs = []string{"127.0.0.1"}
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
	message := generateMessage(config.Cfg.PayloadSize)
	message.MessageId = "Test1000000"
	message.Workload = config.Cfg.DefaultWorkloadSize
	message.Hops = []NodeData{{NodeId: InstanceName, ExitTime: time.Now().UnixMicro()}}
	err := sendTestMessage(message)
	for err != nil {
		log("Can't reach clients yet, retrying")
		err = sendTestMessage(message)
	}

	//loglines := []string{}
	for i := 0; i < config.Cfg.Messages; i++ {
		go func() {
			message := generateMessage(config.Cfg.PayloadSize)
			message.MessageId = fmt.Sprintf("%s%d", InstanceName, id)
			message.Workload = config.Cfg.DefaultWorkloadSize

			//message.SortSize = config.Cfg.DefaultWorkloadSize

			if config.Cfg.ServiceMode {
				sendRESTMessage(message)
			} else {
				sendMqttMessage(message)
			}
		}()

		id++
		time.Sleep(time.Duration(frequency) * time.Millisecond)
	}
	/*f, err := os.OpenFile("/usr/bin/output.txt", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		panic(err)
	}

	defer f.Close()
	for _, logline := range loglines {
		if _, err = f.WriteString(logline); err != nil {
			panic(err)
		}
	}*/
}

func execCmdBash(dfCmd string) (string, error) {
	log(fmt.Sprintf("Executing %s\n", dfCmd))
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
	for _, targetIP := range TargetIPs {
		go func(targetIP string) {
			serviceUrl := fmt.Sprintf(config.Cfg.PushServiceURL, targetIP)
			message.Hops = []NodeData{{NodeId: InstanceName, ExitTime: time.Now().UnixMicro()}}
			jsonData, err := json.Marshal(message)
			_, err = http.Post(serviceUrl, "application/json",
				bytes.NewBuffer(jsonData))

			if err != nil {
				log(fmt.Sprintf("Failed to write to service %s\n", serviceUrl))
			}
		}(targetIP)
	}

	log(fmt.Sprintf("Message id %d sent\n", message.MessageId))

	return "" //logline
}

func sendTestMessage(message Message) error {
	//message.History = []string{InstanceName}
	jsonData, err := json.Marshal(message)
	//fmt.Println(string(jsonData))
	for _, targetIP := range TargetIPs {
		serviceUrl := fmt.Sprintf(config.Cfg.PushServiceURL, targetIP)
		_, err = http.Post(serviceUrl, "application/json",
			bytes.NewBuffer(jsonData))

		if err != nil {
			return err
		}
	}
	return nil
}

func sendMqttMessage(message Message) {
	log("Sending content to MQTT")
	client := *getClient()
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		log("Connect failed")
		log(token.Error().Error())
	}

	payload, err := json.Marshal(message)
	if err != nil {
		log(err.Error())
	}

	log("Publishing")
	token := client.Publish(config.Cfg.MqttTopicWrite, 0, false, payload)
	token.Wait()
	if token.Error() != nil {
		log("Publish failed")
		log(token.Error().Error())
	}

	log("Published")
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

	log(fmt.Sprintf("Connecting to %s\n", config.Cfg.MqttBroker))

	client := mqtt.NewClient(opts)
	return &client
}

func log(line string) {
	f, err := os.OpenFile("/usr/bin/output.txt", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		panic(err)
	}

	defer f.Close()
	if _, err = f.WriteString(line); err != nil {
		panic(err)
	}
}
