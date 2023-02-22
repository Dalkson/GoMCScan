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

// options holds configuration options for the scan
type options struct {
	threads int       // the number of threads to use for the scan
	timeout int       // the timeout for each ping in seconds
	addressList []string // the list of IP addresses to scan
	portList []uint16    // the list of ports to scan
	outputPath string    // the path to save the output of the scan
	saveFavicon bool     // whether to save the server's favicon
}

// statistics holds various statistics about the scan
type statistics struct {
	pinged int       // the number of pings sent
	completed int    // the number of pings that have completed
	found int        // the number of servers that responded to the ping
	total int        // the total number of servers to scan
	startTime time.Time // the time the scan started
}

var stats statistics   // holds the current statistics for the scan
var conf options       // holds the current configuration options for the scan

var pool *gopool.GoPool // the goroutine pool for the scan

// main is the entry point of the program
func main() {
	conf = getFlags()  // parse command-line flags to set the values of conf
	stats.total = totalToSend()  // calculate the total number of servers to scan
	pool = gopool.NewPool(conf.threads)  // initialize the goroutine pool
	fmt.Println("GoMCScan starting. for more information, use -h")
	fmt.Println("Total to scan:", stats.total)
	fmt.Println("- - -")
	go logLoop(2 * time.Second)  // start the log loop
	loopBlock()  // start the scan
	pool.Wait()  // wait for the scan to complete
	fmt.Println("Scan Complete!")
}

// loopBlock scans the list of IP addresses and ports specified in conf
func loopBlock() {
	stats.startTime = time.Now()  // record the start time of the scan
	for _, port := range conf.portList {  // loop through the list of ports
		for _, address := range conf.addressList {  // loop through the list of addresses
			if !strings.Contains(address, "/") { // if the address is not in CIDR notation, put it in CIDR notation
				address = fmt.Sprintf("%v/32", address)
			}
			ip, ipnet, err := net.ParseCIDR(address)  // parse the address and determine the range of addresses to scan
			if err != nil {
				log.Fatal(err)  // if an error occurs, log the error and exit
			}
			for ip := ip.Mask(ipnet.Mask); ipnet.Contains(ip); incrementIP(ip) {  // iterate through the range
				pool.Add(1)  // add a new goroutine to the pool
				go pingIt(string(net.IP.String(ip)), port)  // start a goroutine to ping the server
				stats.pinged++  // increment the ping counter
			}
		}
	}
}

// incrementIP increments the given IP address
func incrementIP(ip net.IP) {
	for j := len(ip) - 1; j >= 0; j-- {
		ip[j]++
		if ip[j] > 0 {
			break
		}
	}
}

// formattedOutput holds the output of a successful ping in a formatted way
type formattedOutput struct {
	Timestamp string  // the time the ping was sent
	IP string  // the IP address and port of the server
	Version string  // the version of the server
	Motd string   // the message of the day of the server
	PlayersCount types.PlayerCount  // the number of players online on the server
	Sample []types.PlayerSample  // a sample of the players online on the server
}

// pingIt pings the given server and logs the response
func pingIt(ip string, port uint16) {
	defer pool.Done()  // mark the goroutine as done when it returns
	data, _, err := mcping.PingWithTimeout(ip, port, time.Duration(conf.timeout)*time.Second)  // ping the server
	stats.completed++  // increment the completed counter
	if err == nil {  // if the ping was successful
		stats.found++  // increment the found counter
		printStatus(fmt.Sprintf("%v:%v | %v  online | %v", ip, port, data.PlayerCount.Online, data.Motd))  // log the response
		formatted := formattedOutput{time.Now().Format("2006-01-02 15:04:04"), ip+":"+fmt.Sprint(port), data.Version, data.Motd, data.PlayerCount, data.Sample}  // create a formattedOutput struct
		record(formatted)  // record the output
		if conf.saveFavicon && data.Favicon != "" {  // if the saveFavicon option is set and the server has a favicon
				saveFavicon(data.Favicon, ip, port)  // save the favicon
		}
	}
}
