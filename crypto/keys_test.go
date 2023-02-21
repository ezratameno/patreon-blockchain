package crypto

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGeneratePrivateKey(t *testing.T) {

	privateKey := GeneratePrivateKey()
	assert.Equal(t, privKeyLen, len(privateKey.Bytes()))

	pubKey := privateKey.Public()
	assert.Equal(t, pubKeyLen, len(pubKey.key))
}

func TestPrivateKeyFromString(t *testing.T) {

	var (
		seed       = "2679fa26cf6eadc52d4fd629d76848babdf347f216df6df51bf96f5cff49f50f"
		privateKey = NewPrivateKeyFromString(seed)
		addressStr = "5b88b4cab1d8d9ec9041a9343bf9ff50f76c66d2"
	)

	assert.Equal(t, privKeyLen, len(privateKey.Bytes()))

	address := privateKey.Public().Address()

	assert.Equal(t, addressStr, address.String())
}

func TestPrivateKeySign(t *testing.T) {

	privateKey := GeneratePrivateKey()
	pubKey := privateKey.Public()
	msg := []byte("foo bar baz")

	sig := privateKey.Sign(msg)

	assert.True(t, sig.Verify(pubKey, msg))

	// Test with invalid msg.
	assert.False(t, sig.Verify(pubKey, []byte("foo")))

	// Test with invalid public key.
	invalidPrivateKey := GeneratePrivateKey()
	assert.False(t, sig.Verify(invalidPrivateKey.Public(), msg))
}

func TestPublicKeyToAddress(t *testing.T) {
	privateKey := GeneratePrivateKey()
	pubKey := privateKey.Public()
	address := pubKey.Address()

	assert.Equal(t, addressLen, len(address.Bytes()))
}
