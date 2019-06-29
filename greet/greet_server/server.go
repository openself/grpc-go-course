package main

import (
	"../greetpb"
	"context"
	"fmt"
	"google.golang.org/grpc"
	"log"
	"net"
	"time"
)

type server struct{}

func (*server) Greet(ctx context.Context, req *greetpb.GreetRequest) (*greetpb.GreetResponse, error) {
	firstName := req.GetGreeting().GetFirstName()
	lastName := req.GetGreeting().GetLastName()

	log.Printf("Greet func was invoked. Params: %s %s", firstName, lastName)

	result := &greetpb.GreetResponse{
		Result: fmt.Sprintf("Hello %s %s!\n", firstName, lastName),
	}
	return result, nil
}

func (*server) GreetManyTimes(req *greetpb.GreetManyTimesRequest,
	stream greetpb.GreetService_GreetManyTimesServer) error {

	firstName := req.GetGreeting().GetFirstName()
	lastName := req.GetGreeting().GetLastName()

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
