package main

import (
	"encoding/binary"
	"fmt"
)

type ByteWriter struct {
	data []byte
}

func NewByteWriter() *ByteWriter {
	return &ByteWriter{}
}

func (w *ByteWriter) Bytes() []byte {
	return w.data
}

func (w *ByteWriter) WriteBytes(value []byte) {
	w.data = append(w.data, value...)
}

func (w *ByteWriter) WriteBool(value bool) {
	if value {
		w.data = append(w.data, 0x01)
	} else {
		w.data = append(w.data, 0x00)
	}
}

func (w *ByteWriter) WriteUInt8(value uint8) {
	w.data = append(w.data, value)
}

func (w *ByteWriter) WriteUInt16(value uint16) {
	var valueBytes [2]byte
	binary.LittleEndian.PutUint16(valueBytes[:], value)
	w.data = append(w.data, valueBytes[:]...)
}

func (w *ByteWriter) WriteUInt16BE(value uint16) {
	var valueBytes [2]byte
	binary.BigEndian.PutUint16(valueBytes[:], value)
	w.data = append(w.data, valueBytes[:]...)
}

func (w *ByteWriter) WriteInt32(value int32) {
	w.WriteUInt32(uint32(value))
}

func (w *ByteWriter) WriteUInt32(value uint32) {
	var valueBytes [4]byte
	binary.LittleEndian.PutUint32(valueBytes[:], value)
	w.data = append(w.data, valueBytes[:]...)
}

func (w *ByteWriter) WriteInt64(value int64) {
	w.WriteUInt64(uint64(value))
}

func (w *ByteWriter) WriteUInt64(value uint64) {
	var valueBytes [8]byte
	binary.LittleEndian.PutUint64(valueBytes[:], value)
	w.data = append(w.data, valueBytes[:]...)
}

func (w *ByteWriter) WriteFixedSizeString(value string, size int) error {
	if len(value) > size {
		return fmt.Errorf("%q has a size of %d which is greater than %d", value, len(value), size)
	}

	valueBytes := make([]byte, size)
	copy(valueBytes, value)

	w.data = append(w.data, valueBytes...)

	return nil
}

func (w *ByteWriter) WriteCompactSize(value uint64) {
	switch {
	case value < 0xfd:
		w.WriteUInt8(uint8(value))

	case value <= 0xffff:
		w.WriteUInt8(0xfd)
		w.WriteUInt16(uint16(value))

	case value <= 0xffffffff:
		w.WriteUInt8(0xfe)
		w.WriteUInt32(uint32(value))

	default:
		w.WriteUInt8(0xff)
		w.WriteUInt64(value)
	}
}

func (w *ByteWriter) WriteVarSizeString(value string) {
	w.WriteCompactSize(uint64(len(value)))
	valueBytes := make([]byte, len(value))
	copy(valueBytes, value)
	w.data = append(w.data, valueBytes...)
}
