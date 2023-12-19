package main

import (
	"fmt"
	"os"

	"github.com/lesnuages/userland-ping/ping"
)

func main() {
	// Get IP from command line args
	if len(os.Args) != 2 {
		fmt.Printf("Usage: %s <IP>\n", os.Args[0])
		return
	}
	ipAddr := os.Args[1]

	p := ping.Pinger{}
	if p.CheckHost(ipAddr) {
		fmt.Printf("%s is up!\n", ipAddr)
	} else {
		fmt.Printf("%s is down!\n", ipAddr)
	}
}
