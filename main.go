package main

import (
	"GoMCScan/mcping"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"strings"
	"time"

	"github.com/aherve/gopool"
)

var threads int
var timeout int
var outputDir string
var addressList []string
var portList []uint16

var pinged int
var completed int
var found int
var total int

var startTime time.Time
var pool *gopool.GoPool

func main() {
	getFlags()
	total = totalToSend()
	pool = gopool.NewPool(threads)
	fmt.Println("Scanning ports:", portList)
	fmt.Println("Total to scan:", total)
	go logLoop(2 * time.Second)
	loopBlock()
	pool.Wait()	
	fmt.Println("Scan Complete!")
}

func loopBlock() {
	startTime = time.Now()
	for _, port := range portList {
		for _, address := range addressList {
			if !strings.Contains(address, "/") {
				address = fmt.Sprintf("%v/32", address)
			}
			ip, ipnet, err := net.ParseCIDR(address)
			if err != nil {
				log.Fatal(err)
			}
			for ip := ip.Mask(ipnet.Mask); ipnet.Contains(ip); inc(ip) {
				pool.Add(1)
				go pingIt(string(net.IP.String(ip)), port)
				pinged++
			}
		}
	}
}

func inc(ip net.IP) {
	for j := len(ip) - 1; j >= 0; j-- {
		ip[j]++
		if ip[j] > 0 {
			break
		}
	}
}

func pingIt(ip string, port uint16) {
	defer pool.Done()

	data, _, err := mcping.PingWithTimeout(ip, port, time.Duration(timeout)*time.Second)
	completed++
	if err == nil {
		sampleBytes, _ := json.Marshal(data.Sample)
		sample := string(sampleBytes)
		if sample == "null" {
			sample = "[]"
		}
		formatted := fmt.Sprintf("{\"Timestamp\":%q, \"Ip\":\"%v:%v\", \"Version\":%q, \"Motd\":%q, \"Players:%v/%v\", \"Sample\":%v}", time.Now().Format("2006-01-02 15:04:05"), ip, port, data.Version, data.Motd, data.PlayerCount.Online, data.PlayerCount.Max, sample)
		found++
		printStatus(fmt.Sprintf("%v:%v | %v", ip, port, data.Motd))
		record(formatted)
	} else {
		//fmt.Println(err)
	}

}
