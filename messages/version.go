package main

import (
	"crypto/rand"
	"encoding/binary"
	"fmt"
	"time"
)

type VersionPayload struct {
	Version   int32
	Services  uint64
	Timestamp int64
	AddrRecv  NetworkAddress
	AddrFrom  NetworkAddress
	Nonce     uint64
	UserAgent string
	Height    int32
	Relay     bool
}

func (v VersionPayload) String() string {
	return fmt.Sprintf(
		"VersionPayload{Version: %d, Services: %d, Timestamp: %d, AddrRecv: %v, AddrFrom: %v, Nonce: %d, UserAgent: %q, Height: %d, Relay: %t}",
		v.Version,
		v.Services,
		v.Timestamp,
		v.AddrRecv,
		v.AddrFrom,
		v.Nonce,
		v.UserAgent,
		v.Height,
		v.Relay,
	)
}

func NewVersionPayload(
	version int32,
	services uint64,
	addrRecv NetworkAddress,
	addrFrom NetworkAddress,
	userAgent string,
	height int32,
	relay bool,
) (VersionPayload, error) {
	nonce, err := randomNonce()
	if err != nil {
		return VersionPayload{}, fmt.Errorf("generate version nonce: %w", err)
	}

	return VersionPayload{
		Version:   version,
		Services:  services,
		Timestamp: time.Now().Unix(),
		AddrRecv:  addrRecv,
		AddrFrom:  addrFrom,
		Nonce:     nonce,
		UserAgent: userAgent,
		Height:    height,
		Relay:     relay,
	}, nil
}

func ReadVersionPayload(r *ByteReader) (VersionPayload, error) {
	version, err := r.ReadInt32()
	if err != nil {
		return VersionPayload{}, err
	}
	services, err := r.ReadUInt64()
	if err != nil {
		return VersionPayload{}, err
	}
	timestamp, err := r.ReadInt64()
	if err != nil {
		return VersionPayload{}, err
	}
	addrRecv, err := ReadNetworkAddress(r)
	if err != nil {
		return VersionPayload{}, err
	}
	addrFrom, err := ReadNetworkAddress(r)
	if err != nil {
		return VersionPayload{}, err
	}
	nonce, err := r.ReadUInt64()
	if err != nil {
		return VersionPayload{}, err
	}
	userAgent, err := r.ReadVarSizeString()
	if err != nil {
		return VersionPayload{}, err
	}
	height, err := r.ReadInt32()
	if err != nil {
		return VersionPayload{}, err
	}
	relay, err := r.ReadBool()
	if err != nil {
		return VersionPayload{}, err
	}
	return VersionPayload{
		Version:   version,
		Services:  services,
		Timestamp: timestamp,
		AddrRecv:  addrRecv,
		AddrFrom:  addrFrom,
		Nonce:     nonce,
		UserAgent: userAgent,
		Height:    height,
		Relay:     relay,
	}, nil
}

func (p VersionPayload) Serialize(w *ByteWriter) {
	w.WriteInt32(p.Version)
	w.WriteUInt64(p.Services)
	w.WriteInt64(p.Timestamp)
	p.AddrRecv.Serialize(w)
	p.AddrFrom.Serialize(w)
	w.WriteUInt64(p.Nonce)
	w.WriteVarSizeString(p.UserAgent)
	w.WriteInt32(p.Height)
	w.WriteBool(p.Relay)
}

func randomNonce() (uint64, error) {
	var nonce uint64

	if err := binary.Read(rand.Reader, binary.LittleEndian, &nonce); err != nil {
		return 0, err
	}

	return nonce, nil
}
