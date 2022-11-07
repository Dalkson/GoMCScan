package main

import (
	"fmt"
	"strings"
)

func totalToSend() int {
	var total int
	totalAddresses := 0
	for _, address := range conf.addressList {
		if !strings.Contains(address, "/") { //puts single addresses in CIDR notation
			address = fmt.Sprintf("%v/32", address)
		}
		addressSuffix := strings.Split(address, "/")
		switch addressSuffix[1] { // counts total addresses based off CIDR suffix, doesnt check for repeats so nested IP ranges will be counted twice.
		case "0":
			totalAddresses = 4294967296
			total = len(conf.portList) * totalAddresses
			return total
		case "1":
			totalAddresses += 2147483648
		case "2":
			totalAddresses += 1073741824
		case "3":
			totalAddresses += 536870912
		case "4":
			totalAddresses += 268435456
		case "5":
			totalAddresses += 134217728
		case "6":
			totalAddresses += 67108864
		case "7":
			totalAddresses += 33554432
		case "8":
			totalAddresses += 16777216
		case "9":
			totalAddresses += 8388608
		case "10":
			totalAddresses += 4194304
		case "11":
			totalAddresses += 2097152
		case "12":
			totalAddresses += 1048576
		case "13":
			totalAddresses += 524288
		case "14":
			totalAddresses += 262144
		case "15":
			totalAddresses += 131072
		case "16":
			totalAddresses += 65536
		case "17":
			totalAddresses += 32768
		case "18":
			totalAddresses += 16384
		case "19":
			totalAddresses += 8192
		case "20":
			totalAddresses += 4096
		case "21":
			totalAddresses += 2048
		case "22":
			totalAddresses += 1024
		case "23":
			totalAddresses += 512
		case "24":
			totalAddresses += 256
		case "25":
			totalAddresses += 128
		case "26":
			totalAddresses += 64
		case "27":
			totalAddresses += 32
		case "28":
			totalAddresses += 16
		case "29":
			totalAddresses += 8
		case "30":
			totalAddresses += 2
		case "31":
			totalAddresses += 2
		case "32":
			totalAddresses++
		}
	}

	total = len(conf.portList) * totalAddresses
	return total
}
