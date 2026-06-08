package main

import (
	"context"
	calculatorpb "grpcstreamsclient/proto/gen"
	"io"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	conn, err := grpc.NewClient("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		log.Fatalln("error occurred when trying to make grpc client", err)
	}
	defer conn.Close()

	client := calculatorpb.NewCalculatorServiceClient(conn)
	ctx := context.Background()

	req := &calculatorpb.GenerateFibonacciRequest{}
	req.SetN(10)

	stream, err := client.GenerateFibonacci(ctx, req)
	if err != nil {
		log.Fatalln("error occurred when calling GenerateFibonacci RPC", err)
	}

	for {
		res, err := stream.Recv()
		if err == io.EOF {
			log.Println("stream ended")
			break
		} else if err != nil {
			log.Fatalln("error occurred while receiving response from stream", err)
		}

		log.Println("Res: ", res.GetNumber())
	}

	sendNumbersStream, err := client.SendNumbers(ctx)
	if err != nil {
		log.Fatalln("error occurred when calling SendNumbers RPC", err)
	}

	for i := range 5 {
		req := &calculatorpb.SendNumbersRequest{}
		req.SetNumber(int32(i + 10))

		err = sendNumbersStream.Send(req)
		if err != nil {
			log.Fatalln("Error when sending stream: ", err)
		}
		time.Sleep(300 * time.Millisecond)
		log.Println("Sent number: ", req.GetNumber())
	}

	res, err := sendNumbersStream.CloseAndRecv()
	if err != nil {
		log.Fatalln("Error when closing and receiving: ", err)
	}

	log.Println("Sum: ", res.GetSum())
	log.Println("Count: ", res.GetNumber())
}
