package config

import (
	"encoding/json"
	"fmt"
	"os"
)

var Cfg *Config

type Config struct {
	SubnetBridgeIPs map[string]string
	IPRouteMap      map[string]string
	NumServices     int
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

	return err
}
