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

func main() {
	threads := flag.Int("threads", 1000, "number of threads to use")
	timeout := flag.Int("timeout", 1, "timeout in seconds")
	flag.Parse()
	pool := gopool.NewPool(*threads)
	var g uint8 = 176
	var h uint8 = 9
	ports := []uint16{25565}
	loopBlock(timeout, g, h, ports, pool)
	pool.Wait()
}

func loopBlock(timeout *int, a uint8, b uint8, ports []uint16, pool *gopool.GoPool) {
	for _, port := range ports {
		for j := 0; j < 255; j++ {
			for k := 0; k < 255; k++ {
				var ip = fmt.Sprintf("%v.%v.%v.%v", a, b, j, k)
				pool.Add(1)
				go pingIt(ip, port, timeout, pool)
			}
		}
	}
}

func pingIt(ip string, port uint16, timeout *int, pool *gopool.GoPool) {
	defer pool.Done()
	data, _, err := mcping.PingWithTimeout(ip, port, time.Duration(*timeout)*time.Second)
	if err == nil {
		fmt.Println("Found")
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
	f, err := os.OpenFile("out/scan.log",
		os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Println(err)
	}
	defer f.Close()
	if _, err := f.WriteString(data + "\n"); err != nil {
		log.Println(err)
	}
}
