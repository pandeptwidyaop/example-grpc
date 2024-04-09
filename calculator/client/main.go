package main

import (
	"context"
	"io"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	pb "learn-grpc/calculator/proto"
)

func main() {
	address := "0.0.0.0:9000"

	cred, err := credentials.NewClientTLSFromFile("ssl/ca.crt", "localhost")
	if err != nil {
		log.Fatalf("Failed to load certificates: %v\n", err)
	}

	conn, err := grpc.Dial(address, grpc.WithTransportCredentials(cred))

	if err != nil {
		log.Fatalf("Failed to connect: %v\n", err)
	}

	defer conn.Close()

	c := pb.NewCalculatorServiceClient(conn)

	// doSum(c)
	// doPrime(c)
	// doAverage(c)
	doMax(c)
}

func doSum(c pb.CalculatorServiceClient) {
	req := &pb.SumRequest{
		FirstNumber:  10,
		SecondNumber: 20,
	}

	res, err := c.Sum(context.Background(), req)

	if err != nil {
		log.Fatalf("Failed to call Sum: %v\n", err)
	}

	log.Printf("Response from Sum: %v\n", res)
}

func doPrime(c pb.CalculatorServiceClient) {
	req := &pb.PrimeRequest{
		Number: 120,
	}
	resStream, err := c.Prime(context.Background(), req)
	if err != nil {
		log.Fatalf("Failed to call Prime: %v\n", err)
	}
	for {
		msg, err := resStream.Recv()

		if err == io.EOF {
			break
		}

		if err != nil {
			break
		}
		log.Printf("Response from Prime: %v\n", msg.GetPrime())
	}
}

func doAverage(c pb.CalculatorServiceClient) {

	stream, err := c.Average(context.Background())

	if err != nil {
		log.Fatalf("Failed to call Average: %v\n", err)
	}

	for i := 0; i < 1000; i++ {
		req := &pb.AverageRequest{
			Number: float32(i),
		}

		stream.Send(req)
		log.Println("Sent request: ", req)
		time.Sleep(1 * time.Second)
	}

	res, err := stream.CloseAndRecv()

	if err != nil {
		log.Fatalf("Failed to receive response from Average: %v\n", err)
	}

	log.Printf("Response from Average: %v\n", res)

}

func doMax(c pb.CalculatorServiceClient) {
	stream, err := c.Max(context.Background())

	if err != nil {
		log.Fatalf("Failed to call Max: %v\n", err)
	}

	waitc := make(chan struct{})

	go func() {
		for {
			res, err := stream.Recv()

			if err == io.EOF {
				break
			}

			if err != nil {
				log.Fatalf("Failed to receive response from Max: %v\n", err)
			}

			log.Printf("Response from Max: %v\n", res.GetMax())
		}
		close(waitc)
	}()

	numbers := []int32{1, 5, 3, 6, 2, 20}

	for _, number := range numbers {
		req := &pb.MaxRequest{
			Number: number,
		}

		stream.Send(req)
		log.Println("Sent request: ", req)
		time.Sleep(1 * time.Second)
	}

	stream.CloseSend()

	<-waitc
}
