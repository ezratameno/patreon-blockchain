package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"time"

	"github.com/ezratameno/pateon-bloackchain/node"
	"github.com/ezratameno/pateon-bloackchain/proto"
	"google.golang.org/grpc"
)

func main() {
	if err := run(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func run() error {
	node := node.NewNode()

	opts := []grpc.ServerOption{}
	grpcServer := grpc.NewServer(opts...)
	ln, err := net.Listen("tcp", ":3001")
	if err != nil {
		return err
	}

	// register the node so we can serve it.
	proto.RegisterNodeServer(grpcServer, node)

	fmt.Println("node running on port:", ":3001")

	go func() {
		for {
			time.Sleep(2 * time.Second)
			makeTransaction()
		}

	}()
	err = grpcServer.Serve(ln)
	if err != nil {
		return err
	}
	return nil
}

func makeTransaction() {
	client, err := grpc.Dial(":3001", grpc.WithInsecure())
	if err != nil {
		log.Fatal(err)
	}

	c := proto.NewNodeClient(client)

	v := &proto.Version{
		Version: "blocker-1.0.1",
		Height:  28,
	}
	_, err = c.Handshake(context.TODO(), v)
	if err != nil {
		log.Fatal(err)
	}
}
