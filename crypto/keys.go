package crypto

import (
	"crypto/ed25519"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
)

const (
	// the first 32 bytes of the private key will be the private key, the remaining 32 will be the public key.
	privKeyLen = 64
	pubKeyLen  = 32
	seedLen    = 32

	addressLen = 20
)

// PrivateKey will be used to sign the transactions.
type PrivateKey struct {
	key ed25519.PrivateKey
}

func NewPrivateKeyFromString(s string) *PrivateKey {

	b, err := hex.DecodeString(s)
	if err != nil {
		panic(err)
	}

	return NewPrivateKeyFromSeed(b)
}

func NewPrivateKeyFromSeed(seed []byte) *PrivateKey {
	if len(seed) != seedLen {
		panic(fmt.Sprintf("invalid seed length, must be %d", seedLen))
	}

	return &PrivateKey{
		key: ed25519.NewKeyFromSeed(seed),
	}
}

// GeneratePrivateKey will generate a random private key.
func GeneratePrivateKey() *PrivateKey {

	seed := make([]byte, seedLen)
	_, err := io.ReadFull(rand.Reader, seed)
	if err != nil {
		panic(err)
	}

	return &PrivateKey{
		key: ed25519.NewKeyFromSeed(seed),
	}
}

func (p *PrivateKey) Bytes() []byte {
	return p.key
}

// Sign will sign the transaction.
func (p *PrivateKey) Sign(msg []byte) *Signature {
	return &Signature{
		value: ed25519.Sign(p.key, msg),
	}
}

func (p *PrivateKey) Public() *PublicKey {
	b := make([]byte, pubKeyLen)

	// copy the public key into the buffer.
	copy(b, p.key[32:])

	return &PublicKey{
		key: b,
	}
}

type PublicKey struct {
	key ed25519.PublicKey
}

func (p *PublicKey) Address() Address {
	// create an address from the last 20 bytes of the key.
	return Address{
		value: p.key[len(p.key)-addressLen:],
	}
}

type Signature struct {
	value []byte
}

func (s *Signature) Bytes() []byte {
	return s.value
}

// Verify will assert the the signature is valid.
// We sign a msg with a private key and we verify with a public key.
func (s *Signature) Verify(pubKey *PublicKey, msg []byte) bool {
	return ed25519.Verify(pubKey.key, msg, s.value)
}

type Address struct {
	value []byte
}

func (a *Address) Bytes() []byte {
	return a.value
}

func (a *Address) String() string {
	return hex.EncodeToString(a.value)
}
