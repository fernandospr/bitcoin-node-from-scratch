// Bitcoin node that connects to multiple peers, performs handshakes, and handles ping/pong messages.
//
// go run . -help
package main

import (
	"flag"
	"fmt"
	"net"
)

func buildOutboundPeerAddresses(connectNodes []string) []PeerAddress {
	peers := make([]PeerAddress, 0, len(connectNodes))
	for _, address := range connectNodes {
		host, port, err := net.SplitHostPort(address)
		if err != nil {
			fmt.Printf("Invalid peer address %q: %s\n", address, err)
			continue
		}

		peers = append(
			peers,
			PeerAddress{
				IP:   host,
				Port: port,
			},
		)
	}
	return peers
}

func printUsage() {
	fmt.Println("Bitcoin node")
	fmt.Println()
	fmt.Println("Usage:")
	fmt.Println("  go run . [options]")
	fmt.Println()
	fmt.Println("Options:")
	flag.PrintDefaults()
}

func main() {
	flag.Usage = printUsage

	// -port
	port := flag.String(
		"port",
		MainnetPort,
		"TCP port to listen on",
	)

	// -connect
	var connectNodes []string
	flag.Func(
		"connect",
		"Connect to a specific outbound peer (can be specified multiple times)",
		func(value string) error {
			connectNodes = append(connectNodes, value)
			return nil
		},
	)

	// -max-outbound-peers
	maxOutboundPeers := flag.Int(
		"max-outbound-peers",
		8,
		"Maximum number of outbound peers",
	)

	flag.Parse()

	config := NodeConfig{
		ListenPort:            *port,
		ExplicitOutboundPeers: buildOutboundPeerAddresses(connectNodes),
		MaxOutboundPeers:      *maxOutboundPeers,
	}
	n := NewNode(config)
	n.Bootstrap()
	n.Run()
}

// TODOs
//
// Hacer un diagrama ascii del funcionamiento/arquitectura general
// Ejemplos de uso de flags

// go run . -port 8334
// go run . -connect 127.0.0.1:8334 -max-outbound-peers 1
