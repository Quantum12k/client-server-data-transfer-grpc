package data_transfer_service

import (
	"context"
	"fmt"
	"google.golang.org/protobuf/types/known/timestamppb"
	"log"
	"sync"
	"time"

	"google.golang.org/grpc/examples/helloworld/helloworld"
)

type Service struct {
	ActiveStreams *sync.Map
	helloworld.UnimplementedGreeterServer
}

func (s *Service) PrintMap() {
	log.Print("map: ")
	s.ActiveStreams.Range(func(k,v interface{})bool{
		log.Print(k.(string))
		return true
	})

	log.Print("\n")
}

func (s *Service) Init() {
	s.ActiveStreams = &sync.Map{}
}

func (s Service) GetDataStream(request *RequestStream, server DataTransfer_GetDataStreamServer) error {
	ctx, cancel := context.WithCancel(context.Background())

	ticker := time.NewTicker(time.Duration(request.GetDataReceptionInterval()) * time.Millisecond)
	defer ticker.Stop()

	val := int64(1)

	s.ActiveStreams.Store(request.GetRequestID(), cancel)

	for {
		select {
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

	cancelFunc, ok := s.ActiveStreams.LoadAndDelete(reqID)
	if ok {
		cancelFuncCasted, ok := cancelFunc.(context.CancelFunc)
		if ok {
			cancelFuncCasted()
		}
	}

	return &StopResponse{
		Msg:           "ok",
	}, nil
}

func (s Service) mustEmbedUnimplementedDataTransferServer() {
	return
}
