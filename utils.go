package fetcher

import (
	"encoding/hex"
	"hash"
	"strings"

	"golang.org/x/crypto/sha3"
)

func decodeString(s string) ([]byte, error) {
	if strings.HasPrefix(s, "0x") {
		s = s[2:]
	}
	return hex.DecodeString(s)
}

func encodeToString(b []byte) string {
	return "0x" + hex.EncodeToString(b)
}

func keccak256Hex(data []byte) string {
	hash := keccak256(data)
	return encodeToString(hash)
}

func keccak256(data ...[]byte) []byte {
	type KeccakState interface {
		hash.Hash
		Read([]byte) (int, error)
	}
	b := make([]byte, 32)
	d := sha3.NewLegacyKeccak256().(KeccakState)
	for _, b := range data {
		d.Write(b)
	}
	d.Read(b)
	return b
}
