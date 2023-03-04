package node

import (
	"context"
	"fmt"
	"net"
	"sync"

	"github.com/ezratameno/pateon-bloackchain/proto"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/peer"
)

type Node struct {
	version    string
	listenAddr string
	logger     *zap.SugaredLogger

	peerLock sync.RWMutex
	peers    map[proto.NodeClient]*proto.Version

	// UnimplementedNodeServer must be embedded to have forward compatible implementations.
	proto.UnimplementedNodeServer
}

func NewNode() *Node {

	loggerConfig := zap.NewDevelopmentConfig()
	loggerConfig.EncoderConfig.TimeKey = ""
	logger, _ := loggerConfig.Build()

	return &Node{
		version:  "blocker-0.1",
		peers:    make(map[proto.NodeClient]*proto.Version),
		peerLock: sync.RWMutex{},
		logger:   logger.Sugar(),
	}
}

// Start starts the grpc server.
func (n *Node) Start(listenAddr string, bootstrapNodes []string) error {

	n.listenAddr = listenAddr
	opts := []grpc.ServerOption{}
	grpcServer := grpc.NewServer(opts...)

	ln, err := net.Listen("tcp", listenAddr)
	if err != nil {
		return err
	}

	// register the node so we can serve it.
	proto.RegisterNodeServer(grpcServer, n)
	n.logger.Infof("node running on port: %s", listenAddr)

	// bootstrap the network with a list of already know nodes
	// in the network.

	go n.bootstrapNetwork(bootstrapNodes)

	err = grpcServer.Serve(ln)
	if err != nil {
		return err
	}

	return nil
}

func (n *Node) HandleTransaction(ctx context.Context, tx *proto.Transaction) (*proto.Ack, error) {
	peer, _ := peer.FromContext(ctx)
	n.logger.Debugf("received tx from: %+v", peer.Addr)
	return &proto.Ack{}, nil
}

// Handshake determines if we will receive this connection or not.
func (n *Node) Handshake(ctx context.Context, v *proto.Version) (*proto.Version, error) {

	// Create a client from the sender listen address.
	c, err := makeNodeClient(v.ListenAddr)
	if err != nil {
		return nil, err
	}

	// Add the peer.
	n.addPeer(c, v)

	n.logger.Debugf("[%s] received version from %s", n.listenAddr, v.ListenAddr)
	return n.getVersion(), nil
}

func (n *Node) addPeer(c proto.NodeClient, v *proto.Version) {
	n.peerLock.Lock()
	defer n.peerLock.Unlock()

	// Handle the logic where we decide to accept or drop
	// the incoming node connection.

	n.peers[c] = v

	// When we connect to a node we want to connect to all his known peers.
	go n.bootstrapNetwork(v.PeerList)

	n.logger.Debugf("[%s] new peer connected (%s) - height (%d)", n.listenAddr, v.ListenAddr, v.Height)

}

func (n *Node) deletePeer(c proto.NodeClient) {
	n.peerLock.Lock()
	defer n.peerLock.Unlock()
	delete(n.peers, c)

}

// BootstrapNetwork bootstraps the network by adding nodes to each other.
func (n *Node) bootstrapNetwork(addrs []string) error {

	for _, addr := range addrs {
		if !n.canConnectWith(addr) {
			continue
		}

		n.logger.Debugf("[%s] dialing remote node %s", n.listenAddr, addr)

		c, v, err := n.dialRemoteNode(addr)
		if err != nil {
			return err
		}

		n.addPeer(c, v)
	}
	return nil
}

func (n *Node) dialRemoteNode(addr string) (proto.NodeClient, *proto.Version, error) {
	c, err := makeNodeClient(addr)
	if err != nil {
		return nil, nil, err
	}

	v, err := c.Handshake(context.Background(), n.getVersion())
	if err != nil {
		return nil, nil, fmt.Errorf("handshake error: %w", err)
	}

	return c, v, nil
}

// canConnectWith checks if we already connected.
func (n *Node) canConnectWith(addr string) bool {
	if n.listenAddr == addr {
		return false
	}
	connectedPeers := n.getPeerList()
	for _, connectedAddr := range connectedPeers {
		if addr == connectedAddr {
			return false
		}
	}
	return true

}
func (n *Node) getVersion() *proto.Version {

	return &proto.Version{
		Version:    "blocker-0.1",
		Height:     0,
		ListenAddr: n.listenAddr,
		PeerList:   n.getPeerList(),
	}
}

// getPeerList returns the listen address for the peers.
func (n *Node) getPeerList() []string {
	n.peerLock.Lock()
	defer n.peerLock.Unlock()

	peers := []string{}

	for _, version := range n.peers {
		peers = append(peers, version.ListenAddr)
	}

	return peers
}

func makeNodeClient(listenAddr string) (proto.NodeClient, error) {
	c, err := grpc.Dial(listenAddr, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}

	return proto.NewNodeClient(c), nil
}
