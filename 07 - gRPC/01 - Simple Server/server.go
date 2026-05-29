package main

import (
	"context"
	"fmt"
	"log"
	"net"
	mainpb "simplegrpcserver/gen"

	"google.golang.org/grpc"
)

type server struct {
	mainpb.UnimplementedCalculateServiceServer
}

func (s *server) Add(ctx context.Context, req *mainpb.AddRequest) (*mainpb.AddResponse, error) {
	addResponse := &mainpb.AddResponse{}
	addResponse.SetSum(req.GetA() + req.GetB())

	return addResponse, nil
}

func main() {
	port := "50051"

	listener, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatal("Failed to listen:", err)
	}

	grpcServer := grpc.NewServer()

	mainpb.RegisterCalculateServiceServer(grpcServer, &server{})

	err = grpcServer.Serve(listener)
	if err != nil {
		log.Fatal("Failed to serve:", err)
	}

	fmt.Printf("Server is running on %v\n", port)
}
