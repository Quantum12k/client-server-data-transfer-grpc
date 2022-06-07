package main

import (
	"context"
	"flag"
	"fmt"
	"google.golang.org/grpc"
	"log"
	"time"

	"github.com/Quantum12k/client-server-data-transfer-grpc/data_transfer_service"
)

const (
	DefaultServerPort            = "1000"
	DefaultCancelStreamTime      = 1
	DefaultDataReceptionInterval = 250
)

func main() {
	portFlag := flag.String("port", DefaultServerPort, "server port")
	cancelStreamTimeFlag := flag.Int64("cancel_stream_time", DefaultCancelStreamTime, "cancellation time in seconds")
	dataReceptionIntervalFlag := flag.Int64("data_reception_interval", DefaultDataReceptionInterval, "data reception interval in milliseconds")
	flag.Parse()

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(*cancelStreamTimeFlag)*time.Second)
	defer cancel()

	conn, err := grpc.DialContext(ctx, fmt.Sprintf(":%s", *portFlag), grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("can not connect to server %v", err)
	}
	defer conn.Close()

	client := data_transfer_service.NewDataTransferClient(conn)

	request := &data_transfer_service.Request{
		DataReceptionInterval: *dataReceptionIntervalFlag,
	}

	stream, err := client.GetDataStream(context.Background(), request)
	if err != nil {
		log.Fatalf("open stream error %v", err)
	}

	for {
		select {
		case <-ctx.Done():
			cancel()
			log.Println("context canceled")
			return
		default:
			resp, err := stream.Recv()
			if err != nil {
				log.Fatalf("recv from stream error %v", err)
			}

			go saveValue()
			log.Printf("Value received: %d\n", resp.Value)
		}
	}
}

func saveValue() {

}
