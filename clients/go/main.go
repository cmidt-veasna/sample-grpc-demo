package main

import (
	"flag"
	"fmt"
	"os"

	"google.golang.org/grpc"

	"example.com/pkg/example"
)

var address string
var port int

var command string
var filter string
var data string

func main() {
	flag.StringVar(&address, "address", "127.0.0.1", "The service server address to connect to, default address is 127.0.0.1")
	flag.IntVar(&port, "port", 8080, "The service port to connect to, default port is 8080")
	flag.StringVar(&command, "command", "list", "The command to interact with server. Valid command is list and save")
	flag.StringVar(&filter, "filter", "", "The filter work together with command list to filter the result. Filter is format as JSON string")
	flag.StringVar(&data, "data", "", "The data is data to send to server with command save. Data is format as JSON string")
	flag.Parse()

	client, err := grpc.Dial(fmt.Sprintf("%s:%d", address, port), grpc.WithInsecure())
	if err != nil {
		fmt.Println("unable to create client grpc connection error", err)
		os.Exit(1)
	}

	esc := example.NewElementServiceClient(client)

	switch command {
	case "list":
		example.List(filter, esc)

	case "save":
		example.Save(data, esc)
	}
}
