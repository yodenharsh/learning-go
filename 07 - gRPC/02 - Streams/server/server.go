package main

import (
	calculatorpb "grpcstreams/proto/gen"
	"io"
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
		time.Sleep(time.Millisecond * 200)
	}

	return nil
}

func (s *server) SendNumbers(stream calculatorpb.CalculatorService_SendNumbersServer) error {
	var sum, count int32

	for {

		req, err := stream.Recv()
		if err == io.EOF {
			res := calculatorpb.SendNumbersResponse{}
			res.SetNumber(count)
			res.SetSum(sum)
			return stream.SendAndClose(&res)
		} else if err != nil {
			return err
		}
		log.Println(req.GetNumber())

		count++
		sum += req.GetNumber()
	}
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
