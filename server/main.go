package main

import (
	"flag"
	"fmt"
	"google.golang.org/grpc"
	"log"
	"net"

	"github.com/Quantum12k/client-server-data-transfer-grpc/data_transfer_service"
)

const (
	DefaultServerPort = "1000"
)

func main() {
	portFlag := flag.String("port", DefaultServerPort, "server port")
	flag.Parse()

	dataTransferService := data_transfer_service.Service{}
	grpcServer := grpc.NewServer()

	data_transfer_service.RegisterDataTransferServer(grpcServer, dataTransferService)

	listener, err := net.Listen("tcp", fmt.Sprintf(":%s", *portFlag))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	fmt.Printf("Server started on port: %s\n", *portFlag)
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %s", err)
	}
}
