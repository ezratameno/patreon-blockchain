package node

import (
	"context"
	"fmt"

	"github.com/ezratameno/pateon-bloackchain/proto"
	"google.golang.org/grpc/peer"
)

type Node struct {
	version string

	// UnimplementedNodeServer must be embedded to have forward compatible implementations.
	proto.UnimplementedNodeServer
}

func NewNode() *Node {
	return &Node{
		version: "blocker-0.1",
	}
}

func (n *Node) HandleTransaction(ctx context.Context, tx *proto.Transaction) (*proto.Ack, error) {
	peer, _ := peer.FromContext(ctx)
	fmt.Printf("received tx from: %+v\n", peer.Addr)
	return &proto.Ack{}, nil
}

// Handshake determines if we will receive this connection or not.
func (n *Node) Handshake(ctx context.Context, v *proto.Version) (*proto.Version, error) {
	ourVersion := &proto.Version{
		Version: n.version,
		Height:  100,
	}

	peer, _ := peer.FromContext(ctx)

	fmt.Printf("received version from %s: %+v \n", peer.Addr, v)
	return ourVersion, nil
}
