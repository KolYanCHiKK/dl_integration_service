package ecdsa

import (
	"crypto"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"encoding/base64"
	"fmt"

	"github.com/KolYanCHiKK/dl_integration_service/pkg/crypt/aead"
	"github.com/go-jose/go-jose"
)

const AlgES256 = "ES256"

type CipherEcdsa struct {
	*aead.Cipher
}

type GeneratedSigningKey struct {
	KID           string
	Algorithm     string
	PublicJWK     []byte // JSON для jwks.json
	PrivateKeyEnc []byte // PKCS#8, зашифрованный
}

func NewJWKCipher(key []byte) (*CipherEcdsa, error) {
	cipherImplementation, err := aead.NewSessionCipher(key)
	if err != nil {
		return nil, err
	}
	return &CipherEcdsa{Cipher: cipherImplementation}, nil
}

func NewGenerateSigningKey(cipher *CipherEcdsa) (*GeneratedSigningKey, error) {
	priv, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return nil, fmt.Errorf("generate key: %w", err)
	}

	pub := jose.JSONWebKey{
		Key:       &priv.PublicKey,
		Algorithm: AlgES256,
		Use:       "sig",
	}

	// kid = отпечаток публичного ключа (RFC 7638)
	thumb, err := pub.Thumbprint(crypto.SHA256)
	if err != nil {
		return nil, fmt.Errorf("thumbprint: %w", err)
	}

	kid := base64.RawURLEncoding.EncodeToString(thumb)
	pub.KeyID = kid

	publicJWK, err := pub.MarshalJSON()
	if err != nil {
		return nil, fmt.Errorf("marshal jwk: %w", err)
	}

	der, err := x509.MarshalPKCS8PrivateKey(priv)
	if err != nil {
		return nil, fmt.Errorf("marshal private key: %w", err)
	}

	enc, err := cipher.Encrypt(der, []byte(kid)) // AAD = kid
	if err != nil {
		return nil, fmt.Errorf("encrypt private key: %w", err)
	}

	return &GeneratedSigningKey{
		KID:           kid,
		Algorithm:     AlgES256,
		PublicJWK:     publicJWK,
		PrivateKeyEnc: enc,
	}, nil
}

func DecryptSigningKey(c *CipherEcdsa, enc []byte, kid string) (*ecdsa.PrivateKey, error) {
	der, err := c.Decrypt(enc, []byte(kid))
	if err != nil {
		return nil, fmt.Errorf("decrypt private key: %w", err)
	}

	key, err := x509.ParsePKCS8PrivateKey(der)
	if err != nil {
		return nil, fmt.Errorf("parse private key: %w", err)
	}

	priv, ok := key.(*ecdsa.PrivateKey)
	if !ok {
		return nil, fmt.Errorf("unexpected key type %T", key)
	}
	return priv, nil
}
