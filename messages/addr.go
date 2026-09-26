package main

import (
	"fmt"
	"strings"
)

const maxAddrEntries = 1000

type AddrPayload struct {
	Addresses []AddrItemPayload
}

type AddrItemPayload struct {
	Time    uint32
	Address NetworkAddress
}

func (a AddrPayload) String() string {
	const maxItems = 10

	var parts []string

	for i, item := range a.Addresses {
		if i >= maxItems {
			break
		}

		parts = append(parts, item.String())
	}

	if len(a.Addresses) > maxItems {
		return fmt.Sprintf(
			"AddrPayload{Addresses: [%s, ...] (total: %d)}",
			strings.Join(parts, ", "),
			len(a.Addresses),
		)
	}

	return fmt.Sprintf(
		"AddrPayload{Addresses: [%s] (total: %d)}",
		strings.Join(parts, ", "),
		len(a.Addresses),
	)
}

func (i AddrItemPayload) String() string {
	return fmt.Sprintf("AddrItemPayload{Time: %d, Address: %s}", i.Time, i.Address)
}

func NewAddrPayload(
	items []AddrItemPayload,
) (AddrPayload, error) {
	return AddrPayload{
		Addresses: items,
	}, nil
}

func ReadAddrPayload(r *ByteReader) (AddrPayload, error) {
	addrQuantity, err := r.ReadCompactSize()
	if addrQuantity > maxAddrEntries {
		return AddrPayload{}, fmt.Errorf(
			"too many addresses: %d (max %d)",
			addrQuantity,
			maxAddrEntries,
		)
	}
	if err != nil {
		return AddrPayload{}, err
	}
	addresses := make([]AddrItemPayload, 0, addrQuantity)
	for range addrQuantity {
		time, err := r.ReadUInt32()
		if err != nil {
			return AddrPayload{}, err
		}
		addr, err := ReadNetworkAddress(r)
		if err != nil {
			return AddrPayload{}, err
		}
		item := AddrItemPayload{
			Time:    time,
			Address: addr,
		}
		addresses = append(addresses, item)
	}

	return AddrPayload{
		Addresses: addresses,
	}, nil
}

func (p AddrPayload) Serialize(w *ByteWriter) {
	addrQuantity := len(p.Addresses)
	w.WriteCompactSize(uint64(addrQuantity))

	for _, a := range p.Addresses {
		w.WriteUInt32(a.Time)
		a.Address.Serialize(w)
	}
}
