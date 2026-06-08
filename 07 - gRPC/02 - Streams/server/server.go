package main

import (
	calculatorpb "grpcstreams/proto/gen"
	"log"
	"net"
	"time"

	"google.golang.org/grpc"
)

type server struct {
	calculatorpb.UnimplementedCalculatorServiceServer
}

func (s *server) GenerateFibonacci(req *calculatorpb.GenerateFibonacciRequest, stream calculatorpb.CalculatorService_GenerateFibonacciServer) error {
	n := req.GetN()
	a, b := 0, 1
	for i := 0; i < int(n); i++ {
		generateFibonacciResponse := calculatorpb.GenerateFibonacciResponse{}
		generateFibonacciResponse.SetNumber(int32(a))

		if err := stream.Send(&generateFibonacciResponse); err != nil {
			return err
		}

		a, b = b, a+b
		time.Sleep(time.Millisecond * 700)
	}

	return nil
}

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalln("Failed to listen:", err)
	}

	grpcServer := grpc.NewServer()
	calculatorpb.RegisterCalculatorServiceServer(grpcServer, &server{})

	err = grpcServer.Serve(lis)
	if err != nil {
		log.Fatalln(err)
	}
}
