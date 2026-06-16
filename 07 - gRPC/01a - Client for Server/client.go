package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"log"
	"net/http"
	"os"
	mainpb "simplegrpcclient/proto/gen"
	"simplegrpcclient/proto/gen/genconnect"
	"time"

	"connectrpc.com/connect"
	"golang.org/x/net/http2"
)

func main() {
	caCert, err := os.ReadFile("cert.pem")
	if err != nil {
		log.Fatalln("Failed to read certificate:", err)
	}

	// Add it to a cert pool
	certPool := x509.NewCertPool()
	if ok := certPool.AppendCertsFromPEM(caCert); !ok {
		log.Fatalln("Failed to append certificate")
	}

	// Build a TLS-enabled HTTP/2 client
	tlsConfig := &tls.Config{RootCAs: certPool}
	httpClient := &http.Client{
		Transport: &http2.Transport{
			TLSClientConfig: tlsConfig,
		},
	}

	client1 := genconnect.NewCalculateServiceClient(httpClient, "https://localhost:50051", connect.WithGRPC())

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	defer cancel()

	reqBody := &mainpb.AddRequest{}
	reqBody.SetA(10)
	reqBody.SetB(15)

	req := connect.NewRequest(reqBody)
	req.Header().Set("Authorization", "Bearer placeholder-token")
	req.Header().Set("X-Testing", "testing")

	res1, err := client1.Add(ctx, req)
	if err != nil {
		log.Fatalln("Error while calling Add RPC:", err)
	}

	log.Println("Received header from server: ", res1.Header())
	log.Println("Received trailer from server: ", res1.Trailer())

	fmt.Println("Response from server:", res1.Msg.GetSum())

	client2 := genconnect.NewGreeterServiceClient(httpClient, "https://localhost:50051", connect.WithGRPC())
	greeterReqBody := mainpb.GreetRequest{}
	greeterReqBody.SetName("Harsh Morayya")

	greeterReq := connect.NewRequest(&greeterReqBody)

	res2, err := client2.Greet(ctx, greeterReq)
	if err != nil {
		log.Fatalln("Error when calling Greet RPC:", err)
	}

	fmt.Println("Greet response from server:", res2.Msg.GetMessage())

	client3 := genconnect.NewAufWiedershenServiceClient(httpClient, "https://localhost:50051", connect.WithGRPC())
	farewellReqBody := mainpb.GoodByeRequest{}
	farewellReqBody.SetName("Harsh Morayya")

	farewellReq := connect.NewRequest(&farewellReqBody)

	res3, err := client3.GoodBye(ctx, farewellReq)
	if err != nil {
		log.Fatalln("Error when calling GoodBye RPC:", err)
	}

	fmt.Println("GoodBye response from server:", res3.Msg.GetMessage())
}
