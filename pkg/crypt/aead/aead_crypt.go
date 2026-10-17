package aead

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"errors"
)

type Cipher struct {
	aead cipher.AEAD
}

func NewSessionCipher(key []byte) (*Cipher, error) {
	if len(key) != 32 {
		return nil, errors.New("crypt: key must be exactly 32 bytes")
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	aead, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	return &Cipher{aead: aead}, nil
}

func (c *Cipher) Encrypt(plain, aad []byte) ([]byte, error) {
	nonce := make([]byte, c.aead.NonceSize())
	if _, err := rand.Read(nonce); err != nil {
		return nil, err
	}
	return c.aead.Seal(nonce, nonce, plain, aad), nil
}

func (c *Cipher) Decrypt(data, aad []byte) ([]byte, error) {
	ns := c.aead.NonceSize()
	if len(data) < ns+c.aead.Overhead() {
		return nil, errors.New("crypt: ciphertext too short")
	}
	return c.aead.Open(nil, data[:ns], data[ns:], aad)
}
