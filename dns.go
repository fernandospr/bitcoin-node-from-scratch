// Queries the Bitcoin DNS seeds
// This program queries the seeds and retrieves node IP addresses
//
// Examples:
// go run dns.go							Queries the default mainnet seeds
// go run dns.go seed.bitcoin.sipa.be		Queries the specifc seed
// go run dns.go x9.seed.bitcoin.sipa.be	Filters by service bits
package main

import (
	"fmt"
	"net"
	"os"
)

func main() {
	// Seeds from https://github.com/bitcoin/bitcoin/blob/v31.1/src/kernel/chainparams.cpp
	hosts := []string{
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
		hosts = os.Args[1:]
	}

	uniqueIpSet := map[string]struct{}{}
	ipQuantity := 0

	for _, host := range hosts {
		fmt.Println(host)

		ips, err := net.LookupIP(host)
		if err != nil {
			fmt.Println(err)
			continue
		}

		for _, ip := range ips {
			if ip.To4() != nil {
				fmt.Println(ip, "→ IPv4")
			} else if ip.To16() != nil {
				fmt.Println(ip, "→ IPv6")
			} else {
				continue
			}
			ipQuantity++
			uniqueIpSet[ip.String()] = struct{}{}
		}
		fmt.Println()
	}

	fmt.Println("IPs quantity:", ipQuantity)
	fmt.Println("Unique IPs:", len(uniqueIpSet))
}
