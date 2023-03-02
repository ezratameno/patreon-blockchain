package types

import (
	"crypto/sha256"

	"github.com/ezratameno/pateon-bloackchain/crypto"
	"github.com/ezratameno/pateon-bloackchain/proto"
	pb "github.com/golang/protobuf/proto"
)

// SignBlock - each time a validator commits a block he need to sign it,
// so other can know it's a valid block.
func SignBlock(pk *crypto.PrivateKey, block *proto.Block) *crypto.Signature {

	return pk.Sign(HashBlock(block))

}

// HashBlock returns a SHA256 of the header.
func HashBlock(block *proto.Block) []byte {
	b, err := pb.Marshal(block)
	if err != nil {
		panic(err)
	}

	hash := sha256.Sum256(b)

	return hash[:]
}
