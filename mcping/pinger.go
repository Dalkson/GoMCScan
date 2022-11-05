package mcping

import (
	"GoMCScan/api/types"
	"GoMCScan/dns"
	"net"
	"strconv"
	"time"
)

type pinger struct {
	DnsResolver types.DnsResolver
	Latency     types.Latency
}

// Create a new Minecraft Pinger
func NewPinger() *pinger {
	var resolver types.DnsResolver
	resolver = dns.NewResolver()
	return &pinger{DnsResolver: resolver}
}

// Create a new Minecraft Pinger with a custom DNS resolver
func NewPingerWithDnsResolver(dnsResolver types.DnsResolver) *pinger {
	return &pinger{DnsResolver: dnsResolver}
}

// Ping and get information from an host and port, default timeout: 3s
// Return a pointer to types.PingResponse or an error
//
// # Error is thrown when the host is unreachable or the data received are incorrect
//
// Example: pinger.Ping("play.hypixel.net", 25565)
func (p *pinger) Ping(host string, port uint16) (*types.PingResponse, string, error) {
	return p.PingWithTimeout(host, port, 3*time.Second)
}

// Ping and get information from an host and port with a custom timeout
// Return a pointer to types.PingResponse or an error
//
// # Error is thrown when the host is unreachable or the data received are incorrect
//
// Example: pinger.Ping("play.hypixel.net", 25565, 5 * time.Second)
func (p *pinger) PingWithTimeout(host string, port uint16, timeout time.Duration) (*types.PingResponse, string, error) {
	resolve, hostSRV, portSRV := p.DnsResolver.SRVResolve(host)
	if resolve {
		host = hostSRV
		port = portSRV
	}

	addr := host + ":" + strconv.Itoa(int(port))

	// Open connection to server
	conn, err := net.DialTimeout("tcp", addr, timeout)
	if err != nil {
		return nil, "", err
	}
	//延迟计算修改为与服务器连接过程时间,部分插件可能导致实际延迟过长
	defer conn.Close()

	sendPacket(host, port, &conn)
	response, err := readResponse(&conn)
	if err != nil {
		return nil, "", err
	}
	decoded := decodeResponse(response)
	return decoded, response, nil
}
