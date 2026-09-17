// Queries the Bitcoin DNS seeds, obtains node IP addresses, and connects to each IP via TCP on port 8333
// Saves the connection results to peers.json
//
// Examples:
// go run tcp.go							Queries the default mainnet seeds
// go run tcp.go seed.bitcoin.sipa.be		Queries the specifc seed
// go run tcp.go x9.seed.bitcoin.sipa.be	Filters by service bits
package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"os"
	"time"
)

const PORT = "8333"
const FILE = "peers.json"

type Peer struct {
	Address   string    `json:"address"`
	Port      string    `json:"port"`
	Responded bool      `json:"responded"`
	Failures  int       `json:"failures"`
	LastTried time.Time `json:"last_tried"`
}

type Peers map[string]Peer

type ConnectionStatus int

const (
	ConnectionConnected ConnectionStatus = 0
	ConnectionTimeout   ConnectionStatus = 1
	ConnectionRejected  ConnectionStatus = 2
)

func tcpConnect(address string, port string) ConnectionStatus {
	timeout := 5 * time.Second
	server := net.JoinHostPort(address, port)
	
	fmt.Printf("Connecting to %s... ", server)
	conn, err := net.DialTimeout("tcp", server, timeout)
	if err != nil {
		var netErr net.Error
		if errors.As(err, &netErr) && netErr.Timeout() {
			fmt.Println("❌")
			return ConnectionTimeout
		}

		fmt.Println("❌")
		return ConnectionRejected
	}
	defer conn.Close()

	fmt.Println("✅")

	return ConnectionConnected
}

func resolveSeed(seed string) []string {
	result := []string{}

	ips, err := net.LookupIP(seed)
	if err == nil {
		for _, ip := range ips {
			result = append(result, ip.String())
		}
	}

    return result
}

func unique(strings []string) []string {
    seen := make(map[string]bool)
    result := []string{}

    for _, s := range strings {
        if !seen[s] {
            seen[s] = true
            result = append(result, s)
        }
    }

    return result
}

func resolveSeeds(seeds []string) []string {
	ips := []string{}
	for _, seed := range seeds {
		fmt.Printf("Resolving seed %s...\n", seed)
		ipsForSeed := resolveSeed(seed)
		ips = append(ips, ipsForSeed...)
		fmt.Printf("IPs for seed %s: %s\n", seed, ips)
	}

	return unique(ips)
}

func loadPeers(file string) Peers {
	peers := make(Peers)

	fmt.Printf("Reading %s...\n", file)
	data, err := os.ReadFile(file)
	if err != nil {
		fmt.Printf("File %s not found\n", file)
		return peers
	}

	fmt.Println("Unmarshalling...")
	if err := json.Unmarshal(data, &peers); err != nil {
		fmt.Println("Error unmarshalling")
		return make(Peers)
	}

	fmt.Println("Done")

	return peers
}

func savePeers(file string, peers Peers) {
	fmt.Println("Marshalling...")
	data, err := json.MarshalIndent(peers, "", "  ")
	if err != nil {
		fmt.Println("Error marshalling")
		return
	}

	fmt.Printf("Writing %s...\n", file)
	err = os.WriteFile(file, data, 0644)
	if err != nil {
		fmt.Println("Error writing")
		return
	}

	fmt.Println("Done")
}

func updatePeersWithIps(peers Peers, ips []string) Peers {
	for _, address := range ips {
		outcome := tcpConnect(address, PORT)

		peer, exists := peers[address]
		if !exists {
			peer = Peer{
				Address:   address,
				Port:      PORT,
				Responded: false,
				Failures:  0,
				LastTried: time.Time{},
			}
		}

		peer.Responded = peer.Responded || outcome == ConnectionConnected
		if outcome == ConnectionConnected {
			peer.Failures = 0
		} else {
			peer.Failures++
		}
		peer.LastTried = time.Now()

		peers[address] = peer
		fmt.Printf("Peer %s updated\n", address)
	}

	return peers
}

func main() {
	// Seeds from https://github.com/bitcoin/bitcoin/blob/v31.1/src/kernel/chainparams.cpp
	seeds := []string{
		"seed.bitcoin.sipa.be",          // Pieter Wuille
		"dnsseed.bluematt.me",           // Matt Corallo
		"seed.bitcoin.jonasschnelli.ch", // Jonas Schnelli
		"seed.btc.petertodd.net",        // Peter Todd
		"seed.bitcoin.sprovoost.nl",     // Sjors Provoost
		"dnsseed.emzy.de",               // Stephan Oeste
		"seed.bitcoin.wiz.biz",          // Jason Maurice
		"seed.mainnet.achownodes.xyz",   // Ava Chow
	}

	if len(os.Args) > 1 {
		seeds = os.Args[1:]
	}

	ips := resolveSeeds(seeds)

	peers := loadPeers(FILE)
	peers = updatePeersWithIps(peers, ips)
	savePeers(FILE, peers)
}