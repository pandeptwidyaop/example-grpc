package main

import (
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	pb "learn-grpc/greet/proto"
)

func main() {
	var address string = "0.0.0.0:9000"

	con, err := grpc.Dial(address, grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		log.Fatalf("Failed to connect: %v\n", err)
	}

	defer con.Close()

	c := pb.NewGreetServiceClient(con)

	doGreet(c)
	doGreetManyTimes(c)
}
