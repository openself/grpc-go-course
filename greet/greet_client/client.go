package main

import (
	"../greetpb"
	"context"
	"google.golang.org/grpc"
	"io"
	"log"
)

func main() {
	conn, err := grpc.Dial(":50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	defer conn.Close()

	clnt := greetpb.NewGreetServiceClient(conn)

	callUnaryGreeting(clnt)

	callServerSreamingGreeting(clnt)
}

func callUnaryGreeting(clnt greetpb.GreetServiceClient) {
	req := &greetpb.GreetRequest{
		Greeting: &greetpb.Greeting{
			FirstName: "Serjio",
			LastName:  "Borisov",
		},
	}
	res, err := clnt.Greet(context.Background(), req)
	if err != nil {
		log.Fatalf("Failed to call Greet RPC: %v", err)
	}

	log.Println("Response:", res.GetResult())
}

func callServerSreamingGreeting(clnt greetpb.GreetServiceClient) {
	req := &greetpb.GreetManyTimesRequest{
		Greeting: &greetpb.Greeting{
			FirstName: "openself",
		},
	}
	resStream, err := clnt.GreetManyTimes(context.Background(), req)
	if err != nil {
		log.Fatalf("Failed to call GreetManyTimes RPC: %v", err)
	}

	for {
		msg, err := resStream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("Error reading GreetManyTimes stream : %v", err)
		}
		log.Println("Response:", msg.GetResult())
	}

}
