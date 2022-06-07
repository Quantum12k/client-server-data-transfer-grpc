package data_transfer_service

import (
	"fmt"
	"log"
	"time"

	"google.golang.org/grpc/examples/helloworld/helloworld"
)

type Service struct {
	helloworld.UnimplementedGreeterServer
}

func (s Service) GetDataStream(request *Request, server DataTransfer_GetDataStreamServer) error {
	ctx := server.Context()

	ticker := time.NewTicker(time.Duration(request.DataReceptionInterval) * time.Millisecond)
	defer ticker.Stop()

	val := int64(1)

	for {
		select {
		case <-ctx.Done():
			{
				log.Println("context canceled")
				return nil
			}
		case <-ticker.C:
			if err := server.Send(&Data{
				Value: val,
			}); err != nil {
				return fmt.Errorf("send data to stream, err: %s", err.Error())
			}

			log.Printf("Sended value = %d to stream\n", val)

			val++
		}
	}
}

func (s Service) mustEmbedUnimplementedDataTransferServer() {
	return
}
