package main

import (
	"bufio"
	"fmt"
	calculatorpb "grpcstreams/proto/gen"
	"io"
	"log"
	"net"
	"os"
	"strings"
	"time"

	"google.golang.org/grpc"
	_ "google.golang.org/grpc/encoding/gzip"
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

func (s *server) Chat(stream calculatorpb.CalculatorService_ChatServer) error {
	reader := bufio.NewReader(os.Stdin)

	for {
		// Receiving
		msg, err := stream.Recv()
		if err == io.EOF {
			log.Println("Completing chat")
			return nil
		} else if err != nil {
			return err
		}

		log.Println("Message received: ", msg)

		fmt.Println("Enter response: ")
		str, err := reader.ReadString('\n')
		if err != nil {
			return err
		}

		str = strings.TrimSpace(str)

		// Sending message through the stream
		msgToSend := calculatorpb.ChatMessage{}
		msgToSend.SetMessage(str)
		err = stream.Send(&msgToSend)

		if err != nil {
			return err
		}
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
