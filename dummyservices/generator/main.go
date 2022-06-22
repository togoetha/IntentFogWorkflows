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
	"strconv"
	"strings"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

var InstanceName string
var MessageFrequency int
var Targets []Target
var TotalMessages int

type Target struct {
	IPQuota       map[string]float64
	MessageCounts map[string]int
}

var messageData []byte

func main() {
	argsWithoutProg := os.Args[1:]
	cfgFile := "defaultconfig.json"
	InstanceName = "generator"
	Targets = []Target{}
	if len(argsWithoutProg) > 0 {
		cfgFile = argsWithoutProg[0]
		InstanceName = argsWithoutProg[1]
		MessageFrequency, _ = strconv.Atoi(argsWithoutProg[2])
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

	config.LoadConfig(cfgFile)

	generate()

	i := 0
	for true {
		i++
		time.Sleep(time.Second)
	}
}

//var running bool

func generate() {
	frequency := 1000.0 / float32(MessageFrequency)
	id := 1
	//running = true

	messageData = []byte{}
	for i := 0; i < config.Cfg.PayloadSize; i++ {
		messageData = append(messageData, byte(i%10))
	}

	message := generateMessage()
	message.MessageId = "Test1000000"
	message.Workload = config.Cfg.DefaultWorkloadSize
	message.Hops = []NodeData{{NodeId: InstanceName, ExitTime: time.Now().UnixMicro()}}
	err := sendTestMessage(message)
	for err != nil {
		log(fmt.Sprintf("%d Can't reach clients yet, retrying\n", time.Now().UnixMilli()))
		err = sendTestMessage(message)
	}

	//loglines := []string{}
	for i := 0; i < config.Cfg.Messages; i++ {
		go func(msgId int) {
			message := generateMessage()
			message.MessageId = fmt.Sprintf("%s%d", InstanceName, msgId)
			message.Workload = config.Cfg.DefaultWorkloadSize

			//message.SortSize = config.Cfg.DefaultWorkloadSize

			if config.Cfg.ServiceMode {
				sendRESTMessage(message)
			} else {
				sendMqttMessage(message)
			}
		}(id)

		id++
		time.Sleep(time.Duration(frequency) * time.Millisecond)
	}
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
	for _, target := range Targets {

		tIP := ""
		for ip, quota := range target.IPQuota {
			if quota == 0 || TotalMessages == 0 {
				tIP = ip
			} else {
				current := float64(target.MessageCounts[ip]) / float64(TotalMessages)
				log(fmt.Sprintf("%d Messages processed %d, target %s current %f quota %f\n", time.Now().UnixMilli(), TotalMessages, ip, current, quota))
				if current < quota {
					log(fmt.Sprintf("%d Sending to %s\n", time.Now().UnixMilli(), ip))
					tIP = ip
					break
				}
			}
		}
		go func(target string) {
			serviceUrl := fmt.Sprintf(config.Cfg.PushServiceURL, tIP)
			message.Hops = []NodeData{{NodeId: InstanceName, ExitTime: time.Now().UnixMicro()}}
			jsonData, err := json.Marshal(message)
			_, err = http.Post(serviceUrl, "application/json",
				bytes.NewBuffer(jsonData))
			log(fmt.Sprintf("%d Message id %s to service %s\n", time.Now().UnixMilli(), message.MessageId, serviceUrl))
			if err != nil {
				log(fmt.Sprintf("%d Failed to write to service %s\n", time.Now().UnixMilli(), serviceUrl))
			}

		}(tIP)
		target.MessageCounts[tIP] += 1
	}
	TotalMessages++
	log(fmt.Sprintf("%d Message id %s sent\n", time.Now().UnixMilli(), message.MessageId))
	//fmt.Printf("%d Message id %s sent\n", time.Now().UnixMilli(), message.MessageId)

	return "" //logline
}

func sendTestMessage(message Message) error {
	//message.History = []string{InstanceName}
	jsonData, err := json.Marshal(message)
	//fmt.Println(string(jsonData))
	for _, targetIP := range Targets {
		for ip := range targetIP.IPQuota {
			serviceUrl := fmt.Sprintf(config.Cfg.PushServiceURL, ip)
			_, err = http.Post(serviceUrl, "application/json",
				bytes.NewBuffer(jsonData))

			if err != nil {
				log(fmt.Sprintf("Failed to write to service %s\n", serviceUrl))
				log(err.Error())
				return err
			}
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
