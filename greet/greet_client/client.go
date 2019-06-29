package main

import (
	"../greetpb"
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

	clnt := greetpb.NewGreetServiceClient(conn)

	callUnaryGreeting(clnt)

	callServerSreamingGreeting(clnt)

	callClientSreamingGreeting(clnt)
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

func callClientSreamingGreeting(clnt greetpb.GreetServiceClient) {
	names := []string{
		"Mark",
		"Lucy",
		"John",
		"Anna",
	}
	msg := &greetpb.LongGreetRequest{
		Greeting: &greetpb.Greeting{
			FirstName: "",
		},
	}
	reqStream, err := clnt.LongGreet(context.Background())
	if err != nil {
		log.Fatalf("Failed to call LongGreet RPC: %v", err)
	}

	for _, name := range names {
		msg.Greeting.FirstName = name
		log.Println("Greet", name)
		err := reqStream.Send(msg)
		if err != nil {
			log.Fatalf("Error sending message: %v", err)
		}
		time.Sleep(100 * time.Millisecond)
	}

	resp, err := reqStream.CloseAndRecv()
	if err != nil {
		log.Fatalf("Failed to get response from LongGreet RPC: %v", err)
	}
	log.Println("Response:", resp.GetResult())
}
