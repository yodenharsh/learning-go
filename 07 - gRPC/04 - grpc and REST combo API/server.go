package main

import (
	"comboapi/proto/gen"
	"comboapi/proto/gen/genconnect"
	"context"
	"log"
	"net"
	"net/http"
	"strings"

	"connectrpc.com/connect"
	"connectrpc.com/grpcreflect"
	"connectrpc.com/vanguard"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

type AnimalServiceHandler struct{}

func (AnimalServiceHandler) GetAnimalCry(ctx context.Context, req *connect.Request[gen.GetAnimalCryRequest]) (*connect.Response[gen.GetAnimalCryResponse], error) {
	message := strings.Builder{}
	var animalCry string

	switch req.Msg.GetAnimalType() {
	case gen.AnimalType_ANIMAL_TYPE_CAT:
		animalCry = "meow"
	case gen.AnimalType_ANIMAL_TYPE_DOG:
		animalCry = "bark"
	case gen.AnimalType_ANIMAL_TYPE_DRAGON:
		animalCry = "wrrr"
	default:
		animalCry = "unknown"
	}

	for range req.Msg.GetCount() {
		message.WriteString(animalCry + " ")
	}

	trimmed := strings.TrimSpace(message.String())
	animalCryResponse := gen.GetAnimalCryResponse{}
	animalCryResponse.SetResponse(trimmed)

	return connect.NewResponse(&animalCryResponse), nil
}

func main() {
	animalService := vanguard.NewService(genconnect.NewAnimalServiceHandler(&AnimalServiceHandler{}))
	transcoder, err := vanguard.NewTranscoder([]*vanguard.Service{animalService})

	if err != nil {
		log.Fatalln(err)
		return
	}

	lis, err := net.Listen("tcp", ":50051")

	if err != nil {
		log.Fatalln(err)
		return
	}

	mux := http.NewServeMux()
	reflector := grpcreflect.NewStaticReflector(genconnect.AnimalServiceName)

	mux.Handle(grpcreflect.NewHandlerV1(reflector))
	mux.Handle(grpcreflect.NewHandlerV1Alpha(reflector))
	mux.Handle("/", transcoder)
	err = http.Serve(lis, h2c.NewHandler(mux, &http2.Server{}))

	if err != nil {
		log.Fatalln(err)
	}
}
