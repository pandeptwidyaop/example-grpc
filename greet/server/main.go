package main

import (
	pb "learn-grpc/greet/proto"
	"log"
	"net"

	"google.golang.org/grpc"
)

type Server struct {
	pb.GreetServiceServer
}

func main() {
	var address string = "0.0.0.0:9000"

	lis, err := net.Listen("tcp", address)

	if err != nil {
		log.Fatalf("Failed to listen: %v\n", err)
	}

	log.Printf("Listening on %s, waiting for connections...\n", address)

	s := grpc.NewServer()

	pb.RegisterGreetServiceServer(s, &Server{})

	err = s.Serve(lis)

	if err != nil {
		log.Fatalf("Failed to serve: %v\n", err)
	}

}
