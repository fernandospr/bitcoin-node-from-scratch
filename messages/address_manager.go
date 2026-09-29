package main

import (
	"encoding/json"
	"fmt"
	"math/rand/v2"
	"net"
	"os"
	"sync"
)

type AddressManager struct {
	mu        sync.RWMutex
	addresses map[string]PeerAddress
}

func (am *AddressManager) RandomSample(max int) []PeerAddress {
	all := am.All()

	if len(all) <= max {
		return all
	}

	rand.Shuffle(len(all), func(i, j int) {
		all[i], all[j] = all[j], all[i]
	})

	return all[:max]
}

func NewAddresssManager() *AddressManager {
	return &AddressManager{
		addresses: make(map[string]PeerAddress),
	}
}

func (am *AddressManager) Add(address PeerAddress) {
	key := net.JoinHostPort(address.IP, address.Port)

	am.mu.Lock()
	defer am.mu.Unlock()

	am.addresses[key] = address
}

func (am *AddressManager) AddMany(addresses []PeerAddress) {
	for _, a := range addresses {
		am.Add(a)
	}
}

func (am *AddressManager) All() []PeerAddress {
	am.mu.RLock()
	defer am.mu.RUnlock()

	result := make([]PeerAddress, 0, len(am.addresses))
	for _, address := range am.addresses {
		result = append(result, address)
	}
	return result
}

func (am *AddressManager) Save(path string) (int, error) {
	addresses := am.All()

	data, err := json.MarshalIndent(addresses, "", "  ")
	if err != nil {
		return 0, fmt.Errorf("marshal addresses: %w", err)
	}

	if err := os.WriteFile(path, data, 0644); err != nil {
		return 0, fmt.Errorf("write addresses file: %w", err)
	}

	return len(addresses), nil
}

func (am *AddressManager) Load(path string) (int, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return 0, fmt.Errorf("read addresses file: %w", err)
	}

	var addresses []PeerAddress

	if err := json.Unmarshal(data, &addresses); err != nil {
		return 0, fmt.Errorf("unmarshal addresses: %w", err)
	}

	am.AddMany(addresses)

	return len(addresses), nil
}
