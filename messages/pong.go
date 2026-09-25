package main

import "fmt"

type PongPayload struct {
	Nonce uint64
}

func (v PongPayload) String() string {
	return fmt.Sprintf(
		"PongPayload{Nonce: %d}",
		v.Nonce,
	)
}

func NewPongPayload(
	nonce uint64,
) (PongPayload, error) {
	return PongPayload{
		Nonce: nonce,
	}, nil
}

func ReadPongPayload(r *ByteReader) (PongPayload, error) {
	nonce, err := r.ReadUInt64()
	if err != nil {
		return PongPayload{}, err
	}
	return PongPayload{
		Nonce: nonce,
	}, nil
}

func (p PongPayload) Serialize(w *ByteWriter) {
	w.WriteUInt64(p.Nonce)
}
