package mcping

import (
	"GoMCScan/api/types"
	"net"
	"strconv"
	"time"
)

// Create a new Minecraft Pinger

// Create a new Minecraft Pinger with a custom DNS resolver

// Ping and get information from an host and port, default timeout: 3s
// Return a pointer to types.PingResponse or an error
//
// # Error is thrown when the host is unreachable or the data received are incorrect
//
// Example: pinger.Ping("play.hypixel.net", 25565)
func Ping(host string, port uint16) (*types.PingResponse, string, error) {
	return PingWithTimeout(host, port, 3*time.Second)
}

// Ping and get information from an host and port with a custom timeout
// Return a pointer to types.PingResponse or an error
//
// # Error is thrown when the host is unreachable or the data received are incorrect
//
// Example: pinger.Ping("play.hypixel.net", 25565, 5 * time.Second)
func PingWithTimeout(host string, port uint16, timeout time.Duration) (*types.PingResponse, string, error) {
	addr := host + ":" + strconv.Itoa(int(port))

	// Open connection to server
	conn, err := net.DialTimeout("tcp", addr, timeout)
	if err != nil {
		return nil, "", err
	}
	defer conn.Close()

	sendPacket(host, port, &conn)
	response, err := readResponse(&conn)
	if err != nil {
		return nil, "", err
	}
	decoded := decodeResponse(response)
	return decoded, response, nil
}
