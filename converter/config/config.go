package config

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
)

var Cfg *Config

type Config struct {
	PreferredRuntime  string `json:"preferredRuntime"`
	DeviceName        string `json:"deviceName"`
	ShortDeviceName   string
	DeviceIP          string `json:"deviceIP"`
	ServicePort       string `json:"servicePort"`
	KubeletPort       string `json:"kubeletPort"`
	VkubeServiceURL   string `json:"vkubeServiceURL"`
	UseKubeAPI        bool   `json:"useKubeAPI"`
	FledgeAPIPort     int    `json:"fledgeAPIPort"`
	IgnoreKubeProxy   string `json:"ignoreKubeProxy"`
	ExternalInterface string `json:"interface"`
	HeartbeatTime     int    `json:"heartbeatTime"`
	UseCLVers         bool   `json:"useCLVers"`
}

func LoadConfig(filename string) error {
	fmt.Printf("Loading config %s\n", filename)
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

	fmt.Printf("VkubeServiceURL check %s\n", Cfg.VkubeServiceURL)
	if Cfg.VkubeServiceURL == "" {
		fmt.Printf("Loading from env instead")
		//Cfg.Runtime = os.Getenv("FLEDGE_RUNTIME")
		Cfg.DeviceName = os.Getenv("FLEDGE_DEVICE_NAME")
		Cfg.DeviceIP = os.Getenv("FLEDGE_DEVICE_IP")
		Cfg.ServicePort = os.Getenv("FLEDGE_SERVICE_PORT")
		Cfg.KubeletPort = os.Getenv("FLEDGE_KUBELET_PORT")
		Cfg.VkubeServiceURL = os.Getenv("FLEDGE_VKUBE_URL")
		Cfg.IgnoreKubeProxy = os.Getenv("FLEDGE_IGNORE_KPROXY")
		Cfg.ExternalInterface = os.Getenv("FLEDGE_INET_INTERFACE")
		Cfg.HeartbeatTime, _ = strconv.Atoi(os.Getenv("HEARTBEAT_TIME"))
	}

	return err
}
