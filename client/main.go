package main

import (
	"context"
	"flag"
	"fmt"
	"google.golang.org/grpc"
	"log"
	"sync"
	"time"

	"github.com/Quantum12k/client-server-data-transfer-grpc/data_transfer_service"
)

const (
	DefaultServerPort            = "1000"
	DefaultCancelStreamTime      = 2
	DefaultDataReceptionInterval = 100
	DefaultBufferMaxSize         = 4
)

func main() {
	portFlag := flag.String("port", DefaultServerPort, "server port")
	cancelStreamTimeFlag := flag.Int64("cancel_stream_time", DefaultCancelStreamTime, "cancellation time in seconds")
	intervalFlag := flag.Int64("interval", DefaultDataReceptionInterval, "data reception interval in milliseconds")
	bufferMaxSizeFlag := flag.Int64("buffer", DefaultBufferMaxSize, "buffer of elements max size")
	flag.Parse()

	conn, err := grpc.Dial(fmt.Sprintf(":%s", *portFlag), grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("can not connect to server %v", err)
	}
	defer conn.Close()

	client := data_transfer_service.NewDataTransferClient(conn)

	getData(client, *cancelStreamTimeFlag, *intervalFlag, *bufferMaxSizeFlag)
}

func getData(client data_transfer_service.DataTransferClient, timeout, interval, bufferMaxSize int64) {
	ctx, cancel := context.WithCancel(context.Background())

	request := &data_transfer_service.Request{
		DataReceptionInterval: interval,
	}

	stream, err := client.GetDataStream(ctx, request)
	if err != nil {
		log.Fatalf("open stream error %v", err)
	}

	timeoutTicker := time.NewTicker(time.Duration(timeout) * time.Second)
	defer timeoutTicker.Stop()

	bufferOfElements := make([]int64, 0, bufferMaxSize)
	valueQueue := make(chan int64)

	var wg sync.WaitGroup
	wg.Add(1)
	go bufferController(bufferOfElements, valueQueue, &wg)

	for {
		select {
		case <-timeoutTicker.C:
			cancel()
			close(valueQueue)
			wg.Wait()
			return
		default:
			resp, err := stream.Recv()
			if err != nil {
				log.Fatalf("recv from stream error %v", err)
			}

			valueQueue <- resp.GetValue()
			log.Printf("Value received: %d\n", resp.GetValue())
		}
	}
}

func bufferController(buffer []int64, valueQueue chan int64, wg *sync.WaitGroup) {
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
