package main

import (
	"../greetpb"
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

	clnt := greetpb.NewGreetServiceClient(conn)

	callUnaryGreeting(clnt)
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
