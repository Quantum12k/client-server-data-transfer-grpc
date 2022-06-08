package main

import (
	"flag"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
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

	listener, err := net.Listen("tcp", fmt.Sprintf(":%s", *portFlag))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	creds, err := credentials.NewServerTLSFromFile("../certs/localhost.crt", "../certs/localhost.key")
	if err != nil {
		log.Fatalf("could not load TLS keys: %s", err)
	}

	opts := []grpc.ServerOption{grpc.Creds(creds), grpc.StreamInterceptor(streamInterceptor), grpc.UnaryInterceptor(unaryInterceptor)}

	dataTransferService := data_transfer_service.Service{}
	dataTransferService.Init()

	grpcServer := grpc.NewServer(opts...)

	data_transfer_service.RegisterDataTransferServer(grpcServer, dataTransferService)

	fmt.Printf("Server started on port: %s\n", *portFlag)
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %s", err)
	}
}
