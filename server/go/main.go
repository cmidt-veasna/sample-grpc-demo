package main

import (
	"flag"
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"

	"example.com/envoygrpc/pkg/example"
)

var address string
var port int

func main() {
	flag.StringVar(&address, "address", "0.0.0.0", "provide listening address, default to 0.0.0.0 listening on all interface.")
	flag.IntVar(&port, "port", 8090, "provide listening port, default to 8090")
	flag.Parse()

	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%d", address, port))
	if err != nil {
		log.Fatal(err)
	}

	server := grpc.NewServer()
	example.RegisterElementServiceServer(server, example.New())

	log.Printf("Start grpc server at %s on %d\n", address, port)
	server.Serve(lis)
}
