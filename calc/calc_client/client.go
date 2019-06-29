package main

import (
	"../calcpb"
	"context"
	"google.golang.org/grpc"
	"log"
)

func main() {
	conn, err := grpc.Dial(":50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	defer conn.Close()

	clnt := calcpb.NewCalculatorServiceClient(conn)

	callAdditionService(clnt)
}

func callAdditionService(clnt calcpb.CalculatorServiceClient) {
	args := []int32{2, 4, 6, 8}
	log.Printf("Add numbers: %v", args)
	req := &calcpb.CalcRequest{
		Args: &calcpb.AdditionArgs{
			Arg: args,
		},
	}
	res, err := clnt.Add(context.Background(), req)
	if err != nil {
		log.Fatalf("Failed to call Add service: %v", err)
	}
	log.Println("Sum =", res.GetSum())
}
