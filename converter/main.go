package main

import (
	"converter/config"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"

	yaml "gopkg.in/yaml.v3"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func main() {
	argsWithoutProg := os.Args[1:]
	cfgFile := "defaultconfig.json"
	if len(argsWithoutProg) == 0 {
		//	cfgFile = argsWithoutProg[0]
		fmt.Println("Use: converter <inputyamlfilename>")
		//os.Exit(1)
	}

	//fmt.Printf("Loading config file %s\n", cfgFile)
	config.LoadConfig(cfgFile)

	inputYamlFile := "test.yml"
	if len(argsWithoutProg) > 0 {
		inputYamlFile = argsWithoutProg[0]
	}
	convertToDeployments(inputYamlFile)
}

var pod *corev1.Pod

func convertToDeployments(inputYamlFile string) {
	exp := ExperimentInfo{}

	data, err := os.ReadFile(inputYamlFile)
	if err != nil {
		fmt.Printf("Couldn't read file %s", inputYamlFile)
		panic(err)
	}
	err = yaml.Unmarshal(data, &exp)
	if err != nil {
		fmt.Printf("Couldn't parse file %s", inputYamlFile)
		panic(err)
	}

	nodeSubnets := []string{exp.Network[0].Msubnet, exp.Network[0].Dsubnet}
	firstSubnet := exp.Network[0].Igresssubnet
	nodeSubnets = append(nodeSubnets, firstSubnet)
	lastSubnet := exp.Network[0].Egresssubnet
	if firstSubnet != lastSubnet {
		nodeSubnets = append(nodeSubnets, lastSubnet)
	}

	nodeServices := make(map[string][]corev1.Pod)
	nodeConfigs := make(map[string]NodeConfig)

	for counter, nodeInfo := range exp.Nodes {
		//handle first and last differently
		subnets := make(map[string]string)
		for _, subnet := range nodeSubnets {
			//subMask, _ := strconv.Atoi(strings.Split(subnet, "/")[1])
			baseSubnetIP, _ := IPStringToInt(strings.Split(subnet, "/")[0])
			subnets[subnet], _ = IPIntToString(baseSubnetIP + 254 - counter) //int(math.Pow(2, float64(subMask))) - counter - 1)
		}

		nodeConfigs[nodeInfo.Name] = NodeConfig{
			NodeName:        nodeInfo.Name,
			SubnetBridgeIPs: subnets,
			IPRouteMap:      routeMap(nodeInfo, exp.Nodes),
			NumServices:     len(nodeInfo.Services),
		}

		nodesvcs := []corev1.Pod{}
		for _, podInfo := range nodeInfo.Services {
			args := []string{"defaultconfig.json", podInfo.Name}
			addresses := getNextAddresses(podInfo, exp.Links)
			for _, address := range addresses {
				args = append(args, strings.Split(address, "/")[0])
			}
			pod := corev1.Pod{
				ObjectMeta: v1.ObjectMeta{
					Name:      podInfo.Name,
					Namespace: "default",
					Labels: map[string]string{
						"MAddress":   podInfo.Mipaddr,
						"DAddresses": strings.Join(podInfo.Dipaddr, ","),
						"Bandwidth":  podInfo.Bandwidth,
						"Latency":    strconv.Itoa(podInfo.Latency),
					},
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:            podInfo.Name,
							Image:           podInfo.Image,
							ImagePullPolicy: corev1.PullIfNotPresent,
							Command:         []string{"./service"},
							Args:            args,
							Resources: corev1.ResourceRequirements{
								Limits: corev1.ResourceList{
									corev1.ResourceCPU:    resource.MustParse(podInfo.Cpu),
									corev1.ResourceMemory: resource.MustParse(podInfo.Memory),
								},
							},
						},
					},
				},
			}
			nodesvcs = append(nodesvcs, pod)
		}
		nodeServices[nodeInfo.Name] = nodesvcs
	}

	err = os.Mkdir("configs", 0777)
	for node, cfg := range nodeConfigs {
		err = os.MkdirAll(fmt.Sprintf("configs/%s", node), 0777)
		contents, err := json.Marshal(cfg)
		if err == nil {
			os.WriteFile(fmt.Sprintf("configs/%s/defaultconfig.json", node), contents, 0777)
		}
	}

	for node, svcs := range nodeServices {
		for id, svc := range svcs {
			//os.Mkdir(fmt.Sprintf("configs/%s", node), os.ModeDir)
			contents, err := json.Marshal(svc)
			if err == nil {
				os.WriteFile(fmt.Sprintf("configs/%s/svc%d.json", node, id), contents, 0777)
			}
		}
	}
}

func routeMap(local NodeInfo, nodes []NodeInfo) map[string]string {
	routes := make(map[string]string)

	for _, node := range nodes {
		if node.Name != local.Name {
			for _, svc := range node.Services {
				for _, ip := range svc.Dipaddr {
					routes[ip] = node.Publicip
				}
				routes[svc.Mipaddr] = node.Publicip
			}
		}
	}

	return routes
}

func getNextAddresses(svcInfo ServiceInfo, links []LinkInfo) []string {
	addresses := []string{}

	for _, link := range links {
		for _, dAddr := range svcInfo.Dipaddr {
			for _, inAddr := range link.In {
				if dAddr == inAddr {
					addresses = append(addresses, link.Out...)
				}
			}
		}
	}
	return addresses
}

func IPIntToString(ip int) (string, error) {
	var ipStr string = ""

	for ip > 0 {
		ipPart := ip % 256
		ipStr = strconv.Itoa(ipPart) + "." + ipStr
		ip = (ip - ipPart) / 256
	}

	return ipStr[0 : len(ipStr)-1], nil
}

func IPStringToInt(ipStr string) (int, error) {
	parts := strings.Split(ipStr, ".")
	var ip int = 0
	for i := 0; i < len(parts); i++ {
		p, _ := strconv.Atoi(parts[i])
		ip += p << uint(24-i*8)
	}
	return ip, nil
}

type NodeConfig struct {
	NodeName        string
	SubnetBridgeIPs map[string]string
	IPRouteMap      map[string]string
	NumServices     int
}

type ExperimentInfo struct {
	Network []NetworkInfo
	Nodes   []NodeInfo
	Links   []LinkInfo
}

type NetworkInfo struct {
	Name         string
	Msubnet      string
	Dsubnet      string
	P2pdsubnets  []string
	Igresssubnet string
	Egresssubnet string
	Numservices  int
	Numedges     int
	Numnodes     int
}

type NodeInfo struct {
	Name     string
	Publicip string
	Services []ServiceInfo
}

type ServiceInfo struct {
	Name      string
	Image     string
	Bandwidth string
	Memory    string
	Cpu       string
	Latency   int
	Mipaddr   string
	Dipaddr   []string
}

type LinkInfo struct {
	Name string
	In   []string
	Out  []string
}
