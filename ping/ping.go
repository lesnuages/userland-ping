package ping

import (
	"context"
	"fmt"
	"net"
	"strings"
	"sync"
	"time"
)

var commonPorts = []uint16{
	21, 22, 23, 25, 53, 80, 110, 111, 135, 139, 143, 389, 443, 445, 993, 995, 1723, 3306, 3389, 5900, 8080,
}

type Pinger struct {
	// Ports is the list of ports to check on the host.
	Ports []uint16
	// Timeout is the timeout for each connection attempt.
	Timeout time.Duration
}

// isUp checks if a host is up by trying to establish a TCP connection
// to the given IP and port.
func (p *Pinger) isUp(ctx context.Context, wg *sync.WaitGroup, ip string, port uint16) bool {
	defer wg.Done()
	d := net.Dialer{
		Timeout: p.Timeout,
	}
	// DialContext will return an error on timeout or if we call the cancel function
	// triggering the wg.Done call.
	if c, err := d.DialContext(ctx, "tcp", fmt.Sprintf("%s:%d", ip, port)); err == nil {
		// Close the establish connection to avoid
		// dangling file descriptors
		c.Close()
		return true
	} else {
		// Check if we got a RST packet, meaning someone responded,
		// which means it's highly likely the host is up.
		return strings.Contains(err.Error(), "connection refused")
	}
}

// CheckHost checks if a host is up by trying to establish a TCP connection on
// each of the Ports. Any TCP response (SYN/ACK or RST) is interpreted as
// a positive result. If no response is received, the host is considered down.
// The default timeout for each connection is 5 seconds.
func (p *Pinger) CheckHost(ip string) bool {
	var hostUP bool = false

	if p.Ports == nil {
		p.Ports = commonPorts
	}

	if p.Timeout == 0 {
		p.Timeout = 5 * time.Second
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	wg := sync.WaitGroup{}

	for _, port := range p.Ports {
		wg.Add(1)
		go func(pt uint16) {
			if p.isUp(ctx, &wg, ip, pt) {
				hostUP = true
				// Cancel the context to stop the other goroutines
				cancel()
			}
		}(port)
	}
	wg.Wait()
	return hostUP
}
