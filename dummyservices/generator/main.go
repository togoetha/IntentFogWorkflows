package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"generator/config"
	"generator/message"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	easyjson "github.com/mailru/easyjson"
)

var InstanceName string
var MessageFrequency int
var Targets []Target
var TotalMessages int

type Target struct {
	IPQuota       map[string]float64
	MessageCounts map[string]int
}

var messageData string

func main() {
	/*fmt.Println("JSON test")
	data := []byte{}
	for i := 0; i < 30000; i++ {
		data = append(data, byte(70+i%10))
	}

	msg := message.Message{
		Payload: string(data),
	}
	msg.MessageId = fmt.Sprintf("%s%d", "testinstance", 1000)
	msg.Workload = 1000
	msg.Hops = []message.NodeData{
		{NodeId: InstanceName, ExitTime: time.Now().UnixMicro()},
		{NodeId: InstanceName, ExitTime: time.Now().UnixMicro()},
		{NodeId: InstanceName, ExitTime: time.Now().UnixMicro()},
	}

	jsonData, _ := json.Marshal(msg)
	fmt.Println(string(jsonData))

	start := time.Now()

	for i := 0; i < 1000; i++ {
		json.Marshal(msg)
	}

	fmt.Printf("100 marshals took %dms\n", time.Since(start).Milliseconds())
	*/
	//jsonStrData := []byte{}
	/*start = time.Now()
	for i := 0; i < 1000; i++ {
		bytes := toFancyJson(msg)
		if i == 0 {
			//jsonStrData = bytes
			fmt.Println(string(bytes))
		}
	}*/

	//fmt.Printf("100 manual marshals took %dms\n", time.Since(start).Milliseconds())

	/*start = time.Now()

	for i := 0; i < 1000; i++ {
		umsg := message.Message{}
		json.Unmarshal(jsonData, &umsg)
	}

	fmt.Printf("100 unmarshals took %dms\n", time.Since(start).Milliseconds())

	start = time.Now()

	for i := 0; i < 1000; i++ {
		jsonStr := ""
		json.Unmarshal(jsonStrData, &jsonStr)
		fromJson(jsonStr[0:300])
	}

	fmt.Printf("100 manual unmarshals took %dms\n", time.Since(start).Milliseconds())
	*/
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

/*func toFancyJson(msg message.Message) []byte {
	//{"hops":[{"history":"","startTime":0,"endTime":1656010631547025},{"history":"","startTime":0,"endTime":1656010631547025},{"history":"","startTime":0,"endTime":1656010631547025}],"workload":1000,"payload":"","messageId":"testinstance1000"}
	hops := []string{}
	for _, data := range msg.Hops {
		hops = append(hops, fmt.Sprintf("{\x22history\x22:\x22%s\x22,\x22startTime\x22:%d,\x22endTime\x22:%d}", data.NodeId, data.EntryTime, data.ExitTime))
	}
	str := fmt.Sprintf("{\x22hops\x22:[%s],\x22workload\x22:%d,\x22payload\x22:\x22%s\x22,\x22messageId\x22:\x22%s\x22}", strings.Join(hops, ","), msg.Workload, msg.Payload, msg.MessageId)
	jsonData, _ := json.Marshal(str)

	return jsonData
}*/

//var running bool

func generateMessage() message.Message {
	message := message.Message{
		Payload: messageData,
	}

	return message
}

var client *http.Client

func generate() {
	tr := http.DefaultTransport.(*http.Transport).Clone()
	tr.MaxIdleConns = 0
	tr.MaxIdleConnsPerHost = 0
	tr.MaxConnsPerHost = 0
	client = &http.Client{Transport: tr}

	frequency := 1000.0 / float32(MessageFrequency)
	id := 1
	//running = true

	data := []byte{}
	for i := 0; i < config.Cfg.PayloadSize; i++ {
		data = append(data, byte(70+i%10))
	}
	messageData = string(data)

	msg := generateMessage()
	msg.MessageId = "Test1000000"
	msg.Workload = config.Cfg.DefaultWorkloadSize
	msg.Hops = []message.NodeData{{NodeId: InstanceName, ExitTime: time.Now().UnixMicro()}}
	err := sendTestMessage(msg)
	for err != nil {
		log(fmt.Sprintf("%d Can't reach clients yet, retrying\n", time.Now().UnixMilli()))
		err = sendTestMessage(msg)
	}

	//loglines := []string{}
	for i := 0; i < config.Cfg.Messages; i++ {
		go func(msgId int) {
			msg := generateMessage()
			msg.MessageId = fmt.Sprintf("%s%d", InstanceName, msgId)
			msg.Workload = config.Cfg.DefaultWorkloadSize

			sendRESTMessage(msg)
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

func sendRESTMessage(msg message.Message) string {
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
			msg.Hops = []message.NodeData{{NodeId: InstanceName, ExitTime: time.Now().UnixMicro()}}
			jsonData, err := json.Marshal(msg)
			resp, err := http.Post(serviceUrl, "application/json",
				bytes.NewBuffer(jsonData))
			resp.Body.Close()

			log(fmt.Sprintf("%d Message id %s to service %s\n", time.Now().UnixMilli(), msg.MessageId, serviceUrl))
			if err != nil {
				log(fmt.Sprintf("%d Failed to write to service %s\n", time.Now().UnixMilli(), serviceUrl))
			}

		}(tIP)
		target.MessageCounts[tIP] += 1
	}
	TotalMessages++
	log(fmt.Sprintf("%d Message id %s sent\n", time.Now().UnixMilli(), msg.MessageId))
	//fmt.Printf("%d Message id %s sent\n", time.Now().UnixMilli(), message.MessageId)

	return "" //logline
}

func sendTestMessage(msg message.Message) error {
	//message.History = []string{InstanceName}
	jsonData, err := easyjson.Marshal(msg)
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
