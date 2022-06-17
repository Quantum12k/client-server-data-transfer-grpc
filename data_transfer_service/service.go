package data_transfer_service

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"google.golang.org/grpc/examples/helloworld/helloworld"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Service struct {
	ActiveStreams *sync.Map
	helloworld.UnimplementedGreeterServer
}

func (s *Service) Init() {
	s.ActiveStreams = &sync.Map{}
}

func (s Service) GetDataStream(request *RequestStream, server DataTransfer_GetDataStreamServer) error {
	ctx, cancel := context.WithCancel(context.Background())
	srvCtx := server.Context()

	ticker := time.NewTicker(time.Duration(request.GetDataReceptionInterval()) * time.Millisecond)
	defer ticker.Stop()

	val := int64(1)

	s.ActiveStreams.Store(request.GetRequestID(), cancel)

	for {
		select {
		case <-srvCtx.Done():
			s.cancelStream(request.RequestID)
			log.Println("client disconnected")
			return nil
		case <-ctx.Done():
			log.Println("stream canceled by client")
			return nil
		case valueGenerationTime := <-ticker.C:
			if err := server.Send(&Data{
				Value: val,
				Time:  timestamppb.New(valueGenerationTime),
			}); err != nil {
				return fmt.Errorf("send data to stream, err: %s", err.Error())
			}

			log.Printf("Sended value = %d to stream\n", val)

			val++
		}
	}
}

func (s Service) StopStream(ctx context.Context, request *RequestStopStream) (*StopResponse, error) {
	reqID := request.GetRequestID()

	s.cancelStream(reqID)

	return &StopResponse{
		Msg:           "ok",
	}, nil
}

func (s Service) cancelStream(reqID string) {
	if cancelFunc, ok := s.ActiveStreams.LoadAndDelete(reqID); ok {
		if cancelFuncCasted, ok := cancelFunc.(context.CancelFunc); ok {
			log.Println("canceling stream by reqID: ", reqID)
			cancelFuncCasted()
		}
	}
}

func (s Service) mustEmbedUnimplementedDataTransferServer() {
	return
}
