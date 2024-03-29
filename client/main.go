package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	"github.com/Quantum12k/client-server-data-transfer-grpc/data_transfer_service"
)

const (
	DefaultServerPort            = "1000"
	DefaultCancelStreamTime      = 2
	DefaultDataReceptionInterval = 100
	DefaultBufferMaxSize         = 4
	DefaultLogin                 = "example login"
	DefaultPassword              = "example password"
)

//структура получаемых данных с сервера
type data struct {
	Value     int64
	Timestamp time.Time
}

func main() {
	portFlag := flag.String("port", DefaultServerPort, "server port")
	cancelStreamTimeFlag := flag.Int64("cancel_stream_time", DefaultCancelStreamTime, "cancellation time in seconds")
	intervalFlag := flag.Int64("interval", DefaultDataReceptionInterval, "data reception interval in milliseconds")
	bufferMaxSizeFlag := flag.Int64("buffer", DefaultBufferMaxSize, "buffer of elements max size")
	loginFlag := flag.String("login", DefaultLogin, "login")
	passwordFlag := flag.String("password", DefaultPassword, "password")
	flag.Parse()

	creds, err := credentials.NewClientTLSFromFile("../certs/localhost.crt", "")
	if err != nil {
		log.Fatalf("could not load tls cert: %s", err)
	}

	conn, err := grpc.Dial(
		fmt.Sprintf(":%s", *portFlag),
		grpc.WithTransportCredentials(creds),
		grpc.WithPerRPCCredentials(&Authentication{
			Login:    *loginFlag,
			Password: *passwordFlag,
		}))
	if err != nil {
		log.Fatalf("can not connect to server %v", err)
	}
	defer conn.Close()

	client := data_transfer_service.NewDataTransferClient(conn)

	getData(client, *cancelStreamTimeFlag, *intervalFlag, *bufferMaxSizeFlag)
}

func getData(client data_transfer_service.DataTransferClient, timeout, interval, bufferMaxSize int64) {
	ctx := context.Background()

	reqID := uuid.New().String()

	request := &data_transfer_service.RequestStream{
		RequestID:             reqID,
		DataReceptionInterval: interval,
	}

	stream, err := client.GetDataStream(ctx, request)
	if err != nil {
		log.Fatalf("open stream error %v", err)
	}

	timeoutTicker := time.NewTicker(time.Duration(timeout) * time.Second)
	defer timeoutTicker.Stop()

	dataQueue := make(chan data)

	var wg sync.WaitGroup
	wg.Add(1)
	go printData(bufferMaxSize, dataQueue, &wg)

	for {
		select {
		case <-timeoutTicker.C:
			resp, err := client.StopStream(ctx, &data_transfer_service.RequestStopStream{RequestID: reqID})
			if err != nil {
				log.Fatalf("stop stream request error %v", err)
			}

			log.Println("Response to cancel request:", resp.Msg)

			close(dataQueue)

			wg.Wait()

			return
		default:
			resp, err := stream.Recv()
			if err != nil {
				log.Fatalf("recv from stream error %v", err)
			}

			dataQueue <- data{
				Value:     resp.GetValue(),
				Timestamp: resp.GetTime().AsTime(),
			}
		}
	}
}

func printData(bufferMaxSize int64, valueQueue chan data, wg *sync.WaitGroup) {
	buffer := make([]data, 0, bufferMaxSize)

	defer wg.Done()

	for value := range valueQueue {
		buffer = append(buffer, value)

		if len(buffer) == cap(buffer) {
			log.Printf("Buffer full, elements to print: %v, (memory address: %p)\n", buffer, &buffer)
			buffer = buffer[:0]
		}
	}

	log.Printf("Remaining buffer, elements to print: %v, (memory address: %p)\n", buffer, &buffer)
}
