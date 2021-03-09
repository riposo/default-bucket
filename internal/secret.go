package internal

import (
	"encoding/hex"
	"io"

	"golang.org/x/crypto/blake2b"
)

// HashSecret is used for hashing values.
type HashSecret []byte

// Decode implements envconfig.Decoder interface.
func (s *HashSecret) Decode(val string) error {
	key, err := hex.DecodeString(val)
	if err != nil {
		return err
	}

	if _, err := blake2b.New256(key); err != nil {
		return err
	}

	*s = HashSecret(key)
	return nil
}

// Encode hashes a string value.
func (s HashSecret) Encode(value string) []byte {
	h, _ := blake2b.New256([]byte(s))
	_, _ = io.WriteString(h, value)
	return h.Sum(nil)
}
