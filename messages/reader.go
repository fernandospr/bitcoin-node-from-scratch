package main

import (
	"encoding/binary"
	"fmt"
)

type ByteReader struct {
	data []byte
	pos  int
}

func NewByteReader(data []byte) *ByteReader {
	return &ByteReader{data: data}
}

func (r *ByteReader) remaining() int {
	return len(r.data) - r.pos
}

func (r *ByteReader) read(n int) ([]byte, error) {
	if n < 0 || r.pos+n > len(r.data) {
		return nil, fmt.Errorf(
			"unexpected end of payload: need %d bytes, have %d",
			n,
			r.remaining(),
		)
	}

	result := r.data[r.pos : r.pos+n]
	r.pos += n

	return result, nil
}

func (r *ByteReader) ReadBytes(n int) ([]byte, error) {
	return r.read(n)
}

func (r *ByteReader) ReadBool() (bool, error) {
	data, err := r.read(1)
	if err != nil {
		return false, err
	}

	return data[0] != 0x00, nil
}

func (r *ByteReader) ReadUInt8() (uint8, error) {
	data, err := r.read(1)
	if err != nil {
		return 0, err
	}

	return data[0], nil
}

func (r *ByteReader) ReadUInt16() (uint16, error) {
	data, err := r.read(2)
	if err != nil {
		return 0, err
	}

	return binary.LittleEndian.Uint16(data), nil
}

func (r *ByteReader) ReadUInt16BE() (uint16, error) {
	data, err := r.read(2)
	if err != nil {
		return 0, err
	}

	return binary.BigEndian.Uint16(data), nil
}

func (r *ByteReader) ReadInt32() (int32, error) {
	data, err := r.ReadUInt32()
	if err != nil {
		return 0, err
	}

	return int32(data), nil
}

func (r *ByteReader) ReadUInt32() (uint32, error) {
	data, err := r.read(4)
	if err != nil {
		return 0, err
	}

	return binary.LittleEndian.Uint32(data), nil
}

func (r *ByteReader) ReadInt64() (int64, error) {
	data, err := r.ReadUInt64()
	if err != nil {
		return 0, err
	}

	return int64(data), nil
}

func (r *ByteReader) ReadUInt64() (uint64, error) {
	data, err := r.read(8)
	if err != nil {
		return 0, err
	}

	return binary.LittleEndian.Uint64(data), nil
}

func (r *ByteReader) ReadCompactSize() (uint64, error) {
	n, err := r.ReadUInt8()
	if err != nil {
		return 0, err
	}

	switch n {
	case 0xfd:
		number, err := r.ReadUInt16()
		if err != nil {
			return 0, err
		}
		if number < 253 {
			return 0, fmt.Errorf(
				"non-canonical compact size: value %d encoded with 0xfd",
				number,
			)
		}
		return uint64(number), nil

	case 0xfe:
		number, err := r.ReadUInt32()
		if err != nil {
			return 0, err
		}
		if number < 65536 {
			return 0, fmt.Errorf(
				"non-canonical compact size: value %d encoded with 0xfe",
				number,
			)
		}
		return uint64(number), nil

	case 0xff:
		number, err := r.ReadUInt64()
		if err != nil {
			return 0, err
		}
		if number < 4294967296 {
			return 0, fmt.Errorf(
				"non-canonical compact size: value %d encoded with 0xff",
				number,
			)
		}
		return number, nil

	default:
		return uint64(n), nil
	}
}

func (r *ByteReader) ReadVarSizeString() (string, error) {
	size, err := r.ReadCompactSize()
	if err != nil {
		return "", err
	}
	strBytes, err := r.read(int(size))
	if err != nil {
		return "", err
	}

	return string(strBytes), nil
}

func (r *ByteReader) ReadFixedSizeString(size int) (string, error) {
	strBytes, err := r.read(size)
	if err != nil {
		return "", err
	}

	for i, b := range strBytes {
		if b == 0x00 {
			// check padding
			for _, padding := range strBytes[i:] {
				if padding != 0x00 {
					return "", fmt.Errorf(
						"invalid padding: expected 0x00, got 0x%02x",
						padding,
					)
				}
			}

			return string(strBytes[:i]), nil
		}
	}

	return string(strBytes), nil
}
