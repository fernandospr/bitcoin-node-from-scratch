package main

import (
	"errors"
	"fmt"
)

var ErrUnsupportedCommand = errors.New("unsupported command")

type Message struct {
	Header  Header
	Payload any
	payload []byte
}

func (m Message) String() string {
	if m.Payload == nil {
		return fmt.Sprintf(
			"Message{Header: %v, Payload: nil}",
			m.Header,
		)
	}

	return fmt.Sprintf(
		"Message{Header: %v, Payload: %v}",
		m.Header,
		m.Payload,
	)
}

func (m Message) Serialize() []byte {
	w := NewByteWriter()

	m.Header.Serialize(w)
	w.WriteBytes(m.payload)

	return w.Bytes()
}

func NewVersionMessage(
	magic uint32,
	version int32,
	services uint64,
	addrRecv NetworkAddress,
	userAgent string,
	height int32,
	relay bool,
) (Message, error) {
	payload, err := NewVersionPayload(
		version,
		services,
		addrRecv,
		NetworkAddress{},
		userAgent,
		height,
		relay,
	)
	if err != nil {
		return Message{}, err
	}

	payloadWriter := NewByteWriter()
	payload.Serialize(payloadWriter)
	payloadBytes := payloadWriter.Bytes()

	header := NewHeader(
		magic,
		"version",
		payloadBytes,
	)
	return Message{
		Header:  header,
		Payload: payload,
		payload: payloadBytes,
	}, nil
}

func NewVerAckMessage(
	magic uint32,
) (Message, error) {
	header := NewHeader(
		magic,
		"verack",
		nil,
	)
	return Message{
		Header:  header,
		Payload: nil,
		payload: nil,
	}, nil
}

func NewPingMessage(
	magic uint32,
	nonce uint64,
) (Message, error) {
	payload, err := NewPingPayload(
		nonce,
	)
	if err != nil {
		return Message{}, err
	}

	payloadWriter := NewByteWriter()
	payload.Serialize(payloadWriter)
	payloadBytes := payloadWriter.Bytes()

	header := NewHeader(
		magic,
		"ping",
		payloadBytes,
	)
	return Message{
		Header:  header,
		Payload: payload,
		payload: payloadBytes,
	}, nil
}

func NewPongMessage(
	magic uint32,
	nonce uint64,
) (Message, error) {
	payload, err := NewPongPayload(
		nonce,
	)
	if err != nil {
		return Message{}, err
	}

	payloadWriter := NewByteWriter()
	payload.Serialize(payloadWriter)
	payloadBytes := payloadWriter.Bytes()

	header := NewHeader(
		magic,
		"pong",
		payloadBytes,
	)
	return Message{
		Header:  header,
		Payload: payload,
		payload: payloadBytes,
	}, nil
}

func Deserialize(data []byte) (Message, error) {
	r := NewByteReader(data)

	header, err := ReadHeader(r)
	if err != nil {
		return Message{}, err
	}

	switch header.Command {

	case "verack":
		payloadChecksum := Checksum(nil)
		if header.Checksum != payloadChecksum {
			return Message{}, fmt.Errorf("checksum mismatch. Header checksum is %s but checksum(payload) is %s", header.ChecksumString(), HexToString(payloadChecksum[:]))
		}
		return Message{
			Header:  header,
			Payload: nil,
			payload: nil,
		}, nil

	case "version":
		payloadBytes, err := r.ReadBytes(int(header.Length))
		if err != nil {
			return Message{}, err
		}
		payloadChecksum := Checksum(payloadBytes)
		if header.Checksum != payloadChecksum {
			return Message{}, fmt.Errorf("checksum mismatch. Header checksum is %s but checksum(payload) is %s", header.ChecksumString(), HexToString(payloadChecksum[:]))
		}
		payloadReader := NewByteReader(payloadBytes)
		payload, err := ReadVersionPayload(payloadReader)
		if err != nil {
			return Message{}, err
		}
		return Message{
			Header:  header,
			Payload: payload,
			payload: payloadBytes,
		}, nil
	case "ping":
		payloadBytes, err := r.ReadBytes(int(header.Length))
		if err != nil {
			return Message{}, err
		}
		payloadChecksum := Checksum(payloadBytes)
		if header.Checksum != payloadChecksum {
			return Message{}, fmt.Errorf("checksum mismatch. Header checksum is %s but checksum(payload) is %s", header.ChecksumString(), HexToString(payloadChecksum[:]))
		}
		payloadReader := NewByteReader(payloadBytes)
		payload, err := ReadPingPayload(payloadReader)
		if err != nil {
			return Message{}, err
		}
		return Message{
			Header:  header,
			Payload: payload,
			payload: payloadBytes,
		}, nil
	case "pong":
		payloadBytes, err := r.ReadBytes(int(header.Length))
		if err != nil {
			return Message{}, err
		}
		payloadChecksum := Checksum(payloadBytes)
		if header.Checksum != payloadChecksum {
			return Message{}, fmt.Errorf("checksum mismatch. Header checksum is %s but checksum(payload) is %s", header.ChecksumString(), HexToString(payloadChecksum[:]))
		}
		payloadReader := NewByteReader(payloadBytes)
		payload, err := ReadPongPayload(payloadReader)
		if err != nil {
			return Message{}, err
		}
		return Message{
			Header:  header,
			Payload: payload,
			payload: payloadBytes,
		}, nil

	default:
		return Message{}, fmt.Errorf(
			"%w: %q",
			ErrUnsupportedCommand,
			header.Command,
		)
	}
}
