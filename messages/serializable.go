package main

type Serializable interface {
	Serialize(w *ByteWriter)
}
