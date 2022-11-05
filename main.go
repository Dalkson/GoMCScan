package main

import (
	"GoMCScan/mcping"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/aherve/gopool"
)

const usage = `Usage of MCScan:
    MCScan [-T Threads] [-t Timeout] [-p PortRange] [-o output]
Options:
    -T, --threads number of threads to use
    -t, --timeout timeout in seconds
    -h, --help prints help information
    -o, --output output location for scan file
`

var threads int
var timeout int
var output string
var portRange string
var pool *gopool.GoPool

func main() {
	flags()
	fmt.Println(portRange)
	pool = gopool.NewPool(threads)
	var a uint8 = 176
	var b uint8 = 9
	ports := []uint16{25565}
	loopBlock(a, b, ports)
	pool.Wait()
}

func flags() {
	flag.IntVar(&threads, "T", 1000, "number of threads to use")
	flag.IntVar(&threads, "threads", 1000, "number of threads to use")
	flag.IntVar(&timeout, "t", 1, "timeout in seconds")
	flag.IntVar(&timeout, "timeout", 1, "timeout in seconds")
	flag.StringVar(&output, "output", "out/scan.log", "output location for scan file")
	flag.StringVar(&output, "o", "out/scan.log", "output location for scan file")
	flag.StringVar(&portRange, "p", "25565-25570", "output location for scan file")
	flag.Usage = func() { fmt.Print(usage) }
	flag.Parse()
}

func loopBlock(a uint8, b uint8, ports []uint16) {
	for _, port := range ports {
		for c := 0; c < 255; c++ {
			for d := 0; d < 255; d++ {
				var ip = fmt.Sprintf("%v.%v.%v.%v", a, b, c, d)
				pool.Add(1)
				go pingIt(ip, port)
			}
		}
	}
}

func pingIt(ip string, port uint16) {
	defer pool.Done()
	data, _, err := mcping.PingWithTimeout(ip, port, time.Duration(timeout)*time.Second)
	if err == nil {
		sampleBytes, _ := json.Marshal(data.Sample)
		sample := string(sampleBytes)
		if sample == "null" {
			sample = "[]"
		}
		formatted := fmt.Sprintf("{\"Ip\":\"%v:%v\", \"Version\":%q, \"Motd\":%q, \"Players:%v/%v\", \"Sample\":%v}", ip, port, data.Version, data.Motd, data.PlayerCount.Online, data.PlayerCount.Max, sample)
		fmt.Println(formatted)
		record(formatted)
	} else {
		//fmt.Println(err)
	}
}

func record(data string) {
	f, err := os.OpenFile(output,
		os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Println(err)
	}
	defer f.Close()
	if _, err := f.WriteString(data + "\n"); err != nil {
		log.Println(err)
	}
}
