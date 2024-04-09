package main

import (
	"context"
	"io"
	pb "learn-grpc/greet/proto"
	"log"
)

func doGreet(c pb.GreetServiceClient) {
	req := &pb.GreetRequest{
		FirstName: "John",
	}

	res, err := c.Greet(context.Background(), req)

	if err != nil {
		log.Fatalf("Failed to call Greet: %v\n", err)
	}

	log.Printf("Response from Greet: %s\n", res.Result)
}
func doGreetManyTimes(c pb.GreetServiceClient) {
	req := &pb.GreetRequest{
		FirstName: "John",
	}

	stream, err := c.GreetManyTimes(context.Background(), req)

	if err != nil {
		log.Fatalf("Failed to call GreetManyTimes: %v\n", err)
	}

	for {
		res, err := stream.Recv()

		if err == io.EOF {
			break
		}

		if err != nil {
			log.Fatalf("Failed to receive response: %v\n", err)
		}

		log.Printf("Response from GreetManyTimes: %s\n", res.Result)
	}
}
