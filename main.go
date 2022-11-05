package main

import (
	"GoMCScan/mcping"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/aherve/gopool"
)

func main() {
	threads, _ := strconv.Atoi(os.Args[1])
	pool := gopool.NewPool(int(threads))
	g := 176
	h := 9
	var port uint16 = 25565

	for j := 0; j < 255; j++ {
		for k := 0; k < 255; k++ {
			var ip = fmt.Sprintf("%v.%v.%v.%v", g, h, j, k)
			pool.Add(1)
			go pingIt(ip, port, pool)
		}
	}
	pool.Wait()
}

func pingIt(ip string, port uint16, pool *gopool.GoPool) {
	defer pool.Done()
	timeout, _ := strconv.Atoi(os.Args[2])
	dec, _, err := mcping.PingWithTimeout(ip, port, time.Duration(timeout)*time.Second)
	if err == nil {
		fmt.Println("Found")
		sampleBytes, _ := json.Marshal(dec.Sample)
		sample := string(sampleBytes)
		if sample == "null" {
			sample = "[]"
		}
		formatted := fmt.Sprintf("{\"Ip\":\"%v:%v\", \"Version\":%q, \"Motd\":%q, \"Players:%v/%v\",  \"Sample\":%v}", ip, port, dec.Version, dec.Motd, dec.PlayerCount.Online, dec.PlayerCount.Max, sample)
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
