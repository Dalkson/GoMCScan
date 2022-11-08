package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	"os"
	"path"
	"time"
)

func logLoop(interval time.Duration) {
	for t := range time.Tick(1 * time.Second) {
		logsleep(t, interval)
	}

}

func logsleep(tick time.Time, interval time.Duration) {
	printStatus("")
	time.Sleep(interval)

}

func printStatus(announce string) {
	percentage := math.Round(100 * float64(stats.completed) / float64(stats.total))
	elapsed := math.Round(time.Since(stats.startTime).Minutes()*10) / 10
	elapsedSec := float64(time.Since(stats.startTime).Seconds())
	findRate := math.Round(float64(stats.found) / float64(elapsedSec))
	pingRate := math.Round(float64(stats.pinged) / float64(elapsedSec))
	completedRate := math.Round(float64(stats.completed) / float64(elapsedSec))
	remaining := math.Round((float64(stats.total-stats.completed) / float64(completedRate))/6)/10
	fmt.Printf("%v%% | Found: %v at %v/s | Pinged: %v at %v/s | Time: %vm, %vm rem | %v \n", percentage, stats.found, findRate, stats.pinged, pingRate, elapsed, remaining, announce)
}

var outputExists bool = false

func record(dataJSON formattedOutput) {
	if !outputExists {
		_, err := os.Stat(path.Dir(conf.outputPath))
		if os.IsNotExist(err) {
			os.Mkdir(path.Dir(conf.outputPath), 0750)
		}
		outputExists = true
	}
	dataBytes, _ := json.Marshal(dataJSON)
	dataString := string(dataBytes)
	f, err := os.OpenFile(conf.outputPath,
		os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		handleError("Could not open output file")
	}
	defer f.Close()
	if _, err := f.WriteString(dataString + "\n"); err != nil {
		log.Println(err)
	}
}

func handleError(err string) {
	fmt.Fprintf(os.Stderr, "error: %v\n", err)
	os.Exit(1)
}
