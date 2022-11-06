package main

import (
	"flag"
	"fmt"
)

func getFlags() {
	const usage = `Usage of MCScan:
    MCScan [-T Threads] [-t Timeout] [-p PortRange] [-o output]
Options:
    -T, --threads number of threads to use
    -t, --timeout timeout in seconds
    -h, --help prints help information
    -o, --output output location for scan file
`
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

	flag.Usage = func() { fmt.Print(usage) }
	flag.Parse()

	expandAddress(addressRange)
	portList = expandPort(portRange)
}