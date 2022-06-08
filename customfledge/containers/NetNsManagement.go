package containers

import (
	"customfledge/config"
	"customfledge/utils"
	"fmt"
	"strings"

	v1 "k8s.io/api/core/v1"
)

func GetNetNs(namespace string, pod string) string {
	return namespace + "-" + pod
}

func BindNetNamespace(namespace string, pod string, pid int, bandwidth int64, latency int, podips []string) {
	netNs := GetNetNs(namespace, pod)

	ipsPerSubnet := make(map[string][]string)
	for subnet, gwip := range config.Cfg.SubnetBridgeIPs {
		subnetParts := strings.Split(subnet, "/")
		fmt.Printf("Matching pod IPs for subnet %s\n", subnetParts[0])
		mask := subnetParts[1]
		ipParts := strings.Split(subnetParts[0], ".")

		ips := []string{}
		for _, podip := range podips {
			podIpParts := strings.Split(podip, ".")
			fmt.Printf("Matching IP %s\n", podip)
			if ipParts[0] == podIpParts[0] && ipParts[1] == podIpParts[1] && ipParts[2] == podIpParts[2] {
				ips = append(ips, podip)
			}
		}

		subnetIdx := fmt.Sprintf("%s/%s", gwip, mask)
		ipsPerSubnet[subnetIdx] = ips
	}
	//ip, _ := RequestIP(namespace, pod)
	//This used to set up net namespace + pid IP address and some routing shit
	//cmd := fmt.Sprintf("sh -x ./setupcontainercni.sh %s %d eth1 %s %d %s", netNs, pid, ip, subnetMask, gatewayIP)
	//Here we just need to set up net namespace + array of IP addresses, routing is done separately after starting all containers
	counter := 1
	for subnet, ips := range ipsPerSubnet {
		if len(ips) > 0 {
			gwip := strings.Split(subnet, "/")[0]
			subnetMask := strings.Split(subnet, "/")[1]
			cniif := fmt.Sprintf("eth%d", counter)
			cmd := fmt.Sprintf("sh -x ./setupcontainercni.sh %s %d %s %s %s %d %d %s", netNs, pid, cniif, gwip, subnetMask, bandwidth, latency, strings.Join(ips, " "))
			utils.ExecCmdBash(cmd)
		}
		counter++
	}
}

func GetNetworkNamespace(namespace string, pod *v1.Pod) string {
	nsName := namespace + "-" + pod.ObjectMeta.Name
	//cmd := fmt.Sprintf("ip netns add %s", nsName)
	//ExecCmdBash(cmd)

	nsPath := fmt.Sprintf("/var/run/netns/%s", nsName)
	return nsPath
}

/*func SetupNetNamespace(namespace string, pod string) string {
	netNs := GetNetNs(namespace, pod)
	fmt.Printf("Setting up network namespace %s", netNs)
	ip, _ := RequestIP(namespace, pod)
	fmt.Printf("Setting up pod veth netns %s ip %s subnet %d gateway %s", netNs, ip, subnetMask, gatewayIP)
	cmd := fmt.Sprintf("sh -x /setupcontainerveth.sh %s eth1 %s %d %s", netNs, ip, subnetMask, gatewayIP)
	ExecCmdBash(cmd)
	return ip
}*/

func RemoveNetNamespace(namespace string, pod string) {
	netNs := GetNetNs(namespace, pod)
	cmd := fmt.Sprintf("sh -x ./shutdowncontainercni.sh %s eth1", netNs)
	utils.ExecCmdBash(cmd)
}
