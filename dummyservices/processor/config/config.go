package config

import (
	"encoding/json"
	"fmt"
	"os"
)

var Cfg *Config

type Config struct {
	PayloadSize int `json:"payloadSize"`
	//DefaultWorkloadSize int    `json:"defaultWorkloadSize"`
	ServiceMode     bool   `json:"serviceMode"`
	PushServiceURL  string `json:"pushServiceURL"`
	MqttBroker      string `json:"mqttBroker"`
	MqttTopicRead   string `json:"mqttTopicRead"`
	MqttTopicWrite  string `json:"mqttTopicWrite"`
	MessageTemplate string `json:"messageTemplate"`
	MqttClientId    string `json:"mqttClientId"`
	MqttUser        string `json:"mqttUser"`
	MqttPass        string `json:"mqttPass"`
}

func LoadConfig(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		//return err
	}
	decoder := json.NewDecoder(file)
	Cfg = &Config{}
	err = decoder.Decode(Cfg)
	if err != nil {
		fmt.Println(err.Error())
		//return err
	}

	return err
}
