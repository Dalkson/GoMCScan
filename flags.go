package main

import (
	"flag"
	"fmt"
	"os"
)

func getFlags() {
	// sets help page for MCScan
	const usage = `Usage of MCScan:
     MCScan (-a targets) [-p PortRange] [-T Threads] [-t Timeout] [-o output]
 Options:
     -h, --help prints help information
     -a, --targets IP address range to scan. Ex: 116.202.0.0/16,116.203.0.0/16
     -p, --ports port range to scan. (25565) Ex: 1337,25565-25570
     -T, --threads number of threads to use. (1000)
     -t, --timeout timeout in seconds. (3)
	 -I, --input location for target file. (targets.txt)
     -o, --output output location for scan file. (out/scan.log)
	 `
	flag.IntVar(&threads, "T", 1000, "number of threads to use")
	flag.IntVar(&threads, "threads", 1000, "number of threads to use")
	flag.IntVar(&timeout, "t", 1, "timeout in seconds")
	flag.IntVar(&timeout, "timeout", 1, "timeout in seconds")
	var addressRange string
	flag.StringVar(&addressRange, "a", "", "IP address range to scan")
	flag.StringVar(&addressRange, "targets", "", "IP address range to scan")
	flag.StringVar(&outputDir, "output", "out/scan.log", "output location for scan file")
	flag.StringVar(&outputDir, "o", "out/scan.log", "output location for scan file")
	var inputDir string
	flag.StringVar(&inputDir, "input", "targets.txt", "input location for target file")
	flag.StringVar(&inputDir, "i", "targets.txt", "input location for target file")
	var portRange string
	flag.StringVar(&portRange, "p", "25565", "port range to scan")
	flag.StringVar(&portRange, "port", "25565", "port range to scan")

	flag.Usage = func() { fmt.Print(usage) }
	flag.Parse()

	// if input fill is passed, set addressRange to the contents of the file.
	if isFlagPassed("a") || isFlagPassed("targets") {
		expandAddress(addressRange)
	} else if _, err := os.Stat(inputDir); err == nil {
		b, err := os.ReadFile(inputDir)
		if err != nil {
			fmt.Print(err)
		}
		str := string(b)
		if str != "" {
			addressRange = str
			expandAddress(addressRange)
		} else {
			handleError("Input file is empty.")
		}
	} else {
		handleError("targets have not been set and no input file found.")
	}

	portList = expandPort(portRange)
}

func isFlagPassed(name string) bool {
	found := false
	flag.Visit(func(f *flag.Flag) {
		if f.Name == name {
			found = true
		}
	})
	return found
}
