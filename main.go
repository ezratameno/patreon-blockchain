package main

import (
	"context"
	"fmt"
	"log"
	"os"

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
	makeNode(":3001", []string{})
	makeNode(":4001", []string{":3001"})

	// go func() {
	// 	for {
	// 		time.Sleep(2 * time.Second)
	// 		makeTransaction()
	// 	}

	// }()
	// err := node.Start(":3001")
	// if err != nil {
	// 	return err
	// }

	select {}
	return nil
}

func makeNode(listenAddr string, bootstrapNodes []string) *node.Node {
	n := node.NewNode()

	go n.Start(listenAddr)

	if len(bootstrapNodes) > 0 {
		err := n.BootstrapNetwork(bootstrapNodes)
		if err != nil {
			log.Fatal(err)
		}
	}

	return n
}

func makeTransaction() {
	client, err := grpc.Dial(":3001", grpc.WithInsecure())
	if err != nil {
		log.Fatal(err)
	}

	c := proto.NewNodeClient(client)

	v := &proto.Version{
		Version:    "blocker-1.0.1",
		Height:     28,
		ListenAddr: ":4001",
	}
	_, err = c.Handshake(context.TODO(), v)
	if err != nil {
		log.Fatal(err)
	}
}
