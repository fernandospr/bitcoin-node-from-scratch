package main

import (
	"fmt"
	"math/rand/v2"
	"net"
	"time"
)

const addressSaveInterval = 1 * time.Minute
const addressDatabasePath = "addresses.json"

type Node struct {
	addressManager *AddressManager

	explicitOutboundPeers []PeerAddress
	listenPort            string
	maxOutboundPeers      int
}

func (n *Node) AddExplicitOutboundPeer(address string) error {
	host, port, err := net.SplitHostPort(address)
	if err != nil {
		return fmt.Errorf("invalid peer address: %w", err)
	}

	n.explicitOutboundPeers = append(
		n.explicitOutboundPeers,
		PeerAddress{
			IP:   host,
			Port: port,
		},
	)

	return nil
}

func NewNode(config NodeConfig) *Node {
	return &Node{
		addressManager:        NewAddresssManager(),
		listenPort:            config.ListenPort,
		explicitOutboundPeers: config.ExplicitOutboundPeers,
		maxOutboundPeers:      config.MaxOutboundPeers,
	}
}

func (n *Node) Run() {
	go n.listenForPeers(n.listenPort)
	go n.connectToPeers()
	go n.persistAddresses()

	select {}
}

// Bootstrap

func (n *Node) Bootstrap() {
	fmt.Println("Bootstrapping node...")

	loaded, err := n.addressManager.Load(addressDatabasePath)
	if err != nil {
		fmt.Printf(
			"Failed to load address database: %s\n",
			err,
		)
	} else {
		fmt.Printf(
			"Loaded %d known addresses\n",
			loaded,
		)

		return
	}

	fmt.Println("Address database not found")
	fmt.Println("Resolving DNS seeds...")

	addresses := ResolveSeeds(MainnetSeeds)
	n.addressManager.AddMany(addresses)

	fmt.Printf("Discovered %d addresses from DNS seeds\n",
		len(addresses),
	)

	saved, err := n.addressManager.Save(addressDatabasePath)
	if err != nil {
		fmt.Printf(
			"Failed to save address database: %s\n",
			err,
		)
		return
	}

	fmt.Printf(
		"Saved %d addresses to %s\n",
		saved,
		addressDatabasePath,
	)
}

// Inbound

func (n *Node) listenForPeers(port string) {
	fmt.Printf("Listening on port %s...\n", port)

	listener, err := net.Listen("tcp", fmt.Sprintf(":%s", port))
	if err != nil {
		fmt.Printf(
			"Listening error: %s\n",
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

		go n.acceptPeer(conn)
	}
}

func (n *Node) acceptPeer(conn net.Conn) {
	remote := conn.RemoteAddr()
	fmt.Println("Client connected:", remote)

	p := Peer{
		Address:        remote.String(),
		Conn:           conn,
		Direction:      Inbound,
		AddressManager: n.addressManager,
	}

	p.run()
}

// Outbound

func (n *Node) connectToPeers() {
	connected := 0

	// 1. Explicit peers first.
	for _, address := range n.explicitOutboundPeers {
		if connected >= n.maxOutboundPeers {
			break
		}

		if err := n.connectToPeer(address); err != nil {
			fmt.Printf(
				"[%s] Explicit connection failed: %s\n",
				net.JoinHostPort(address.IP, address.Port),
				err,
			)
			continue
		}

		connected++

		fmt.Printf(
			"Outbound peers: %d/%d\n",
			connected,
			n.maxOutboundPeers,
		)
	}

	// 2. Fill remaining slots with discovered addresses.
	if connected >= n.maxOutboundPeers {
		return
	}

	candidates := n.addressManager.All()

	rand.Shuffle(len(candidates), func(i, j int) {
		candidates[i], candidates[j] =
			candidates[j], candidates[i]
	})

	for _, address := range candidates {
		if connected >= n.maxOutboundPeers {
			break
		}

		if err := n.connectToPeer(address); err != nil {
			fmt.Printf(
				"[%s] Connection failed: %s\n",
				net.JoinHostPort(address.IP, address.Port),
				err,
			)
			continue
		}

		connected++

		fmt.Printf(
			"Outbound peers: %d/%d\n",
			connected,
			n.maxOutboundPeers,
		)
	}
}

func (n *Node) connectToPeer(pa PeerAddress) error {
	server := net.JoinHostPort(pa.IP, pa.Port)

	fmt.Printf("[%s] Connecting...\n", server)

	conn, err := net.DialTimeout(
		"tcp",
		server,
		5*time.Second,
	)
	if err != nil {
		return fmt.Errorf(
			"connection error: %w",
			err,
		)
	}

	fmt.Printf(
		"[%s] TCP connection established\n",
		server,
	)

	p := Peer{
		Address:        server,
		Conn:           conn,
		Direction:      Outbound,
		AddressManager: n.addressManager,
	}

	go p.run()

	return nil
}

// Address persistance

func (n *Node) persistAddresses() {
	ticker := time.NewTicker(addressSaveInterval)
	defer ticker.Stop()

	for range ticker.C {
		saved, err := n.addressManager.Save(addressDatabasePath)
		if err != nil {
			fmt.Printf(
				"Failed to save address database: %s\n",
				err,
			)
			continue
		}

		fmt.Printf(
			"Saved %d addresses to %s\n",
			saved,
			addressDatabasePath,
		)
	}
}
