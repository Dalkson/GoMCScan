package main

import (
	"strconv"
	"strings"
	"unicode"
)

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