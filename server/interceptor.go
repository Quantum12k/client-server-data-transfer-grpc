package main

import (
	"context"
	"errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"log"
)

func streamInterceptor(srv interface{}, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
	if err := printCredentials(stream.Context()); err != nil {
		return err
	}

	return handler(srv, stream)
}

func printCredentials(ctx context.Context) error {
	if incMetadata, ok := metadata.FromIncomingContext(ctx); ok {
		if len(incMetadata["login"]) > 0 && len(incMetadata["password"]) > 0 &&
			incMetadata["login"][0] != "" && incMetadata["password"][0] != "" {
			log.Printf("Request from %s, password: %s", incMetadata["login"][0], incMetadata["password"][0])
			return nil
		}

		return errors.New("access denied")
	}

	return errors.New("empty metadata")
}
