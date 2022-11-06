package main

import (
	"GoMCScan/mcping"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
	"time"
	"unicode"

	"github.com/aherve/gopool"
)

// const usage = `Usage of MCScan:
//     MCScan [-T Threads] [-t Timeout] [-p PortRange] [-o output]
// Options:
//     -T, --threads number of threads to use
//     -t, --timeout timeout in seconds
//     -h, --help prints help information
//     -o, --output output location for scan file
// `

var threads int
var timeout int
var outputDir string
var addressList []string
var portList []uint16

var pinged int
var completed int
var found int

var startTime time.Time
var pool *gopool.GoPool

func main() {
	flags()
	pool = gopool.NewPool(threads)
	fmt.Println("Scanning ports:", portList)
	loopBlock()
	pool.Wait()
	fmt.Println("Scan Complete!")
}

func flags() {
	flag.IntVar(&threads, "T", 1000, "number of threads to use")
	flag.IntVar(&threads, "threads", 1000, "number of threads to use")
	flag.IntVar(&timeout, "t", 1, "timeout in seconds")
	flag.IntVar(&timeout, "timeout", 1, "timeout in seconds")
	flag.StringVar(&outputDir, "output", "out/scan.log", "output location for scan file")
	flag.StringVar(&outputDir, "o", "out/scan.log", "output location for scan file")
	var portRange string
	flag.StringVar(&portRange, "p", "25565", "port range to scan")
	flag.StringVar(&portRange, "port", "25565", "port range to scan")
	var addressRange string
	flag.StringVar(&addressRange, "a", "", "IP address range to scan")
	flag.StringVar(&addressRange, "targets", "", "IP address range to scan")

	// flag.Usage = func() { fmt.Print(usage) }
	flag.Parse()

	expandAddress(addressRange)
	portList = expandPort(portRange)
}

func expandAddress(input string) {
	// Example string input: "176.9.0.0/16,116.202.0.0/16", output: [176.9.0.0/16,116.202.0.0/16]
	if input == "" {
		handleError("no target set, set one with -a or --targets")
	}
	for _, a := range input {
		if !(unicode.IsNumber(a) || a == ',' || a == '.' || a == '/') {
			handleError("Invalid characters in ports list. Valid characters include: \"123456789,./\"")
		}
	}
	addressList = strings.Split(input, ",")
}

func expandPort(input string) []uint16 {
	// Example input: "123,456-458,11111", output: [123,456,567,458,11111]
	for _, a := range input {
		if !(unicode.IsNumber(a) || a == ',' || a == '-') {
			handleError("Invalid characters in ports list. Valid characters include: \"123456789,-\"")
		}
	}
	var output []uint16
	for _, o := range strings.Split(input, ",") {
		if strings.Contains(o, "-") {
			test := strings.Split(o, "-")
			startPort, err1 := strconv.ParseInt(test[0], 10, 16)
			endPort, err2 := strconv.ParseInt(test[1], 10, 16)
			if err1 != nil || err2 != nil {
				handleError("Port could not be parsed to integer")
			}
			for port := uint16(startPort); port < uint16(endPort+1); port++ {
				output = append(output, port)
			}
		} else {
			port, err := strconv.ParseUint(o, 10, 16)
			if err != nil {
				handleError("Port could not be parsed to integer")
			}
			output = append(output, uint16(port))
		}
	}
	return output
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
		fmt.Println(formatted)
		found++
		fmt.Printf("%v/%v, %v percent complete\n", completed, pinged, uint8(100*float64(completed)/float64(pinged)))
		fmt.Printf("Time Elapsed: %v min, finding rate: %v servers per second", time.Since(startTime).Minutes(), int(float64(found)/float64(time.Since(startTime).Seconds())))
		record(formatted)
	} else {
		//fmt.Println(err)
	}
}

func record(data string) {
	f, err := os.OpenFile(outputDir,
		os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Println(err)
	}
	defer f.Close()
	if _, err := f.WriteString(data + "\n"); err != nil {
		log.Println(err)
	}
}

func handleError(err string) {
	fmt.Fprintf(os.Stderr, "error: %v\n", err)
	os.Exit(1)
}
