package main

import (
	"encoding/hex"
	"fmt"
)

type Header struct {
	Magic    uint32
	Command  string
	Length   uint32
	Checksum [4]byte
}

func NewHeader(magic uint32, command string, payload []byte) Header {
	return Header{
		Magic:    magic,
		Command:  command,
		Length:   uint32(len(payload)),
		Checksum: checksum(payload),
	}
}

func ReadHeader(r *ByteReader) (Header, error) {
	magic, err := r.ReadUInt32()
	if err != nil {
		return Header{}, err
	}
	cmd, err := r.ReadFixedSizeString(12)
	if err != nil {
		return Header{}, err
	}
	length, err := r.ReadUInt32()
	if err != nil {
		return Header{}, err
	}
	checksum, err := r.ReadBytes(4)
	if err != nil {
		return Header{}, err
	}
	return Header{
		Magic:    magic,
		Command:  cmd,
		Length:   length,
		Checksum: [4]byte(checksum),
	}, nil
}

func (h Header) Serialize(w *ByteWriter) {
	w.WriteUInt32(h.Magic)
	w.WriteFixedSizeString(h.Command, 12)
	w.WriteUInt32(h.Length)
	w.WriteBytes(h.Checksum[:])
}

func (h Header) ChecksumString() string {
	return "0x" + hex.EncodeToString(h.Checksum[:])
}

func (h Header) String() string {
	return fmt.Sprintf(
		"Header{Magic: 0x%08x, Command: %s, Length: %d, Checksum: %s}",
		h.Magic,
		h.Command,
		h.Length,
		h.ChecksumString(),
	)
}

func checksum(data []byte) [4]byte {
	hash := DoubleSha256(data)

	var result [4]byte
	copy(result[:], hash[:4])

	return result
}
