package main

import (
	"fmt"
	"net"
)

func uniquePeerAddresses(addresses []PeerAddress) []PeerAddress {
	seen := make(map[string]struct{})
	result := make([]PeerAddress, 0, len(addresses))

	for _, address := range addresses {
		key := net.JoinHostPort(
			address.IP,
			address.Port,
		)

		if _, exists := seen[key]; exists {
			continue
		}
		seen[key] = struct{}{}
		result = append(result, address)
	}

	return result
}

func resolveSeed(seed string) ([]PeerAddress, error) {
	ips, err := net.LookupIP(seed)
	if err != nil {
		return nil, fmt.Errorf(
			"resolve DNS seed %s: %w",
			seed,
			err,
		)
	}
	addresses := make([]PeerAddress, 0, len(ips))
	for _, ip := range ips {
		address := PeerAddress{
			IP:   ip.String(),
			Port: MainnetPort,
		}
		addresses = append(addresses, address)
	}
	return addresses, nil
}

func ResolveSeeds(seeds []string) []PeerAddress {
	var addresses []PeerAddress

	for _, seed := range seeds {
		fmt.Printf("Resolving seed %s...\n", seed)
		addressesForSeed, err := resolveSeed(seed)
		if err != nil {
			fmt.Printf("Seed %s failed: %s\n", seed, err)
			continue
		}
		addresses = append(addresses, addressesForSeed...)
	}

	return uniquePeerAddresses(addresses)
}
