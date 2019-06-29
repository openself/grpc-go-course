package main

import (
	"../calcpb"
	"context"
	"google.golang.org/grpc"
	"io"
	"log"
	"net"
	"time"
)

type server struct{}

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

func (*server) CalcSum(ctx context.Context, req *calcpb.CalcSumRequest) (*calcpb.CalcSumResponse, error) {
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

func (*server) CalcAvg(stream calcpb.CalculatorService_CalcAvgServer) error {
	var total, count int32
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			// we have finished reading the client stream
			if count == 0 {
				count = 1
			}
			res := &calcpb.CalcAvgResponse{
				Avg: float32(total) / float32(count),
			}
			return stream.SendAndClose(res)
		}
		if err != nil {
			log.Fatalf("Error reading stream: %v", err)
		}
		total += req.GetNumber()
		count++
	}
	return nil
}

func (*server) CalcMax(stream calcpb.CalculatorService_CalcMaxServer) error {
	var max int32
	res := &calcpb.CalcMaxResponse{}
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			log.Fatalf("Error reading stream: %v", err)
		}
		number := req.GetNumber()
		if number <= max {
			continue
		}
		max = number
		res.Max = max
		err = stream.Send(res)
		if err != nil {
			log.Fatalf("Error sending message: %v", err)
		}
	}
	return nil
}
