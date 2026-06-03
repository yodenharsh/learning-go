package main

import (
	"context"
	"fmt"
	"log"
	mainpb "simplegrpcclient/gen"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	conn, err := grpc.NewClient("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalln("Failed to connect to server:", err)
	}

	defer conn.Close()

	client := mainpb.NewCalculateServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	defer cancel()

	req := mainpb.AddRequest{}
	req.SetA(10)
	req.SetB(15)

	res, err := client.Add(ctx, &req)
	if err != nil {
		log.Fatalln("Error while calling Add RPC:", err)
	}

	fmt.Println("Response from server:", res.GetSum())
}
