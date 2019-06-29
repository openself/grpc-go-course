package main

import (
	"../calcpb"
	"context"
	"google.golang.org/grpc"
	"log"
	"net"
)

type server struct{}

func (*server) Add(ctx context.Context, req *calcpb.CalcRequest) (*calcpb.CalcResponse, error) {
	args := req.GetArgs()
	var sum, a int32
	for _, a = range args.GetArg() {
		sum += a
	}
	res := &calcpb.CalcResponse{
		Sum: sum,
	}
	return res, nil
}

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	srv := grpc.NewServer()
	calcpb.RegisterCalculatorServiceServer(srv, &server{})
	log.Println("Listen on :50051")
	if err := srv.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
