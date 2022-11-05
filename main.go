package main

import (
	"GoMCScan/mcping"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/aherve/gopool"
)

type PingResponse struct {
	//Latency     uint // Latency between you and the server
	//PlayerCount PlayerCount // Players count information of the server
	//Protocol    int // Protocol number of the server
	//Favicon     string // Favicon in base64 of the server
	Motd    string         // Motd of the server without color
	Version string         // Version of the server
	Sample  []PlayerSample // List of connected players on the server
}

type PlayerCount struct {
	Online int // Number of connected players on the server
	Max    int // Number of maximum players on the server
}

type PlayerSample struct {
	UUID string // UUID of the player
	Name string // Name of the player
}

func main() {
	threads, _ := strconv.Atoi(os.Args[1])
	pool := gopool.NewPool(int(threads))
	g := 176
	h := 9
	var port uint16
	port = 25565

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
		strDec := fmt.Sprintf("%#v", dec)
		strDec = strings.ReplaceAll(strDec, "types.PlayerSample", "")
		strDec = strings.ReplaceAll(strDec, "&types.PingResponse", "")
		strDec = strings.ReplaceAll(strDec, "(nil)", "")
		strDec = strings.ReplaceAll(strDec, "[]{", "{")
		strport := fmt.Sprintf("%#v", int(port))
		address := "IP:\"" + ip + ":" + strport + "\", "
		data := "{" + address + strDec[1:]
		fmt.Println(data)

		if strings.Contains(strings.ToLower(strDec), "minecraft") || strings.Contains(strings.ToLower(strDec), "long_int") {
			fmt.Println("Found")
			record(data)
		}
	} else {
		//fmt.Println(err)
	}
}

func record(data string) {
	f, err := os.OpenFile("./out/scan.log",
		os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Println(err)
	}
	defer f.Close()
	if _, err := f.WriteString(data + "\n"); err != nil {
		log.Println(err)
	}
}
