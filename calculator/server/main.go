package main

import (
	"context"
	"io"
	pb "learn-grpc/calculator/proto"
	"log"
	"math"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/status"
)

type Server struct {
	pb.CalculatorServiceServer
}

func main() {
	address := "0.0.0.0:9000"

	lis, err := net.Listen("tcp", address)

	if err != nil {
		log.Fatalf("Failed to listen: %v\n", err)
	}

	log.Printf("Listening on %s, waiting for connections...\n", address)

	cred, err := credentials.NewServerTLSFromFile("ssl/server.crt", "ssl/server.pem")
	if err != nil {
		log.Fatalf("Failed to load certificates: %v\n", err)
	}
	s := grpc.NewServer(grpc.Creds(cred))

	pb.RegisterCalculatorServiceServer(s, &Server{})

	err = s.Serve(lis)

	if err != nil {
		log.Fatalf("Failed to serve: %v\n", err)
	}

}

func (s *Server) Sum(ctx context.Context, req *pb.SumRequest) (*pb.SumResponse, error) {
	a := req.FirstNumber
	b := req.SecondNumber
	result := a + b
	res := &pb.SumResponse{
		Result: result,
	}
	return res, nil
}

func (s *Server) Prime(req *pb.PrimeRequest, stream pb.CalculatorService_PrimeServer) error {
	k := int32(2)
	n := req.Number
	for n > 1 {
		if n%k == 0 {
			res := &pb.PrimeResponse{
				Prime: k,
			}
			stream.Send(res)
			n = n / k
		} else {
			k = k + 1
		}
	}
	return nil
}

func (s *Server) Average(stream pb.CalculatorService_AverageServer) error {
	var sum float32
	var count float32
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			avg := sum / count
			return stream.SendAndClose(&pb.AverageResponse{
				Average: avg,
			})
		}
		if err != nil {
			log.Printf("Failed to receive: %v\n", err)
		}
		sum += req.Number
		count++

		log.Printf("Received: %v, sum: %v, count: %v\n", req.Number, sum, count)
	}
}

func (s *Server) Max(stream pb.CalculatorService_MaxServer) error {
	var max int32
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			log.Printf("Failed to receive: %v\n", err)
		}
		if req.Number > max {
			max = req.Number
			err = stream.Send(&pb.MaxResponse{
				Max: max,
			})
			if err != nil {
				log.Printf("Failed to send: %v\n", err)
			}
		}
	}
}

func (s *Server) Sqrt(ctx context.Context, req *pb.SqrtRequest) (*pb.SqrtResponse, error) {
	if req.Number < 0 {
		return nil, status.Errorf(codes.InvalidArgument, "Number must be positive")
	}

	return &pb.SqrtResponse{
		Sqrt: int32(math.Sqrt(float64(req.Number))),
	}, nil
}
