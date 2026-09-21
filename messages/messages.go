// Serializer/Deserializer
//
// Serialize examples:
// -------------------
// go run . s bool true|false
// go run . s uint8 <<number>>
// go run . s uint16 <<number>>
// go run . s uint16BE <<number>>
// go run . s int32 <<number>>
// go run . s uint32 <<number>>
// go run . s int64 <<number>>
// go run . s uint64 <<number>>
// go run . s compactsize <<number>>
// go run . s varsizestr <<str>>
// go run . s fixsizestr <<str>>
// go run . s message demoVersion
// go run . s message demoVerack
//
// Deserialize examples:
// ---------------------
// go run . d bool <<hexstr>>
// go run . d uint8 <<hexstr>>
// go run . d uint16 <<hexstr>>
// go run . d uint16BE <<hexstr>>
// go run . d int32 <<hexstr>>
// go run . d uint32 <<hexstr>>
// go run . d int64 <<hexstr>>
// go run . d uint64 <<hexstr>>
// go run . d compactsize <<hexstr>>
// go run . d varsizestr <<hexstr>>
// go run . d fixsizestr <<hexstr>>
// go run . d message <<hexstr>>
package main

import (
	"encoding/hex"
	"fmt"
	"os"
	"strconv"
)

const (
	MainnetMagic uint32 = 0xD9B4BEF9
	TestnetMagic uint32 = 0x0709110B
)

type serializer func(*ByteWriter, string) error
type deserializer func(*ByteReader) (any, error)

func demoSerializeVersionMessage() ([]byte, error) {
	addrRecv, err := NewNetworkAddress(
		0,
		"1.2.3.4",
		8333,
	)
	if err != nil {
		return nil, fmt.Errorf("error creating demo version message: %s", err)
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
		return nil, fmt.Errorf("error creating demo version message: %s", err)
	}

	data := message.Serialize()

	return data, nil
}

func demoSerializeVerackMessage() ([]byte, error) {
	message, err := NewVerAckMessage(
		MainnetMagic,
	)
	if err != nil {
		return nil, fmt.Errorf("error creating demo verack message: %s", err)
	}

	return message.Serialize(), nil
}

func serialize(args []string) {
	serializers := map[string]serializer{
		"bool": func(w *ByteWriter, value string) error {
			v, err := strconv.ParseBool(value)
			if err != nil {
				return fmt.Errorf("invalid bool: %q", value)
			}
			w.WriteBool(v)
			return nil
		},

		"uint8": func(w *ByteWriter, value string) error {
			v, err := strconv.ParseUint(value, 10, 8)
			if err != nil {
				return fmt.Errorf("invalid uint8: %q", value)
			}
			w.WriteUInt8(uint8(v))
			return nil
		},

		"uint16": func(w *ByteWriter, value string) error {
			v, err := strconv.ParseUint(value, 10, 16)
			if err != nil {
				return fmt.Errorf("invalid uint16: %q", value)
			}
			w.WriteUInt16(uint16(v))
			return nil
		},

		"uint16BE": func(w *ByteWriter, value string) error {
			v, err := strconv.ParseUint(value, 10, 16)
			if err != nil {
				return fmt.Errorf("invalid uint16: %q", value)
			}
			w.WriteUInt16BE(uint16(v))
			return nil
		},

		"int32": func(w *ByteWriter, value string) error {
			v, err := strconv.ParseInt(value, 10, 32)
			if err != nil {
				return fmt.Errorf("invalid int32: %q", value)
			}
			w.WriteInt32(int32(v))
			return nil
		},

		"uint32": func(w *ByteWriter, value string) error {
			v, err := strconv.ParseUint(value, 10, 32)
			if err != nil {
				return fmt.Errorf("invalid uint32: %q", value)
			}
			w.WriteUInt32(uint32(v))
			return nil
		},

		"int64": func(w *ByteWriter, value string) error {
			v, err := strconv.ParseInt(value, 10, 64)
			if err != nil {
				return fmt.Errorf("invalid int64: %q", value)
			}
			w.WriteInt64(v)
			return nil
		},

		"uint64": func(w *ByteWriter, value string) error {
			v, err := strconv.ParseUint(value, 10, 64)
			if err != nil {
				return fmt.Errorf("invalid uint64: %q", value)
			}
			w.WriteUInt64(v)
			return nil
		},

		"compactsize": func(w *ByteWriter, value string) error {
			v, err := strconv.ParseUint(value, 10, 64)
			if err != nil {
				return fmt.Errorf("invalid number: %q", value)
			}
			w.WriteCompactSize(v)
			return nil
		},

		"varsizestr": func(w *ByteWriter, value string) error {
			w.WriteVarSizeString(value)
			return nil
		},

		"fixsizestr": func(w *ByteWriter, value string) error {
			w.WriteFixedSizeString(value, 12)
			return nil
		},

		"message": func(w *ByteWriter, value string) error {
			switch value {
			case "demoVersion":
				bytes, err := demoSerializeVersionMessage()
				if err != nil {
					return err
				}
				w.WriteBytes(bytes)
				return nil
			case "demoVerack":
				bytes, err := demoSerializeVerackMessage()
				if err != nil {
					return err
				}
				w.WriteBytes(bytes)
				return nil
			default:
				return fmt.Errorf("invalid message: %q", value)
			}
		},
	}

	if len(args) < 3 {
		fmt.Println("Missing arguments")
		os.Exit(1)
	}

	serialize, ok := serializers[args[1]]
	if !ok {
		fmt.Printf("Unknown type: %q\n", args[1])
		os.Exit(1)
	}

	w := NewByteWriter()

	if err := serialize(w, args[2]); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Printf("%x\n", w.Bytes())
}

func deserialize(args []string) {
	var deserializers = map[string]deserializer{
		"bool": func(r *ByteReader) (any, error) {
			return r.ReadBool()
		},

		"uint8": func(r *ByteReader) (any, error) {
			return r.ReadUInt8()
		},

		"uint16": func(r *ByteReader) (any, error) {
			return r.ReadUInt16()
		},

		"uint16BE": func(r *ByteReader) (any, error) {
			return r.ReadUInt16BE()
		},

		"int32": func(r *ByteReader) (any, error) {
			return r.ReadInt32()
		},

		"uint32": func(r *ByteReader) (any, error) {
			return r.ReadUInt32()
		},

		"int64": func(r *ByteReader) (any, error) {
			return r.ReadInt64()
		},

		"uint64": func(r *ByteReader) (any, error) {
			return r.ReadUInt64()
		},
		"compactsize": func(r *ByteReader) (any, error) {
			return r.ReadCompactSize()
		},
		"varsizestr": func(r *ByteReader) (any, error) {
			return r.ReadVarSizeString()
		},
		"fixsizestr": func(r *ByteReader) (any, error) {
			return r.ReadFixedSizeString(12)
		},
		"message": func(r *ByteReader) (any, error) {
			return Deserialize(r.data)
		},
	}

	if len(args) < 3 {
		fmt.Println("Missing arguments")
		os.Exit(1)
	}

	deserialize, ok := deserializers[args[1]]
	if !ok {
		fmt.Printf("Unknown type: %s\n", args[1])
		os.Exit(1)
	}

	data, err := hex.DecodeString(args[2])
	if err != nil {
		fmt.Printf("Invalid hex: %v\n", err)
		os.Exit(1)
	}

	r := NewByteReader(data)

	value, err := deserialize(r)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Printf("%v\n", value)
}

func main() {
	args := os.Args
	if len(args) < 2 {
		fmt.Println("Missing arguments")
		os.Exit(1)
	}
	switch args[1] {
	case "s":
		serialize(args[1:])
	case "d":
		deserialize(args[1:])
	default:
		fmt.Printf("Unrecognized command %q\n", args[1])
		os.Exit(1)
	}
}
