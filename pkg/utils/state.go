package utils

import (
	"crypto/rand"
	"encoding/base64"
	"io"
)

// NewStateString generates a new random string from rand.Reader of n length.
func NewStateString(n int) (string, error) {
	data := make([]byte, n)

	if _, err := io.ReadFull(rand.Reader, data); err != nil {
		return "", err
	}

	return base64.URLEncoding.EncodeToString(data), nil
}
