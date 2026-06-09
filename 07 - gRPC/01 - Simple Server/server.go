package main

import (
	"context"
	"fmt"
	"log"
	"net"
	mainpb "simplegrpcserver/proto/gen"
	farewellpb "simplegrpcserver/proto/gen/farewell"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
)

type server struct {
	mainpb.UnimplementedCalculateServiceServer
	mainpb.UnimplementedGreeterServiceServer
	farewellpb.UnimplementedAufWiedershenServiceServer
}

func (s *server) Add(ctx context.Context, req *mainpb.AddRequest) (*mainpb.AddResponse, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		log.Println("No metadata received")
	} else {
		log.Println("Metadata: ", md)
		if val, ok := md["a uthorization"]; ok {
			log.Println("Authorization: ", val)
		}
	}

	addResponse := &mainpb.AddResponse{}
	addResponse.SetSum(req.GetA() + req.GetB())

	responseMd := metadata.Pairs("x-testing", "from-server")
	grpc.SendHeader(ctx, responseMd)

	trailerMd := metadata.Pairs("after-finishing", "done")
	grpc.SetTrailer(ctx, trailerMd)

	return addResponse, nil
}

func (s *server) Greet(ctx context.Context, req *mainpb.GreetRequest) (*mainpb.GreetResponse, error) {
	greetResponse := &mainpb.GreetResponse{}
	greetResponse.SetMessage("Hello " + req.GetName() + "!")

	return greetResponse, nil
}

func (s *server) GoodBye(ctx context.Context, req *farewellpb.GoodByeRequest) (*farewellpb.GoodByeResponse, error) {
	farewellResponse := &farewellpb.GoodByeResponse{}
	farewellResponse.SetMessage("Auf Wiedersehen " + req.GetName() + "!")

	return farewellResponse, nil
}

func main() {
	cert := "cert.pem"
	key := "key.pem"

	port := "50051"

	listener, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatal("Failed to listen:", err)
	}

	creds, err := credentials.NewServerTLSFromFile(cert, key)
	if err != nil {
		log.Fatal("Failed to load TLS credentials:", err)
	}
	grpcServer := grpc.NewServer(grpc.Creds(creds))

	mainpb.RegisterCalculateServiceServer(grpcServer, &server{})
	mainpb.RegisterGreeterServiceServer(grpcServer, &server{})
	farewellpb.RegisterAufWiedershenServiceServer(grpcServer, &server{})

	err = grpcServer.Serve(listener)
	if err != nil {
		log.Fatal("Failed to serve:", err)
	}

	fmt.Printf("Server is running on %v\n", port)
}
