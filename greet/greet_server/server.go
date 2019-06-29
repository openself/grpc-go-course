package main

import (
	"../greetpb"
	"context"
	"fmt"
	"google.golang.org/grpc"
	"io"
	"log"
	"net"
	"strings"
	"time"
)

type server struct{}

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	srv := grpc.NewServer()
	greetpb.RegisterGreetServiceServer(srv, &server{})
	log.Println("Listen on :50051")
	if err := srv.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}

}

func (*server) Greet(ctx context.Context, req *greetpb.GreetRequest) (*greetpb.GreetResponse, error) {
	firstName := req.GetGreeting().GetFirstName()
	lastName := req.GetGreeting().GetLastName()
	correctNames(&firstName, &lastName)

	log.Printf("Greet func was invoked. Params: %s %s", firstName, lastName)

	result := &greetpb.GreetResponse{
		Result: fmt.Sprintf("Hello %s %s!\n", firstName, lastName),
	}
	return result, nil
}

func (*server) GreetManyTimes(req *greetpb.GreetManyTimesRequest,
	stream greetpb.GreetService_GreetManyTimesServer) error {
	log.Println("GreetManyTimes func was invoked.")

	firstName := req.GetGreeting().GetFirstName()
	lastName := req.GetGreeting().GetLastName()
	correctNames(&firstName, &lastName)

	log.Printf("GreetManyTimes func was invoked. Params: %s %s", firstName, lastName)

	resp := &greetpb.GreetManyTimesResponse{}
	for i := 1; i < 6; i++ {
		resp.Result = fmt.Sprintf("%d Hello %s %s!\n", i, firstName, lastName)
		err := stream.Send(resp)
		if err != nil {
			log.Fatalf("Error sending message: %v", err)
		}
		time.Sleep(500 * time.Millisecond)
	}
	return nil
}

func (*server) LongGreet(stream greetpb.GreetService_LongGreetServer) error {
	log.Println("LongGreet func was invoked.")
	result := []string{}
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			// we have finished reading the client stream
			res := &greetpb.LongGreetResponse{
				Result: strings.Join(result, " "),
			}
			return stream.SendAndClose(res)
		}
		if err != nil {
			log.Fatalf("Error reading stream: %v", err)
		}
		firstName := req.GetGreeting().GetFirstName()
		lastName := req.GetGreeting().GetLastName()
		correctNames(&firstName, &lastName)

		text := fmt.Sprintf("Hello %s %s!", firstName, lastName)
		result = append(result, text)
	}
	return nil
}

func correctNames(firstName *string, lastName *string) {
	if *firstName == "" {
		*firstName = "\b"
	}
	if *lastName == "" {
		*lastName = "\b"
	}
}

func (*server) BiDiGreet(stream greetpb.GreetService_BiDiGreetServer) error {
	log.Println("BiDiGreet func was invoked.")
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			log.Fatalf("Error reading stream: %v", err)
		}
		firstName := req.GetGreeting().GetFirstName()
		lastName := req.GetGreeting().GetLastName()
		correctNames(&firstName, &lastName)

		text := fmt.Sprintf("Hello %s %s!", firstName, lastName)
		res := &greetpb.BiDiGreetResponse{
			Result: text,
		}
		err = stream.Send(res)
		if err != nil {
			log.Fatalf("Error sending message: %v", err)
		}
	}
	return nil
}
