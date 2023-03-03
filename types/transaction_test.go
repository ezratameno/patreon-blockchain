package types

import (
	"fmt"
	"testing"

	"github.com/ezratameno/pateon-bloackchain/crypto"
	"github.com/ezratameno/pateon-bloackchain/proto"
	"github.com/ezratameno/pateon-bloackchain/util"
	"github.com/stretchr/testify/assert"
)

func TestNewTransaction(t *testing.T) {

	fromPrivateKey := crypto.GeneratePrivateKey()
	fromAddress := fromPrivateKey.Public().Address().Bytes()

	toPrivateKey := crypto.GeneratePrivateKey()
	toAddress := toPrivateKey.Public().Address().Bytes()

	input := &proto.TxInput{
		PrevTxHash:   util.RandomHash(),
		PrevOutIndex: 0,
		PubicKey:     fromPrivateKey.Public().Bytes(),
	}

	// We need to spend all the amount of the transaction.
	// If we have a 100 and we want to send 5, so we send 5 to the dentation and 95 to ourself.

	output1 := &proto.TxOutput{
		Amount:  5,
		Address: toAddress,
	}

	// send back 95.
	output2 := &proto.TxOutput{
		Amount:  95,
		Address: fromAddress,
	}

	tx := &proto.Transaction{
		Version: 1,
		Inputs:  []*proto.TxInput{input},
		Outputs: []*proto.TxOutput{output1, output2},
	}

	sig := SignTransaction(fromPrivateKey, tx)

	input.Signature = sig.Bytes()

	assert.True(t, VerifyTransaction(tx))

	fmt.Printf("%+v\n", tx)

}
