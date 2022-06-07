package containers

import (
	"customfledge/config"
	"customfledge/utils"
	"fmt"
	"strconv"
	"strings"
)

//var baseSubnetIP int
//var maxSubnetIP int
//var subnetMask int
//var nodeSubnetsMasks []string
//var gatewayIPs []string
//var gatewayIP string

//var usedAddresses map[int]string

func InitContainerNetworking() {
	cmd := fmt.Sprintf("sh -x ./startcni.sh")
	_, err := utils.ExecCmdBash(cmd)
	if err != nil {
		fmt.Println("Could not set up cni0")
		panic(err)
	}

	for subnetMask, brip := range config.Cfg.SubnetBridgeIPs {
		subnet := strings.Split(subnetMask, "/")[0]
		mask := strings.Split(subnetMask, "/")[1]
		cmd = fmt.Sprintf("sh -x ./addcniiip.sh %s %s %s", subnet, mask, brip)
		_, err = utils.ExecCmdBash(cmd)
	}
	//nodeSubnetsMasks = subnetsMasks
	//subnetMask, _ = strconv.Atoi(subMask)
	//baseSubnetIP, _ = IPStringToInt(nodeSubnet)
	//maxSubnetIP = baseSubnetIP + int(math.Pow(2, float64(subnetMask)))
	/*gatewayIPs = []string{}
	  for _, iprange := range nodeSubnetsMasks {
	  	subMask, _ := strconv.Atoi(strings.Split(iprange, "/")[1])
	  	baseSubnetIP, _ := IPStringToInt(strings.Split(iprange, "/")[0])
	  	gatewayIP, _ := IPIntToString(baseSubnetIP + int(math.Pow(2, float64(subMask))) - 1)
	  	gatewayIPs = append(gatewayIPs, gatewayIP)
	  }*/
	//usedAddresses = make(map[int]string)
}

func SetupRoutes(ipRouteMap map[string]string) {
	// ip route add $nodeip/$mask via $routeip dev $extif

	for svcIP, publicIP := range config.Cfg.IPRouteMap {
		cmd := fmt.Sprintf("sh -x ./addroute.sh %s %s", svcIP, publicIP)
		_, err := utils.ExecCmdBash(cmd)
		if err != nil {
			fmt.Printf("Failed to set up route %s to %s", publicIP, svcIP)
		}
	}
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
