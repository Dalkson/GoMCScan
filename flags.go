package main

import (
	"flag"
	"fmt"
)

func getFlags() {
	// sets help page for MCScan
	const usage = `Usage of MCScan:
     MCScan (-a targets) [-p PortRange] [-T Threads] [-t Timeout] [-o output]
 Options:
     -h, --help prints help information
     -a, --targets IP address range to scan
     -p, --ports port range to scan (25565)
     -T, --threads number of threads to use (1000)
     -t, --timeout timeout in seconds (3)
     -o, --output output location for scan file (out/scan.log)
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
