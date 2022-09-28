package main

import (
	"fmt"
	"net"
	"os"

	"github.com/hashicorp/go-hclog"
	protos "github.com/pmohanj/go-microservices/grpc/protos/currency"
	"github.com/pmohanj/go-microservices/grpc/server"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	log := hclog.Default()

	gs := grpc.NewServer()

	// kind of like a handler
	c := server.NewCurrency(log)

	protos.RegisterCurrencyServer(gs, c)
	// register the reflection service which allows clients to determine the methods
	// for this gRPC service
	reflection.Register(gs)

	// create a TCP socket for inbound server connections
	l, err := net.Listen("tcp", fmt.Sprintf(":%d", 9092))
	if err != nil {
		log.Error("Unable to create listener", "error", err)
		os.Exit(1)
	}

	// listen for requests
	gs.Serve(l)
}
