package main

import (
	"errors"
	"fmt"
	"io"
	"net"
	"strconv"
	"time"
)

type Peer struct {
	Server    string
	Conn      net.Conn
	PingNonce uint64
}

const maxPayloadSize = 4 * 1024 * 1024
const pingInterval = 1 * time.Minute

func (p *Peer) log(format string, args ...any) {
	fmt.Printf(
		"[%s] [%s] %s\n",
		time.Now().Format("2006-01-02 15:04:05.000"),
		p.Server,
		fmt.Sprintf(format, args...),
	)
}

func (p *Peer) handshake() error {
	p.log("Sending version message...")

	ip, port, err := net.SplitHostPort(p.Server)
	if err != nil {
		return fmt.Errorf("parse peer address: %w", err)
	}

	portInt, err := strconv.ParseUint(port, 10, 16)
	if err != nil {
		return fmt.Errorf("parse peer port: %w", err)
	}
	versionMessage, err := BuildVersionMessage(ip, uint16(portInt))
	if err != nil {
		return fmt.Errorf(
			"build version message: %w",
			err,
		)
	}

	if err := p.sendMessage(versionMessage); err != nil {
		return fmt.Errorf(
			"send version message: %w",
			err,
		)
	}

	p.log("Version sent ✅")

	receivedVersion := false
	receivedVerack := false
	sentVerack := false

	for !(receivedVersion && receivedVerack) {
		message, err := p.readMessage()
		if err != nil {
			return fmt.Errorf(
				"read handshake message: %w",
				err,
			)
		}

		p.log(
			"Received %s",
			message.Header.Command,
		)

		switch message.Header.Command {

		case "version":
			receivedVersion = true

			p.log(
				"Received Version: %s",
				message,
			)

			if !sentVerack {
				p.log("Sending verack message...")

				verackMessage, err := BuildVerackMessage()
				if err != nil {
					return fmt.Errorf(
						"build verack message: %w",
						err,
					)
				}

				if err := p.sendMessage(
					verackMessage,
				); err != nil {
					return fmt.Errorf(
						"send verack message: %w",
						err,
					)
				}

				p.log("Verack sent ✅")

				sentVerack = true
			}

		case "verack":
			receivedVerack = true

			p.log(
				"Received Verack: %s",
				message,
			)

		case "ping":
			p.log(
				"Received ping during handshake: %s",
				message,
			)

			if err := p.handlePing(message); err != nil {
				return err
			}

		default:
			p.log(
				"Ignoring %s during handshake",
				message.Header.Command,
			)
		}
	}

	return nil
}

func (p *Peer) messageLoop() error {
	messages := make(chan Message)
	errorsCh := make(chan error, 1)

	go func() {
		for {
			message, err := p.readMessage()
			if err != nil {
				errorsCh <- err
				return
			}

			messages <- message
		}
	}()

	ticker := time.NewTicker(pingInterval)
	defer ticker.Stop()

	for {
		select {
		case message := <-messages:
			p.log(
				"Received %s",
				message.Header.Command,
			)

			if err := p.handleMessage(message); err != nil {
				return err
			}

		case err := <-errorsCh:
			return err

		case <-ticker.C:
			if err := p.sendPing(); err != nil {
				return fmt.Errorf(
					"send ping: %w",
					err,
				)
			}
		}
	}
}

func (p *Peer) sendMessage(message Message) error {
	p.log(
		"Sending %s",
		message.Header.Command,
	)

	data := message.Serialize()

	_, err := p.Conn.Write(data)
	if err != nil {
		return fmt.Errorf(
			"send %s: %w",
			message.Header.Command,
			err,
		)
	}

	return nil
}

func (p *Peer) handlePing(message Message) error {
	pingPayload, ok := message.Payload.(PingPayload)
	if !ok {
		return fmt.Errorf("invalid ping payload")
	}

	p.log(
		"Ping nonce: %d",
		pingPayload.Nonce,
	)

	pongMessage, err := BuildPongMessage(
		pingPayload.Nonce,
	)
	if err != nil {
		return fmt.Errorf(
			"build pong message: %w",
			err,
		)
	}

	if err := p.sendMessage(pongMessage); err != nil {
		return fmt.Errorf(
			"send pong message: %w",
			err,
		)
	}

	return nil
}

func (p *Peer) handlePong(message Message) error {
	pongPayload, ok := message.Payload.(PongPayload)
	if !ok {
		return fmt.Errorf("invalid pong payload")
	}

	p.log(
		"Received pong nonce: %d",
		pongPayload.Nonce,
	)

	if p.PingNonce == 0 {
		return fmt.Errorf(
			"unexpected pong: no ping is pending",
		)
	}

	if pongPayload.Nonce != p.PingNonce {
		return fmt.Errorf(
			"unexpected pong nonce: got %d, expected %d",
			pongPayload.Nonce,
			p.PingNonce,
		)
	}

	p.log("Pong matches our ping! ✅")

	p.PingNonce = 0

	return nil
}

func (p *Peer) handleAddr(message Message) error {
	payload, ok := message.Payload.(AddrPayload)
	if !ok {
		return fmt.Errorf("invalid addr payload")
	}
	p.log("addr payload: %s", payload)
	// TODO Save payload.Addresses

	return nil
}

func (p *Peer) handleGetaddr(message Message) error {
	// TODO Build items with known addresses
	items := []AddrItemPayload{}
	message, err := BuildAddrMessage(items)
	if err != nil {
		return fmt.Errorf(
			"build addr message: %w",
			err,
		)
	}

	if err := p.sendMessage(message); err != nil {
		return fmt.Errorf(
			"send addr message: %w",
			err,
		)
	}

	return nil
}

func (p *Peer) readMessage() (Message, error) {
	for {
		headerBytes := make([]byte, 24)

		if _, err := io.ReadFull(
			p.Conn,
			headerBytes,
		); err != nil {
			return Message{}, err
		}

		r := NewByteReader(headerBytes)

		header, err := ReadHeader(r)
		if err != nil {
			return Message{}, fmt.Errorf(
				"read header: %w",
				err,
			)
		}

		p.log(
			"Received Header: %s",
			header,
		)

		if header.Length > maxPayloadSize {
			return Message{}, fmt.Errorf(
				"payload too large: %d bytes",
				header.Length,
			)
		}

		payload := make([]byte, header.Length)

		if _, err := io.ReadFull(
			p.Conn,
			payload,
		); err != nil {
			return Message{}, fmt.Errorf(
				"read payload: %w",
				err,
			)
		}

		data := make(
			[]byte,
			0,
			24+len(payload),
		)

		data = append(data, headerBytes...)
		data = append(data, payload...)

		message, err := Deserialize(data)
		if err != nil {
			if errors.Is(err, ErrUnsupportedCommand) {
				p.log(
					"Ignoring unsupported command: %s",
					header.Command,
				)

				continue
			}

			return Message{}, fmt.Errorf(
				"deserialize message: %w",
				err,
			)
		}

		return message, nil
	}
}

func (p *Peer) handleMessage(message Message) error {
	switch message.Header.Command {

	case "ping":
		return p.handlePing(message)

	case "pong":
		return p.handlePong(message)

	case "getaddr":
		return p.handleGetaddr(message)

	case "addr":
		return p.handleAddr(message)

	case "version":
		p.log("Unexpected version message")

	case "verack":
		p.log("Unexpected verack message")

	default:
		p.log(
			"Unsupported command: %s",
			message.Header.Command,
		)
	}

	return nil
}

func (p *Peer) sendPing() error {
	nonce, err := RandomNonce()
	if err != nil {
		return fmt.Errorf(
			"generate ping nonce: %w",
			err,
		)
	}

	message, err := BuildPingMessage(nonce)
	if err != nil {
		return fmt.Errorf(
			"build ping message: %w",
			err,
		)
	}

	if err := p.sendMessage(message); err != nil {
		return err
	}

	p.PingNonce = nonce

	return nil
}

func (p *Peer) sendGetaddr() error {
	message, err := BuildGetaddrMessage()
	if err != nil {
		return fmt.Errorf(
			"build getaddr message: %w",
			err,
		)
	}

	if err := p.sendMessage(message); err != nil {
		return err
	}

	return nil
}
