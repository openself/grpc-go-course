package main

import (
	"../calcpb"
	"context"
	"google.golang.org/grpc"
	"log"
	"net"
	"time"
)

type server struct{}

func (*server) Add(ctx context.Context, req *calcpb.CalcSumRequest) (*calcpb.CalcSumResponse, error) {
	args := req.GetArgs()
	var sum, a int32
	for _, a = range args.GetArg() {
		sum += a
	}
	res := &calcpb.CalcSumResponse{
		Sum: sum,
	}
	return res, nil
}

func (*server) CalcPND(req *calcpb.CalcPNDRequest,
	stream calcpb.CalculatorService_CalcPNDServer) error {
	number := req.GetNumber()
	log.Printf("PrimeNumberDecomposition func was invoked. Param: %d", number)

	resp := &calcpb.CalcPNDResponse{}
	var k int32
	k = 2
	for number > 1 {
		if number%k == 0 {
			resp.PrimeFactor = k // this is a factor
			err := stream.Send(resp)
			if err != nil {
				log.Fatalf("Error sending message: %v", err)
			}
			time.Sleep(500 * time.Millisecond)

			number /= k // divide N by k so that we have the rest of the number left.
			continue
		} // if k evenly divides into N
		k = k + 1
	}
	return nil
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
