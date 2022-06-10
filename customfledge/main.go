package main

import (
	"customfledge/config"
	"customfledge/containers"
	"encoding/json"
	"fmt"
	"os"

	v1 "k8s.io/api/core/v1"
)

func main() {
	argsWithoutProg := os.Args[1:]
	cfgFile := "defaultconfig.json"
	if len(argsWithoutProg) > 0 {
		cfgFile = argsWithoutProg[0]
	}

	config.LoadConfig(cfgFile)

	cri := (&containers.ContainerdRuntimeInterface{}).Init()

	containers.InitCgroups()
	containers.InitContainerNetworking()

	for i := 0; i < config.Cfg.NumServices; i++ {
		fmt.Printf("Reading svc%d.json\n", i)
		jsonBytes, err := os.ReadFile(fmt.Sprintf("svc%d.json", i))
		if err != nil {
			fmt.Printf("Failed to read svc%d.json", i)
		}
		pod := &v1.Pod{}
		err = json.Unmarshal(jsonBytes, pod)
		if err != nil {
			fmt.Printf("Failed to parse svc%d.json", i)
		}
		cri.DeployPod(pod)
	}

	containers.SetupRoutes(config.Cfg.IPRouteMap)
}
