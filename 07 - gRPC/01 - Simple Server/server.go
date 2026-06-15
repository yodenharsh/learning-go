package main

import (
	"context"
	"log"
	"net/http"
	gen "simplegrpcserver/proto/gen"
	"simplegrpcserver/proto/gen/genconnect"

	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type server struct{}

func (s *server) Add(ctx context.Context, req *gen.AddRequest) (*gen.AddResponse, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		log.Println("No metadata received")
	} else {
		log.Println("Metadata: ", md)
		if val, ok := md["a uthorization"]; ok {
			log.Println("Authorization: ", val)
		}
	}

	addResponse := &gen.AddResponse{}
	addResponse.SetSum(req.GetA() + req.GetB())

	responseMd := metadata.Pairs("x-testing", "from-server")
	grpc.SendHeader(ctx, responseMd)

	trailerMd := metadata.Pairs("after-finishing", "done")
	grpc.SetTrailer(ctx, trailerMd)

	return addResponse, nil
}

func (s *server) Greet(ctx context.Context, req *gen.GreetRequest) (*gen.GreetResponse, error) {
	greetResponse := &gen.GreetResponse{}
	greetResponse.SetMessage("Hello " + req.GetName() + "!")

	return greetResponse, nil
}

func (s *server) GoodBye(ctx context.Context, req *gen.GoodByeRequest) (*gen.GoodByeResponse, error) {
	farewellResponse := &gen.GoodByeResponse{}
	farewellResponse.SetMessage("Auf Wiedersehen " + req.GetName() + "!")

	return farewellResponse, nil
}

func main() {
	cert := "cert.pem"
	key := "key.pem"

	port := "50051"

	mux := http.NewServeMux()
	path, handler := genconnect.NewGreeterServiceHandler(&server{})
	mux.Handle(path, handler)

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: h2c.NewHandler(mux, &http2.Server{}),
	}

	srv.ListenAndServeTLS(cert, key)
}
