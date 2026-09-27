// Bitcoin node that connects to multiple peers, performs handshakes, and handles ping/pong messages.
//
// Usage:
// go run . ip1:port1 ip2:port2
package main

import (
	"fmt"
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

	fmt.Printf("[%s] Connecting...\n", server)

	conn, err := net.DialTimeout(
		"tcp",
		server,
		5*time.Second,
	)
	if err != nil {
		fmt.Printf("[%s] Connection error: %s\n", server, err)
		return
	}

	p := Peer{
		Address:   server,
		Conn:      conn,
		Direction: Outbound,
	}

	p.run()
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

func buildPeerAddresses(args []string) []PeerAddress {
	pas := make([]PeerAddress, 0, len(args))

	for _, arg := range args {
		ip, port, err := net.SplitHostPort(arg)
		if err != nil {
			fmt.Printf("Invalid peer address %q: %s\n", arg, err)
			continue
		}

		pas = append(pas, PeerAddress{
			IP:   ip,
			Port: port,
		})
	}
	return pas
}

func listenForPeers(port int) {
	fmt.Printf("Listening on port %d...\n", port)

	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		fmt.Printf(
			"[%s] Error: %s\n",
			time.Now().Format("2006-01-02 15:04:05.000"),
			err,
		)
		return
	}
	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Printf(
				"[%s] Accept error: %s\n",
				time.Now().Format("2006-01-02 15:04:05.000"),
				err,
			)
			continue
		}

		go acceptPeer(conn)
	}
}

func acceptPeer(conn net.Conn) {
	remote := conn.RemoteAddr()
	fmt.Println("Client connected:", remote)

	p := Peer{
		Address:   remote.String(),
		Conn:      conn,
		Direction: Inbound,
	}

	p.run()
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

	go listenForPeers(8333)

	pas := buildPeerAddresses(args)
	if len(pas) > 0 {
		connectToPeers(pas)
	}
}
