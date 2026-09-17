package wireguard

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"

	"golang.org/x/crypto/curve25519"
)

// KeyLen is the byte length of a WireGuard key (Curve25519).
const KeyLen = 32

var ErrInvalidKey = errors.New("invalid wireguard key")

// KeyPair holds a base64-encoded Curve25519 key pair.
type KeyPair struct {
	PrivateKey string
	PublicKey  string
}

// GenerateKeyPair creates a new WireGuard private/public key pair.
func GenerateKeyPair() (KeyPair, error) {
	priv := make([]byte, KeyLen)
	if _, err := rand.Read(priv); err != nil {
		return KeyPair{}, fmt.Errorf("generate private key: %w", err)
	}
	clampPrivateKey(priv)

	pub, err := curve25519.X25519(priv, curve25519.Basepoint)
	if err != nil {
		return KeyPair{}, fmt.Errorf("derive public key: %w", err)
	}

	return KeyPair{
		PrivateKey: encodeKey(priv),
		PublicKey:  encodeKey(pub),
	}, nil
}

// PublicKeyFromPrivate derives the base64 public key from a base64 private key.
func PublicKeyFromPrivate(privateKey string) (string, error) {
	priv, err := decodeKey(privateKey)
	if err != nil {
		return "", err
	}

	pub, err := curve25519.X25519(priv, curve25519.Basepoint)
	if err != nil {
		return "", fmt.Errorf("derive public key: %w", err)
	}
	return encodeKey(pub), nil
}

// GeneratePresharedKey creates a random 32-byte preshared key.
func GeneratePresharedKey() (string, error) {
	psk := make([]byte, KeyLen)
	if _, err := rand.Read(psk); err != nil {
		return "", fmt.Errorf("generate preshared key: %w", err)
	}
	return encodeKey(psk), nil
}

// ValidateKey checks that the given string is a well-formed base64 WireGuard key.
func ValidateKey(key string) error {
	_, err := decodeKey(key)
	return err
}

func clampPrivateKey(key []byte) {
	key[0] &= 248
	key[31] &= 127
	key[31] |= 64
}

func encodeKey(b []byte) string {
	return base64.StdEncoding.EncodeToString(b)
}

func decodeKey(s string) ([]byte, error) {
	b, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrInvalidKey, err)
	}
	if len(b) != KeyLen {
		return nil, fmt.Errorf("%w: expected %d bytes, got %d", ErrInvalidKey, KeyLen, len(b))
	}
	return b, nil
}
