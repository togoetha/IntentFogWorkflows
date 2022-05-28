package config

import (
	"encoding/json"
	"fmt"
	"os"
)

var Cfg *Config

type Config struct {
	InstanceName        string  `json:"instanceName"`
	MessageFrequency    float32 `json:"messageFrequency"`
	DefaultWorkloadSize int     `json:"defaultWorkloadSize"`
	PayloadSize         int     `json:"payloadSize"`
	ServiceMode         bool    `json:"serviceMode"`
	PushServiceURL      string  `json:"pushServiceURL"`
	MqttBroker          string  `json:"mqttBroker"`
	MqttTopicWrite      string  `json:"mqttTopicWrite"`
	MessageTemplate     string  `json:"messageTemplate"`
	MqttClientId        string  `json:"mqttClientId"`
	MqttUser            string  `json:"mqttUser"`
	MqttPass            string  `json:"mqttPass"`
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
