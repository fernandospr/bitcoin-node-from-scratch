package main

import "crypto/sha256"

func DoubleSha256(data []byte) [32]byte {
	firstSha256 := sha256.Sum256(data)
	return sha256.Sum256(firstSha256[:])
}
