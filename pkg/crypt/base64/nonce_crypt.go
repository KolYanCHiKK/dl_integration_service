package base64

import (
	crand "crypto/rand"
	"encoding/base64"
)

func RandomNonce(n int) (string, error) {
	b := make([]byte, n)
	if _, err := crand.Read(b); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
}
