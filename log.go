package main

import (
	"fmt"
	"log"
	"math"
	"os"
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
	const percent = "%"
	percentage := math.Round(100 * float64(completed) / float64(total))
	elapsed := math.Round(time.Since(startTime).Minutes()*10) / 10
	elapsedSec := float64(time.Since(startTime).Seconds())
	findRate := math.Round(float64(found) / float64(elapsedSec))
	pingRate := math.Round(float64(pinged) / float64(elapsedSec))
	fmt.Printf("%v%v | Found: %v at %v/s | Pinged: %v at %v/s | Time: %vm | %v \n", percentage, percent, found, findRate, pinged, pingRate, elapsed, announce)
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
