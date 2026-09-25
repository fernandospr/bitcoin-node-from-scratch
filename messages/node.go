// Bitcoin node that connects to multiple peers, performs handshakes, and handles ping/pong messages.
//
// Usage:
// go run . ip1:port1 ip2:port2
package main

import (
	"errors"
	"fmt"
	"io"
	"net"
	"os"
	"sync"
	"time"
)

type PeerAddress struct {
	IP   string
	Port string
}

func connectToPeer(pa PeerAddress) {
	server := net.JoinHostPort(pa.IP, pa.Port)

	fmt.Printf("[%s] Connecting... ", server)

	conn, err := net.DialTimeout(
		"tcp",
		server,
		5*time.Second,
	)
	if err != nil {
		fmt.Println("❌")
		fmt.Printf("[%s] Connection error: %s\n", server, err)
		return
	}

	fmt.Println("✅")

	defer conn.Close()

	peer := Peer{
		Server: server,
		Conn:   conn,
	}

	peer.log("Starting handshake...")

	if err := peer.handshake(); err != nil {
		peer.log(
			"Handshake failed: %s",
			err,
		)
		return
	}

	peer.log("Handshake success! ✅")

	peer.log("Starting message loop...")

	if err := peer.messageLoop(); err != nil {
		if errors.Is(err, io.EOF) {
			peer.log("Peer closed the connection")
		} else {
			peer.log(
				"Message loop ended: %s",
				err,
			)
		}

		return
	}
}

func connectToPeers(pas []PeerAddress) {
	var wg sync.WaitGroup

	for _, pa := range pas {
		wg.Add(1)

		go func(pa PeerAddress) {
			defer wg.Done()

			connectToPeer(pa)
		}(pa)
	}

	wg.Wait()
}

func printUsage() {
	fmt.Println("Usage:")
	fmt.Println("  go run . <ip:port> [ip:port...]")
	fmt.Println()
	fmt.Println("Example:")
	fmt.Println("  go run . 108.36.121.109:8333 176.126.71.51:8333")
}

func main() {
	args := os.Args[1:]

	if len(args) == 0 {
		printUsage()
		os.Exit(1)
	}
	pas := make([]PeerAddress, 0, len(args))

	for _, arg := range args {
		ip, port, err := net.SplitHostPort(arg)
		if err != nil {
			fmt.Printf("Invalid peer address %q: %s\n", arg, err)
			os.Exit(1)
		}

		pas = append(pas, PeerAddress{
			IP:   ip,
			Port: port,
		})
	}

	connectToPeers(pas)
}
