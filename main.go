package main

import (
	"GoMCScan/mcping"
	"GoMCScan/mcping/types"
	"fmt"
	"log"
	"net"
	"strings"
	"time"

	"github.com/aherve/gopool"
)

var threads int
var timeout int
var outputPath string
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
			if !strings.Contains(address, "/") { //puts single addresses in CIDR notation
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

type formattedOutput struct {
	Timestamp string
	Ip string
	Version string
	Motd string 
	PlayersCount types.PlayerCount
	Sample []types.PlayerSample
}

func pingIt(ip string, port uint16) {
	defer pool.Done()
	data, _, err := mcping.PingWithTimeout(ip, port, time.Duration(timeout)*time.Second)
	completed++ // this is somewhat broken because of concurency
	if err == nil {

		found++ // also would be broken
		printStatus(fmt.Sprintf("%v:%v | %v", ip, port, data.Motd))
		formatted := formattedOutput{time.Now().Format("2006-01-02 15:04:04"), ip+":"+fmt.Sprint(port), data.Version, data.Motd, data.PlayerCount, data.Sample}
		record(formatted)
	}
}
