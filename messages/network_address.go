package main

import (
	"bytes"
	"fmt"
	"net"
	"net/netip"
)

type NetworkAddress struct {
	Services uint64
	IP       [16]byte
	Port     uint16
}

var ipv4Prefix = [12]byte{
	0, 0, 0, 0,
	0, 0, 0, 0,
	0, 0, 0xff, 0xff,
}

func NewNetworkAddress(services uint64, ip string, port uint16) (NetworkAddress, error) {
	checkedIp, err := netip.ParseAddr(ip)
	if err != nil {
		return NetworkAddress{}, fmt.Errorf("IPv4/IPv6 expected but got: %q", ip)
	}

	var ipBytes [16]byte
	if checkedIp.Is4() {
		ipv4 := checkedIp.As4()
		copy(ipBytes[:], ipv4Prefix[:])
		copy(ipBytes[12:], ipv4[:])
	} else {
		ipv6 := checkedIp.As16()
		copy(ipBytes[:], ipv6[:])
	}

	return NetworkAddress{
		Services: services,
		IP:       ipBytes,
		Port:     port,
	}, nil
}

func ReadNetworkAddress(r *ByteReader) (NetworkAddress, error) {
	services, err := r.ReadUInt64()
	if err != nil {
		return NetworkAddress{}, err
	}

	ip, err := r.ReadBytes(16)
	if err != nil {
		return NetworkAddress{}, err
	}

	port, err := r.ReadUInt16BE()
	if err != nil {
		return NetworkAddress{}, err
	}

	return NetworkAddress{
		Services: services,
		IP:       [16]byte(ip),
		Port:     port,
	}, nil
}

func (a NetworkAddress) IPString() string {
	if bytes.Equal(a.IP[0:12], ipv4Prefix[:]) {
		return net.IP(a.IP[12:16]).To4().String()
	}
	return net.IP(a.IP[:]).To16().String()
}

func (a NetworkAddress) String() string {
	return fmt.Sprintf(
		"NetworkAddress{Services: %d, IP:%s, Port:%d}",
		a.Services,
		a.IPString(),
		a.Port,
	)
}

func (a NetworkAddress) Serialize(w *ByteWriter) {
	w.WriteUInt64(a.Services)
	w.WriteBytes(a.IP[:])
	w.WriteUInt16BE(a.Port)
}
