package main

import (
	"context"
	"fmt"
	"log"
	mainpb "simplegrpcclient/proto/gen"
	farewellpb "simplegrpcclient/proto/gen/farewell"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
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

	client1 := mainpb.NewCalculateServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	defer cancel()

	ctxWithMd := metadata.AppendToOutgoingContext(ctx,
		"Authorization", "Bearer placeholder-token",
		"X-Testing", "testing")

	req := mainpb.AddRequest{}
	req.SetA(10)
	req.SetB(15)

	var resHeader, resTrailer metadata.MD
	res1, err := client1.Add(
		ctxWithMd, &req, grpc.Header(&resHeader), grpc.Trailer(&resTrailer))
	if err != nil {
		log.Fatalln("Error while calling Add RPC:", err)
	}

	log.Println("Received header from server: ", resHeader)
	log.Println("Received trailer from server: ", resTrailer)

	fmt.Println("Response from server:", res1.GetSum())

	client2 := mainpb.NewGreeterServiceClient(conn)
	greeterReq := mainpb.HelloRequest{}
	greeterReq.SetName("Harsh Morayya")

	res2, err := client2.Greet(ctxWithMd, &greeterReq)
	if err != nil {
		log.Fatalln("Error when calling Greet RPC:", err)
	}

	fmt.Println("Greet response from server:", res2.GetMessage())

	client3 := farewellpb.NewAufWiedershenServiceClient(conn)
	farewellReq := farewellpb.GoodByeRequest{}
	farewellReq.SetName("Harsh Morayya")

	res3, err := client3.GoodBye(ctxWithMd, &farewellReq)
	if err != nil {
		log.Fatalln("Error when calling GoodBye RPC:", err)
	}

	fmt.Println("GoodBye response from server:", res3.GetMessage())

	state := conn.GetState()
	log.Println("Connection state:", state)
}
