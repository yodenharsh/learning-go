package main

import (
	"context"
	"log"
	"net/http"
	gen "simplegrpcserver/proto/gen"
	"simplegrpcserver/proto/gen/genconnect"

	"connectrpc.com/connect"
)

type server struct{}

func (s *server) Add(ctx context.Context, req *connect.Request[gen.AddRequest]) (*connect.Response[gen.AddResponse], error) {
	log.Println("Authorization:", req.Header().Get("Authorization"))
	log.Println("X-Testing:", req.Header().Get("X-Testing"))

	addResponse := &gen.AddResponse{}
	addResponse.SetSum(req.Msg.GetA() + req.Msg.GetB())

	res := connect.NewResponse(addResponse)

	res.Header().Set("x-testing", "from-server")
	res.Trailer().Set("after-finishing", "done")

	return res, nil
}

func (s *server) Greet(ctx context.Context, req *connect.Request[gen.GreetRequest]) (*connect.Response[gen.GreetResponse], error) {
	greetResponse := &gen.GreetResponse{}
	greetResponse.SetMessage("Hello " + req.Msg.GetName() + "!")

	res := connect.NewResponse(greetResponse)

	return res, nil
}

func (s *server) GoodBye(ctx context.Context, req *connect.Request[gen.GoodByeRequest]) (*connect.Response[gen.GoodByeResponse], error) {
	farewellResponse := &gen.GoodByeResponse{}
	farewellResponse.SetMessage("Auf Wiedersehen " + req.Msg.GetName() + "!")

	res := connect.NewResponse(farewellResponse)

	return res, nil
}

func main() {
	cert := "cert.pem"
	key := "key.pem"

	port := "50051"

	mux := http.NewServeMux()
	path1, handler1 := genconnect.NewGreeterServiceHandler(&server{})
	path2, handler2 := genconnect.NewCalculateServiceHandler(&server{})
	path3, handler3 := genconnect.NewAufWiedershenServiceHandler(&server{})
	mux.Handle(path2, handler2)
	mux.Handle(path1, handler1)
	mux.Handle(path3, handler3)

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	err := srv.ListenAndServeTLS(cert, key)
	if err != nil {
		log.Fatalln()
	}
}
