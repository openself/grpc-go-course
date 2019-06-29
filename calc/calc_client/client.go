package main

import (
	"../calcpb"
	"context"
	"google.golang.org/grpc"
	"io"
	"log"
	"time"
)

func main() {
	conn, err := grpc.Dial(":50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	defer conn.Close()

	clnt := calcpb.NewCalculatorServiceClient(conn)

	callAdditionService(clnt)

	callDecompositionService(clnt)

	callAvgService(clnt)

	callMaxService(clnt)
}

func callAdditionService(clnt calcpb.CalculatorServiceClient) {
	args := []int32{2, 4, 6, 8}
	log.Printf("Add numbers: %v", args)
	req := &calcpb.CalcSumRequest{
		Args: &calcpb.AdditionArgs{
			Arg: args,
		},
	}
	res, err := clnt.CalcSum(context.Background(), req)
	if err != nil {
		log.Fatalf("Failed to call Add service: %v", err)
	}
	log.Println("Sum =", res.GetSum())
}

func callDecompositionService(clnt calcpb.CalculatorServiceClient) {
	var number int32
	number = 120
	log.Printf("Ask PDN for number: %d", number)
	req := &calcpb.CalcPNDRequest{
		Number: number,
	}
	res, err := clnt.CalcPND(context.Background(), req)
	if err != nil {
		log.Fatalf("Failed to call Add service: %v", err)
	}

	for {
		factor, err := res.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("Error reading stream: %v", err)
		}
		log.Printf("factor = %d; ", factor.GetPrimeFactor())
	}
	log.Println("The end")
}

func callAvgService(clnt calcpb.CalculatorServiceClient) {
	numbers := []int32{1, 2, 3, 4}
	msg := &calcpb.CalcAvgRequest{}
	reqStream, err := clnt.CalcAvg(context.Background())
	if err != nil {
		log.Fatalf("Failed to call CalcAvg RPC: %v", err)
	}

	for _, number := range numbers {
		msg.Number = number
		log.Println("number", number)
		err := reqStream.Send(msg)
		if err != nil {
			log.Fatalf("Error sending message: %v", err)
		}
		time.Sleep(100 * time.Millisecond)
	}

	resp, err := reqStream.CloseAndRecv()
	if err != nil {
		log.Fatalf("Failed to get response from CalcAvg RPC: %v", err)
	}
	log.Println("Avg = ", resp.GetAvg())
}

func callMaxService(clnt calcpb.CalculatorServiceClient) {
	// 1. Create a stream by invoking the client
	bdStream, err := clnt.CalcMax(context.Background())
	if err != nil {
		log.Fatalf("Failed to call CalcMax RPC: %v", err)
	}

	waitChan := make(chan struct{})

	// 2. Send a bunch of messages to the server (goroutine)
	numbers := []int32{1, 5, 3, 6, 2, 20, 1}
	msg := &calcpb.CalcMaxRequest{}
	go func() {
		for _, number := range numbers {
			msg.Number = number
			log.Println("current number", number)
			err := bdStream.Send(msg)
			if err != nil {
				log.Fatalf("Error sending message: %v", err)
			}
			time.Sleep(100 * time.Millisecond)
		}
		err := bdStream.CloseSend()
		if err != nil {
			close(waitChan)
		}
	}()

	// 3. Receive a bunch of messages from the server (goroutine)
	go func() {
		for {
			msg, err := bdStream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatalf("Error reading BiDiGreet stream : %v", err)
			}
			log.Println("Maximum =", msg.GetMax())
		}
		close(waitChan)
	}()

	// 4. Block until everything is done
	<-waitChan
}
