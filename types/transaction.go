package types

import (
	"crypto/sha256"

	"github.com/ezratameno/pateon-bloackchain/crypto"
	"github.com/ezratameno/pateon-bloackchain/proto"
	pb "github.com/golang/protobuf/proto"
)

func SignTransaction(pk *crypto.PrivateKey, tx *proto.Transaction) *crypto.Signature {
	return pk.Sign(HashTransaction(tx))
}

func HashTransaction(tx *proto.Transaction) []byte {
	b, err := pb.Marshal(tx)
	if err != nil {
		panic(err)
	}

	hash := sha256.Sum256(b)
	return hash[:]

}

// VerifyTransaction ensures all the inputs are verified.
// The output is only valid if all the inputs are valid.
func VerifyTransaction(tx *proto.Transaction) bool {
	for _, input := range tx.Inputs {

		sig := crypto.SignatureFromBytes(input.Signature)
		pubKey := crypto.PublicKeyFromBytes(input.PubicKey)

		// TODO: make sure we don't run into problems after verification because we set the signature to nil.
		input.Signature = nil

		if !sig.Verify(pubKey, HashTransaction(tx)) {
			return false
		}

	}

	return true
}
