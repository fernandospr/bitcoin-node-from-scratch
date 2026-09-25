package main

import (
	"fmt"
)

const (
	MainnetMagic uint32 = 0xD9B4BEF9
	TestnetMagic uint32 = 0x0709110B
)

func BuildVersionMessage(ip string, port uint16) (Message, error) {
	addrRecv, err := NewNetworkAddress(
		0,
		ip,
		port,
	)
	if err != nil {
		return Message{}, fmt.Errorf("error creating version message: %s", err)
	}

	message, err := NewVersionMessage(
		MainnetMagic,
		70016,
		0,
		addrRecv,
		"/satoshilib:0.1/",
		0,
		false,
	)
	if err != nil {
		return Message{}, fmt.Errorf("error creating version message: %s", err)
	}
	return message, nil
}

func BuildVerackMessage() (Message, error) {
	message, err := NewVerAckMessage(
		MainnetMagic,
	)
	if err != nil {
		return Message{}, fmt.Errorf("error creating verack message: %s", err)
	}

	return message, nil
}

func BuildPingMessage(nonce uint64) (Message, error) {
	message, err := NewPingMessage(
		MainnetMagic,
		nonce,
	)
	if err != nil {
		return Message{}, fmt.Errorf("error creating ping message: %s", err)
	}

	return message, nil
}

func BuildPongMessage(nonce uint64) (Message, error) {
	message, err := NewPongMessage(
		MainnetMagic,
		nonce,
	)
	if err != nil {
		return Message{}, fmt.Errorf("error creating pong message: %s", err)
	}

	return message, nil
}
