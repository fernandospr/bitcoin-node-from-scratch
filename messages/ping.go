package main

import "fmt"

type PingPayload struct {
	Nonce uint64
}

func (v PingPayload) String() string {
	return fmt.Sprintf(
		"PingPayload{Nonce: %d}",
		v.Nonce,
	)
}

func NewPingPayload(
	nonce uint64,
) (PingPayload, error) {
	return PingPayload{
		Nonce: nonce,
	}, nil
}

func ReadPingPayload(r *ByteReader) (PingPayload, error) {
	nonce, err := r.ReadUInt64()
	if err != nil {
		return PingPayload{}, err
	}
	return PingPayload{
		Nonce: nonce,
	}, nil
}

func (p PingPayload) Serialize(w *ByteWriter) {
	w.WriteUInt64(p.Nonce)
}
