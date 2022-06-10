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
	output, err := utils.ExecCmdBash(cmd)
	fmt.Println(output)
	if err != nil {
		fmt.Println("Could not set up cni0")
	}

	for subnetMask, brip := range config.Cfg.SubnetBridgeIPs {
		subnet := strings.Split(subnetMask, "/")[0]
		mask := strings.Split(subnetMask, "/")[1]
		cmd = fmt.Sprintf("sh -x ./addcniip.sh %s %s %s", subnet, mask, brip)
		output, err = utils.ExecCmdBash(cmd)
		fmt.Println(output)
	}
	//init root class on phys ifs for delay per IP
	configDone := make(map[string]bool)

	counter := 0
	if config.Cfg.Delays != nil {
		for extIp, delay := range config.Cfg.Delays {
			publicIf := GetRouteDev(extIp)
			_, contains := configDone[publicIf]
			if !contains {
				cmd = fmt.Sprintf("sh -x ./setupdelay.sh %s", publicIf)
				configDone[publicIf] = true
			}
			class := 5 + counter
			classId := fmt.Sprintf("1:%d", class)
			handle := 10 * counter
			cmd = fmt.Sprintf("sh -x ./addipdelay.sh %s %s %s %d %s", publicIf, classId, delay, handle, extIp)
			counter++
		}
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

	for svcIPMask, publicIP := range config.Cfg.IPRouteMap {
		routeDev := GetRouteDev(publicIP)

		svcIP := strings.Split(svcIPMask, "/")[0]
		cmd := fmt.Sprintf("sh -x ./addroute.sh %s %s %s", svcIP, publicIP, routeDev)
		_, err := utils.ExecCmdBash(cmd)
		if err != nil {
			fmt.Printf("Failed to set up route %s to %s\n", publicIP, svcIP)
		}
	}
}

func GetRouteDev(publicIP string) string {
	cmd := fmt.Sprintf("ip route get %s | grep -E -o '[0-9\\.]* dev [a-z0-9]*'", publicIP)
	route, err := utils.ExecCmdBash(cmd)
	if err != nil {
		fmt.Println("Failed to determine public dev")
	}
	routeDev := strings.Split(route, " ")[2]
	return routeDev
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
