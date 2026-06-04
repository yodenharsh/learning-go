package main

import (
	"context"
	"fmt"
	"log"
	mainpb "simplegrpcclient/proto/gen"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

func main() {
	cert := "cert.pem"
	creds, err := credentials.NewClientTLSFromFile(cert, "")
	if err != nil {
		log.Fatalln("Failed to load TLS credentials:", err)
		return
	}

	conn, err := grpc.NewClient("localhost:50051", grpc.WithTransportCredentials(creds))
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
	state := conn.GetState()
	log.Println("Connection state:", state)
}
