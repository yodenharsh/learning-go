package main

import (
	"context"
	calculatorpb "grpcstreamsclient/proto/gen"
	"io"
	"log"

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
	req.SetN(20)

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

}
