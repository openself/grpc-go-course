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

	callBiDiSreamingGreeting(clnt)
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

func callBiDiSreamingGreeting(clnt greetpb.GreetServiceClient) {
	// 1. Create a stream by invoking the client
	bdStream, err := clnt.BiDiGreet(context.Background())
	if err != nil {
		log.Fatalf("Failed to call BiDiGreet RPC: %v", err)
	}

	waitChan := make(chan struct{})

	// 2. Send a bunch of messages to the server (goroutine)
	names := []string{
		"Mark",
		"Lucy",
		"John",
		"Anna",
	}
	msg := &greetpb.BiDiGreetRequest{
		Greeting: &greetpb.Greeting{
			FirstName: "",
		},
	}
	go func() {
		for _, name := range names {
			msg.Greeting.FirstName = name
			log.Println("Greet", name)
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
			log.Println("Response:", msg.GetResult())
		}
		close(waitChan)
	}()

	// 4. Block until everything is done
	<-waitChan
}
