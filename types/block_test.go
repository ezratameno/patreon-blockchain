package types

import (
	"encoding/hex"
	"fmt"
	"testing"

	"github.com/ezratameno/pateon-bloackchain/crypto"
	"github.com/ezratameno/pateon-bloackchain/util"
	"github.com/stretchr/testify/assert"
)

func TestSignBlock(t *testing.T) {
	block := util.RandomBlock()
	privateKey := crypto.GeneratePrivateKey()
	publicKey := privateKey.Public()

	sig := SignBlock(privateKey, block)

	assert.Equal(t, 64, len(sig.Bytes()))

	assert.True(t, sig.Verify(publicKey, HashBlock(block)))

}

func TestHashBlock(t *testing.T) {

	block := util.RandomBlock()
	hash := HashBlock(block)
	fmt.Println(hex.EncodeToString(hash))
	assert.Equal(t, 32, len(hash))
}
