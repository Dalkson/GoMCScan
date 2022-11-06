package mcping

import (
	"GoMCScan/mcping/types"
	"net"
	"strconv"
	"time"
)

// Ping and get information from an host and port with a custom timeout
// Return a pointer to types.PingResponse or an error
//
// # Error is thrown when the host is unreachable or the data received are incorrect
//
func PingWithTimeout(host string, port uint16, timeout time.Duration) (*types.PingResponse, string, error) {
	addr := host + ":" + strconv.Itoa(int(port))

	// Open connection to server
	conn, err := net.DialTimeout("tcp", addr, timeout)
	if err != nil {
		return nil, "", err
	}
	defer conn.Close()

	// Send packet to server
	sendPacket(host, port, &conn)
	response, err := readResponse(&conn)
	if err != nil {
		return nil, "", err
	}

	// Decide response
	decoded := decodeResponse(response)
	return decoded, response, nil
}
